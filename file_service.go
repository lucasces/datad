package main

import "encoding/json"
import "errors"
import "fmt"
import "hash/adler32"
import "log"

import "github.com/boltdb/bolt"

type FileService struct {
	ctx   *Context
	files map[string]File
}

var filesBucketName = []byte("files")

func (self *FileService) AddFile(file File) error {

	writer := adler32.New()
	job := NewJob(file, writer, 2048, func(job *Job) error {
		var data []byte
		var err error

		switch {
		case file.Type() == FILE_DEFAULT:
			// Unchecked conversion. May cause problems
			defaultFile, _ := file.(DefaultFile)
			defaultFile.SetId(fmt.Sprintf("%x", writer.Sum32()))
			log.Printf("Adding file %s", defaultFile.Id())
			data, err = json.Marshal(defaultFile.Data())
			if err != nil {
				return err
			}
		}

		err = self.ctx.Database.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(filesBucketName)
			return b.Put([]byte(file.Id()), data)
		})
		return nil
	})
	self.ctx.JobService.Channel() <- job
	return nil
}

func (self *FileService) FileExists(id string) (bool, error) {
	exists := false
	err := self.ctx.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(filesBucketName)

		data := b.Get([]byte(id))
		exists = (data != nil)
		return nil
	})

	return exists, err
}

func (self *FileService) RemoveFile(id string) {
	delete(self.files, id)
}

func (self *FileService) GetFile(id string) (File, error) {
	fileData := FileData{}
	err := self.ctx.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(filesBucketName)
		data := b.Get([]byte(id))
		if data == nil {
			return errors.New("File does not exist")
		}

		json.Unmarshal(data, &fileData)
		return nil
	})

	file := NewFile(fileData)

	return file, err
}

func (self *FileService) ListFiles(offset int, limit int) ([]File, error) {
	files := []File{}
	err := self.ctx.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(filesBucketName)
		n := 0
		err := b.ForEach(func(k, v []byte) error {
			switch {
			case n > (offset + limit - 1):
				return NewLimitReachedError()
			case n >= offset:
				fileData := FileData{}
				err := json.Unmarshal(v, &fileData)
				if err != nil {
					return err
				}
				files = append(files, NewFile(fileData))
			}
			n += 1
			return nil
		})
		if _, ok := err.(LimitReachedError); ok {
			return nil
		}
		return err
	})
	return files, err
}

func CreateFileService(ctx *Context) {
	fileService := FileService{ctx, make(map[string]File)}

	err := ctx.Database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(filesBucketName)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	ctx.FileService = &fileService
}
