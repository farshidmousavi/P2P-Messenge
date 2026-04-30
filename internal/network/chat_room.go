package network

import (
	"context"
	"fmt"
	"log"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type ChatRoom struct {
	ps       *pubsub.PubSub
	topic    *pubsub.Topic
	sub      *pubsub.Subscription
	roomName string
	selfID   peer.ID
}

func JoinChatRoom(ctx context.Context, h host.Host, roomName string) (*ChatRoom, error) {
	// GossipSub: PubSub protocol for large-scale networks
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}

	topic, err := ps.Join(fmt.Sprintf("chat-room:%s", roomName))
	if err != nil {
		return nil, err
	}

	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	room := &ChatRoom{
		ps:       ps,
		topic:    topic,
		sub:      sub,
		roomName: roomName,
		selfID:   h.ID(),
	}

	log.Printf("Joined room: %s (ID: %s)", roomName, h.ID())
	return room, nil
}

// Publish sends a message to the entire chat room
func (cr *ChatRoom) Publish(message string) error {
	return cr.topic.Publish(context.Background(), []byte(message))
}

// Listen listens for new messages (run in a separate goroutine)
func (cr *ChatRoom) Listen(messageHandler func(peerID string, message string)) {
	for {
		msg, err := cr.sub.Next(context.Background())
		if err != nil {
			log.Println("Subscription error:", err)
			return
		}

		// Send message to handler (without filtering our own messages)
		messageHandler(msg.ReceivedFrom.String(), string(msg.Data))
	}
}