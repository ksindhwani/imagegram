package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		Name          string
		Expected      *Config
		ExpectedError error
	}{
		{
			Name: "Test Default Config",
			Expected: &Config{
				Addr:                 defaultAddr,
				ServerReadTimeout:    defaultServerReadTimeout,
				ServerWriteTimeout:   defaultServerWriteTimeout,
				ServerIdleTimeout:    defaultServerIdleTimeout,
				DBUserID:             defaultDBUserID,
				DBPassword:           defaultDBPassword,
				DBHostName:           defaultDBHostName,
				DBPort:               defaultDBPort,
				DBDatabaseName:       defaultDBDatabaseName,
				DBMaxIdleConnections: defaultDBMaxIdleConnections,
				DBMaxOpenConnections: defaultDBMaxOpenConnections,
				DBMaxConnLifetime:    defaultDBMaxConnLifeTime,
				HostImageDirectory:   defaultHostImageDirectory,
				LocalImageDirectory:  defaultLocalImageDirectory,
			},
		},
	}

	for _, test := range tests {
		config, err := New()
		assert.Equal(t, test.ExpectedError, err, test.Name)
		assert.Equal(t, test.Expected, config, test.Name)
	}
}
