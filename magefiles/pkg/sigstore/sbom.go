package sigstore

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/knabben/tutorial-istio-sec/magefiles/writter"
)

var file = "netshoot.sbom.json"

func SBOM(container string) error {
	PushContainer(container)

	cmd := "syft " + container + " -o cyclonedx-json"

	fmt.Println(cmd)
	script.Exec(cmd).WriteFile(file)
	p := script.Exec("head -n 40 " + file)
	o, _ := p.Bytes()
	fmt.Println(string(o))

	return nil
}

func Attest(container string) error {
	return writter.Cosign("attest", "--predicate", file, "--type", "cyclonedx", container)
}
