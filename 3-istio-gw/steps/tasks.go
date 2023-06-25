package steps

import (
	"github.com/magefile/mage/sh"
	"os"
	"path"
)

var (
	SpecsFolder string

	kubectl  = sh.RunCmd("kubectl")
	istioctl = sh.RunCmd("istioctl")
	kind     = sh.RunCmd("kind")
)

// ServiceMeshI is an interface for fresh new Istio installation
type ServiceMeshI interface {
	InstallKind(name, config string, withLB bool) error
	InstallIstio(path, namespace string) error

	DeployApplication(namespace string) error
	ApplyPolicies(namespace string) error

	DeleteKind(name string) error
	DeleteIstio() error
}

type ServiceMesh struct{}

func NewServiceMesh() ServiceMeshI {
	return &ServiceMesh{}
}

func init() {
	wd, _ := os.Getwd()
	SpecsFolder = path.Join(wd, "specs")
}
