package network

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/d4rkpc/p2p-messenger/internal/config"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

var BootstrapPeers = []string{
	// Add your default bootstrap peers here if needed
}

func ConnectToBootstrap(ctx context.Context, h host.Host) {
	// Load bootstrap settings from config file
	bootstrapConfig, err := config.LoadBootstrapConfig()
	if err != nil {
		log.Println("Failed to load bootstrap config:", err)
		return
	}

	log.Printf("Loaded %d bootstrap peers from config", len(bootstrapConfig.Peers))

	// Connect to default bootstrap peers
	for _, addr := range BootstrapPeers {
		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			log.Println("Invalid bootstrap addr:", err)
			continue
		}

		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Println("AddrInfo error:", err)
			continue
		}

		if info.ID == h.ID() {
			continue
		}

		h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
		
		log.Printf("Attempting to connect to bootstrap: %s", info.ID)

		if err := h.Connect(ctx, *info); err != nil {
			log.Printf("Bootstrap connect failed to %s: %v", info.ID, err)
			continue
		}

		log.Printf("Connected to bootstrap: %s", info.ID)
	}

	// Connect to bootstrap peers stored in config file
	for _, addr := range bootstrapConfig.Peers {
		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			log.Println("Invalid bootstrap addr from config:", err)
			continue
		}

		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Println("AddrInfo error from config:", err)
			continue
		}

		if info.ID == h.ID() {
			log.Println("Skipping self-connection")
			continue
		}

		h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

		log.Printf("Attempting to connect to bootstrap (config): %s", info.ID)

		if err := h.Connect(ctx, *info); err != nil {
			log.Printf("Bootstrap connect failed to %s: %v", info.ID, err)
			continue
		}

		log.Printf("Connected to bootstrap (config): %s", info.ID)
	}
}

// ConnectToBootstrapPeers connects to all saved bootstrap peers (manual command)
func ConnectToBootstrapPeers(ctx context.Context, h host.Host) {
	// Load bootstrap settings from config file
	bootstrapConfig, err := config.LoadBootstrapConfig()
	if err != nil {
		log.Printf("Failed to load bootstrap config: %v", err)
		return
	}

	if len(bootstrapConfig.Peers) == 0 {
		fmt.Println("No bootstrap peers found in config")
		return
	}

	fmt.Printf("Found %d bootstrap peers in config\n", len(bootstrapConfig.Peers))

	connectedCount := 0
	for i, addr := range bootstrapConfig.Peers {
		fmt.Printf("   [%d/%d] Trying to connect to: %s\n", i+1, len(bootstrapConfig.Peers), addr)

		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			log.Printf("   Invalid address: %v", err)
			continue
		}

		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Printf("   AddrInfo error: %v", err)
			continue
		}

		// Skip if it's our own node
		if info.ID == h.ID() {
			fmt.Println("   Skipping self-connection")
			continue
		}

		// Skip if already connected
		if h.Network().Connectedness(info.ID) == network.Connected {
			fmt.Printf("   Already connected to: %s\n", info.ID)
			connectedCount++
			continue
		}

		// Attempt to connect
		connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		err = h.Connect(connCtx, *info)
		cancel()

		if err != nil {
			log.Printf("   Connection failed to %s: %v", info.ID, err)
			continue
		}

		fmt.Printf("   Connected to: %s\n", info.ID)
		connectedCount++
	}

	fmt.Printf("\nSummary: Connected to %d/%d bootstrap peers\n", connectedCount, len(bootstrapConfig.Peers))
}

// AddSelfAsBootstrap adds the current node's address to bootstrap list
func AddSelfAsBootstrap(ctx context.Context, h host.Host) error {
	// Get all addresses of the current node
	if len(h.Addrs()) == 0 {
		log.Println("No addresses found for self")
		return nil
	}

	// Prefer using TCP address
	var bestAddr string
	for _, addr := range h.Addrs() {
		// Prefer TCP protocol
		if len(addr.Protocols()) > 0 && addr.Protocols()[0].Code == multiaddr.P_TCP {
			// Only select non-localhost addresses
			addrStr := addr.String()
			if !strings.Contains(addrStr, "127.0.0.1") && !strings.Contains(addrStr, "::1") {
				fullAddr := addrStr + "/p2p/" + h.ID().String()
				bestAddr = fullAddr
				break
			}
		}
	}

	// If no public address found, try any non-localhost address
	if bestAddr == "" {
		for _, addr := range h.Addrs() {
			addrStr := addr.String()
			if !strings.Contains(addrStr, "127.0.0.1") && !strings.Contains(addrStr, "::1") {
				bestAddr = addrStr + "/p2p/" + h.ID().String()
				break
			}
		}
	}

	// If no public address found, fall back to localhost (for testing)
	if bestAddr == "" {
		bestAddr = h.Addrs()[0].String() + "/p2p/" + h.ID().String()
	}

	log.Printf("Adding self as bootstrap: %s", bestAddr)
	
	// Add self to bootstrap list
	if err := config.AddBootstrapPeer(bestAddr); err != nil {
		return err
	}
	
	// Automatically clean up inactive peers
	fmt.Println("Running automatic cleanup of inactive bootstrap peers...")
	CleanInactiveBootstrapPeers(ctx, h)
	
	return nil
}
// PrintBootstrapList prints the list of saved bootstrap peers
func PrintBootstrapList() {
	bootstrapConfig, err := config.LoadBootstrapConfig()
	if err != nil {
		log.Printf("Failed to load bootstrap config: %v", err)
		return
	}

	fmt.Println("Current bootstrap peers:")
	if len(bootstrapConfig.Peers) == 0 {
		fmt.Println("   (no bootstrap peers registered)")
		return
	}

	for i, peer := range bootstrapConfig.Peers {
		fmt.Printf("   %d. %s\n", i+1, peer)
	}
}

// CleanInactiveBootstrapPeers checks and removes inactive bootstrap peers
func CleanInactiveBootstrapPeers(ctx context.Context, h host.Host) {
	bootstrapConfig, err := config.LoadBootstrapConfig()
	if err != nil {
		log.Printf("Failed to load bootstrap config: %v", err)
		return
	}

	if len(bootstrapConfig.Peers) == 0 {
		fmt.Println("No bootstrap peers to clean")
		return
	}

	fmt.Println("Checking bootstrap peers for inactivity...")
	
	activePeers := []string{}
	inactiveCount := 0

	for i, addr := range bootstrapConfig.Peers {
		fmt.Printf("   [%d/%d] Checking: %s\n", i+1, len(bootstrapConfig.Peers), addr)
		
		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			fmt.Printf("      Invalid address format, removing\n")
			inactiveCount++
			continue
		}

		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			fmt.Printf("      Cannot parse address, removing\n")
			inactiveCount++
			continue
		}

		// Keep if it's our own node
		if info.ID == h.ID() {
			fmt.Printf("      This is self, keeping\n")
			activePeers = append(activePeers, addr)
			continue
		}

		// Check connectivity with short timeout
		checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		err = h.Connect(checkCtx, *info)
		cancel()

		if err != nil {
			fmt.Printf("      Inactive or unreachable, removing\n")
			inactiveCount++
			continue
		}

		fmt.Printf("      Active, keeping\n")
		activePeers = append(activePeers, addr)
		
		// Close test connection (optional)
		go h.Network().ClosePeer(info.ID)
	}

	// Update bootstrap peers list
	if inactiveCount > 0 {
		bootstrapConfig.Peers = activePeers
		if err := config.SaveBootstrapConfig(bootstrapConfig); err != nil {
			log.Printf("Failed to save cleaned bootstrap config: %v", err)
			return
		}
		fmt.Printf("\nCleanup complete: %d inactive peer(s) removed, %d active peer(s) remain\n", 
			inactiveCount, len(activePeers))
	} else {
		fmt.Printf("\nCleanup complete: No inactive peers found (%d active peers)\n", len(activePeers))
	}
}