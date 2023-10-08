package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	InitOK string = "init_ok"
	EchoOK string = "echo_ok"
)

type Body struct {
	Type      string   `json:"type"`
	NodeId    string   `json:"node_id,omitempty"`
	NodeIds   []string `json:"node_ids,omitempty"`
	MsgId     uint     `json:"msg_id"`
	InReplyTo uint     `json:"in_reply_to"`
	Message   string   `json:"echo,omitempty"`
}

type Message struct {
	Source string `json:"src"`
	Dest   string `json:"dest"`
	Body   Body   `json:"body"`
}

func ScanMessage(msg string) (*Message, error) {
	decoder := json.NewDecoder(strings.NewReader(msg))

	var s Message
	err := decoder.Decode(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

type Server struct {
	NodeId        string
	NextMessageId uint
}

func (s *Server) Reply(req *Message, msgType string, message string) {
	msgId := s.NextMessageId + 1

	res := Message{
		Source: s.NodeId,
		Dest:   req.Source,
		Body: Body{
			MsgId:     msgId,
			Type:      msgType,
			InReplyTo: req.Body.MsgId,
			Message:   message,
		},
	}

	out, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, string(out))
}

func main() {
	server := Server{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintln(os.Stderr, "received:", line)

		msg, err := ScanMessage(line)
		if err != nil {
			log.Fatalf("failed to scan message: %v", err)
		}

		switch msg.Body.Type {
		case "init":
			server.NodeId = msg.Body.NodeId
			fmt.Fprintln(os.Stderr, "initialized node", server.NodeId)
			server.Reply(msg, InitOK, "")
		case "echo":
			fmt.Fprintln(os.Stderr, "Echoing", msg.Body)
			server.Reply(msg, EchoOK, msg.Body.Message)
		default:
			log.Fatalf("unsupported message type")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
	}
}
