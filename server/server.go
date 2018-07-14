package server

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
	"sync/atomic"
	"fmt"
	"github.com/kataras/iris/context"
	"a-list-music/utilities"
	"net/http"
)

var fileUploadedChan chan utilities.Action
var Client = func() AListServerClient {
	return AListServerClient{
		FileUploaded: fileUploadedChan,
	}
}

type AListServer interface {
	BuildServer() (server iris.Application)
}

type AListServerClient struct {
	FileUploaded chan utilities.Action
}

func BuildServer() (server *iris.Application) {
	fileUploadedChan = make(chan utilities.Action)
	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html"))

	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	app.Get("/style-sheet", func(ctx context.Context) {
		ctx.ServeFile("./views/main.style.css", false)
	})

	app.Get("/alist-service", func(ctx context.Context) {
		ctx.ServeFile("./workers/alist-service.worker.js", false)
	})

	app.Get("/admin", func(ctx context.Context) {
		ctx.View("admin.html")
	})

	mvc.Configure(app.Party("/websocket"), configureMVC)

	return app
}

func configureMVC(m *mvc.Application) {
	ws := websocket.New(websocket.Config{
		CheckOrigin: func(r *http.Request) bool {
			fmt.Println(r)
			return true
		},
		IDGenerator: func(ctx context.Context) string {
			var count= int(0)
			var name= "ClientID" + string(count+1)
			fmt.Println(name)
			return name
		},
	})
	// http://localhost:8080/websocket/iris-ws.js
	m.Router.Any("/iris-ws.js", websocket.ClientHandler())

	// This will bind the result of ws.Upgrade which is a websocket.Connection
	// to the controller(s) served by the `m.Handle`.
	m.Register(ws.Upgrade)
	m.Handle(new(websocketController))
}

var visits uint64

func increment() uint64 {
	return atomic.AddUint64(&visits, 1)
}

func decrement() uint64 {
	return atomic.AddUint64(&visits, ^uint64(0))
}

type websocketController struct {
	Conn websocket.Connection
}

func (c *websocketController) onLeave(roomName string) {
	newCount := decrement()
	c.Conn.To(websocket.Broadcast).Emit("visit", newCount)
}

func (c *websocketController) update() {
	newCount := increment()
	c.Conn.Join("sample")
	c.Conn.To(websocket.All).Emit("visit", newCount)
}

func (c *websocketController) Get( /* websocket.Connection could be lived here as well, it doesn't matter */ ) {
	c.Conn.OnLeave(c.onLeave)
	c.Conn.On("visit", c.update)
	c.fileSockets()
	c.Conn.Wait()
}

func(c *websocketController) fileSockets() {
	const RoomName = "file_upload"

	if c.Conn.Join(RoomName); c.Conn.IsJoined(RoomName) {
		var onUpload = func(data []byte) {
			fmt.Println("file received", data)
		}
		c.Conn.On("upload", onUpload)
	}
	c.Conn.Emit("FileUpload::Done", nil)
}