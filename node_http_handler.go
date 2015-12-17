package main

import "encoding/json"
import "net/http"

type NodeHttpHandler struct {
	ctx *Context
}

func (self NodeHttpHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	nodes, err := self.ctx.NodeService.ListNodes(0, 10)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	data, err := json.Marshal(nodes)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.WriteHeader(http.StatusOK)
	response.Write(data)
}

func NewNodeHttpHandler(ctx *Context) NodeHttpHandler {
	return NodeHttpHandler{ctx}
}
