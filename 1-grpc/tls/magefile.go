package tls

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"path"
)

type TLS mg.Namespace

const (
	certificateFolder = "./tls/certs"

	certCAKeyPair = "root_keypair.pem"
	certCACSR     = "root_csr.pem"
	certCACert    = "root_cert.pem"

	// valid algorithm options RSA, RSA-PSS, EC, X25519, X448, ED25519 and ED448.
	certificateAlgo = "ED448"
)

var (
	fullCertKeyPair = path.Join(certificateFolder, certCAKeyPair)
	fullRootCsr     = path.Join(certificateFolder, certCACSR)
	fullRootCert    = path.Join(certificateFolder, certCACert)
)

// GenCAKeyPair 1 - Create a new public/private key pair under the certs folder
func (TLS) GenCAKeyPair() error {
	// generate a new key pair
	args := []string{"genpkey", "-algorithm", certificateAlgo, "-out", fullCertKeyPair}
	if err := sh.Run("openssl", args...); err != nil {
		return err
	}

	if err := printCert(fullCertKeyPair, "pkey"); err != nil {
		return err
	}

	return nil
}

// GenCARootCSR 2 - Creates a new CSR file and print
func (TLS) GenCARootCSR() error {
	subject := "/CN=Root CA" // uses DN notation
	extension := "basicConstraints=critical,CA:TRUE"

	// generate a new CSR (Certificate Signing Request)
	args := []string{"req", "-new",
		"-subj", subject,
		"-addext", extension,
		"-key", fullCertKeyPair,
		"-out", fullRootCsr,
	}
	if err := sh.Run("openssl", args...); err != nil {
		return err
	}

	if err := printCert(fullRootCsr, "req"); err != nil {
		return err
	}

	return nil
}

// GenCARootCert 3 - Generate CA root certificate from CSR
func (TLS) GenCARootCert() error {
	// generate a new CA certificate (Certificate Signing Request)
	args := []string{"x509",
		"-req",
		"-in", fullRootCsr,
		"-copy_extensions", "copyall",
		"-key", fullCertKeyPair,
		"-days", "365",
		"-out", fullRootCert,
	}
	if err := sh.Run("openssl", args...); err != nil {
		return err
	}

	if err := printCert(fullRootCert, "x509"); err != nil {
		return err
	}

	return nil
}

// printCert prints out the keypair values
func printCert(certPath, ptype string) error {
	args := []string{ptype, "-in", certPath, "-noout", "-text"}
	if err := sh.RunV("openssl", args...); err != nil {
		return err
	}
	return nil
}
