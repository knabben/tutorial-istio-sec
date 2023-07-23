package spire

import (
	"fmt"
	"github.com/knabben/tutorial-istio-sec/magefiles/writter"
	"github.com/magefile/mage/sh"
	"strings"
	"time"
)

func Bootstrap(specPath string) error {
	bootstrap := writter.AppendFolder(specPath, "bootstrap")
	if err := writter.Kubectl("apply", "-f", writter.AppendFolder(bootstrap, "spire-quickstart.yaml")); err != nil {
		return err
	}

	args := []string{"wait", "--for=condition=Ready", "-n", "spire", "pod", "-l", "app=spire-server", "--timeout", "600s"}
	if err := writter.Kubectl(args...); err != nil {
		return err
	}

	if err := writter.Kubectl("apply", "-f", writter.AppendFolder(bootstrap, "clusterspiffee.yaml")); err != nil {
		return err
	}

	writter.Kubectl("label", "namespace", "default", "istio-injection=enabled")

	if err := writter.Istioctl("install", "--skip-confirmation", "-f", writter.AppendFolder(specPath, "istio.yaml")); err != nil {
		return err
	}

	return nil
}

// InstallSpire installs and set the SPIRE server manually
func InstallSpire(specPath, specApps string) error {
	if err := writter.Kubectl("apply", "-f", specPath); err != nil {
		return err
	}

	for {
		// wait for server run, if pod does not exist keep trying.
		args := []string{"wait", "--for=condition=Ready", "-n", "spire", "pod", "-l", "app=spire-server", "--timeout", "300s"}
		if err := writter.Kubectl(args...); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	// create the SPIRE node attestation
	args := []string{
		"exec", "-n", "spire",
		"spire-server-0", "--",
		"/opt/spire/bin/spire-server", "entry", "create",
		"-spiffeID", "spiffe://opssec.in/ns/spire/sa/spire-agent",
		"-selector", "k8s_sat:cluster:kind",
		"-selector", "k8s_sat:agent_ns:spire",
		"-selector", "k8s_sat:agent_sa:spire-agent",
		"-node",
	}
	if err := writter.Kubectl(args...); err != nil {
		return err
	}
	// create client and server wl attestation
	for _, arg := range []string{"client", "server"} {
		if err := writter.Kubectl("create", "serviceaccount", arg); err != nil {
			return err
		}
		args = []string{
			"exec", "-n", "spire",
			"spire-server-0", "--",
			"/opt/spire/bin/spire-server", "entry", "create",
			"-spiffeID", fmt.Sprintf("spiffe://opssec.in/ns/default/sa/%s", arg),
			"-parentID", "spiffe://opssec.in/ns/spire/sa/spire-agent",
			"-selector", "k8s:ns:default",
			"-selector", "k8s:sa:" + arg,
		}
		if err := writter.Kubectl(args...); err != nil {
			return err
		}
	}
	// Install client/server deployment
	return writter.Kubectl("apply", "-f", specApps)
}

func Deploy(specPath string) error {
	if err := writter.Kubectl("apply", "-f", writter.AppendFolder(specPath, "deploy.yaml")); err != nil {
		return err
	}
	args := []string{"wait", "--for=condition=Ready", "pod", "-l", "app=sleep", "--timeout", "300s"}
	if err := writter.Kubectl(args...); err == nil {
		return err
	}
	return nil
}

func Check() error {
	pod, _ := sh.Output("kubectl", "get", "pod", "-l", "app=spire-server", "-n", "spire", "-o", "jsonpath=\"{.items[0].metadata.name}\"")
	fmt.Println(pod)

	out, _ := sh.Output("kubectl", "exec", "-t", strings.Trim(pod, "\""), "-n", "spire", "-c", "spire-server", "--", "./bin/spire-server", "entry", "show")
	fmt.Println(out)
	return nil
}
