package main

import (
	"fmt"
	"github.com/rjzak/tinywww"
	"os"
)

const LISTEN_STRING = "127.0.0.1:8080"

func main() {
	server, err := tinywww.NewTinyWWW(LISTEN_STRING)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not listen on %s: %s\n", LISTEN_STRING, err)
		os.Exit(1)
	}
	defer server.Close()

	server.HandleGet("/", func(resp *tinywww.TinyResponse, req *tinywww.TinyRequest) {
		fmt.Printf("Received %s request: %s\n", req.HTTP, req.Request)
		fmt.Printf("\tMethod: %s\n", req.Method)
		fmt.Printf("\tHeaders:\n")
		for key, value := range req.Headers {
			fmt.Printf("\t\t%s: %s\n", key, value)
		}
		fmt.Printf("\tForm:\n")
		for key, value := range req.Form {
			fmt.Printf("\t\t%s: %s\n", key, value)
		}
		fmt.Printf("\tBody: %s\n", string(req.Body))
		resp.Write([]byte("Hello, World!"))
	})

	if err := server.Listen(); err != nil {
		fmt.Fprintf(os.Stderr, "Error serving HTTP: %s\n", err)
	}
}
