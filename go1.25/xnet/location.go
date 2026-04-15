package xnet

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
	"github.com/Deimvis/go-ext/go1.25/xptr"
)

// NOTE: DEPRECATED! Use xurlcfg instead

// TODO: move to smth like xnetcfg, config structs should not interfere with functionality

type LocationTLS struct {
	Location `yaml:",inline"`
	TLS      TLS `yaml:",inline"`
}

type LocationPF struct {
	Location   `yaml:",inline"`
	PathFormat `yaml:",inline"`
}

// DEPRECATED. Use xurlcfg instead
type Location struct {
	Scheme string `yaml:"scheme" json:"scheme" validate:"required"`
	Host   string `yaml:"host" json:"host" validate:"required"`
	// TODO: use uint16 for port
	Port *int `yaml:"port" json:"port" validate:"omitnil,lte=65535"`
}

type TLS struct {
	SkipVerify         *bool                    `yaml:"skip_verify" json:"skip_verify"`
	RootCACertificates []X509Certificate        `yaml:"root_ca_certificates" json:"root_ca_certificates"`
	AuthCertificates   []X509CertificateKeyPair `yaml:"auth_certificates" json:"auth_certificates"`
}

type X509Certificate struct {
	PEMFilePath string `yaml:"pem_file_path" validate:"file"`
}

type X509PrivateKey struct {
	PEMFilePath string `yaml:"pem_file_path" validate:"file"`
}

type X509CertificateKeyPair struct {
	Certificate X509Certificate `yaml:"certificate"`
	PrivateKey  X509PrivateKey  `yaml:"private_key"`
}

type PathFormat struct {
	PathFormat string `yaml:"path_format"`
}

type TimeoutS struct {
	TimeoutS int `yaml:"timeout_s"`
}

func (l *LocationPF) URL(pathFormatArgs ...any) *url.URL {
	u := l.Location.URL()
	u.Path = fmt.Sprintf(l.PathFormat.PathFormat, pathFormatArgs...)
	return u
}

func (l *Location) URL() *url.URL {
	u := &url.URL{}
	u.Scheme = l.Scheme
	if l.Port != nil {
		u.Host = fmt.Sprintf("%s:%d", l.Host, *l.Port)
	} else {
		u.Host = l.Host
	}
	return u
}

func (t *TLS) TLSClientConfig() (*tls.Config, error) {
	tlsCfg := &tls.Config{}

	if t.SkipVerify != nil && *t.SkipVerify {
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

func (t *TimeoutS) Timeout() time.Duration {
	return time.Duration(t.TimeoutS) * time.Second
}

func (loc Location) FromURL(u *url.URL) Location {
	var port *int
	if len(u.Port()) > 0 {
		port = xptr.T(xmust.Do(strconv.Atoi(u.Port())))
	}
	return Location{
		Scheme: u.Scheme,
		Host:   u.Hostname(),
		Port:   port,
	}
}
