
# P2P Messenger - Decentralized Peer-to-Peer Messenger

A fully decentralized messenger based on libp2p protocol that works without any central server. This application enables group chat, automatic peer discovery, and NAT traversal.

## ✨ Features

- **No Central Server** - Direct peer-to-peer communication
- **Resilient Network** - Network stays alive even if some nodes go offline
- **NAT Traversal** - Connect through home and corporate firewalls
- **Group Chat** - All users can communicate in a public room
- **Automatic Discovery** - Nodes automatically find each other
- **Dynamic Bootstrap** - Register any node as a network entry point

## 📋 Prerequisites

- **Go 1.21 or higher** - [Download from golang.org](https://golang.org/dl/)
- No internet connection required for running (offline-ready)

## 🚀 Quick Start

### Running the Application (Offline Ready)

All dependencies are already included in the `vendor` folder, so **no internet connection is required** to run the program!

```bash
# Clone the project
git clone https://github.com/d4rkpc/p2p-messenger.git
cd p2p-messenger

# Run directly (using included vendor dependencies)
go run -mod=vendor ./cmd/node/main.go

```

Note: The vendor folder contains all the required Go packages. You don't need to download anything - just run the command above!

First Time Setup (If vendor folder doesn't exist)
If you're building from source and the vendor folder is missing:

```bash
# Download dependencies (requires internet, use Iranian proxy if needed)
go env -w GOPROXY=https://goproxy.ir,direct
go env -w GOSUMDB=off
go mod download

# Create vendor folder for offline distribution
go mod vendor
```

## 📖 Usage Guide
Starting the Program
After running the program, it will automatically connect to saved Bootstrap peers and join the general chat room:

## Output
🚀 Starting P2P Messenger...
✅ Node ID: 12D3KooWBwUZ5juuwos9MzXByjGwzBw87d3qYAyAgdHmHEXTV4kN
📡 Listening on: /ip4/192.168.1.1/tcp/63730/p2p/12D3KooWB...
🌐 DHT bootstrapped
✅ Joined room: general-chat

✅ Connected! Type your message or commands:
📝 Commands: 'bootstrap', 'list-bootstrap', 'connect-bootstrap', 'clean-bootstrap'
✏️ You:

-------------------------------------------------------
## Available Commands
Command	                    Description	Example
- "bootstrap":	            Register current node as a Bootstrap (network entry point)	bootstrap
- "list-bootstrap"	        Show list of saved Bootstrap peers	list-bootstrap
- "connect-bootstrap"	    Connect to all saved Bootstrap peers	ct bootstrap
- "clean-bootstrap"	        Automatically remove inactive Bootstrap peers	clean-bootstrap
- "connect <address>"	    Manually connect to a specific peer	connect /ip4/192.168.88.254/tcp/63730/p2p/12D3KooW...
- "<message text>"           Send message to the public chat room	Hello everyone!
< br / >
📖 To use commands you simply need to type the command below and hit the ienter

## Sending Messages
Simply type your message and press Enter:

text
✏️ You: Hello this is a test
💬 [HEXTV4kN]: Hello this is a test
🏗️ Setting Up a Personal Bootstrap Node
To have a stable node (like a server) act as a network entry point:

On your server (with public IP):
bash
go run -mod=vendor ./cmd/node/main.go
# After startup, type:
bootstrap
On client machines:
bash
go run -mod=vendor ./cmd/node/main.go
# The program will automatically connect to the Bootstrap
🔧 Troubleshooting
1. connection refused or dial backoff errors
This means a Bootstrap peer is offline. To fix:

bash
# Remove inactive Bootstrap peers
clean-bootstrap

# Reconnect to active Bootstrap peers
ct bootstrap
2. Two nodes can't find each other
Solutions:

Wait a few minutes for DHT to propagate information

Use the ct bootstrap command

Manually connect using connect <address>

3. How to share the program with others
Simply zip the entire project folder (including the vendor directory) and share it. The recipient only needs Go installed and can run:

bash
go run -mod=vendor ./cmd/node/main.go
No internet connection required!

📁 Project Structure
text
p2p-messenger/
├── cmd/
│   └── node/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Bootstrap configuration management
│   ├── crypto/                  # Cryptography module (currently disabled)
│   │   ├── identity.go
│   │   └── signature.go
│   └── network/
│       ├── bootstrap.go         # Bootstrap peer management
│       ├── chat_room.go         # Chat room with GossipSub
│       ├── connection_manager.go # Connection management
│       ├── dht.go               # DHT setup
│       ├── discovery.go         # Automatic peer discovery
│       ├── node.go              # libp2p node setup
│       └── peers.go             # Peer listing utilities
├── vendor/                      # All dependencies (offline-ready)
├── go.mod                       # Project dependencies
└── README.md                    # This file
💾 Configuration Storage
Bootstrap settings are stored in:

Windows: C:\Users\[Username]\.p2p-messenger\bootstrap.json

Linux/Mac: ~/.p2p-messenger/bootstrap.json

Example file content:

json
{
  "peers": [
    "/ip4/192.168.88.254/tcp/63730/p2p/12D3KooWBwUZ5juuwos9MzXByjGwzBw87d3qYAyAgdHmHEXTV4kN"
  ]
}
🤝 Contributing
Fork the repository

Create a feature branch (git checkout -b feature/amazing-feature)

Commit your changes (git commit -m 'Add amazing feature')

Push to the branch (git push origin feature/amazing-feature)

Open a Pull Request

📜 License
This project is licensed under the MIT License.

⚠️ Privacy & Security Notes
Privacy: Your IP is only visible to directly connected peers

Security: Message encryption is currently disabled (can be added)

Persistence: Network depends on stable Bootstrap peers

For enhanced privacy, consider using a VPN or running your own relay node.

