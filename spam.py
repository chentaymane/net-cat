import socket
import time

HOST = "127.0.0.1"
PORT = 8080
MAX_CLIENTS = 15

def connect_user(i):
    name = f"user{i}"
    print(f"Connecting {name}")

    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((HOST, PORT))

        # Receive server prompt
        data = s.recv(1024).decode()
        print(f"{name} received:", data.strip())

        if "Server is full" in data:
            s.close()
            return

        # Send name
        s.sendall((name + "\n").encode())
        time.sleep(0.2)

        # Send a message
        s.sendall((f"hello from {name}\n").encode())

        return s  # keep connection alive

    except Exception as e:
        print(f"{name} error:", e)
        return None


sockets = []

for i in range(1, MAX_CLIENTS + 1):
    sock = connect_user(i)
    if sock:
        sockets.append(sock)

    time.sleep(0.5)  # force sequential join

print("\nAll clients attempted. Keeping connections alive...")
time.sleep(30)

for s in sockets:
    s.close()
