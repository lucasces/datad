package defs

type MessageService interface {
	Announce(Message) error
	Send(Message, string) error
	Channel() chan Message
}
