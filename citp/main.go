package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const (
	// CITP multicast address and port
	multicastAddress = "224.0.0.180:4809"
	serverPort       = ":4809"
)

// CITP Header structure
type CITPHeader struct {
	Protocol string
	Version  uint16
	Message  string
}

// CITP Info message structure
type CITPInfoMessage struct {
	Header  CITPHeader
	Details string
}

// Create a CITP header
func createCITPHeader(message string) CITPHeader {
	return CITPHeader{
		Protocol: "CITP",
		Version:  1,
		Message:  message,
	}
}

// Create an Info message
func createInfoMessage(details string) CITPInfoMessage {
	return CITPInfoMessage{
		Header:  createCITPHeader("PINF"),
		Details: details,
	}
}

// Serialize the CITP Info message to bytes
func serializeInfoMessage(info CITPInfoMessage) ([]byte, error) {
	var buf bytes.Buffer

	// Write Protocol
	if err := binary.Write(&buf, binary.BigEndian, []byte(info.Header.Protocol)); err != nil {
		return nil, err
	}

	// Write Version
	if err := binary.Write(&buf, binary.BigEndian, info.Header.Version); err != nil {
		return nil, err
	}

	// Write Message Type
	if err := binary.Write(&buf, binary.BigEndian, []byte(info.Header.Message)); err != nil {
		return nil, err
	}

	// Write Details
	if err := binary.Write(&buf, binary.BigEndian, []byte(info.Details)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Send UDP multicast message
func sendMulticastMessage(conn *net.UDPConn, message []byte) error {
	_, err := conn.Write(message)
	return err
}

// Send the multicast message
func sendMulticast() error {
	// Resolve the multicast address
	addr, err := net.ResolveUDPAddr("udp", multicastAddress)
	if err != nil {
		return err
	}

	// Create the UDP socket
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create and serialize the CITP Info message
	infoMessage := createInfoMessage("Media Server Available")
	serializedMessage, err := serializeInfoMessage(infoMessage)
	if err != nil {
		return err
	}

	// Send the message over the UDP socket
	return sendMulticastMessage(conn, serializedMessage)
}

// Handle incoming TCP connections from clients
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New TCP connection established.")

	// Send a response when a client connects
	infoMessage := createInfoMessage("Media Server Connected")
	serializedMessage, err := serializeInfoMessage(infoMessage)
	if err != nil {
		log.Fatalf("Error serializing message: %v", err)
	}

	// Send response message to client
	_, err = conn.Write(serializedMessage)
	if err != nil {
		log.Fatalf("Error sending response to client: %v", err)
	}

	// Read incoming messages from the client (for now, just log them)
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error reading from client: %v", err)
			break
		}
		fmt.Printf("Received from client: %s\n", string(buf[:n]))
	}
}

// Start listening for incoming TCP connections
func startTCPServer() {
	// Listen on TCP port
	listener, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer listener.Close()
	fmt.Println("Media Server listening for TCP connections on port", serverPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting TCP connection: %v", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

func main() {
	// Start sending multicast messages
	err := sendMulticast()
	if err != nil {
		log.Fatalf("Error sending multicast message: %v", err)
	}
	fmt.Println("Multicast message sent to announce Media Server availability.")

	// Start listening for incoming TCP connections
	go startTCPServer()

	// Keep the server running
	select {}
}
