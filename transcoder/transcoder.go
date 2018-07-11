package transcoder

import (
	"os"
	"os/exec"
	"github.com/kjk/betterguid"
	"net/http"
	"fmt"
	"io"
)

type SoundFileMeta struct {
	id string
	uri string
	encoding string
	codex string
	size int
}

func buildFFMPEGCMD(sourceMeta SoundFileMeta) *exec.Cmd {
	return exec.Command("ffmpeg", "-i", sourceMeta.uri, "-vn", "-ar 44100", "-ac 2", "-ab 192k", "-f mp3", sourceMeta.id + ".mp3")
}

type Transcode interface {
	StoreToMediaLibrary()(SoundFileMeta, err error)
	NewJob(file os.File, targetMime []string)(source SoundFileMeta)
	RunTranscode()(data byte, err error)
	exitChan() chan error
}


type TranscodeJob struct {
	id string
	ready bool
	done bool
	sourceMeta SoundFileMeta
	targetMeta SoundFileMeta
	ffmpegCMD exec.Cmd
}

//type TranscodesQueue struct {
//
//}

type TranscoderClient struct {
	TranscodeJobs chan map[string]TranscodeJob
	TranscodeC chan Transcode
}


func BuildTranscodeClient(transcoderClient *TranscoderClient) (TranscoderClient) {
	jobs := make(chan map[string] TranscodeJob)
	transcoderC := make(chan Transcode)
	client := TranscoderClient{TranscodeJobs: jobs, TranscodeC: transcoderC}
	if transcoderClient != nil {
		transcoderClient = &client
	}
	return client
}


//func TransStore() {
//	initializeFFMPEG()

	//exec.Command("ffmpeg", )
	// catch STDOUT
//}
//

//func (c *TranscoderClient) StoreToMediaLibrary()(SoundFileMeta, err error) {
//	todo...
//}

func (c *TranscoderClient) NewJob(_file *os.File, targetMime ...string)(SoundFileMeta, error) {
	buffer := make([]byte, 1024)
	testBuffer := make([]byte, 512)
	result := SoundFileMeta{}
	id := betterguid.New()
	result.id = id
	fmt.Println("Starting to build meta")
	cwd, err := os.Getwd()
	if err != nil {
		return result, err
	}
	libDir := string(cwd + "/sound-files" )

	n, err := _file.Read(testBuffer)
	encoding := http.DetectContentType(testBuffer[:n])
	fmt.Println("encoding is", EncExtMap[encoding])

	fileExt := string(EncExtMap[encoding])
	fileName := string("/" + id + "." + fileExt)

	err = os.Chmod(libDir, 755)
	if err != nil {
		return result, err
	}
	dir := string(cwd + "/sound-files" + "/" + id + "/source" + "/" + fileExt + "/")
	result.uri = dir + fileName

	err = os.MkdirAll(dir, 755)
	if err != nil {
		return result, err
	}

	_newFile, err := os.Create(result.uri)
	if err != nil  {
		fmt.Println(err, _newFile)
		return result, err
	}

	fmt.Println("about to for")
	for {
		n, err := _file.Read(buffer)
		if err != nil && err != io.EOF {
			return result, err
		}

		if _, err := _newFile.Write(buffer[:n]); err != nil {
			return result, err
		}

	}
	return result, nil
}

var EncExtMap = map[string] string {
	"audio/wave": "wav",
}


func (c *TranscoderClient) RunJobs() (map[string] SoundFileMeta, error) {
	result := make(map[string] SoundFileMeta)
	//transcodedFile := make([]byte, 512)
	//for job :=  range c.TranscodesJobs {
	//	_out := os.Stdout
	//	job.ffmpegCMD.Stdout = _out
	//
	//	fmt.Printf("stdout %s", _out)
	//	job.ffmpegCMD.Run()
		//result[job.id] = SoundFileMeta{}
	//}
	//
	return result, nil
}

func (c *TranscoderClient) exitChan() chan error {
	return c.exitChan()
}
