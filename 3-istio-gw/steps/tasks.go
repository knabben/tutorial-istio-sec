package steps

import (
	"github.com/magefile/mage/sh"
	"os"
	"path"
)

var (
	SpecsFolder string
)

func init() {
	wd, _ := os.Getwd()
	SpecsFolder = path.Join(wd, "specs")
}

// ServiceMeshI is an interface for fresh new Istio installation
type ServiceMeshI interface {
	InstallKind(name, config string, withLB bool) error
	InstallIstio(mode, path string)
	DeployApplication(enableAmbient bool)
	ApplyLayer4()
	ApplyLayer7()
	ApplyControlTraffic()
	DeleteKind(name string) error
	DeleteIstio()
}

type ServiceMesh struct {
	Layer4Spec string
	Layer7Spec string
}

const (
	METALLB_URL = "https://raw.githubusercontent.com/metallb/metallb/v0.13.10/config/manifests/metallb-native.yaml"
	METALLB_CR  = "metallb_cr.yaml"
)

func (s *ServiceMesh) InstallIstio(mode, path string) {
	// kubectl get crd gateways.gateway.networking.k8s.io &> /dev/null || \
	//  { kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v0.6.2" | kubectl apply -f -; }
	//TODO implement me
	//istioctl install --set profile=minimal -y
	panic("implement me")
}

func (s *ServiceMesh) ApplyLayer4() {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceMesh) ApplyLayer7() {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceMesh) ApplyControlTraffic() {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceMesh) DeployApplication(enableAmbient bool) {
	//TODO implement me
	panic("implement me")
}

// InstallKind installs kind as a base K8s cluster for local development
func (s *ServiceMesh) InstallKind(name string, config string, withLB bool) error {
	// Delete old cluster if exists
	if err := sh.RunV("kind", "delete", "cluster", "--name", name); err != nil {
		return err
	}

	// Create a new cluster using predefined configuration file
	if err := sh.RunV("kind", "create", "cluster", "--name", name, "--config", path.Join(SpecsFolder, config)); err != nil {
		return err
	}

	if withLB {
		for _, spec := range []string{METALLB_URL, path.Join(SpecsFolder, METALLB_CR)} {
			_ = sh.Run("kubectl", "-n", "metallb-system", "wait", "--for=condition=Ready", "pod", "-l", "app=metallb")
			// Create a new cluster using predefined configuration file
			if err := sh.RunV("kubectl", "apply", "-f", spec); err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteKind delete an existent kind cluster
func (s *ServiceMesh) DeleteKind(name string) error {
	// Delete kind cluster by name
	if err := sh.RunV("kind", "delete", "cluster", "--name", name); err != nil {
		return err
	}
	return nil
}

func (s *ServiceMesh) DeleteIstio() {
	//TODO implement me
	panic("implement me")
}

func NewServiceMesh() ServiceMeshI {
	return &ServiceMesh{}
}
