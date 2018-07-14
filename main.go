package main

import (
	"fmt"
	"a-list-music/server"
	"github.com/kataras/iris"
	"a-list-music/store"
	"a-list-music/transcoder"
	"a-list-music/utilities"
)

const ServerHost = string("localhost:8888")

func main() {

	StartServer()
	Nexus()
}

func StartServer() {
	fmt.Println("starting Server")
		serv := server.BuildServer()
		serv.Logger().Info("on connection")
		fmt.Println("starting server on ", ServerHost)
		err := serv.Run(iris.TLS(ServerHost, "./alist.cert", "./alist.key"))
		if err != nil {
			utilities.ErrorHandler(err)
		}
}

// This will coordinate the different modules into routines.
func Nexus() {
	for job := range store.Client().Jobs {
		fmt.Println(job)
	}

	for file := range server.Client().FileUploaded {
		fmt.Println(file)
	}

	for transcodes := range transcoder.Client().Jobs {
		fmt.Println(transcodes)
	}
}