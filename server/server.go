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
		// load templates.
		app.RegisterView(iris.HTML("./views", ".html"))

		app.Get("/alist-service", func(context context.Context) {
			context.ServeFile("./workers/alist-service.worker.js", false)
		})
		// render the ./views/index.html.
		app.Get("/", func(ctx iris.Context) {
			ctx.View("index.html")
		})

		mvc.Configure(app.Party("/websocket"), configureMVC)
		// Or mvc.New(app.Party(...)).Configure(configureMVC)

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
	// Note that you could use an anonymous field as well, it doesn't matter, binder will find it.
	//
	// This is the current websocket connection, each client has its own instance of the *websocketController.
	Conn websocket.Connection
}

func (c *websocketController) onLeave(roomName string) {
	// visits--
	newCount := decrement()
	// This will call the "visit" event on all clients, except the current one,
	// (it can't because it's left but for any case use this type of design)
	c.Conn.To(websocket.Broadcast).Emit("visit", newCount)
}

func (c *websocketController) update() {
	// visits++
	newCount := increment()

	// This will call the "visit" event on all clients, including the current
	// with the 'newCount' variable.
	//
	// There are many ways that u can do it and faster, for example u can just send a new visitor
	// and client can increment itself, but here we are just "showcasing" the websocket controller.
	c.Conn.To(websocket.All).Emit("visit", newCount)
}

func (c *websocketController) Get( /* websocket.Connection could be lived here as well, it doesn't matter */ ) {
	c.Conn.OnLeave(c.onLeave)
	c.Conn.On("visit", c.update)

	c.Conn.On("get_service_worker", func() {

	})
	// call it after all event callbacks registration.
	c.Conn.Wait()
}
