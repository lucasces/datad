package main

import "log"

type MessageProcessor struct {
	ctx *Context
}

func NewMessageProcessor(ctx *Context) MessageProcessor {
	return MessageProcessor{ctx}
}

func (self *MessageProcessor) MessageProcessorRun() {
	for {
		msg := <-self.ctx.MessageService.Channel()

		switch {
		case msg.Kind == TYPE_ANNOUNCE:
			self.processAnnounce(msg)
		}
	}
}

func (self *MessageProcessor) processAnnounce(msg Message) {
	detail := msg.Detail.(AnnounceDetail)
	if detail.Respond {
		self.ctx.MessageService.Send(NewAnnounceMessage(self.ctx.Config.NodeInfo, false), msg.Source)
	}

	exists, err := self.ctx.NodeService.NodeExists(detail.Id)
	if err != nil {
		log.Println(err)
		return
	}

	if !exists {
		self.ctx.NodeService.AddNode(NewNode(detail.Id, detail.Name, msg.Source))
	}
}
