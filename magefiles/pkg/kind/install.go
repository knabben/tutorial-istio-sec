package kind

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/writter"
	"github.com/magefile/mage/sh"
	"path"
)

// InstallKind installs kind as a base K8s cluster for local development
func InstallKind(name, specsFolder string, withLB bool) error {
	// Delete old cluster if exists
	if err := writter.Kind("delete", "cluster", "--name", name); err != nil {
		return err
	}
	// Create a new cluster using predefined configuration file
	configFile := writter.AppendFolder(specsFolder, "kind.yaml")
	if err := writter.Kind("create", "cluster", "--name", name, "--config", configFile); err != nil {
		return err
	}
	nodes := []string{name + "-worker2", name + "-control-plane", name + "-worker"}
	// Set sysctl parameters on each node
	for _, node := range nodes {
		args := []string{"exec", node, "sysctl", "-w"}
		config := []string{"fs.inotify.max_user_instances=1024", "fs.inotify.max_user_watches=1048576"}
		for _, s := range config {
			if err := sh.RunV("docker", append(args, s)...); err != nil {
				return nil
			}
		}
	}
	if withLB {
		if err := installMetalLB(specsFolder); err != nil {
			return err
		}
	}
	return nil
}

// installMetalLB Installs MetalLB in the cluster
func installMetalLB(specsFolder string) error {
	for _, spec := range []string{METALLB_URL, path.Join(specsFolder, METALLB_CR)} {
		_ = writter.Kubectl("-n", "metallb-system", "wait", "--for=condition=Ready", "pod", "-l", "app=metallb", "--timeout", "300s")
		// Create a new cluster using predefined configuration file
		if err := writter.Kubectl("apply", "-f", spec); err != nil {
			return err
		}
	}
	return nil
}
