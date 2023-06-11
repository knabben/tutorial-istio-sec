package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffegrpc/grpccredentials"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	host     = flag.String("host", "172.18.0.2:50051", "Host gRPC server")
	sockPath = flag.String("sock", "unix:///run/spire/sockets/agent.sock", "SPIRE Agent socket")
	allowID  = flag.String("allowID", "spiffe://opssec.in/ns/default/sa/server", "Server SPIFFEID")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	log.Printf("connecting on %s", *host)

	source, err := workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithAddr(*sockPath)))
	if err != nil {
		return fmt.Errorf("unable to create X509Source: %w", err)
	}
	defer source.Close()

	// Allowed SERVER SPIFFE ID
	serverID := spiffeid.RequireFromString(*allowID)

	// Dial the server with credentials that do mTLS and verify that presented certificate has SPIFFE ID `spiffe://example.org/server`
	conn, err := grpc.DialContext(ctx, *host, grpc.WithTransportCredentials(
		grpccredentials.MTLSClientCredentials(source, source, tlsconfig.AuthorizeID(serverID)),
	))
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	for {
		client := pb.NewGreeterClient(conn)
		reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "world"})
		if err != nil {
			return fmt.Errorf("failed issuing RPC to server: %w", err)
		}

		log.Print(reply.Message)
		time.Sleep(time.Second)
	}

	return nil
}
