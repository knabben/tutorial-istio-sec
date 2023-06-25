## Istio Gateway Example

Use `INSTALL_KIND=y` for `sm:install` and `sm:delete` targets if you want to operate on Kind

```
Targets:
  sm:install   -- install resources (kind/istio) into the cluster
  sm:delete      cleans up resources from cluster
  sm:deploy      creates the pre-defined topology for tests
  sm:policies    create a VirtualService and define application Authorization files
```
