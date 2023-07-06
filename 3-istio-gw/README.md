## Istio Gateway Example

### EKS production cluster

Create a new EKS cluster (default is enough)

```shell
eksctl create cluster --name eks-demo2 --region us-east-1 
```

### Local development environment

Mage will provide targets for running local clusters

```shell
sM3:install         installs kind and metallb into the cluster
sM3:delete          cleans up kind from cluster
```

## Installing Istio Ambient

```shell
sM3:installIstio    install ambient
sM3:deleteIstio     cleans up resources from cluster
```

## Deploying the DEMO app

Deploying the DEMO app and policies can be achieved with

```shell
sM3:deploy          creates the pre-defined application topology
sM3:policies        create a VirtualService and define application Authorization files
```
