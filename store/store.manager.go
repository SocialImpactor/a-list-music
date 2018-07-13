package store

import (
	"os"
	"time"
	"path"
	"github.com/kataras/iris/core/router"
		"encoding/json"
	"a-list-music/utilities"
	)

var StoreBasePath = path.Join(utilities.CWD(), "sound-files")
var errorH = utilities.ErrorHandler

type Manifest struct {
	entries map[string]ManifestEntry
}

type ManifestOptions struct {
	removal []ManifestEntry
	update 	[]ManifestEntry
}

type ManifestEntry struct {
	Id       string
	URI      string
	Encoding string
	Size     int
}

func withManifest(options ManifestOptions) {
	buffer := make([]byte, 1024)
	manifest := Manifest{}

	_manifest, err := os.OpenFile(
		path.Join(StoreBasePath, "store_manifest.json"),
		os.O_WRONLY,
		os.FileMode(utilities.PermissionsCodes["rw--"]))
	errorH(err)
	defer func() {
		jBuff := make([]byte, 1024)
		jByte, err := json.Marshal(jBuff)
		errorH(err)
		_manifest.Write(jByte)
		_manifest.Close()
	}()
	json.Unmarshal(buffer, manifest)
	updates := options.update
	removals := options.removal

	for i := 0; i < len(updates); i++ {
		manifest.entries[updates[i].Id] = updates[i]
	}

	for i := 0; i < len(removals); i++ {
		delete(manifest.entries, removals[i].Id)
	}


}

type FileMeta struct {
	URIs map[string] string	`json:"uris"`
	OwnerId string			`json:"owner_id"`
	Size int				`json:"size"`
	StoredOn time.Time 		`json:"stored_on"`
}

type StoreOptions struct {
	File *os.File
	Name string
	Id string
}

type Job struct {
	Callback *func(data []byte, err error)
	Meta FileMeta
	Ready bool
	Done bool
}

type StoreClient struct {
	Jobs chan utilities.Action
}

type StoreManager interface {
	FetchFile(options StoreOptions) ([]byte, error)
	//WriteToFS(file *os.File, meta transcoder.SoundFileMeta)
}

func InitSoundLib() (string, error) {
	libDir := string(utilities.CWD() + "/sound-files" )
	if !router.DirectoryExists(libDir) {
		err := os.MkdirAll(libDir, os.FileMode(utilities.PermissionsCodes["rw--"]))
		if err != nil {
			return  libDir, err
		}
	}
	return libDir, nil
}

func SetClient(_client *StoreClient) {
	jobC := make(chan utilities.Action)
	storeClient := StoreClient{Jobs: jobC}
	_client = &storeClient
}

func (client *StoreClient) ProcessJobs() {
	//uri := manifest()
	//
	//for fj := range client.FetchJob {
	//	// get job
	//
	//	os.OpenFile(path.Join(StoreBasePath), fj.Meta.)
	//	// get / update meta
	//
	//
	//	// stream file
	//
	//	// check
	//
	//	//callback
	//
	//}
	//
	//for wj := range client.FetchJob {
	//	// get job
	//
	//	// get / update meta
	//
	//	// stream file
	//
	//	// check
	//
	//	//callback
	//
	//}
}