package main

import (
	"fmt"
	"github.com/kataras/iris"
	"a-list-music/server"
	"a-list-music/transcoder"
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
	readyJobs := make(chan map[string] transcoder.TranscodeJob)
	transcoded := make(chan map[string] transcoder.TranscodeJob)
	tclient.ReadyTranscodes = readyJobs

	if soundFile, err := os.Open("sound-files/sound-demo/gtr-nylon22.wav"); err == nil {

		go tclient.NewJob(soundFile, "mp3")

		jmap := <- readyJobs
		for key, val := range jmap {
			fmt.Println(key, val)
		}
		fmt.Println("running readyJobs", )
		tclient.Transcoded = transcoded

		go tclient.RunTranscodes(jmap)

		done := <- transcoded
		println(done)
	}  else {
		panic(err)
	}

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
	server.Run(iris.Addr("localhost:2820"))
}

