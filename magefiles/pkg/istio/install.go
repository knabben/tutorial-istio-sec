package istio

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/writter"
	"path"
)

const (
	ISTIO_CONFIG = "istio.yaml"
)

func InstallIstio(specsFolder, namespace string) error {
	// Apply the Gateway API custom resources.
	argsList := [][]string{
		{"kustomize", "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v0.6.2", "-o", "/tmp/kustomized"},
		{"apply", "-f", "/tmp/kustomized"},
	}
	for _, arg := range argsList {
		if err := writter.Kubectl(arg...); err != nil {
			return err
		}
	}
	// Install Istio with custom ambient
	config := writter.AppendFolder(specsFolder, ISTIO_CONFIG)
	if err := writter.Istioctl("install", "-y", "--set", "values.global.proxy.logLevel=info", "-f", config); err != nil {
		return err
	}
	// Apply otel addons
	otelFolder := path.Join(specsFolder, "otel/")
	if err := writter.Kubectl("apply", "-f", otelFolder); err != nil {
		return err
	}
	// Enable ambient mode on default namespace and wait Kiali for completion
	argsList = [][]string{
		{
			"label", "namespace", namespace, "istio.io/dataplane-mode=ambient",
		},
		{
			"-n", "istio-system", "wait", "--for=condition=Ready", "pod", "-l", "app=kiali", "--timeout", "300s",
		},
	}
	for _, a := range argsList {
		if err := writter.Kubectl(a...); err != nil {
			return err
		}
	}
	// Start Kiali dashboard
	if err := writter.Istioctl("dashboard", "kiali"); err != nil {
		return err
	}
	return nil
}
