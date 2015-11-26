package message

import "log"

import "datad/defs"
import "datad/context"
import "datad/node"

type MessageProcessor struct {
	ctx context.Context
}

func NewMessageProcessor(ctx context.Context) MessageProcessor {
	return MessageProcessor{ctx}
}

func (self MessageProcessor) MessageProcessorRun() {
	for {
		msg := <-self.ctx.MessageService.Channel()

		switch {
		case msg.Kind == defs.TYPE_ANNOUNCE:
			self.processAnnounce(msg)
		}
	}
}

func (self MessageProcessor) processAnnounce(msg defs.Message) {
	detail := msg.Detail.(AnnounceDetail)
	if detail.Respond {
		self.ctx.MessageService.Send(NewAnnounceMessage(self.ctx.NodeService.NodeInfo(), false), msg.Source)
	}

	exists, err := self.ctx.NodeService.NodeExists(detail.Id)
	if err != nil {
		log.Println(err)
		return
	}

	if !exists {
		self.ctx.NodeService.AddNode(node.NewNode(detail.Id, detail.Name, msg.Source))
	}
}
