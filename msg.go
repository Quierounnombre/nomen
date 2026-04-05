package main

import (
	"encoding/json"
	"net"
)

type Msg struct {
	author		string
	encoder		*json.Encoder
	content		map[string]interface{}
}

func (m *Msg) clean_content() {
	m.content = make(map[string]interface{})
}

func (m *Msg) reply() {
	var err error

	err = m.encoder.Encode(m.content)
	if err != nil {
		fmt.Println("Error socket not working")
		fmt.Println(err)
		fmt.Println("Target conn -> ", m.author)
	}
}
