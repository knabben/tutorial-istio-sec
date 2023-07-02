package kind

import "github.com/knabben/tutorial-istio-sec/magefiles/writter"

// DeleteKind delete an existent kind cluster
func DeleteKind(name string) error {
	return writter.Kind("delete", "cluster", "--name", name)
}
