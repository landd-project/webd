package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"webd/internal/config"
	"webd/internal/ipc"
)

func main() {
	config, err := config.GetConfig();
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to get config: %v\n", err);
		os.Exit(1);	
	}

	listener, err := net.Listen("tcp", config.TcpHost);
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

	response := ipc.ParseRequest(req);

	bt, err := json.MarshalIndent(response, "", "	");
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to marshal json: %v\n", err) // TODO: maybe send this to the client.
	}

	conn.Write(bt);
}
