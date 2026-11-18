package server

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func ListenAndServe(address string) error {
	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		return err
	}

	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("failed to add CA certificate from certs/ca.crt")
	}

	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS13,
	}

	listener, err := tls.Listen("tcp", address, serverTLSConfig)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("Server is listening on %s...\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Println("Rejected non-TLS connection from", conn.RemoteAddr())
		return
	}
	if err := tlsConn.Handshake(); err != nil {
		log.Println("mTLS handshake failed:", err)
		return
	}
	clientID := conn.RemoteAddr().String()
	if certs := tlsConn.ConnectionState().PeerCertificates; len(certs) > 0 {
		clientID = fmt.Sprintf("CN=%s (%s)", certs[0].Subject.CommonName, conn.RemoteAddr())
	}

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading client message:", err)
		return
	}
	incident := strings.TrimSpace(message)
	if incident == "" {
		log.Println("Empty message from", clientID, "- ignored")
		return
	}
	
	registeredAt := time.Now()
	fmt.Printf("[%s] Received incident from %s: %s\n",
		registeredAt.Format("2006-01-02 15:04:05"), clientID, incident)

	response := "[SERVER]: ALERT RECEIVED AND VERIFIED."
	if _, err = fmt.Fprintln(conn, response); err != nil {
		log.Println("Error writing response:", err)
	}

	if err := SendIncidentNotification(incident, clientID, registeredAt); err != nil {
		log.Printf("Failed to notify %s: %v\n", DutyEngineerEmail, err)
		return
	}
	log.Printf("Incident notification sent to %s\n", DutyEngineerEmail)
}
