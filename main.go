package certrenewal

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/vault/api"
)

func writeFile(filename string, data string) error {
	dataBytes := []byte(strings.TrimRight(data, "\n\r") + "\n")
	return ioutil.WriteFile(filename, dataBytes, 0600)
}

// LoginApprole gets a token from the approle config.
func LoginApprole(config *Config, client *api.Client) error {
	tokenResp, err := client.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   config.RoleID,
		"secret_id": config.SecretID,
	})
	if err != nil {
		return fmt.Errorf("%w: %s", ErrLogin, err)
	}
	token := tokenResp.Auth.ClientToken
	client.SetToken(token)
	return nil
}

// Run runs the renewal process for the given config.
// Returns 0 if no certifiate was renewed.
// 1 if at least one certificate was renewed and 2 if an error occurred.
func Run(config *Config) (int, error) {
	vaultConfig := api.DefaultConfig()
	if config.VaultAddr != "" {
		vaultConfig.Address = config.VaultAddr
		vaultConfig.ConfigureTLS(&api.TLSConfig{
			Insecure: config.Insecure,
		})
	}
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return 2, fmt.Errorf("%w: cannot get vault client: %s", ErrConfig, err)
	}
	err = LoginApprole(config, client)
	if err != nil {
		return 2, err
	}
	returnCode := 0
	for _, cert := range config.Certs {
		renewed, err := cert.CheckAndRenew(config, client)
		if err != nil {
			return 2, err
		}

		if renewed {
			returnCode = 1
		}
	}

	for _, cert := range config.SSH {
		renewed, err := cert.CheckAndRenew(config, client)
		if err != nil {
			return 2, err
		}
		if renewed {
			returnCode = 1
		}
	}
	return returnCode, nil
}
