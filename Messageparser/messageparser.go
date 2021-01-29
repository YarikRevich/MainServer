package Messageparser

import (
	"encoding/json"
	"log"
)

type Parser interface{

	Parse([]*Message)[]byte
	Unparse([]byte)([]*Message, error)
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

func (m *Message) Unparse(message []byte)([]*Message, error){
	var unparsed []*Message
	//fmt.Println(string(message))
	err := json.Unmarshal(message, &unparsed)
	if err != nil{
		return nil, err
		//_, _ = fmt.Fprintln(os.Stderr, string(message))
	}
	return unparsed, nil
}

func NewMessage()[]*Message{
	return []*Message{
		{
			Type:       "OK",
		},
	}
}