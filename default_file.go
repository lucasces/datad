package main

import "io"
import "os"

type DefaultFile struct {
	FileData
}

func (self DefaultFile) Id() string {
	return self.FileData.Id
}

func (self *DefaultFile) SetId(id string) {
	self.FileData.Id = id
}

func (self DefaultFile) Name() string {
	return self.FileData.Name
}

func (self *DefaultFile) SetName(name string) {
	self.FileData.Name = name
}

func (self DefaultFile) Path() string {
	return self.FileData.Path
}

func (self DefaultFile) Size() (int64, error) {
	reader, err := os.Open(self.Path())
	if err != nil {
		return 0, err
	}

	info, err := reader.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func (self DefaultFile) Type() FileType {
	return self.FileData.Type
}

func (self DefaultFile) NewReader() (io.Reader, error) {
	return os.Open(self.Path())
}

func (self DefaultFile) Data() FileData {
	return self.FileData
}

func NewDefaultFile(data FileData) DefaultFile {
	return DefaultFile{data}
}
