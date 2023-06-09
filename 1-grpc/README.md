# 1 - gRPC secure service 

## Generating the TLS CA PKI

For self-signed certificates use `mage` with TLS namespace.

```shell
tls:genRootCA        1 - Generates a new Root CA and CA key pair
tls:genServerCert    2 - Generates a server CSR and server cert
tls:genClientCert    3 - Generates a client CSR and client cert
```

## Running the golang gRPC

To run the gRPC examples in the folder run, server runs on 50051 TCP by default:

```shell
grpc:runClient       runs a TLS based gRPC client
grpc:runServer       runs a TLS based gRPC server
```