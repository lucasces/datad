package main

import "log"
import "net/http"
import "time"

type HttpService struct {
	ctx    *Context
	server *http.Server
}

func (self *HttpService) run() {
	log.Fatal(self.server.ListenAndServe())
}

func CreateHttpService(ctx *Context) {
	server := &http.Server{
		Addr:         ctx.Config.HttpService.BindAddress,
		ReadTimeout:  time.Duration(ctx.Config.HttpService.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(ctx.Config.HttpService.WriteTimeout) * time.Second,
	}

	httpService := HttpService{ctx, server}

	http.Handle("/nodes", NewNodeHttpHandler(ctx))
	http.Handle("/files", NewFileHttpHandler(ctx))

	go httpService.run()
}
