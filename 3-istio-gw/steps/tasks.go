package steps

import (
	"github.com/magefile/mage/sh"
	"os"
	"path"
)

// ServiceMeshI is an interface for fresh new Istio installation
type ServiceMeshI interface {
	InstallKind(name, config string, withLB bool) error
	InstallIstio(path string) error
	DeployApplication(enableAmbient bool)
	ApplyLayer4()
	ApplyLayer7()
	ApplyControlTraffic()
	DeleteKind(name string) error
	DeleteIstio() error
}

type ServiceMesh struct {
	Layer4Spec string
	Layer7Spec string
}

const (
	METALLB_URL = "https://raw.githubusercontent.com/metallb/metallb/v0.13.10/config/manifests/metallb-native.yaml"
	METALLB_CR  = "metallb_cr.yaml"
)

var (
	SpecsFolder string

	kubectl  = sh.RunCmd("kubectl")
	istioctl = sh.RunCmd("istioctl")
)

func init() {
	wd, _ := os.Getwd()
	SpecsFolder = path.Join(wd, "specs")
}

func (s *ServiceMesh) InstallIstio(p string) error {
	// Apply the Gateway API custom resources.
	args := [][]string{
		{
			"kustomize", "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v0.6.2", "-o", "/tmp/kustomized",
		},
		{
			"apply", "-f", "/tmp/kustomized",
		},
	}
	for _, a := range args {
		if err := kubectl(a...); err != nil {
			return err
		}
	}

	// Install Istio with custom ambient
	argss := []string{"install", "-y", "--set", "values.global.proxy.logLevel=debug", "-f", path.Join(SpecsFolder, p)}
	if err := istioctl(argss...); err != nil {
		return err
	}

	// Apply otel addons
	if err := kubectl("apply", "-f", path.Join(SpecsFolder, "otel/")); err != nil {
		return err
	}

	// Enable ambient mode on default namespace and wait Kiali for completion
	args = [][]string{
		{
			"label", "namespace", "default", "istio.io/dataplane-mode=ambient",
		},
		{
			"-n",
			"istio-system",
			"wait",
			"--for=condition=Ready",
			"pod",
			"-l",
			"app=kiali",
			"--timeout",
			"300s",
		},
	}
	for _, a := range args {
		if err := kubectl(a...); err != nil {
			return err
		}
	}

	// Start Kiali dashboard
	if err := istioctl("dashboard", "kiali"); err != nil {
		return err
	}

	return nil
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
			_ = kubectl("-n", "metallb-system", "wait", "--for=condition=Ready", "pod", "-l", "app=metallb", "--timeout", "300s")
			// Create a new cluster using predefined configuration file
			if err := kubectl("apply", "-f", spec); err != nil {
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

func (s *ServiceMesh) DeleteIstio() error {
	// Uninstall istio
	if err := istioctl("uninstall", "-y", "--purge"); err != nil {
		return err
	}

	// Uninstall otel addons
	if err := kubectl("delete", "-f", path.Join(SpecsFolder, "otel/")); err != nil {
		return err
	}
	return nil
}

func NewServiceMesh() ServiceMeshI {
	return &ServiceMesh{}
}
