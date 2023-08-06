package sigstore

import "github.com/knabben/tutorial-istio-sec/magefiles/writter"

func PushContainer(container string) {
	writter.Docker("tag", "nicolaka/netshoot:latest", container)
	writter.Docker("push", container)
}

func Sign(container string) error {
	PushContainer(container)
	return writter.Cosign("sign", "-a", "attribute=newattr", container)
}

func Verify(container, identity, oidcIssuer string) error {
	return writter.Cosign("verify",
		container, "--certificate-identity="+identity,
		"--certificate-oidc-issuer="+oidcIssuer)

}
