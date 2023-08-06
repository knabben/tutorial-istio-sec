package helm

import "github.com/knabben/tutorial-istio-sec/magefiles/writter"

func InstallPC(namespace string) error {
	writter.Helm("repo", "add", "sigstore", "https://sigstore.github.io/helm-charts")
	writter.Helm("repo", "update")

	writter.Kubectl("create", "namespace", "cosign-system")
	writter.Helm("install", "policy-controller", "-n", "cosign-system", "sigstore/policy-controller", "--devel")

	writter.Kubectl("label", "namespace", namespace, "policy.sigstore.dev/include=true")
	return nil
}
