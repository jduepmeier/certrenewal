package certrenewal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config contains the configuration.
type Config struct {
	RoleID    string    `yaml:"role_id"`
	SecretID  string    `yaml:"secret_id"`
	VaultAddr string    `yaml:"vault_addr"`
	Certs     []Cert    `yaml:"certs"`
	SSH       []SSHCert `yaml:"ssh"`
	PkiPath   string    `yaml:"pki_path"`
	SSHPath   string    `yaml:"ssh_path"`
	Insecure  bool      `yaml:"insecure"`
}

// ReadConfig reads the configuration from the given file.
func ReadConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return config, fmt.Errorf("%w: %s", ErrConfig, err)
	}
	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(config)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrConfig, err)
	}
	return config, err
}
