package socket

import (
	"fmt"
	"net"
	"os"
)

type Client struct {
	conn net.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect() error {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	instanceSig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")

	socketPath := fmt.Sprintf("%s/hypr/%s/.socket2.sock", runtimeDir, instanceSig)

	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") == "" {
		socketPath = "/tmp/hypr/hyprland.sock2"
	}

	fmt.Printf("Connecting to socket: %s\n", socketPath)

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return fmt.Errorf("socket connection failed: %w", err)
	}

	c.conn = conn
	fmt.Println("Connected to Hyprland socket, listening for events...")
	fmt.Println("Press Ctrl+C to stop\n")

	return nil
}

func (c *Client) GetConnection() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
