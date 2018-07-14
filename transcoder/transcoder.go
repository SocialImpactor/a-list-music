package transcoder

import (
	"os"
	"os/exec"
	"net/http"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"strings"
	"bytes"
	"path"
	"github.com/kjk/betterguid"
	"io"
	"a-list-music/utilities"
)

var FFMPEGPath = string(path.Join(utilities.CWD(), "..", "..", "..", "Desktop", "ffmpeg", "ffmpeg"))

var EncExtMap = map[string] string {
	"audio/wave": "wav",
}

var Client = func() TranscodeClient {
	return TranscodeClient{}
}

type AListTranscoder interface {

	// meta is for the library IA, it relates to URLs
	MetaBuilder(file *os.File) (SoundFileMeta)

	// New Jobs set the source file...
	NewJob(file *os.File, targetMime []string)

	ExitChan() chan error
}

type SoundFileMeta struct {
	Id        string
	Name      string
	URI       string
	BaseDir   string
	SourceDir string
	Encoding  string
	Codex     string
	Size      int
}

type TranscodeJob struct {
	Id         string
	Ready      bool
	Done       bool
	SourceMeta SoundFileMeta
	TargetMeta SoundFileMeta
	FFMPEGCmd  *exec.Cmd
}

type TranscodeClient struct {
	Transcode 	*AListTranscoder
	Transcoded  *map[string] TranscodeJob
	Jobs 		chan utilities.Action
	exitChan 	chan error
}

func (c TranscodeClient) ExitChan() chan error  {
	return c.exitChan
}

func (c *TranscodeClient) MakeTranscodeJob(_file *os.File, targetEncode ...string) {
	var err error
	buffer := make([]byte, 1024)
	id := betterguid.New()

	// BUILDING RESPONSE OBJECT //

	responseObj := SoundFileMeta{}
	responseObj.Id = id
	responseObj.Encoding, err = DetectEncoding(_file)
	responseObj.Name = string(id + "." + responseObj.Encoding)
	responseObj.BaseDir = path.Join(utilities.CWD(), "sound-files", id)
	responseObj.SourceDir = path.Join(utilities.CWD(), "sound-files", id,  "source" , "/" , responseObj.Encoding, "/")
	responseObj.URI = path.Join(responseObj.SourceDir,  responseObj.Name)

	// CREATE SOURCE FILE AND FOLDERS

	err = os.MkdirAll(responseObj.BaseDir, os.FileMode(utilities.PermissionsCodes["rw--"]))
	_newFile, err := os.Create(responseObj.URI)

	if err != nil {
		c.exitChan <- err
	}

	if err != nil  {
		c.exitChan <- err
	}

	defer _newFile.Close()

	// write to source file

	for {
		n, err := _file.Read(buffer)
		if err != nil && err != io.EOF {
			c.exitChan <- err
		}

		if _, err := _newFile.Write(buffer[:n]); err != nil {
			c.exitChan <- err
		}

		if n == 0 {
			break
		}
	}

	//  /BUILDING RESPONSE OBJECT //

	// ATTACH FFMPEG CMD //
	encodingCount := len(targetEncode)
	for i := 0;  i < encodingCount; i++ {
		cmd := buildFFMPEGCMD(responseObj, targetEncode[i])
		job := TranscodeJob {
			Id:         responseObj.Id,
			Ready:      true,
			Done:       false,
			SourceMeta: responseObj,
			TargetMeta: SoundFileMeta{},
			FFMPEGCmd:  cmd,
		}
		fmt.Println("New Job Success?", job)
		payload := []byte(fmt.Sprintf("%v", responseObj))
		action := utilities.Action{Type: "transcode", Payload: []byte(fmt.Sprintf("%v", payload))}

		c.Jobs <- action
	}

	fmt.Println("Closing")
}

// Sniffs out a files encoding
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

	// convert any video/audio to mp3 audio
	case "mp3":
		{
			return exec.Command(
				FFMPEGPath,
				"-i",
				sourceMeta.URI,
				// removes video
				"-vn",
				// sets sample rate
				"-ar 44100",
				// something with a 2 ...
				"-ac 2",
				// sets stream rate
				"-ab 192k",
				// forces mp3 encoding
				"-f mp3",
				sourceMeta.Id+".mp3",
			)
		}
		//case "flac":
		//case "wav":
	default:
	}
	return nil
}

// this will replace the other methods, mostly a forever loop that's
// concurrent and takes input through channels

func (c *TranscodeClient)ProcessJobs() {
	// range will keep going until channels is closed
	for action := range c.Jobs {
		buffer := make([]byte, 1024)
		payload := action.Payload
		bReader := bytes.NewReader(buffer)
		n, err := bReader.Read(payload)

		_file, err := os.Create("temp")

		for {
			_file.Write(payload[:n])

		}
		fmt.Println(_file, n)

		//err := j.FFMPEGCmd.Start()
			//utilities.ErrorHandler(err)
		//

		if err != nil {
			c.exitChan <- err
		}

		//err = j.ffmpegCMD.Wait()
		if err != nil {
			c.exitChan <- err
		}

		fmt.Printf("This ")

	}
}
