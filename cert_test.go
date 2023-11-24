package certrenewal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
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

type IssueRequestData struct {
	CommonName string `json:"common_name"`
	AltNames   string `json:"alt_names"`
	IPSans     string `json:"ip_sans"`
}

func TestIssue(t *testing.T) {

	pkiPathFromConfig := "path-from-config"
	testRole := "test-role"
	config := &Config{
		PkiPath: pkiPathFromConfig,
	}

	tests := []struct {
		name string
		cert *Cert
		err  error
	}{
		{
			name: "empty",
			cert: &Cert{
				Role: testRole,
			},
			err: nil,
		},
		{
			name: "common-name",
			cert: &Cert{
				Role: testRole,
				CN:   "test.example.com",
			},
			err: nil,
		},
		{
			name: "sans",
			cert: &Cert{
				Role: testRole,
				SANS: []string{
					"test1.example.com",
					"test2.example.com",
				},
			},
			err: nil,
		},
		{
			name: "ip-sans",
			cert: &Cert{
				Role: testRole,
				IPS: []string{
					"127.0.0.1",
					"1.2.3.4",
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/v1/%s/issue/%s", test.cert.pkiPath(config), test.cert.Role) {
					t.Logf("called path: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
					return
				}

				if !assert.Equal(t, "PUT", r.Method) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				data := IssueRequestData{}
				decoder := json.NewDecoder(r.Body)
				decoder.Decode(&data)

				if !assert.Equal(t, test.cert.CN, data.CommonName) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if !assert.Equal(t, strings.Join(test.cert.SANS, ","), data.AltNames) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if !assert.Equal(t, strings.Join(test.cert.IPS, ","), data.IPSans) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				responseData := api.Secret{
					RequestID:     "1234",
					LeaseID:       "1234",
					LeaseDuration: 1234,
					Renewable:     false,
					Data:          map[string]interface{}{},
				}

				encoder := json.NewEncoder(w)
				encoder.Encode(&responseData)
			}))
			defer ts.Close()

			apiConfig := api.DefaultConfig()
			apiConfig.Address = ts.URL
			apiConfig.HttpClient = ts.Client()
			client, err := api.NewClient(apiConfig)
			if err != nil {
				t.Fatalf("got error from NewClient: %s", err)
			}

			err = test.cert.Issue(config, client)
			assert.ErrorIs(t, err, test.err)
		})
	}
}
