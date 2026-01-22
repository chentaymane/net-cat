package functions

import "fmt"

func Add(c *Client) string {
	mu.Lock()
	defer mu.Unlock()

	if len(clients) >= MAX_CLIENT {
		return "capacity"
	}

	if _, exists := clients[c.name]; exists {
		return "name"
	}

	clients[c.name] = c
	return ""
}

// Remove client from map
func Remove(c *Client) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, c.name)
}

// Rename a client safely
func Rename(c *Client, newName string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := clients[newName]; exists {
		return fmt.Errorf("name already exists")
	}
	oldName := c.name
	c, ok := clients[oldName]
	if !ok {
		return fmt.Errorf("client not found")
	}

	delete(clients, oldName)
	c.name = newName
	clients[newName] = c
	return nil
}
