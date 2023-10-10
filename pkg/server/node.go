package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"sync"
)

type HandlerFunc = func(msg *Message)

type Node struct {
	sync.Mutex
	Id            string
	NodeIds       []interface{}
	NextMessageId uint
	handlers      map[string]HandlerFunc
	scanner       *bufio.Scanner
	logLock       *sync.Mutex
}

func NewNode() *Node {
	n := Node{
		scanner:  bufio.NewScanner(os.Stdin),
		handlers: make(map[string]HandlerFunc, 0),
		logLock:  &sync.Mutex{},
	}

	n.On("init", func(msg *Message) {
		n.Log("starting node initialization")
		n.Id = msg.Body["node_id"].(string)
		n.NodeIds = msg.Body["node_ids"].([]interface{})

		n.Reply(msg, map[string]interface{}{"type": "init_ok"})
		n.Log("node %s initialized", n.Id)
	})

	return &n
}

func (n *Node) Log(message string, args ...any) {
	n.logLock.Lock()
	defer n.logLock.Unlock()

	message += "\n"
	fmt.Fprintf(os.Stderr, message, args...)
}

func (n *Node) Send(destination string, body map[string]interface{}) {
	msg := map[string]interface{}{
		"dest": destination,
		"src":  n.Id,
		"body": body,
	}

	n.Lock()
	defer n.Unlock()

	n.Log(fmt.Sprintf("Sending message: %v", msg))

	out, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}

func (n *Node) Reply(request *Message, body map[string]interface{}) {
	newBody := maps.Clone(body)
	newBody["in_reply_to"] = request.Body["msg_id"]
	n.Send(request.Source, newBody)
}

func (n *Node) On(msgType string, handler HandlerFunc) {
	if _, ok := n.handlers[msgType]; ok {
		panic(fmt.Sprintf("handler already registerd for %s", msgType))
	}

	n.Log("registering new handler for %s", msgType)
	n.handlers[msgType] = handler
}

func (n *Node) Run() {
	for n.scanner.Scan() {
		line := n.scanner.Text()
		n.Log("received msg: %v", line)

		msg := ParseMessage(line)
		n.Log("processing '%s' message", msg.Body["type"])

		n.Lock()
		f, ok := n.handlers[msg.Body["type"].(string)]
		if !ok {
			panic(fmt.Sprintf("handler not registerd for %s", msg.Body["type"]))
		}
		n.Unlock()

		n.Log("found handler for '%s'", msg.Body["type"])
		f(msg)
	}
}
