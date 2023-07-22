package apps

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/writter"
	"log"
)

// DeployApplication install the application and objects
func DeployApplication(specsFolder, namespace string, installCM, gateway bool, serviceAccount string) error {
	if installCM {
		certFolder := writter.AppendFolder(specsFolder, "cert-manager/")
		if err := writter.Kubectl("apply", "-f", certFolder); err != nil {
			return err
		}
	}

	// deploy the application specs
	if err := writter.Kubectl("apply", "-n", namespace, "-f", writter.AppendFolder(specsFolder, "apps/")); err != nil {
		return err
	}
	_ = writter.Kubectl("wait", "--for=condition=Ready", "pod", "-l", "app="+serviceAccount+"b", "--timeout", "300s")
	// create a waypoint for service accounts

	for _, s := range []string{"a", "b"} {
		if err := writter.Istioctl("x", "-n", namespace, "waypoint", "apply", "--service-account", serviceAccount+s); err != nil {
			return err
		}
	}

	if gateway {
		_ = writter.Kubectl("wait", "--for=condition=Ready", "pod", "-l", "istio.io/gateway-name=gateway", "--timeout", "300s")
		printGwListener("deploy/gateway-istio")
	}

	return nil
}

// ApplyPolicies creates a new AuthorizationPolicy for the appb service and VirtualService for control
func ApplyPolicies(specsFolder, namespace string) error {
	policyFolder := writter.AppendFolder(specsFolder, "policies")
	return writter.Kubectl("apply", "-n", namespace, "-f", policyFolder)
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
