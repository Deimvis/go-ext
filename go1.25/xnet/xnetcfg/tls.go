package xnetcfg

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/Deimvis/go-ext/go1.25/xoptional"
)

type TLS struct {
	SkipVerify         xoptional.T[bool]        `json:"skip_verify" yaml:"skip_verify" `
	RootCACertificates []X509Certificate        `json:"root_ca_certificates" yaml:"root_ca_certificates" `
	AuthCertificates   []X509CertificateKeyPair `json:"auth_certificates" yaml:"auth_certificates" `
}

type X509CertificateKeyPair struct {
	Certificate X509Certificate `json:"certificate" yaml:"certificate"`
	PrivateKey  X509PrivateKey  `json:"private_key" yaml:"private_key"`
}

type X509Certificate struct {
	PEMFilePath string `json:"pem_file_path" yaml:"pem_file_path" validate:"file"`
}

type X509PrivateKey struct {
	PEMFilePath string `json:"pem_file_path" yaml:"pem_file_path" validate:"file"`
}

func (t *TLS) StdConfig() (*tls.Config, error) {
	tlsCfg := &tls.Config{}

	if t.SkipVerify.HasValue() && t.SkipVerify.Value() {
		tlsCfg.InsecureSkipVerify = true
	}

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	for _, cert := range t.RootCACertificates {
		certContent, err := os.ReadFile(cert.PEMFilePath)
		if err != nil {
			return nil, err
		}
		ok := rootCAs.AppendCertsFromPEM(certContent)
		if !ok {
			return nil, fmt.Errorf("failed to add certificate: %s", cert.PEMFilePath)
		}
	}
	tlsCfg.RootCAs = rootCAs

	var authCerts []tls.Certificate
	for _, entry := range t.AuthCertificates {
		cert, err := entry.TLSCertificate()
		if err != nil {
			return nil, err
		}
		authCerts = append(authCerts, cert)
	}
	tlsCfg.Certificates = authCerts

	return tlsCfg, nil
}

func (ckp *X509CertificateKeyPair) TLSCertificate() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(ckp.Certificate.PEMFilePath, ckp.PrivateKey.PEMFilePath)
}
