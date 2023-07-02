package steps

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/writter"
	"log"
	"path"
)

// DeployApplication install the application and objects
func DeployApplication(namespace, specFolder string) error {
	appsFolder := path.Join(specFolder, "apps")
	// deploy the application and Kubernetes gateway object
	if err := writter.Kubectl("apply", "-n", namespace, "-f", appsFolder); err != nil {
		return err
	}
	_ = writter.Kubectl("wait", "--for=condition=Ready", "pod", "-l", "app=appb")
	// create a waypoint in the namespace
	if err := writter.Istioctl("x", "-n", namespace, "waypoint", "apply", "--service-account", "appb"); err != nil {
		return err
	}
	_ = writter.Kubectl("wait", "--for=condition=Ready", "pod", "-l", "istio.io/gateway-name=gateway")
	return printGwListener("deploy/gateway-istio")
}

// ApplyPolicies creates a new AuthorizationPolicy for the appb service and VirtualService for control
func ApplyPolicies(namespace, specFolder string) error {
	policyFolder := path.Join(specFolder, "policies")
	if err := writter.Kubectl("apply", "-n", namespace, "-f", policyFolder); err != nil {
		return err
	}
	return nil
}

func printGwListener(gw string) error {
	log.Printf("\n\nProxy-config listener configuration for %v\n\n", gw)
	if err := writter.Istioctl("proxy-config", "listener", gw); err != nil {
		return err
	}
	log.Printf("\n\nProxy-config route configuration for %v\n\n", gw)
	if err := writter.Istioctl("proxy-config", "route", gw); err != nil {
		return err
	}
	return nil
}
