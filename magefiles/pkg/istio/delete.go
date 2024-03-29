package istio

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/writter"
)

// DeleteIstio uninstall istio and local waypoint
func DeleteIstio(specFolder string, handleCM bool) error {
	if err := writter.Istioctl("x", "waypoint", "delete", "appb"); err != nil {
		return err
	}
	if err := writter.Istioctl("x", "waypoint", "delete", "appa"); err != nil {
		return err
	}

	if err := writter.Istioctl("uninstall", "-y", "--purge"); err != nil {
		return err
	}

	// Uninstall otel addons
	otelFolder := writter.AppendFolder("specs", "otel/")
	appsFolder := writter.AppendFolder(specFolder, "apps/")
	if handleCM {
		certFolder := writter.AppendFolder(specFolder, "cert-manager/")
		if err := writter.Kubectl("delete", "-f", certFolder); err != nil {
			return err
		}
	}
	if err := writter.Kubectl("delete", "-f", otelFolder, "-f", appsFolder); err != nil {
		return err
	}

	return nil
}
