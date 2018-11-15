package configuration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/fabric8-services/admin-console/configuration"
)

func TestConfiguration(t *testing.T) {

	t.Run("default config", func(t *testing.T) {

		t.Run("error when using default values", func(t *testing.T) {
			// given
			// no specific ENV VAR set
			// when
			config := configuration.New()
			// then
			err := config.DefaultConfigurationError()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "Auth service url is empty")
		})

		t.Run("error when using empty auth service URL only", func(t *testing.T) {
			// given
			unsetenvs := setenvs(envvars{
				"ADMIN_SENTRY_DSN":             "sentry",
				"ADMIN_DEVELOPER_MODE_ENABLED": "false",
				"ADMIN_POSTGRES_PASSWORD":      "anothersecretpassword",
			})
			defer unsetenvs()
			// when
			config := configuration.New()
			// then
			err := config.DefaultConfigurationError()
			require.Error(t, err)
			assert.Equal(t, err.Error(), "Auth service url is empty")
		})

		t.Run("error when using localhost auth service URL", func(t *testing.T) {
			// given
			unsetenvs := setenvs(envvars{
				"ADMIN_AUTH_URL":               "http://localhost",
				"ADMIN_SENTRY_DSN":             "sentry",
				"ADMIN_POSTGRES_PASSWORD":      "anothersecretpassword",
				"ADMIN_DEVELOPER_MODE_ENABLED": "false",
			})
			defer unsetenvs()
			// when
			config := configuration.New()
			// then
			err := config.DefaultConfigurationError()
			require.Error(t, err)
			assert.Equal(t, err.Error(), "Auth service url is localhost")
		})

		t.Run("error when using invalid/localhost auth service URL", func(t *testing.T) {
			// given
			unsetenvs := setenvs(envvars{
				"ADMIN_AUTH_URL":               "localhost",
				"ADMIN_SENTRY_DSN":             "sentry",
				"ADMIN_POSTGRES_PASSWORD":      "anothersecretpassword",
				"ADMIN_DEVELOPER_MODE_ENABLED": "false",
			})
			defer unsetenvs()
			// when
			config := configuration.New()
			// then
			err := config.DefaultConfigurationError()
			require.Error(t, err)
			assert.Equal(t, err.Error(), "invalid Auth service url (missing scheme?)")
		})

		t.Run("error when using empty Sentry DSN only", func(t *testing.T) {
			// given
			unsetenvs := setenvs(envvars{
				"ADMIN_AUTH_URL":               "localhost",
				"ADMIN_POSTGRES_PASSWORD":      "anothersecretpassword",
				"ADMIN_DEVELOPER_MODE_ENABLED": "false",
			})
			defer unsetenvs()
			// when
			config := configuration.New()
			// then
			err := config.DefaultConfigurationError()
			require.Error(t, err)
			assert.Equal(t, err.Error(), "Sentry DSN is empty")
		})

		t.Run("no error when all envs set", func(t *testing.T) {
			// given
			unsetenvs := setenvs(envvars{
				"ADMIN_AUTH_URL":               "localhost",
				"ADMIN_SENTRY_DSN":             "sentry",
				"ADMIN_POSTGRES_PASSWORD":      "anothersecretpassword",
				"ADMIN_DEVELOPER_MODE_ENABLED": "false",
			})
			defer unsetenvs()
			// when
			config := configuration.New()
			// then
			err := config.DefaultConfigurationError()
			assert.NoError(t, err)
		})

	})

	t.Run("developer mode", func(t *testing.T) {

		t.Run("enabled", func(t *testing.T) {
			// given
			unsetenvs := setenvs(envvars{
				"ADMIN_DEVELOPER_MODE_ENABLED": "true",
			})
			defer unsetenvs()
			// when
			config := configuration.New()
			// then
			assert.True(t, config.IsDeveloperModeEnabled())
		})
	})

}

type envvars map[string]string

// setenvs sets the given env vars and provides a function to unsetenvs them all in a `defer` call
func setenvs(envs envvars) func() {
	for k, v := range envs {
		os.Setenv(k, v)
	}
	return func() {
		for k, _ := range envs {
			os.Unsetenv(k)
		}
	}
}
