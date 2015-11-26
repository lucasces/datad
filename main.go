package main

import "log"

import "datad/context"

import "datad/config"
import "datad/message"
import "datad/node"

type MessageService interface {
	Announce() error
	Send() error
}

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, err := context.NewContext(config)
	if err != nil {
		log.Fatal(err)
	}

	nodeService := node.NewNodeService(ctx)
	ctx.NodeService = nodeService

	messageService, err := message.NewUDPMessageService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ctx.MessageService = messageService

	err = ctx.MessageService.Announce(message.NewAnnounceMessage(ctx.NodeService.NodeInfo(), true))
	if err != nil {
		log.Fatal(err)
	}

	config.SaveConfig()

	ctx.WaitGroup.Wait()

	log.Printf("hello")
}
