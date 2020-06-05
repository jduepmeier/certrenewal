package certrenewal

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// SSHData contains the encoded ssh data string.
type SSHData struct {
	PrivateKey string `mapstructure:","`
	PublicKey  string `mapstructure:","`
	Cert       string `mapstructure:"signed_key"`
}

// SSHCert holds the ssh data and metadata from ssh cert.
type SSHCert struct {
	data           SSHData
	PrivateKeyPath string   `yaml:"private_key"`
	PublicKeyPath  string   `yaml:"public_key"`
	CertPath       string   `yaml:"cert"`
	Role           string   `yaml:"role"`
	Hosts          []string `yaml:"hosts"`
	Hooks          []string `yaml:"hooks"`
}

// Issue renews the certificate.
func (cert *SSHCert) Issue(config *Config, client *api.Client) error {
	err := cert.issue(config, client)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrIssue, err)
	}

	return err
}

func (cert *SSHCert) issue(config *Config, client *api.Client) error {
	err := cert.load()
	if err != nil {
		return err
	}
	sshClient := client.SSHWithMountPoint(config.SSHPath)

	data := map[string]interface{}{
		"cert_type":        "host",
		"valid_principals": strings.Join(cert.Hosts, ","),
		"public_key":       cert.data.PublicKey,
	}

	signedData, err := sshClient.SignKey(cert.Role, data)
	if err != nil {
		return err
	}

	return mapstructure.Decode(&signedData.Data, &cert.data)
}

// WriteFiles writes the certificate files to disks.
func (cert *SSHCert) WriteFiles() (err error) {
	err = writeFile(cert.CertPath, cert.data.Cert)
	return err
}

func (cert *SSHCert) load() error {
	if cert.data.PrivateKey == "" {
		content, err := ioutil.ReadFile(cert.PrivateKeyPath)
		if err != nil {
			return err
		}

		cert.data.PrivateKey = string(content)
	}

	if cert.data.PublicKey == "" {
		content, err := ioutil.ReadFile(cert.PublicKeyPath)
		if err != nil {
			return err
		}
		cert.data.PublicKey = string(content)
	}

	return nil
}

// NeedsRenewal checks if the certificate needs renewal.
// Will return ErrCert error if the certificate cannot be read correctly.
// If an error will be returned the boolean value is always true.
func (cert *SSHCert) NeedsRenewal(config *Config) (bool, error) {
	renewal, err := cert.needsRenewal(config)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrCert, err)
	}

	return renewal, err
}

func (cert *SSHCert) needsRenewal(config *Config) (bool, error) {
	err := cert.load()
	if err != nil {
		return false, err
	}
	content, err := ioutil.ReadFile(cert.CertPath)
	if err != nil {
		return true, nil
	}

	certPub, _, _, _, err := ssh.ParseAuthorizedKey(content)
	if err != nil {
		return true, fmt.Errorf("parse failure for %s: %s", cert.CertPath, err)
	}

	sshCert, ok := certPub.(*ssh.Certificate)
	if !ok {
		return true, fmt.Errorf("%s is not a ssh certificate", cert.CertPath)
	}

	expectedPrincipals := sort.StringSlice(cert.Hosts)
	realPrincipals := sort.StringSlice(sshCert.ValidPrincipals)
	if len(expectedPrincipals) != len(realPrincipals) {
		return true, fmt.Errorf("missing principal %d != %d", len(expectedPrincipals), len(realPrincipals))
	}
	for i, principal := range expectedPrincipals {
		if principal != realPrincipals[i] {
			return true, fmt.Errorf("missing principal %s", principal)
		}
	}
	certTime := time.Unix(int64(sshCert.ValidBefore), 0)
	expectedAfter := time.Now().Add(10 * 24 * time.Hour)
	test := expectedAfter.After(certTime)
	if test {
		return test, fmt.Errorf("certTime (%s) < expectedAfter (%s)", certTime.String(), expectedAfter.String())
	}

	return test, nil
}

// CheckAndRenew checks if the cert needs renewal and renews the certs if needed.
// If renewal is needed the configured hooks will run after the renewal.
func (cert *SSHCert) CheckAndRenew(config *Config, client *api.Client) (bool, error) {
	renewal, err := cert.NeedsRenewal(config)
	if !renewal {
		logrus.Infof("ssh cert %s needs no renewal", cert.CertPath)
		return false, nil
	}
	logrus.Infof("Need renewal with the following error: %s", err)
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

// RunHooks runs the configured hooks.
func (cert *SSHCert) RunHooks() (err error) {
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
