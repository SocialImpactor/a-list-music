package main

import (
	"fmt"
	"a-list-music/transcoder"
	"a-list-music/store"
	"a-list-music/utilities"
	"github.com/kataras/iris/websocket"
	"os"
	"a-list-music/server"
	"github.com/kataras/iris"
)

type Client struct {
	*transcoder.TranscodeClient
	*store.StoreClient
}

func main() {
	var st = string("calling AListTranscoder")
	fmt.Println(st)

	// AListTranscoder Client
	_transcoderClient := initTranscoderClient()
	_storeClient := initStoreClient()

	fmt.Println(_transcoderClient)
	fmt.Println(_storeClient)
	// ServerHandlers

	StartServer()
}

func initTranscoderClient() transcoder.TranscodeClient {
	action := make(chan utilities.Action)
	transcoded := make(map[string] transcoder.TranscodeJob)
	_t := &transcoded
	transcodeClient := transcoder.TranscodeClient{ Transcoded: _t, }
	transcodeClient.Jobs = action
	go transcoder.SetClient(&transcodeClient)
	go transcodeClient.ProcessJobs()
	return transcodeClient
}

func initStoreClient() *store.StoreClient {
	store.InitSoundLib()
	jobs :=  make(chan utilities.Action)
	client := store.StoreClient{Jobs: jobs}
	_client := &client
	store.SetClient(_client)
	client.ProcessJobs()
	return &client
}

func fileSockets(socket websocket.Connection, client transcoder.TranscodeClient) {
	const RoomName = "file_upload"

	if socket.Join(RoomName); socket.IsJoined(RoomName) {
			var onUpload = func(data []byte) {
			fmt.Println("file received", data)

			// todo Form Data

			_file, err := os.Create("temp")
			utilities.ErrorHandler(err)
			_file.Write(data)
			client.MakeTranscodeJob( _file, "mp3")
		}
		socket.On("upload", onUpload)
	}
	socket.Emit("FileUpload::Done", nil)
}

func StartServer() {
	fmt.Println("starting Server")
	server := server.BuildServer()
	server.Run(iris.Addr("localhost:9004"))
}

