package network

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
)

// PrintKnownPeers prints all peers that the node knows
func PrintKnownPeers(h host.Host) {
	peers := h.Peerstore().Peers()

	fmt.Println("Known peers:")
	for _, p := range peers {
		if p == h.ID() {
			continue
		}
		fmt.Println(" -", p.String())
	}
}
