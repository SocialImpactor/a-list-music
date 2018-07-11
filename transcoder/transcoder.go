package transcoder

import (
	"os"
	"os/exec"
	"github.com/kjk/betterguid"
	"net/http"
	"fmt"
	"io"
	"github.com/kataras/iris/core/errors"
	"github.com/kataras/iris/core/router"
	"strings"
)

var EncExtMap = map[string] string {
	"audio/wave": "wav",
}

type SoundFileMeta struct {
	id string
	name string
	uri string
	encoding string
	codex string
	size int
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
	ffmpegCMD *exec.Cmd
}


type TranscoderClient struct {
	TranscodeJobs chan TranscodeJob
}

//type TranscodesQueue struct {
//
//}

func InitSoundLib() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "",  err
	}
	libDir := string(cwd + "/sound-files" )
	if !router.DirectoryExists(libDir) {
		err := os.MkdirAll(libDir, 0744)
		if err != nil {
			return  libDir, err
		}
	}
	return libDir, nil
}

func BuildTranscoderClient(transcoderClient *TranscoderClient) {
	jobs := make(chan TranscodeJob)

	client := TranscoderClient{TranscodeJobs: jobs}
	if transcoderClient != nil {
		transcoderClient = &client
	}
}

func (c *TranscoderClient) NewJob(_file *os.File, targetEncode ...string) {
	fmt.Println("Adding New Job")
	// TODO Clean up..

	buffer := make([]byte, 1024)
	result := SoundFileMeta{}
	id := betterguid.New()

	// build result
	result.id = id
	cwd, err := os.Getwd()
	if err != nil {
		c.exitChan() <- err
	}


	if result.encoding, err = DetectEncoding(_file); err != nil {
		c.exitChan() <- err
	}
	result.name = string(id + "." + result.encoding)

	dir := string(cwd + "/sound-files" + "/" + id + "/source" + "/" + result.encoding + "/")
	result.uri = dir + result.name

	// set source file and folders
	err = os.MkdirAll(dir, 744)
	if err != nil {
		c.exitChan() <- err
	}
	_newFile, err := os.Create(result.uri)
	if err != nil  {
		c.exitChan() <- err
	}
	defer _newFile.Close()

	// write to source file
	for {
		n, err := _file.Read(buffer)
		if err != nil && err != io.EOF {
			c.exitChan() <- err
		}

		if _, err := _newFile.Write(buffer[:n]); err != nil {
			c.exitChan() <- err
		}

		if n == 0 {
			break
		}
	}

	// build FFMPEG CMD
	encodingCount := len(targetEncode)
	for i := 0;  i < encodingCount; i++ {
		cmd := buildFFMPEGCMD(result, targetEncode[i])
		job := TranscodeJob{
			id: result.id,
			ready: true,
			done: false,
			sourceMeta: result,
			targetMeta: SoundFileMeta{},
			ffmpegCMD: cmd,
		}
		fmt.Println("New Job Success?", job)
		c.TranscodeJobs <- job
	}
	fmt.Println("Closing")
	close(c.TranscodeJobs)
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

func (c *TranscoderClient) RunJobs() (map[string] SoundFileMeta, error) {
	result := make(map[string] SoundFileMeta)
	for jobs := range c.TranscodeJobs {
		_out := os.Stdout
		jobs.ffmpegCMD.Stdout = _out
		jobs.ffmpegCMD.Run()
		println("STDOUT", _out)
	}

	return result, nil
}

func (c *TranscoderClient) exitChan() chan error {
	return c.exitChan()
}


func DetectEncoding(_file *os.File) (string, error) {
	testBuffer := make([]byte, 512)
	n, err := _file.Read(testBuffer)
	if err != nil {
		return "", err
	}
	encoding := http.DetectContentType(testBuffer[:n])
	fmt.Println("encoding is", encoding)
	if EncExtMap[encoding] == "" {
		return "", errors.New("Encoding not indexed" + encoding)
	}
	return EncExtMap[encoding], nil
}

func buildFFMPEGCMD(sourceMeta SoundFileMeta, targetEncode string) *exec.Cmd {
	switch strings.ToLower(targetEncode) {
	case "mp3":
		return exec.Command(
			"ffmpeg",
			"-i",
			sourceMeta.uri,
			"-vn",
			"-ar 44100",
			"-ac 2",
			"-ab 192k",
			"-f mp3",
			sourceMeta.id + ".mp3",
		)
	default:
		return nil
	}
}
