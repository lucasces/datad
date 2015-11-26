package context

import "sync"

import "github.com/boltdb/bolt"

import "datad/defs"
import "datad/config"

type Context struct {
	WaitGroup      *sync.WaitGroup
	MessageService defs.MessageService
	NodeService    defs.NodeService
	Config         config.Config
	Database       *bolt.DB
}

func NewContext(config config.Config) (Context, error) {
	var waitGroup sync.WaitGroup

	db, err := bolt.Open(config.DatabaseFile, 0600, nil)

	if err != nil {
		return Context{}, err
	}

	return Context{&waitGroup, nil, nil, config, db}, nil
}
