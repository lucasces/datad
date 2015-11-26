package message

import "encoding/json"
import "reflect"

import "datad/defs"

type AnnounceDetail struct {
	Id      string
	Name    string
	Respond bool
}

func Encode(message defs.Message) ([]byte, error) {
	out, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func Decode(message []byte) (defs.Message, error) {
	out := defs.Message{}
	err := json.Unmarshal(message, &out)
	if err != nil {
		return defs.Message{}, err
	}
	switch {
	case out.Kind == defs.TYPE_ANNOUNCE:
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

func NewAnnounceMessage(node defs.Node, respond bool) defs.Message {
	return defs.Message{0, defs.TYPE_ANNOUNCE, "", AnnounceDetail{node.Id, node.Name, respond}}
}
