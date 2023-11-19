package certrenewal

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkFileContent(t *testing.T, filename, content string) {

	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("cannot read file: %s", err)
		return
	}
	if string(fileContent) != content {
		t.Errorf("Expected %q instead of %q", content, string(fileContent))
	}
}

func TestWriteFiles(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "certrenewal-test")
	if err != nil {
		t.Skipf("cannot create tmpdir: %s", err)
	}
	defer os.RemoveAll(tmpdir)
	certData := CertData{
		Certificate: "test-cert",
		PrivateKey:  "test-cert",
		Chain:       []string{"test-chain"},
	}
	cert := Cert{
		PrivateKey: path.Join(tmpdir, "test.key"),
		CertFile:   path.Join(tmpdir, "cert.pem"),
		ChainFile:  path.Join(tmpdir, "chain.pem"),
		data:       certData,
	}

	err = cert.WriteFiles()

	if err != nil {
		t.Errorf("WriteFiles() should not return an error. Got %s", err)
	}
	checkFileContent(t, cert.PrivateKey, certData.PrivateKey+"\n")
	checkFileContent(t, cert.CertFile, certData.Certificate+"\n")
	checkFileContent(t, cert.ChainFile, certData.Chain[0]+"\n")
}

func TestHooks(t *testing.T) {
	tests := []struct {
		name string
		err  error
		cert *Cert
	}{
		{
			name: "no-hooks",
			cert: &Cert{},
			err:  nil,
		},
		{
			name: "false",
			cert: &Cert{
				Hooks: []string{
					"/bin/false",
				},
			},
			err: ErrHookFailed,
		},
		{
			name: "no-error",
			cert: &Cert{
				Hooks: []string{
					"/bin/true",
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.cert.RunHooks()
			assert.ErrorIs(t, err, test.err)
		})
	}
}

func TestPkiPath(t *testing.T) {
	pathFromConfig := "path-from-config"
	pathFromCert := "path-from-cert"
	config := &Config{
		PkiPath: pathFromConfig,
	}

	tests := []struct {
		name     string
		expected string
		cert     *Cert
	}{
		{
			name:     "path-from-config",
			expected: pathFromConfig,
			cert:     &Cert{},
		},
		{
			name:     "path-from-cert",
			expected: pathFromCert,
			cert: &Cert{
				PkiPath: pathFromCert,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualPath := test.cert.pkiPath(config)
			assert.Equal(t, test.expected, actualPath)
		})
	}
}
