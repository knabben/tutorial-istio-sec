package steps

import (
	"log"
	"path"
)

// DeployApplication install the application and objects
func (s *ServiceMesh) DeployApplication(namespace string) error {
	// deploy the application and Kubernetes gateway object
	if err := kubectl("apply", "-n", namespace, "-f", path.Join(SpecsFolder, "apps")); err != nil {
		return err
	}
	_ = kubectl("wait", "--for=condition=Ready", "pod", "-l", "app=appb")
	// create a waypoint in the namespace
	if err := istioctl("x", "-n", namespace, "waypoint", "apply", "--service-account", "appb"); err != nil {
		return err
	}
	_ = kubectl("wait", "--for=condition=Ready", "pod", "-l", "istio.io/gateway-name=gateway")
	return printGwListener("deploy/gateway-istio")
}

// ApplyPolicies creates a new AuthorizationPolicy for the appb service and VirtualService for control
func (s *ServiceMesh) ApplyPolicies(namespace string) error {
	if err := kubectl("apply", "-n", namespace, "-f", path.Join(SpecsFolder, "policies")); err != nil {
		return err
	}
	return nil
}

func printGwListener(gw string) error {
	log.Printf("\n\nProxy-config listener configuration for %v\n\n", gw)
	if err := istioctl("proxy-config", "listener", gw); err != nil {
		return err
	}
	log.Printf("\n\nProxy-config route configuration for %v\n\n", gw)
	if err := istioctl("proxy-config", "route", gw); err != nil {
		return err
	}
	return nil
}
