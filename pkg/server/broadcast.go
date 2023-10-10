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

	node.On("topology", func(msg *Message) {
		node.Log("registering node neighbors")
		topo := msg.Body["topology"].(map[string]interface{})

		bcast.neighbors = topo[node.Id].([]interface{})
		node.Log("received neighbors %v", bcast.neighbors)
		node.Reply(msg, map[string]interface{}{"type": "topology_ok"})
	})

	node.On("read", func(msg *Message) {
		bcast.Lock()
		defer bcast.Unlock()

		node.Reply(msg, map[string]interface{}{
			"type":     "read_ok",
			"messages": bcast.messages.Items(),
		})
	})

	node.On("broadcast", func(msg *Message) {
		m := msg.Body["message"].(T)

		bcast.Lock()
		// Only broadcast new messages
		if !bcast.messages.Contains(m) {
			bcast.messages.Add(m)

			// Broadcast message to all neighbors
			for _, neighbor := range bcast.neighbors {
				if neighbor == msg.Source {
					// Don't broadcast back to the source
					continue
				}

				node.Send(neighbor.(string), map[string]interface{}{
					"type":    "broadcast",
					"message": m,
				})
			}
		}
		bcast.Unlock()

		if _, ok := msg.Body["msg_id"]; ok {
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
