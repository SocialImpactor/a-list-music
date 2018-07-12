package main

import (
	"fmt"
	"github.com/kataras/iris"
	"a-list-music/server"
	"a-list-music/transcoder"
)

func main() {
	fmt.Println("calling Transcoder")

	// Transcoder Client

	tClient := initTranscoder()

	fmt.Println(tClient)

	// Library Manager

	// ServerHandlers

	StartServer()
}

func demoSoundTranscode(tclient transcoder.TranscoderClient) {
	readyJobs := make(chan map[string] transcoder.TranscodeJob)
	// transcoded := make(chan map[string] transcoder.TranscodeJob)
	tclient.ReadyTranscodes = readyJobs
}

func initTranscoder() transcoder.TranscoderClient{
	transcodeClient := transcoder.TranscoderClient{}
	jobs := make(chan transcoder.TranscodeJob)
	transcodeClient.Jobs = jobs
	go transcoder.SetClient(&transcodeClient)
	go transcodeClient.ProcessJobs()
	return transcodeClient
}

func buildTranscoderClient() transcoder.TranscoderClient {
	transcoder.InitSoundLib()
	tClient := transcoder.TranscoderClient{}
	go transcoder.SetClient(&tClient)
	return tClient
}

func StartServer() {
	fmt.Println("starting Server")
	server := server.BuildServer()
	server.Run(iris.Addr("localhost:2824"))
}

