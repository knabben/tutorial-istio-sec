// Package main implements a client for Greeter service.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")

	rootCA  = flag.String("rootca", "./tls/certs/root-cert.pem", "Root CA certificate")
	cert    = flag.String("cert", "./tls/certs/client-cert.pem", "Client certificate")
	keypair = flag.String("keypair", "./tls/certs/client-keypair.pem", "Client keypair")
)

func main() {
	flag.Parse()

	creds, err := generateCreds()
	if err != nil {
		log.Fatal(err)
	}

	// create client connection
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

func generateCreds() (credentials.TransportCredentials, error) {
	pool := x509.NewCertPool()

	rootBytes, err := os.ReadFile(*rootCA)
	if err != nil {
		return nil, err
	}

	if !pool.AppendCertsFromPEM(rootBytes) {
		return nil, err
	}

	certificate, err := tls.LoadX509KeyPair(*cert, *keypair)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      pool,
	}), nil
}
