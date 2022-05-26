package certrenewal

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateCerts(tmpDir string) error {
	cmdPath := path.Join("test", "generate_test_keys.sh")
	cmd := exec.Command(cmdPath, tmpDir)

	return cmd.Run()
}

func TestNeedRenewal(t *testing.T) {
	tmpDir := t.TempDir()
	sshKeyDir := path.Join(tmpDir, "ssh-keys")
	err := generateCerts(tmpDir)
	if err != nil {
		t.Fatalf("cannot generate test certificates: %s", err)
	}
	brokenCert := path.Join(tmpDir, "broken-cert.pub")
	cert30Days := path.Join(sshKeyDir, "rsa-cert-30days.pub")
	cert1Day := path.Join(sshKeyDir, "rsa-cert-1day.pub")
	privateKey := path.Join(sshKeyDir, "rsa")
	publicKey := path.Join(sshKeyDir, "rsa.pub")

	file, err := os.Create(brokenCert)
	if err != nil {
		t.Fatalf("cannot generate broken certificate: %s", err)
	}
	file.Close()
	tests := []struct {
		name   string
		config *Config
		cert   *SSHCert
		result bool
		err    error
	}{
		{
			name:   "empty",
			config: &Config{},
			cert:   &SSHCert{},
			err:    ErrCert,
			result: false,
		},
		{
			name:   "broken",
			config: &Config{},
			cert: &SSHCert{
				CertPath:       brokenCert,
				PrivateKeyPath: privateKey,
				PublicKeyPath:  publicKey,
			},
			err:    ErrCert,
			result: true,
		},
		{
			name:   "30days",
			config: &Config{},
			cert: &SSHCert{
				CertPath:       cert30Days,
				PrivateKeyPath: privateKey,
				PublicKeyPath:  publicKey,
			},
			err:    nil,
			result: false,
		},
		{
			name:   "1day",
			config: &Config{},
			cert: &SSHCert{
				CertPath:       cert1Day,
				PrivateKeyPath: privateKey,
				PublicKeyPath:  publicKey,
			},
			err:    ErrCert,
			result: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok, err := test.cert.NeedsRenewal(test.config)
			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.result, ok)
		})
	}
}
