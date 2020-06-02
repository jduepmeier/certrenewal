package certrenewal

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// CertData contains the pem encoded certificate.
type CertData struct {
	Certificate string   `mapstructure:"certificate"`
	Chain       []string `mapstructure:"ca_chain"`
	PrivateKey  string   `mapstructure:"private_key"`
}

// Cert contains all infos about a certificate.
type Cert struct {
	PrivateKey string   `yaml:"private_key"`
	CertFile   string   `yaml:"cert_file"`
	ChainFile  string   `yaml:"chain_file"`
	Role       string   `yaml:"role"`
	CN         string   `yaml:"cn"`
	SANS       []string `yaml:"sans"`
	Hooks      []string `yaml:"hooks"`
	data       CertData
}

// Issue renews the certificate.
func (cert *Cert) Issue(config *Config, client *api.Client) error {
	err := cert.issue(config, client)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrIssue, err)
	}

	return err
}

func (cert *Cert) issue(config *Config, client *api.Client) error {
	data := map[string]interface{}{
		"common_name": cert.CN,
		"alt_names":   strings.Join(cert.SANS, ","),
	}
	certResponse, err := client.Logical().Write(fmt.Sprintf("%s/issue/%s", config.PkiPath, cert.Role), data)
	if err != nil {
		return err
	}

	err = mapstructure.Decode(&certResponse.Data, &cert.data)
	if err != nil {
		return err
	}

	return nil
}

// WriteFiles writes the certificate files to disks.
func (cert *Cert) WriteFiles() (err error) {
	err = writeFile(cert.CertFile, cert.data.Certificate)
	if err != nil {
		return err
	}
	err = writeFile(cert.PrivateKey, cert.data.PrivateKey)
	if err != nil {
		return err
	}
	err = writeFile(cert.ChainFile, strings.Join(cert.data.Chain, "\n"))
	if err != nil {
		return err
	}
	return nil
}

// RunHooks runs the configured hooks.
func (cert *Cert) RunHooks() (err error) {
	for _, hook := range cert.Hooks {
		logrus.Infof("run command %s", hook)
		cmd := exec.Command("/bin/sh", "-c", hook)
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// NeedsRenewal checks if the certificate needs renewal.
// Will return ErrCert error if the certificate cannot be read correctly.
// If an error will be returned the boolean value is always true.
func (cert *Cert) NeedsRenewal(config *Config) (bool, error) {
	renewal, err := cert.needsRenewal(config)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrCert, err)
	}

	return renewal, err
}

func (cert *Cert) needsRenewal(config *Config) (bool, error) {
	pemContent, err := ioutil.ReadFile(cert.CertFile)
	if err != nil {
		return true, nil
	}

	block, _ := pem.Decode(pemContent)

	pubFile, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return true, err
	}
	if pubFile.Subject.CommonName != cert.CN {
		return true, nil
	}
	expectedSans := sort.StringSlice(cert.SANS)
	realSans := sort.StringSlice(pubFile.DNSNames)
	if len(expectedSans) != len(realSans) {
		return true, nil
	}
	for i, san := range expectedSans {
		if san != realSans[i] {
			return true, nil
		}
	}

	now := time.Now()
	if now.Add(10 * 24 * time.Hour).After(pubFile.NotAfter) {
		return true, nil
	}

	return false, nil
}

// CheckAndRenew checks if the cert needs renewal and renews the certs if needed.
// If renewal is needed the configured hooks will run after the renewal.
func (cert *Cert) CheckAndRenew(config *Config, client *api.Client) (bool, error) {
	renewal, err := cert.NeedsRenewal(config)
	if err != nil {
		return false, err
	}
	if !renewal {
		logrus.Infof("cert %s needs no renewal", cert.CertFile)
		return false, nil
	}
	err = cert.Issue(config, client)
	if err != nil {
		return true, err
	}

	err = cert.WriteFiles()
	if err != nil {
		return true, err
	}

	return true, cert.RunHooks()
}
