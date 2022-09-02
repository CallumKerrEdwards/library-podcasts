package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CallumKerrEdwards/loggerrific/tlogger"
	"github.com/stretchr/testify/assert"
)

func TestConfig_FromEnvironment(t *testing.T) {

	// given
	envVarCleanup := envVarSetter(t, map[string]string{
		"LIBPODCASTS_APPLICATION_HOST": "http://localhost",
	})
	t.Cleanup(envVarCleanup)

	// when
	config, err := New(tlogger.NewTLogger(t))

	// then
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost", config.Application.Host)
}

func TestConfig_FromFile(t *testing.T) {
	// given
	configFilePath := filepath.Join(t.TempDir(), "config.yaml")
	envVarCleanup := envVarSetter(t, map[string]string{
		"LIBPODCASTS_CONFIG_PATH": configFilePath,
	})
	t.Cleanup(envVarCleanup)

	configYAML := `---
Application:
  Host: "http://localhost"
Podcast:
  Title: "Test Audiobooks"
  Description: "A Test Audiobook feed"
  Explicit: False 
`
	err := os.WriteFile(configFilePath, []byte(configYAML), 0644)
	assert.NoError(t, err)

	// when
	config, err := New(tlogger.NewTLogger(t))

	// then
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost", config.Application.Host)
	assert.Equal(t, "Test Audiobooks", config.Podcast.Title)
	assert.Equal(t, "A Test Audiobook feed", config.Podcast.Description)
	assert.Equal(t, false, config.Podcast.Explicit)
}

func TestConfig_EnvironementOverridesFile(t *testing.T) {
	// given
	configFilePath := filepath.Join(t.TempDir(), "config.yaml")
	envVarCleanup := envVarSetter(t, map[string]string{
		"LIBPODCASTS_CONFIG_PATH":      configFilePath,
		"LIBPODCASTS_APPLICATION_HOST": "http://127.0.0.1",
	})
	t.Cleanup(envVarCleanup)

	configYAML := `---
Application:
  Host: "http://localhost"
`
	err := os.WriteFile(configFilePath, []byte(configYAML), 0644)
	assert.NoError(t, err)

	// when
	config, err := New(tlogger.NewTLogger(t))

	// then
	assert.NoError(t, err)
	assert.Equal(t, "http://127.0.0.1", config.Application.Host)
}

func TestConfig_DefaultsOnly(t *testing.T) {
	// when
	config, err := New(tlogger.NewTLogger(t))

	// then
	assert.NoError(t, err)
	assert.Equal(t, "Audiobooks", config.Podcast.Title)
	assert.Equal(t, true, config.Podcast.Explicit)
	assert.Equal(t, true, config.Podcast.BlockFromITunes)
}

func envVarSetter(t *testing.T, envs map[string]string) (closer func()) {
	originalEnvVars := map[string]string{}

	for name, value := range envs {
		if originalValue, ok := os.LookupEnv(name); ok {
			originalEnvVars[name] = originalValue
		}
		err := os.Setenv(name, value)
		assert.NoError(t, err)
	}

	return func() {
		for name := range envs {
			origValue, has := originalEnvVars[name]
			if has {
				t.Log("Setting env", name, "to", origValue)
				err := os.Setenv(name, origValue)
				assert.NoError(t, err)

			} else {
				t.Log("Unsetting env", name)
				err := os.Unsetenv(name)
				assert.NoError(t, err)
			}
		}
	}
}
