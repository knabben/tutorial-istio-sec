# 1 - gRPC secure service 

## Generating the TLS CA PKI

For self-signed certificates use `mage` with:

```
  tls:genCAKeyPair     1 - Create a new public/private key pair under the certs folder
  tls:genCARootCSR     2 - Creates a new CSR file and print
  tls:genCARootCert    3 - Generate CA root certificate from CSR
```
