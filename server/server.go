package server

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/websocket"
	"sync/atomic"
	"fmt"
	"github.com/kataras/iris/context"
)

type AListServer interface {
	BuildServer()(server iris.Application)
}

func BuildServer() (server *iris.Application){
	fmt.Println("starting")

	app := iris.New()

	// load templates

	app.RegisterView(iris.HTML("./views", ".html"))

	app.Get("/style-sheet", func(ctx context.Context) {
		ctx.ServeFile("./views/main.style.css", false)
	})

	app.Get("/alist-service", func(ctx context.Context) {
		ctx.ServeFile("./workers/alist-service.worker.js", false)
	})

	app.Get("/admin", func(ctx context.Context) {
		ctx.View("admin.html")
	})

	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	mvc.Configure(app.Party("/websocket"), configureMVC)
	// Or
	//mvc.New(app.Party(...)).Configure(configureMVC)
	// http://localhost:8080
	return app
}

func configureMVC(m *mvc.Application) {
	ws := websocket.New(websocket.Config{})
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
	c.Conn.To(websocket.All).Emit("visit", newCount)
}

func (c *websocketController) Get( /* websocket.Connection could be lived here as well, it doesn't matter */ ) {
	c.Conn.OnLeave(c.onLeave)
	c.Conn.On("visit", c.update)
	c.Conn.Wait()
}