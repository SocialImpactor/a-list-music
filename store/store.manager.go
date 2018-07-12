package store

import (
	"os"
	"a-list-music/transcoder"
	"time"
)

type FileMeta struct {
	*transcoder.SoundFileMeta 	`json:"sound_meta"`
	URIs map[string] string		`json:"uris"`
	OwnerId string				`json:"owner_id"`
	Size int					`json:"size"`
	StoredOn time.Time 			`json:"stored_on"`
}

type StoreOptions struct {
	file *os.File
	name string
	id string
}

type StoreManager interface {
	FetchFile(options StoreOptions) ([]byte, error)
	WriteToFS(file *os.File, meta transcoder.SoundFileMeta)
}


func InitializeStoreManager() StoreManager {

}