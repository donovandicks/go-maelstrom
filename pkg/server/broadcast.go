package server

import (
	"sync"

	"github.com/donovandicks/godistsys/pkg/utils"
)

type BroadcastServer[T comparable] struct {
	sync.Mutex
	node      *Node
	neighbors []interface{}
	messages  utils.Set[T]
}

func NewBroadcastServer[T comparable]() *BroadcastServer[T] {
	node := NewNode()

	bcast := BroadcastServer[T]{
		node:     node,
		messages: utils.NewSet[T](),
	}

	node.On("topology", func(msg map[string]interface{}) {
		node.Log("registering node neighbors")
		body := (msg["body"]).(map[string]interface{})
		topo := (body["topology"].(map[string]interface{}))

		bcast.neighbors = topo[node.Id].([]interface{})
		node.Log("received neighbors %v", bcast.neighbors)
		node.Reply(msg, map[string]interface{}{"type": "topology_ok"})
	})

	node.On("read", func(msg map[string]interface{}) {
		bcast.Lock()
		defer bcast.Unlock()

		node.Reply(msg, map[string]interface{}{
			"type":     "read_ok",
			"messages": bcast.messages.Items(),
		})
	})

	node.On("broadcast", func(msg map[string]interface{}) {
		body := (msg["body"]).(map[string]interface{})
		m := body["message"].(T)

		bcast.Lock()
		// Only broadcast new messages
		if !bcast.messages.Contains(m) {
			bcast.messages.Add(m)

			// Broadcast message to all neighbors
			for _, neighbor := range bcast.neighbors {
				node.Send(neighbor.(string), map[string]interface{}{
					"type":    "broadcast",
					"message": m,
				})
			}
		}
		bcast.Unlock()

		if _, ok := body["msg_id"]; ok {
			// Only reply to messages that have an ID, which are only ones
			// from the controller and not from peers
			node.Reply(msg, map[string]interface{}{"type": "broadcast_ok"})
		}
	})

	return &bcast
}

func (s *BroadcastServer[T]) Run() {
	s.node.Run()
}
