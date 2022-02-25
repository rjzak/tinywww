package tinywww

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/alphahorizonio/tinynet/pkg/tinynet"
	"os"
	"strings"
)

const BUFFER_LENGTH = 2048

type endpointFunction func(*TinyResponse, *TinyRequest)

type TinyWWW struct {
	Listener     tinynet.Listener
	UrlMapperGet map[string]endpointFunction
}

func NewTinyWWW(hostIP string) (*TinyWWW, error) {
	listenAddress, err := tinynet.ResolveTCPAddr("tcp", hostIP)
	if err != nil {
		return nil, err
	}

	listenerObject, err := tinynet.ListenTCP("tcp", listenAddress)
	if err != nil {
		return nil, err
	}

	twww := &TinyWWW{
		Listener:     listenerObject,
		UrlMapperGet: make(map[string]endpointFunction),
	}

	return twww, nil
}

func (tw *TinyWWW) Listen() error {
	return tw.ListenMaybeExit(false)
}

func (tw *TinyWWW) ListenMaybeExit(exitOnError bool) error {
	for {
		connection, err := tw.Listener.Accept()
		if err != nil {
			if exitOnError {
				return err
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Listening error: %s\n", err)
				continue
			}
		}
		buffer := make([]byte, BUFFER_LENGTH)
		_, err = connection.Read(buffer)
		if err != nil {
			if exitOnError {
				return err
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Listening error: %s\n", err)
				continue
			}
		}
		//fmt.Printf("Received: %s\n\n\n", buffer)

		requestedURL := string(buffer[:bytes.Index(buffer, []byte("\n"))])
		if strings.Contains(requestedURL, "?") {
			requestedURL = strings.Split(requestedURL, "?")[0]
		}
		if !strings.HasPrefix(requestedURL, "GET") {
			fmt.Printf("Not a GET request: %s\n", requestedURL)
			continue
		}

		requestedURL = strings.Split(requestedURL, " ")[1]
		handler, ok := tw.UrlMapperGet[requestedURL]
		if ok {
			response := NewTinyResponse()
			request := NewTinyRequestFromBuffer(buffer)
			handler(response, request)
			if _, err := connection.Write([]byte("HTTP/1.0 200\n")); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error writing HTTP response to connection: %s\n", err)
			}
			for key, value := range response.Headers {
				if _, err := connection.Write([]byte(fmt.Sprintf("%s: %s\n", key, value))); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Error writing headers to connection: %s\n", err)
				}
			}
			if _, err := connection.Write([]byte("\n\n")); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error writing to response connection: %s\n", err)
			}
			if _, err := connection.Write(response.Buffer); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error writing handler data to response connection: %s\n", err)
			}
		} else {
			if _, err := connection.Write([]byte("HTTP/1.0 404\n\nNot Found.")); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error writing 404 message to response connection: %s\n", err)
			}
			fmt.Printf("No handler for URL %s\n", requestedURL)
		}
		if err := connection.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error closing response connection: %s\n", err)
		}
	}
	return errors.New("server interrupted")
}

func (tw *TinyWWW) Close() error {
	return tw.Listener.Close()
}

func (tw *TinyWWW) HandleGet(pattern string, handler func(resp *TinyResponse, req *TinyRequest)) {
	tw.UrlMapperGet[pattern] = handler
}
