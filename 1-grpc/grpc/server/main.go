package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/knabben/tutorial-istio-sec/1-grpc/grpc/proto"
)

var (
	addr    = flag.String("addr", "0.0.0.0:50051", "Server address")
	rootCA  = flag.String("rootca", "./tls/certs/root-cert.pem", "Root CA certificate")
	cert    = flag.String("cert", "./tls/certs/server-cert.pem", "Server certificate")
	keypair = flag.String("keypair", "./tls/certs/server-keypair.pem", "Server keypair")
)

var (
	total = map[string]float32{}
)

type casino struct {
	pb.UnimplementedCasinoServer
}

func (s *casino) PayIt(ctx context.Context, in *pb.TransactionRequest) (*pb.TransactionReply, error) {
	if in.GetClient() == "" {
		return &pb.TransactionReply{Error: true}, nil
	}

	total[in.GetClient()] += in.GetValue()
	log.Printf("Received %v from %s, total: %v", in.GetValue(), in.GetClient(), total[in.GetClient()])

	return &pb.TransactionReply{Error: false}, nil
}

func main() {
	flag.Parse()

	creds, err := generateCreds()
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterCasinoServer(srv, &casino{})

	listen, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("list port err: %v", err)
	}

	log.Printf("listening at %v", listen.Addr())
	if err := srv.Serve(listen); err != nil {
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
