package main

import "encoding/json"
import "reflect"

type MessageKind int

const (
	TYPE_ANNOUNCE MessageKind = iota
)

type Message struct {
	ID     int
	Kind   MessageKind
	Source string
	Detail interface{}
}

func (self *Message) Encode() ([]byte, error) {
	out, err := json.Marshal(self)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func Decode(message []byte) (Message, error) {
	out := Message{}
	err := json.Unmarshal(message, &out)
	if err != nil {
		return Message{}, err
	}
	switch {
	case out.Kind == TYPE_ANNOUNCE:
		obj := AnnounceDetail{}
		det := out.Detail.(map[string]interface{})
		for key, value := range det {
			vt := reflect.TypeOf(value)
			switch {
			case vt.Name() == "bool":
				reflect.ValueOf(&obj).Elem().FieldByName(key).SetBool(value.(bool))

			case vt.Name() == "string":
				reflect.ValueOf(&obj).Elem().FieldByName(key).SetString(value.(string))
			}
		}
		out.Detail = obj
	}
	return out, nil
}
