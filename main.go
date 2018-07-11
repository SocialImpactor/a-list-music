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

	go transcodeSomething()
	// main thread ends with Server.Run
	StartServer()

}

func transcodeSomething() {
	//soundFile, err := os.Open("sound-files/demo-sound/101358__edge-of-october__distress-signal.wav")
	soundFile, err := os.Open("434602__matdiffusion__crowd-noises-the-puppets.wav")
	if err != nil {
		panic(err)
	}
	val := &soundFile
	fmt.Println(val)
	defer soundFile.Close()
	// initialize transcoder channel
	tClient := transcoder.TranscoderClient{}
	go transcoder.BuildTranscodeClient(&tClient)

	meta, err := tClient.NewJob(soundFile, "mp3")

	if err != nil {
		panic(err)
	}

	fmt.Println(meta)

}
func StartServer() {
	fmt.Println("starting Server")
	server := server.BuildServer()
	server.Run(iris.Addr("localhost:2820"))
}

