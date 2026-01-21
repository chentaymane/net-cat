Perfect! I can create a **highly polished, GitHub-ready README** for **Net-Cat** that’s professional, visually clear, and structured with sections, code blocks, and ASCII diagrams to illustrate the chat flow. Here’s a full version:

````markdown
# Net-Cat

![Net-Cat Logo](https://img.shields.io/badge/TCP-Chat-orange)  

**Net-Cat** is a lightweight terminal-based TCP chat server written in **Go**, designed for real-time messaging with a single chat group. It supports colored messages, private mentions, and essential user commands.

---

## 🚀 Features

- **Real-time messaging** between multiple clients.
- **Single chat group** support (all clients in one room).
- **User commands**:
  - `/rename <new_name>` – Change your display name.
  - `/users` – List all connected users.
- **Private mentions** using `@username`.
- **ANSI-colored interface** for readability in terminal.
- **Chat history** for newly joined clients.
- **Input validation** to prevent invalid or unsafe messages.
- **Lightweight & fast**, using Go goroutines for concurrent client handling.

---

## 💻 Installation

### Prerequisites

- Go (1.20+)
- Git
- Terminal with TCP client (`nc` or `telnet`)

### Clone & Build

```bash
git clone https://github.com/yourusername/net-cat.git
cd net-cat
go build -o net-cat main.go
````

This generates a `net-cat` executable in your directory.

---

## ⚡ Usage

### Start Server

By default, Net-Cat runs on port **8989**. You can specify a custom port:

```bash
./net-cat [port]
```

Example:

```bash
./net-cat 8080
```

### Connect Clients

Clients can connect using **netcat**:

```bash
nc localhost 8080
```

Once connected:

```
[ENTER YOUR NAME]: achent
[2026-01-17 20:50:25][achent]:
```

---

## 📝 Commands

| Command               | Description                            |
| --------------------- | -------------------------------------- |
| `/rename <new_name>`  | Change your display name               |
| `/users`              | List all active users                  |
| `@username <message>` | Send a private mention to another user |

**Example:**

```
@achent Hey! How are you?
/rename RedFox
/users
```

---

## 💾 Chat History

* All messages are stored in memory (`history`) and displayed to new clients upon joining.
* Includes **timestamps**, **join/leave notifications**, and **renames**.

---

## 🖥 Chat Flow Example

```
Welcome to TCP-Chat!

[ENTER YOUR NAME]: achent
[2026-01-17 20:50:25][achent]: Hello everyone!
karim has joined our chat...
[2026-01-17 20:50:48][karim]: @achent Hi there!
[2026-01-17 20:51:00][achent]: /rename RedFox
[2026-01-17 20:51:10][RedFox]: Welcome back!
```

---

## 📦 Code Overview

### Client Structure

```go
type Client struct {
    conn net.Conn // TCP connection
    name string  // Unique username
}
```

### Core Functions

* **HandleClient** – Handles each client in a separate goroutine.
* **broadcast** – Sends messages to all clients (except sender).
* **Tag** – Sends private messages to a specific user.
* **DeleteClient / RenameClient** – Manage clients safely.
* **sendPrompt** – Prints a live prompt with timestamp for the client.
* **validName / validMsg** – Input validation functions.

### Concurrency

* Uses **goroutines** to handle multiple clients simultaneously.
* Shared resources (clients map, history) are protected with **sync.Mutex** for thread safety.

---

## 🔒 Input Validation

* Usernames: Alphanumeric, max 10 characters, unique.
* Messages: Only allowed characters, max 100 characters.

---

## 📂 Project Structure

```
net-cat/
│
├─ main.go         # Entry point, TCP listener
├─ functions/      # Core chat functions
│  ├─ client.go
│  ├─ messages.go
│  └─ utils.go
├─ go.mod
└─ README.md
```

---

## 🤝 Contributing

Contributions welcome!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/myfeature`)
3. Commit changes (`git commit -m 'Add feature'`)
4. Push to branch (`git push origin feature/myfeature`)
5. Open a Pull Request

---

## 📜 License

MIT License – See [LICENSE](LICENSE) for details.

---

## ⚡ Notes

* Only a **single group** chat is supported.
* ANSI escape codes used for coloring; terminal must support it.
* Lightweight and easy to extend.

---

**Author:** achent
**Project:** Net-Cat
**Language:** Go
**Repository:** (https://github.com/chentaymane/net-cat)


