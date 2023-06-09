package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/knabben/tutorial-istio-sec/1-grpc/grpc/proto"
)

var (
	addr    = flag.String("addr", "localhost:50051", "the address to connect to")
	rootCA  = flag.String("rootca", "./tls/certs/root-cert.pem", "Root CA certificate")
	cert    = flag.String("cert", "./tls/certs/client-cert.pem", "Client certificate")
	keypair = flag.String("keypair", "./tls/certs/client-keypair.pem", "Client keypair")
)

func main() {
	flag.Parse()

	client := getCertSubject()
	creds, err := generateCreds()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		c := pb.NewCasinoClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		r, err := c.PayIt(ctx, &pb.TransactionRequest{
			Value: float32(time.Now().Second()), Client: client,
		})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Error = %v", r.GetError())
		cancel()
		time.Sleep(2 * time.Second)
	}
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

func getCertSubject() string {
	data, err := os.ReadFile(*cert)
	if err != nil {
		log.Fatal(err)
	}
	dec, _ := pem.Decode(data)
	crt, _ := x509.ParseCertificate(dec.Bytes)
	return crt.Subject.String()
}
