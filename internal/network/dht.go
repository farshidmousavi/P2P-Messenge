package network

import (
	"context"
	"log"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

func SetupDHT(ctx context.Context, h host.Host) (*dht.IpfsDHT, error) {
	kadDHT, err := dht.New(ctx, h)
	if err != nil {
		return nil, err
	}

	if err := kadDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}

	log.Println("DHT bootstrapped")

	return kadDHT, nil
}
