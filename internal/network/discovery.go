package network

import (
	"context"
	"log"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	routingdiscovery "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	discoveryutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

const DiscoveryNamespace = "p2p-messenger"

func DiscoverPeers(ctx context.Context, h host.Host, kadDHT *dht.IpfsDHT) {
	routingDiscovery := routingdiscovery.NewRoutingDiscovery(kadDHT)

	// Advertise ourselves in the DHT
	discoveryutil.Advertise(ctx, routingDiscovery, DiscoveryNamespace)
	log.Println("Advertised in DHT")
	log.Println("Peer discovery started...")

	// Find other peers
	peerChan, err := routingDiscovery.FindPeers(ctx, DiscoveryNamespace)
	if err != nil {
		log.Println("FindPeers error:", err)
		return
	}

	// Map to prevent reconnecting to the same peer
	connected := make(map[string]bool)

	for peerInfo := range peerChan {
		if peerInfo.ID == "" || peerInfo.ID == h.ID() {
			continue
		}

		// Skip if already connected
		if connected[peerInfo.ID.String()] {
			continue
		}

		// Skip if currently connected
		if h.Network().Connectedness(peerInfo.ID) == network.Connected {
			connected[peerInfo.ID.String()] = true
			continue
		}

		log.Printf("Found peer: %s", peerInfo.ID)

		// Attempt to connect with retries
		go func(p peer.AddrInfo) {
			for attempt := 1; attempt <= 3; attempt++ {
				connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				err := h.Connect(connCtx, p)
				cancel()

				if err == nil {
					log.Printf("Connected to: %s", p.ID)
					connected[p.ID.String()] = true
					return
				}

				log.Printf("Connection attempt %d to %s failed: %v", attempt, p.ID, err)
				time.Sleep(2 * time.Second)
			}
			log.Printf("Failed to connect to %s after 3 attempts", p.ID)
		}(peerInfo)
	}
}
