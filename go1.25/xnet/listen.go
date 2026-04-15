package xnet

import "fmt"

// TODO: move to smth like xnetcfg, config structs should not interfere with functionality

type ListenLocation struct {
	Host *string `yaml:"host" validate:"omitnil"`
	Port int     `yaml:"port" validate:"lte=65535"`
}

type ListenTLS struct {
	Option       `yaml:",inline"`
	Certificates X509CertificateKeyPair `yaml:"certificates"`
}

type Option struct {
	Enabled bool `yaml:"enabled"`
}

func (l *ListenLocation) Address() string {
	if l.Host == nil {
		return fmt.Sprintf(":%d", l.Port)
	}
	return fmt.Sprintf("%s:%d", *l.Host, l.Port)
}
