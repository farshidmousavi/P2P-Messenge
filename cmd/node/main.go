package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/d4rkpc/p2p-messenger/internal/network"
)

func main() {
	// Create context with cancel for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Starting P2P Messenger...")
	node, err := network.NewNode(ctx)
	if err != nil {
		log.Fatal("Failed to create node:", err)
	}

	// Connect to bootstrap peers from config file
	fmt.Println("Connecting to bootstrap peers...")
	network.ConnectToBootstrap(ctx, node.Host)

	// Initialize and bootstrap Kademlia DHT for peer discovery
	kadDHT, err := network.SetupDHT(ctx, node.Host)
	if err != nil {
		log.Fatal("Failed to setup DHT:", err)
	}

	// Start automatic peer discovery via DHT
	go network.DiscoverPeers(ctx, node.Host, kadDHT)
	// Manage and maintain active connections (reconnect if needed)
	go network.MaintainConnections(ctx, node.Host)

	// Optional: Print known peers every 30 seconds
	// Uncomment if you want to see connected peers list
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		case <-time.After(30 * time.Second):
	// 			network.PrintKnownPeers(node.Host)
	// 		}
	// 	}
	// }()

	// Join the global chat room
	room, err := network.JoinChatRoom(ctx, node.Host, "general-chat")
	if err != nil {
		log.Fatal("Failed to join chat room:", err)
	}

	// Listen for incoming messages from other peers
	go room.Listen(func(peerID string, message string) {
		shortID := peerID
		if len(peerID) > 8 {
			shortID = peerID[len(peerID)-8:]
		}
		fmt.Printf("\n[%s]: %s\n", shortID, message)
		fmt.Print("You: ")
	})

	// Handle graceful shutdown on Ctrl+C or SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\nShutting down...")
		cancel()
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	// Display command information and start input loop
	fmt.Println("\nConnected! Type your message or commands:")
	fmt.Println("Commands:")
	fmt.Println("  1. 'bootstrap'        - Register your node as a bootstrap peer")
	fmt.Println("  2. 'list-bootstrap'   - Show all saved bootstrap peers")
	fmt.Println("  3. 'connect-bootstrap'- Connect to all bootstrap peers")
	fmt.Println("  4. 'clean-bootstrap'  - Remove inactive bootstrap peers")
	fmt.Println("  5. 'connect <addr>'   - Manually connect to a peer")
	fmt.Print("You: ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			fmt.Print("You: ")
			continue
		}

		// Register current node as a bootstrap peer
		if text == "bootstrap" {
			fmt.Println("Converting this node to bootstrap...")
			if err := network.AddSelfAsBootstrap(ctx, node.Host); err != nil {
				log.Printf("Failed to add self as bootstrap: %v", err)
			} else {
				fmt.Println("Node registered as bootstrap successfully!")
				fmt.Println("   Other nodes can now discover this node via bootstrap list")
			}
			fmt.Print("You: ")
			continue
		}

		// connect-bootstrap command to connect to all bootstrap peers
		if text == "connect-bootstrap" {
			fmt.Println("Connecting to all bootstrap peers...")
			network.ConnectToBootstrapPeers(ctx, node.Host)
			fmt.Print("You: ")
			continue
		}

		// Display list of all saved bootstrap peers
		if text == "list-bootstrap" {
			network.PrintBootstrapList()
			fmt.Print("You: ")
			continue
		}

		// Remove inactive bootstrap peers from config
		if text == "clean-bootstrap" {
			fmt.Println("Cleaning inactive bootstrap peers...")
			network.CleanInactiveBootstrapPeers(ctx, node.Host)
			fmt.Print("You: ")
			continue
		}
		// Manually connect to a specific peer address
		if strings.HasPrefix(text, "connect ") {
			addr := strings.TrimPrefix(text, "connect ")
			fmt.Println("Connecting to:", addr)
			if err := node.Connect(ctx, addr); err != nil {
				log.Printf("Connection failed: %v", err)
			} else {
				fmt.Println("Connected successfully!")
			}
			fmt.Print("✏️ You: ")
			continue
		}

		// Send regular message to the chat room
		if err := room.Publish(text); err != nil {
			log.Printf("Failed to send: %v", err)
		}
		fmt.Print("You: ")
	}
}