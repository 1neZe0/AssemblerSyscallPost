package main

import (
	"fmt"
	"net"
	"syscall"
)

const telegramAPI = "api.telegram.org:443"

func sendMessage(apiToken string, chatID int64, message string) error {
	// Create the socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	// Set the socket timeout
	tv := syscall.Timeval{
		Sec:  5,
		Usec: 0,
	}
	err = syscall.SetsockoptTimeval(fd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv)
	if err != nil {
		return err
	}

	// Connect to the server
	addr := syscall.SockaddrInet4{
		Port: 848,
		Addr: [4]byte{127, 0, 0, 1},
	}
	copy(addr.Addr[:], net.ParseIP(telegramAPI).To4())
	err = syscall.Connect(fd, &addr)
	if err != nil {
		return err
	}
	// ... code to create the socket and connect to the server

	// Build the HTTP request
	request := fmt.Sprintf("POST /bot%s/sendMessage HTTP/1.1\r\nHost: %s\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n{\"chat_id\":%d,\"text\":\"%s\"}", apiToken, telegramAPI, len(message), chatID, message)
	_, err = syscall.Write(fd, []byte(request))
	if err != nil {
		return err
	}

	// Read the response
	var response [1024]byte
	n, err := syscall.Read(fd, response[:])
	if err != nil {
		return err
	}

	// Check the response status code
	if n < 12 || string(response[:12]) != "HTTP/1.1 200" {
		return fmt.Errorf("invalid response: %s", string(response[:n]))
	}

	return nil
}

func createServer(port int) error {
	// Create the socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	// Bind the socket to the specified port
	addr := syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte{0, 0, 0, 0},
	}
	err = syscall.Bind(fd, &addr)
	if err != nil {
		return err
	}

	// Listen for incoming connections
	err = syscall.Listen(fd, 5)
	if err != nil {
		return err
	}

	// Accept incoming connections in a loop
	for {
		// Accept the incoming connection
		connFd, _, err := syscall.Accept(fd)
		if err != nil {
			return err
		}
		defer syscall.Close(connFd)

		// Read the incoming request
		var request [1024]byte
		n, err := syscall.Read(connFd, request[:])
		if err != nil {
			return err
		}

		// Process the request and send the response
		response := processRequest(string(request[:n]))
		_, err = syscall.Write(connFd, []byte(response))
		if err != nil {
			return err
		}
	}
}

func processRequest(request string) string {
	// Parse the request to get the chat ID and message
	var chatID int64
	var message string
	fmt.Sscanf(request, "{\"chat_id\":%d,\"text\":\"%s\"}", &chatID, &message)

	// Process the message and generate a response
	response := "Your message was: " + message

	fmt.Println(response)

	return response
}

func main() {
	go func() {
		createServer(88)
	}()

	// Replace these with your own values
	apiToken := "5185078921:AAHiAV-xaDHa5UGBbQI91KjdlYxk5emlvso"
	chatID := int64(434057604)
	message := "Hello, world!"

	err := sendMessage(apiToken, chatID, message)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Message sent successfully")
	}
}
