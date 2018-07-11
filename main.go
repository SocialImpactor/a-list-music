package main

import (
	"fmt"
	"github.com/kataras/iris"
	"a-list/server"
	"a-list/transcoder"
	"os"
	)

func main() {
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println("calling Transcoder")
	tClient := buildTranscoderClient()
	demoSoundTranscode(tClient)
	// main thread ends with Server.Run
	StartServer()

}
func demoSoundTranscode(tclient transcoder.TranscoderClient) {
	jobs := make(chan transcoder.TranscodeJob)
	tclient.TranscodeJobs = jobs
	if soundFile, err := os.Open("sound-files/demo-sound/18210__roil-noise__circuitbent-casio-ctk-550-loop1.wav"); err == nil {
		go tclient.NewJob(soundFile, "mp3")
		readyJobs := <- jobs
		fmt.Println("running jobs", readyJobs)
		go tclient.RunJobs()
	}  else {
		panic(err)
	}

}
func buildTranscoderClient() transcoder.TranscoderClient {
	transcoder.InitSoundLib()
	tClient := transcoder.TranscoderClient{}
	go transcoder.BuildTranscoderClient(&tClient)
	return tClient
}

func StartServer() {
	fmt.Println("starting Server")
	server := server.BuildServer()
	server.Run(iris.Addr("localhost:2821"))
}

