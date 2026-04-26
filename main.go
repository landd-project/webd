package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"webd/internal/request"
	"webd/internal/config"
)

func main() {
	config, err := config.GetConfig();
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to get config: %v\n", err);
		os.Exit(1);	
	}

	os.Remove(config.SocketPath);

	listener, err := net.Listen("unix", config.SocketPath);
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to listen: %v\n", err);
		os.Exit(1);
	}

	for {
		conn, err := listener.Accept();
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to accept connection: %v\n", err);
		}

		go handleConnection(conn);
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close();
	reader := bufio.NewReader(conn);
	req, err := reader.ReadString('\n');
	if err != nil {
		if err == io.EOF {
			return
		}
		fmt.Fprintf(os.Stderr, "ERROR: failed to read request: %v\n", err);
	}

	err = request.ParseRequest(req);
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to read request: %v\n", err);
	}
}
