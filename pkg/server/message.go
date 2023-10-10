package server

import (
	"encoding/json"
	"strings"
)

// READ: {"id":32,"src":"c13","dest":"n4","body":{"type":"read","msg_id":1}}
// BROAD: {"id":31,"src":"n4","dest":"n1","body":{"message":0,"type":"broadcast"}}
// TOPO: {"id":10,"src":"c6","dest":"n2","body":{
//        "type":"topology","topology":{"n1":["n4","n2"],"n2":["n5","n3","n1"],"n3":["n2"],"n4":["n1","n5"],"n5":["n2","n4"]},"msg_id":1}}
// INIT: {"id":0,"src":"c0","dest":"n4","body":{"type":"init","node_id":"n4","node_ids":["n1","n2","n3","n4","n5"],"msg_id":1}}

type Message struct {
	Id     uint                   `json:"id"`
	Source string                 `json:"src"`
	Dest   string                 `json:"dest"`
	Body   map[string]interface{} `json:"body"`
}

func ParseMessage(raw string) *Message {
	decoder := json.NewDecoder(strings.NewReader(raw))

	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		panic(err)
	}

	return &msg
}
