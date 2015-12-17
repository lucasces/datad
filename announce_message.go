package main

type AnnounceDetail struct {
	Id      string
	Name    string
	Respond bool
}

func NewAnnounceMessage(node Node, respond bool) Message {
	return Message{0, TYPE_ANNOUNCE, "", AnnounceDetail{node.Id, node.Name, respond}}
}
