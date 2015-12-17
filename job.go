package main

import "io"
import "log"
import "sync"

type CompleteCallback func(*Job) error

type JobType int

const (
	JOB_ID_GENERATION = iota
)

type Job struct {
	Id          string
	Source      File
	Destination io.Writer
	BlockSize   int
	Comm        chan []byte

	size            int64
	readerTaskError error
	writerTaskError error
	status          float32
	onComplete      CompleteCallback
	waitGroup       sync.WaitGroup
}

func (self *Job) Execute() {
	var err error
	if err != nil {
		log.Printf("Job execution terminated with error: %s", err)
		return
	}

	self.size, err = self.Source.Size()
	if err != nil {
		log.Printf("Job execution terminated with error: %s", err)
		return
	}

	self.waitGroup.Add(1)
	go self.ReadTask()

	self.waitGroup.Add(1)
	go self.WriteTask()

	self.waitGroup.Wait()

	err = self.onComplete(self)

	switch {
	case self.readerTaskError != nil:
		log.Printf("Job execution terminated with error: %s", self.readerTaskError)

	case self.writerTaskError != nil:
		log.Printf("Job execution terminated with error: %s", self.writerTaskError)
	case err != nil:
		log.Printf("Job execution terminated with error: %s", self.writerTaskError)
	}
}

func NewJob(source File, destination io.Writer, blockSize int, onComplete CompleteCallback) Job {
	return Job{generateNewId(), source, destination, blockSize, make(chan []byte), 0, nil, nil, 0, onComplete, sync.WaitGroup{}}
}

func (self *Job) ReadTask() {
	read := int64(0)
	buffer := make([]byte, self.BlockSize)

	reader, err := self.Source.NewReader()
	if err != nil {
		self.readerTaskError = err
		return
	}

	for read < self.size && self.writerTaskError == nil {
		n, err := reader.Read(buffer)
		if err != nil {
			self.readerTaskError = err
			break
		}
		read += int64(n)
		self.Comm <- buffer[:n]
	}
	self.waitGroup.Done()
}

func (self *Job) WriteTask() {
	written := int64(0)

	for written < self.size && self.readerTaskError == nil {
		buffer := <-self.Comm
		n, err := self.Destination.Write(buffer)
		if err != nil {
			self.writerTaskError = err
			break
		}
		written += int64(n)
		self.status = (float32(written) / float32(self.size))
	}
	self.waitGroup.Done()
}
