package network

import (
	"context"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

// MaintainConnections checks and maintains stable connections
func MaintainConnections(ctx context.Context, h host.Host) {
	ticker := time.NewTicker(30 * time.Second) // افزایش به 30 ثانیه
	defer ticker.Stop()

	// Map to prevent repeated attempts to failed peers
	failedPeers := make(map[string]struct {
		attempts int
		lastTry  time.Time
	})

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			peers := h.Peerstore().Peers()
			for _, p := range peers {
				if p == h.ID() {
					continue
				}

				// Check connection status
				if h.Network().Connectedness(p) == network.Connected {
					// If connected, clear failed record
					delete(failedPeers, p.String())
					continue
				}

				// Check history of failed attempts
				if record, exists := failedPeers[p.String()]; exists {
					// If more than 3 attempts, remove the peer
					if record.attempts >= 3 {
						log.Printf("Removing unreachable peer %s after %d failed attempts", p, record.attempts)
						h.Peerstore().ClearAddrs(p)
						delete(failedPeers, p.String())
						continue
					}
					
					// If less than 30 seconds since last attempt, wait
					if time.Since(record.lastTry) < 30*time.Second {
						continue
					}
				}

				addrs := h.Peerstore().Addrs(p)
				if len(addrs) == 0 {
					// No addresses, remove peer
					log.Printf("🗑️ Removing peer %s (no addresses)", p)
					h.Peerstore().ClearAddrs(p)
					delete(failedPeers, p.String())
					continue
				}

				// Attempt to connect
				log.Printf("Connecting to peer: %s", p)
				connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				err := h.Connect(connCtx, peer.AddrInfo{ID: p, Addrs: addrs})
				cancel()

				if err != nil {
					// Record failed attempt
					record := failedPeers[p.String()]
					record.attempts++
					record.lastTry = time.Now()
					failedPeers[p.String()] = record
					
					log.Printf("Failed to connect to %s (attempt %d/3): %v", p, record.attempts, err)
				} else {
					log.Printf("Connected to: %s", p)
					delete(failedPeers, p.String())
				}
			}
		}
	}
}