package certrenewal

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
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
