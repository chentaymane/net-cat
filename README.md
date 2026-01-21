🐱 Net-Cat

Net-Cat is a high-performance, terminal-based TCP chat server written in Go. It allows multiple users to communicate in a centralized chat room in real-time, featuring private mentions, command-based interactions, and full message history for late-joiners.
✨ Features

    Concurrent Handling: Built with Go goroutines to handle multiple clients simultaneously.

    Persistent History: New users receive the full chat history upon joining.

    User Commands: * /rename <new_name>: Change your identity on the fly.

        /users: See who else is currently online.

    Direct Mentions: Use @username to highlight a message for a specific user.

    Thread Safety: Utilizes sync.Mutex to ensure data integrity across concurrent connections.

    Visual Feedback: ANSI color coding for timestamps, usernames, and system notifications.

🛠 Installation
Prerequisites

    Go (version 1.20 or higher)

    A terminal with nc (netcat) or telnet installed.

Build from Source
Bash

# Clone the repository
git clone https://github.com/chentaymane/net-cat.git
cd net-cat

# Build the executable
go build -o net-cat main.go

⚡ Usage
1. Start the Server

By default, the server listens on port 8989. You can specify a custom port as an argument.
Bash

./net-cat [port]

# Example
./net-cat 8080

2. Connect as a Client

Open a new terminal window and connect using netcat:
Bash

nc localhost 8080

📋 Interaction Guide
Action	Input / Command	Result
Join	Enter your name when prompted	You enter the chat room
Chat	Type anything and press Enter	Broadcasts to everyone
Mention	@username <message>	Highlights message for that user
Rename	/rename <new_name>	Updates your display name
List Users	/users	Displays all active participants
Example Chat Flow
Plaintext

[ENTER YOUR NAME]: achent
[2026-01-17 20:50:25][achent]: Hello everyone!
karim has joined our chat...
[2026-01-17 20:50:48][karim]: @achent Hi there!
[2026-01-17 20:51:00][achent]: /rename RedFox
[2026-01-17 20:51:10][RedFox]: I am now RedFox!

📂 Project Structure
Plaintext

net-cat/
├── main.go            # Entry point & TCP Listener setup
├── functions/         
│   ├── client.go      # Client struct & lifecycle management
│   ├── messages.go    # Broadcast & Mention logic
│   └── utils.go       # Input validation & formatting
├── go.mod             # Go module definition
└── README.md

⚙️ Technical Overview

    Concurrency: Every client connection triggers a dedicated goroutine, allowing the server to scale effortlessly.

    Validation: * Usernames: Must be unique, alphanumeric, and under 10 characters.

        Messages: Length-restricted to 100 characters to prevent buffer issues.

    Architecture: The server maintains a map[string]Client for O(1) lookups during mentions and a shared buffer for chat history.
