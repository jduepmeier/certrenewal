package certrenewal

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConfigFunc func(t *testing.T, config *Config)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		testConfig testConfigFunc
		err        error
	}{
		{
			name: "missing",
			path: "missing.yaml",
			err:  ErrConfig,
		},
		{
			name: "broken",
			path: "broken_yaml.yaml",
			err:  ErrConfig,
		},
		{
			name: "empty",
			path: "empty.yaml",
		},
	}

	baseDir := path.Join("test", "configs")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := ReadConfig(path.Join(baseDir, test.path))
			if test.err != nil {
				assert.ErrorIs(t, err, test.err)
			} else {
				if test.testConfig != nil {
					test.testConfig(t, config)
				}
			}
		})
	}
}
