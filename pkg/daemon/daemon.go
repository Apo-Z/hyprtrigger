package daemon

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

type Daemon struct {
	listener     net.Listener
	socketPath   string
	reloadChan   chan bool
	shutdownChan chan bool
	stopped      bool
}

type Command struct {
	Type string `json:"type"`
}

func NewDaemon() *Daemon {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = "/tmp"
	}

	socketPath := filepath.Join(runtimeDir, "hyprtrigger.sock")

	return &Daemon{
		socketPath:   socketPath,
		reloadChan:   make(chan bool, 1),
		shutdownChan: make(chan bool, 1),
		stopped:      false,
	}
}

func (d *Daemon) Start() error {
	// Remove existing socket if it exists
	if err := os.Remove(d.socketPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing socket: %w", err)
	}

	listener, err := net.Listen("unix", d.socketPath)
	if err != nil {
		return fmt.Errorf("failed to create daemon socket: %w", err)
	}

	d.listener = listener
	fmt.Printf("Daemon socket created: %s\n", d.socketPath)

	go d.handleConnections()
	return nil
}

func (d *Daemon) handleConnections() {
	for {
		conn, err := d.listener.Accept()
		if err != nil {
			// Check if we're shutting down
			select {
			case <-d.shutdownChan:
				return
			default:
				// Check if listener was closed
				if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
					return
				}
				fmt.Printf("Daemon accept error: %v\n", err)
				continue
			}
		}

		go d.handleConnection(conn)
	}
}

func (d *Daemon) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	var cmd Command

	if err := decoder.Decode(&cmd); err != nil {
		fmt.Printf("Daemon decode error: %v\n", err)
		return
	}

	switch cmd.Type {
	case "reload":
		fmt.Println("Reload command received")
		select {
		case d.reloadChan <- true:
			conn.Write([]byte("OK: Reload initiated\n"))
		default:
			conn.Write([]byte("OK: Reload already in progress\n"))
		}
	case "status":
		conn.Write([]byte("OK: Daemon is running\n"))
	case "shutdown":
		fmt.Println("Shutdown command received")
		conn.Write([]byte("OK: Shutting down\n"))
		d.shutdownChan <- true
	default:
		conn.Write([]byte("ERROR: Unknown command\n"))
	}
}

func (d *Daemon) GetReloadChannel() <-chan bool {
	return d.reloadChan
}

func (d *Daemon) GetShutdownChannel() <-chan bool {
	return d.shutdownChan
}

func (d *Daemon) Stop() {
	if d.stopped {
		return
	}
	d.stopped = true

	fmt.Println("Stopping daemon...")

	// Close listener to stop accepting new connections
	if d.listener != nil {
		d.listener.Close()
	}

	// Clean up socket file
	os.Remove(d.socketPath)

	// Close channels safely
	close(d.shutdownChan)
	close(d.reloadChan)
}

// Check if daemon is already running
func IsDaemonRunning() bool {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = "/tmp"
	}

	socketPath := filepath.Join(runtimeDir, "hyprtrigger.sock")

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Client functions for sending commands
func SendCommand(cmdType string) error {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = "/tmp"
	}

	socketPath := filepath.Join(runtimeDir, "hyprtrigger.sock")

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to daemon (is hyprtrigger running?): %w", err)
	}
	defer conn.Close()

	cmd := Command{Type: cmdType}
	encoder := json.NewEncoder(conn)

	if err := encoder.Encode(cmd); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	// Read response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	response := string(buffer[:n])
	fmt.Print(response)

	return nil
}

func SendReload() error {
	return SendCommand("reload")
}

func SendStatus() error {
	return SendCommand("status")
}

func SendShutdown() error {
	return SendCommand("shutdown")
}
