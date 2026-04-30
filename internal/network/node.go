package network

import (
	"context"
	"fmt"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type Node struct {
	Host host.Host
}

func NewNode(ctx context.Context) (*Node, error) {
	// Settings for NAT and firewall traversal
	opts := []libp2p.Option{
		libp2p.EnableHolePunching(),    // NAT hole punching
		libp2p.EnableRelay(),           // Use relay if direct connection fails
		//EnableAutoNAT is disabled in current version - using default capabilities instead
	}

	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Node ID: %s\n", h.ID())
	
	// Display all addresses (including relay addresses)
	for _, addr := range h.Addrs() {
		fmt.Printf("Listening on: %s/p2p/%s\n", addr, h.ID())
	}

	return &Node{Host: h}, nil
}

func (n *Node) Connect(ctx context.Context, addr string) error {
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}

	if err := n.Host.Connect(ctx, *info); err != nil {
		return err
	}

	log.Printf("Connected to: %s", info.ID)
	return nil
}