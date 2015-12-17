package main

import "github.com/boltdb/bolt"

type Context struct {
	MessageService *MessageService
	NodeService    *NodeService
	JobService     *JobService
	FileService    *FileService
	Config         Config
	Database       *bolt.DB
}

func NewContext() (Context, error) {
	config, err := LoadConfig()
	if err != nil {
		return Context{}, err
	}

	ctx := Context{nil, nil, nil, nil, config, nil}

	err = ctx.InitServices()
	if err != nil {
		return Context{}, err
	}

	return ctx, nil
}

func (self *Context) Close() {
	self.Config.SaveConfig()
	self.Database.Close()
}

func (self *Context) InitServices() error {
	var err error

	self.Database, err = bolt.Open(self.Config.DatabaseFile, 0600, nil)
	if err != nil {
		return err
	}

	CreateNodeService(self)

	_, err = CreateMessageService(self)
	if err != nil {
		return err
	}

	CreateHttpService(self)
	CreateJobService(self)
	CreateFileService(self)

	return nil

}
