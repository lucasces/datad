package main

import "io"

type FileType int

const (
	FILE_DEFAULT = iota
)

type File interface {
	Id() string
	Name() string
	Path() string
	Type() FileType
	Size() (int64, error)
	NewReader() (io.Reader, error)
	Data() FileData
}

type FileData struct {
	Id   string
	Name string
	Path string
	Type FileType
}

func NewFile(data FileData) File {
	switch {
	case data.Type == FILE_DEFAULT:
		return NewDefaultFile(data)
	}
	return nil
}
