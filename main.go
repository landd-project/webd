package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"webd/internal"
)

/*
 GET url TAB 2
 TRUST url
 UNTRUST url

 TAB NEW BLANK
 TAB NEW url
 TAB SEL 0
 TAB DEL 1

{"geminiprotocol.net", "bbs.geminispace.org", "kennedy.gemi.dev/search?abobora"}
*/

func main() {
	socketPath :=  "/tmp/webd.sock";
	os.Remove(socketPath);

	listener, err := net.Listen("unix", socketPath);
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
	request, err := reader.ReadString('\n');
	if err != nil {
		if err == io.EOF {
			return
		}
		fmt.Fprintf(os.Stderr, "ERROR: failed to read request: %v\n", err);
	}

	err = internal.ParseRequest(request);
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to read request: %v\n", err);
	}
}
