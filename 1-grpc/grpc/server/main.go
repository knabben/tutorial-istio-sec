package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type greeterService struct {
	pb.UnimplementedGreeterServer
}

var (
	addr = flag.String("addr", "0.0.0.0:50051", "Server address")

	rootCA  = flag.String("rootca", "./tls/certs/root-cert.pem", "Root CA certificate")
	cert    = flag.String("cert", "./tls/certs/server-cert.pem", "Server certificate")
	keypair = flag.String("keypair", "./tls/certs/server-keypair.pem", "Server keypair")
)

func main() {
	flag.Parse()

	creds, err := generateCreds()
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	// register service into grpc server
	pb.RegisterGreeterServer(grpcServer, &greeterService{})

	// listen port
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("list port err: %v", err)
	}

	log.Printf("listening at %v", lis.Addr())

	// listen port
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve err: %v", err)
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

	// configuration of the certificate what we want to
	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
	}), nil
}
