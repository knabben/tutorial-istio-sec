package tls

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"path"
)

type TLS mg.Namespace

const (
	certificateFolder = "./tls/certs"
	certificateFile   = "root_keypair.pem"

	// valid algorithm options RSA, RSA-PSS, EC, X25519, X448, ED25519 and ED448.
	certificateAlgo = "x448"
)

// GenKeyPair creates new public/private key pair under certs folder
func (TLS) GenKeyPair() error {
	certPath := path.Join(certificateFolder, certificateFile)

	// generate a new key pair
	args := []string{"genpkey", "-algorithm", certificateAlgo, "-out", certPath}
	if err := sh.Run("openssl", args...); err != nil {
		return err
	}

	// print out the keypair values
	args = []string{"pkey", "-in", certPath, "-noout", "-text"}
	if err := sh.RunV("openssl", args...); err != nil {
		return err
	}
	return nil
}

func (TLS)