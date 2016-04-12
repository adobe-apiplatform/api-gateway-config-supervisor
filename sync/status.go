package sync

import (
	"sync"
	"time"
)

type Status struct {
	Status               string    `json:"status"`
	LastSync             time.Time `json:"lastSync"`
	LastFSChangeDetected time.Time `json:"lastChangeDetected"`
	LastReload           time.Time `json:"lastReload"`
}

var instance *Status
var once sync.Once

func GetStatusInstance() *Status {
	once.Do(func() {
		instance = &Status{
			LastReload:           time.Now(), // we assume it has happened already to avoid immediate reloads
			LastFSChangeDetected: time.Now(), // we assume it has happened already to avoid immediate reloads
		}
	})
	return instance
}
