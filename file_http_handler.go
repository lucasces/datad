package main

import "encoding/json"
import "net/http"
import "io/ioutil"

type FileHttpHandler struct {
	ctx *Context
}

func (self FileHttpHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	switch {
	case request.Method == "GET":
		self.handleGet(response, request)
		return
	case request.Method == "POST":
		self.handlePost(response, request)
		return
	}
}

func (self FileHttpHandler) handleGet(response http.ResponseWriter, request *http.Request) {
	files, err := self.ctx.FileService.ListFiles(0, 10)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	data, err := json.Marshal(files)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.WriteHeader(http.StatusOK)
	response.Write(data)
}

func (self FileHttpHandler) handlePost(response http.ResponseWriter, request *http.Request) {
	if request.Body != nil {
		fileData := FileData{}
		data, err := ioutil.ReadAll(request.Body)
		err = json.Unmarshal(data, &fileData)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(err.Error()))
			return
		}
		file := NewFile(fileData)
		self.ctx.FileService.AddFile(file)

		response.WriteHeader(http.StatusCreated)

		return
	}
	response.WriteHeader(http.StatusBadRequest)
}

func NewFileHttpHandler(ctx *Context) FileHttpHandler {
	return FileHttpHandler{ctx}
}
