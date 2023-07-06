package istio

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/writter"
)

// DeleteIstio uninstall istio and local waypoint
func DeleteIstio(specFolder string) error {
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
	otelFolder := writter.AppendFolder(specFolder, "otel/")
	appsFolder := writter.AppendFolder(specFolder, "apps/")
	certFolder := writter.AppendFolder(specFolder, "cert-manager/")
	if err := writter.Kubectl("delete", "-f", otelFolder, "-f", appsFolder, "-f", certFolder); err != nil {
		return err
	}

	return nil
}
