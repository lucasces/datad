package main

import "log"
import "os"
import "os/signal"

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	ctx, err := NewContext()
	if err != nil {
		log.Fatal(err)
	}

	err = ctx.MessageService.Announce(NewAnnounceMessage(ctx.Config.NodeInfo, true))
	if err != nil {
		log.Fatal(err)
	}

	_ = <-c

	ctx.Close()
}
