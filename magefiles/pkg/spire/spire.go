package spire

import (
	"fmt"
	"github.com/knabben/tutorial-istio-sec/magefiles/writter"
	"time"
)

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
