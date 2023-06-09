package grpc

import (
	"os"
)


// A custom install step if you need your bin someplace other than go/bin
func GRPC() error {
	return os.Rename("./MyApp", "/usr/bin/MyApp")
}

