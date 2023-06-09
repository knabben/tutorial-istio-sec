package tls

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"path"
)

type TLS mg.Namespace

const (
	certificateFolder = "./tls/certs"

	certCAKeyPair = "root-keypair.pem"
	certCACert    = "root-cert.pem"

	certSrvKey  = "server-keypair.pem"
	certSrvCsr  = "server-csr.pem"
	certSrvCert = "server-cert.pem"

	certCliKey  = "client-keypair.pem"
	certCliCSR  = "client-csr.pem"
	certCliCert = "client-cert.pem"
)

var (
	fullRootKeyPair = path.Join(certificateFolder, certCAKeyPair)
	fullRootCert    = path.Join(certificateFolder, certCACert)

	fullServerKeyPair = path.Join(certificateFolder, certSrvKey)
	fullServerCSR     = path.Join(certificateFolder, certSrvCsr)
	fullServerCert    = path.Join(certificateFolder, certSrvCert)

	fullClientKeyPair = path.Join(certificateFolder, certCliKey)
	fullClientCSR     = path.Join(certificateFolder, certCliCSR)
	fullClientCert    = path.Join(certificateFolder, certCliCert)
)

type Certificate struct {
	days     string
	request  string
	ca       string
	caKey    string
	certFile string
}

func (x *Certificate) Render() (string, []string) {
	return "openssl", []string{
		"x509", "-req",
		"-in", x.request,
		"-days", x.days,
		"-CAcreateserial",
		"-CA", x.ca,
		"-CAkey", x.caKey,
		"-out", x.certFile,
		"-extfile", "./tls/ext.conf",
	}
}

type CertificateSignRequest struct {
	cypher  string
	subject string
	keyFile string
	csrFile string
}

func (c *CertificateSignRequest) Render() (string, []string) {
	return "openssl", []string{
		"req",
		"-nodes",
		"-newkey", c.cypher,
		"-keyout", c.keyFile,
		"-out", c.csrFile,
		"-subj", c.subject,
	}
}

type RootCACertificate struct {
	cypher  string
	days    string
	subject string
	keyFile string
	certOut string
}

func (r *RootCACertificate) Render() (string, []string) {
	if r.cypher == "" {
		r.cypher = "rsa:4096" // default to RSA
	}
	return "openssl", []string{
		"req",
		"-x509",
		"-nodes",
		"-newkey", r.cypher,
		"-days", r.days,
		"-keyout", r.keyFile,
		"-out", r.certOut,
		"-subj", r.subject,
	}
}

// GenRootCA 1 - Generates a new Root CA and CA key pair
func (TLS) GenRootCA() error {
	cert := &RootCACertificate{
		cypher:  "rsa:4096",
		days:    "365",
		subject: "/CN=Root CA",
		keyFile: fullRootKeyPair,
		certOut: fullRootCert,
	}

	cmd, args := cert.Render()
	if err := sh.Run(cmd, args...); err != nil {
		return err
	}

	if err := printCert(fullRootCert, "x509"); err != nil {
		return err
	}
	return nil
}

// GenServerCert 2 - Generates a server CSR and server cert
func (TLS) GenServerCert() error {
	req := &CertificateSignRequest{
		cypher:  "rsa:4096",
		subject: "/OU=server,CN=localhost",
		keyFile: fullServerKeyPair,
		csrFile: fullServerCSR,
	}

	cmd, args := req.Render()
	if err := sh.Run(cmd, args...); err != nil {
		return err
	}

	cert := &Certificate{
		days:     "365",
		request:  fullServerCSR,
		ca:       fullRootCert,
		caKey:    fullRootKeyPair,
		certFile: fullServerCert,
	}
	cmd, args = cert.Render()
	if err := sh.Run(cmd, args...); err != nil {
		return err
	}

	if err := printCert(fullServerCert, "x509"); err != nil {
		return err
	}

	return nil
}

// GenClientCert 3 - Generates a client CSR and client cert
func (TLS) GenClientCert() error {
	req := &CertificateSignRequest{
		cypher:  "rsa:4096",
		subject: "/OU=client,CN=localhost",
		keyFile: fullClientKeyPair,
		csrFile: fullClientCSR,
	}
	cmd, args := req.Render()
	if err := sh.Run(cmd, args...); err != nil {
		return err
	}

	cert := &Certificate{
		days:     "60",
		request:  fullClientCSR,
		ca:       fullRootCert,
		caKey:    fullRootKeyPair,
		certFile: fullClientCert,
	}
	cmd, args = cert.Render()
	if err := sh.Run(cmd, args...); err != nil {
		return err
	}

	if err := printCert(fullClientCert, "x509"); err != nil {
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
