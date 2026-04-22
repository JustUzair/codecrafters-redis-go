package store

import (
	"os"
	"sync"
)

// No assignment of a particular type to actual cache store,
// to prevent repetitive modification of the store type entries and underlying functions
// Using the concepts of dynamic type inference at relevant helper and utils function, can preserve type where needed
// and still keeps the store generic

// Stream Data struct

type Field struct {
	Key   string
	Value any
}
type StreamEntry struct {
	ID     string
	Fields []Field
}

type Stream struct {
	StreamEntries []StreamEntry
}

//---------------------Config--------------------------//

type Config struct {
	Dir            string
	Appenddirname  string
	Appendonly     bool
	Appendfilename string
	Appendfsync    string
}

//---------------------------------------------------------//

// Generic Base Data Store

type Value[T any] struct {
	Value            T
	Deadline         int64
	IsDeadlineMillis bool
}
type Storage[T any] struct {
	mu        sync.RWMutex
	store     map[string]Value[T]
	notifiers map[string][]chan struct{}
	config    Config
}

// Init Redis store

var Cache = &Storage[any]{
	store:     make(map[string]Value[any]),
	notifiers: make(map[string][]chan struct{}),
	config: Config{
		Dir:            getPwd(),
		Appenddirname:  "appendonlydir",
		Appendfilename: "appendonly.aof",
		Appendfsync:    "everysec",
		Appendonly:     false,
	},
}

func getPwd() string {
	// 1. Check if the environment variable is set
	if path := os.Getenv("REDIS_DATA_DIR"); path != "" {
		return path
	}

	// 2. Fallback to current working directory for local development
	pwd, _ := os.Getwd()
	return pwd
}
