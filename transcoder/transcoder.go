package transcoder

import (
	"os"
	"os/exec"
	"net/http"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/kataras/iris/core/router"
	"strings"
	"bytes"
	"path"
	"github.com/kjk/betterguid"
	"io"
)

func CWD() string {
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}
	return "./"
}

var PermissionsCodes = map[string] int {
	"rwxrr": 744,
	"rwx--": 700,
	"rwrr": 644,
	"rw--": 600,
	"r--": 400,
}
var FFMPEGPath = string(path.Join(CWD(), "..", "..", "..", "Desktop", "ffmpeg", "ffmpeg"))

var EncExtMap = map[string] string {
	"audio/wave": "wav",
}

type SoundFileMeta struct {
	id string
	name string
	uri string
	baseDir string
	encoding string
	codex string
	size int
}

type Transcoder interface {
	StoreToMediaLibrary()(SoundFileMeta, err error)
	// meta is for the library IA, it relates to URLs
	MetaBuilder(file *os.File) (SoundFileMeta)

	// New Jobs set the source file...
	NewJob(file *os.File, targetMime []string)

	//
	//RunTranscodes(jobs map[string] TranscodeJob)

	ExitChan() chan error
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
	ReadyTranscodes chan map[string] TranscodeJob
	Transcoded      chan map[string] TranscodeJob
	Jobs 			chan TranscodeJob
	exitChan 		chan error
}

func InitSoundLib() (string, error) {
	libDir := string(CWD() + "/sound-files" )
	if !router.DirectoryExists(libDir) {
		err := os.MkdirAll(libDir, os.FileMode(PermissionsCodes["rw--"]))
		if err != nil {
			return  libDir, err
		}
	}
	return libDir, nil
}

func SetClient(transcoderClient *TranscoderClient) {
	client := TranscoderClient{}

	if transcoderClient != nil {
		transcoderClient = &client
	}
}

func (c TranscoderClient) ExitChan() chan error  {
	return c.exitChan
}

func (c *TranscoderClient) MakeTranscodeJob(_file *os.File, targetEncode ...string) {
	var err error
	jobs := make(map[string] TranscodeJob)
	buffer := make([]byte, 1024)
	id := betterguid.New()

	// BUILDING RESPONSE OBJECT //

	responseObj := SoundFileMeta{}
	responseObj.id = id
	responseObj.encoding, err = DetectEncoding(_file)
	responseObj.name = string(id + "." + responseObj.encoding)
	responseObj.baseDir = path.Join(CWD(), "sound-files", id)
	dir := path.Join(CWD(), "/sound-files",  "/", id,  "/source" , "/" , responseObj.encoding, "/")
	responseObj.uri = dir + responseObj.name

	// CREATE SOURCE FILE AND FOLDERS

	err = os.MkdirAll(dir, os.FileMode(PermissionsCodes["rw--"]))
	_newFile, err := os.Create(responseObj.uri)

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
			id:         responseObj.id,
			ready:      true,
			done:       false,
			sourceMeta: responseObj,
			targetMeta: SoundFileMeta{},
			ffmpegCMD:  cmd,
		}
		fmt.Println("New Job Success?", job)

		jobs[job.id] = job
	}
	c.ReadyTranscodes <- jobs
	fmt.Println("Closing")
	close(c.ReadyTranscodes)
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
				sourceMeta.uri,
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
				sourceMeta.id+".mp3",
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

func (c *TranscoderClient)ProcessJobs() {

	// range will keep going until channels is closed

	for j := range c.Jobs {
		_out := bytes.Buffer{}
		j.ffmpegCMD.Stdout = &_out

		err := j.ffmpegCMD.Start()


		if err != nil {
			c.exitChan <- err
		}

		err = j.ffmpegCMD.Wait()
		if err != nil {
			c.exitChan <- err
		}



	}
}
