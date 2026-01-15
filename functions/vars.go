package functions

import (
	"net"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
	mu   sync.Mutex
}

var (
	clients = make(map[net.Conn]*Client)
	mu      sync.Mutex
	history []string
)


const (
	MAX_CLIENT   = 10
	allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 .,!?;:'-()[]{}@#$%&*+=/_\"\\\t"
)

const logo = "\033[38;5;208mWelcome to TCP-Chat!\033[0m\n" +
    "\033[1;37m         _nnnn_\033[0m\n" +
    "\033[1;37m        dGGGGMMb\033[0m\n" +
    "\033[1;37m       @p~qp~~qMb\033[0m\n" +
    "\033[1;37m       M|\033[38;5;208m@\033[1;37m||\033[38;5;208m@\033[1;37m) M|\033[0m\n" +
    "\033[1;37m       @,----.JM|\033[0m\n" +
    "\033[1;37m      JS^\\__/  qKL\033[0m\n" +
    "\033[1;37m     dZP        qKRb\033[0m\n" +
    "\033[1;37m    dZP          qKKb\033[0m\n" +
    "\033[1;37m   fZP            SMMb\033[0m\n" +
    "\033[1;37m   HZM            MMMM\033[0m\n" +
    "\033[1;37m   FqM            MMMM\033[0m\n" +
    "\033[38;5;208m __| \".        |\\dS\"qML\033[0m\n" +
    "\033[38;5;208m |    `.       | `' \\Zq\033[0m\n" +
    "\033[38;5;208m_)      \\.___.,|     .'\033[0m\n" +
    "\033[38;5;208m\\____   )MMMMMP|   .'\033[0m\n" +
    "\033[38;5;208m     `-'       `--'\033[0m\n"
