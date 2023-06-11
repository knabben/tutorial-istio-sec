# Kubernetes SPIRE over gRPC

The following `mage` commands are available to manage SPIRE and SPIFFE client/server
go-spiffe example.


```shell
k8S:containerClient    build the SPIFFE gRPC client and push to registry
k8S:containerServer    build the SPIFFE gRPC server and push to registry
k8S:delete             cleans up resources from cluster
k8S:install            installs resources into the cluster
```
