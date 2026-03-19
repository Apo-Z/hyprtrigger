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

func socketPath() string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = "/tmp"
	}
	return filepath.Join(runtimeDir, "hyprtrigger.sock")
}

func NewDaemon() *Daemon {
	return &Daemon{
		socketPath:   socketPath(),
		reloadChan:   make(chan bool, 1),
		shutdownChan: make(chan bool, 1),
	}
}

func (d *Daemon) Start() error {
	if err := os.Remove(d.socketPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing socket: %w", err)
	}

	listener, err := net.Listen("unix", d.socketPath)
	if err != nil {
		return fmt.Errorf("failed to create daemon socket: %w", err)
	}

	d.listener = listener
	fmt.Printf("Daemon socket: %s\n", d.socketPath)
	go d.handleConnections()
	return nil
}

func (d *Daemon) handleConnections() {
	for {
		conn, err := d.listener.Accept()
		if err != nil {
			if d.stopped {
				return
			}
			fmt.Printf("Daemon accept error: %v\n", err)
			continue
		}
		go d.handleConnection(conn)
	}
}

func (d *Daemon) handleConnection(conn net.Conn) {
	defer conn.Close()

	var cmd Command
	if err := json.NewDecoder(conn).Decode(&cmd); err != nil {
		return
	}

	switch cmd.Type {
	case "reload":
		select {
		case d.reloadChan <- true:
			conn.Write([]byte("OK: Reload initiated\n"))
		default:
			conn.Write([]byte("OK: Reload already in progress\n"))
		}
	case "status":
		conn.Write([]byte("OK: Daemon is running\n"))
	case "shutdown":
		conn.Write([]byte("OK: Shutting down\n"))
		d.shutdownChan <- true
	default:
		conn.Write([]byte("ERROR: Unknown command\n"))
	}
}

func (d *Daemon) GetReloadChannel() <-chan bool   { return d.reloadChan }
func (d *Daemon) GetShutdownChannel() <-chan bool { return d.shutdownChan }

func (d *Daemon) Stop() {
	if d.stopped {
		return
	}
	d.stopped = true
	if d.listener != nil {
		d.listener.Close()
	}
	os.Remove(d.socketPath)
}

func IsDaemonRunning() bool {
	conn, err := net.Dial("unix", socketPath())
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func SendCommand(cmdType string) error {
	conn, err := net.Dial("unix", socketPath())
	if err != nil {
		return fmt.Errorf("failed to connect to daemon (is hyprtrigger running?): %w", err)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(Command{Type: cmdType}); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	fmt.Print(string(buf[:n]))
	return nil
}

func SendReload() error   { return SendCommand("reload") }
func SendStatus() error   { return SendCommand("status") }
func SendShutdown() error { return SendCommand("shutdown") }
