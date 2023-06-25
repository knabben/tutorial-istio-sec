package steps

import (
	"path"
)

// DeleteKind delete an existent kind cluster
func (s *ServiceMesh) DeleteKind(name string) error {
	// Delete kind cluster by name
	if err := kind("delete", "cluster", "--name", name); err != nil {
		return err
	}
	return nil
}

// DeleteIstio removes istio installation
func (s *ServiceMesh) DeleteIstio() error {
	// Uninstall istio and local waypoint
	if err := istioctl("x", "waypoint", "delete", "appb"); err != nil {
		return err
	}

	if err := istioctl("uninstall", "-y", "--purge"); err != nil {
		return err
	}

	// Uninstall otel addons
	if err := kubectl("delete", "-f", path.Join(SpecsFolder, "otel/"), "-f", path.Join(SpecsFolder, "apps/")); err != nil {
		return err
	}

	return nil
}
