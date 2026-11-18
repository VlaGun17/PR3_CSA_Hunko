package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	clientCert, err := tls.LoadX509KeyPair("../certs/client.crt", "../certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := os.ReadFile("../certs/ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatal("failed to add CA certificate from certs/ca.crt")
	}

	clientTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS13,
	}

	conn, err := tls.Dial("tcp", "localhost:8443", clientTLSConfig)
	if err != nil {
		log.Fatal(err)
	}
	if err := sendTestMessage(conn); err != nil {
		log.Fatal(err)
	}
}

func sendTestMessage(conn net.Conn) error {
	defer conn.Close()
	response := "CRITICAL_ERROR: Database connection timeout on node-04."
	_, err := fmt.Fprintln(conn, response)
	if err != nil {
		return err
	}
	fmt.Println("Message sent successfully!")

	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("Received reply from server: %s\n", reply)
	return nil
}