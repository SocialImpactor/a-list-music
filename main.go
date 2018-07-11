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
	soundFile, err := os.Open("sound-files/demo-sound/101358__edge-of-october__distress-signal.wav")
	//soundFile, err := os.Open("434626__gis-sweden__electronic-minute-no-136-comparing-x-with-y-left-sloth-1.wav")
	if err != nil {
		panic(err)
	}
	val := &soundFile
	fmt.Println(val)
	defer soundFile.Close()
	// initialize transcoder channel
	tClient := transcoder.TranscoderClient{}
	go transcoder.BuildTranscodeClient(&tClient)

	if meta, err := tClient.NewJob(soundFile, "mp3"); err != nil {
	}

}
func StartServer() {
	fmt.Println("starting Server")
	server := server.BuildServer()
	server.Run(iris.Addr("localhost:2821"))
}

