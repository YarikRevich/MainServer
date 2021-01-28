package Messageparser

import(
	"log"
	"fmt"
	"encoding/json"
)

type Parser interface{

	Parse([]*Message)[]byte
	Unparse([]byte)[]*Message
}

type Message struct{
	Index           int
	Error           string 
	Type            string
	Body            []string
}

func (m *Message) Parse(message []*Message)[]byte{
	
	b, err := json.Marshal(message)
	if err != nil{
		log.Fatalln(err)
	}
	return b
}

func (m *Message) Unparse(message []byte)[]*Message{
	var unparsed []*Message
	fmt.Println(string(message))
	err := json.Unmarshal(message, &unparsed)
	if err != nil{
		log.Fatalln(err)
	}
	return unparsed
}

func NewMessage()[]*Message{
	return []*Message{
		&Message{
			Index:           0,
			Error:           "0",
			Type:            "OK",
			Body:            []string{},
		},
	}
}