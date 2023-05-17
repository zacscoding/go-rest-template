package cfgloader

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type config struct {
	Server struct {
		Port         int           `json:"port"`
		ReadTimeout  time.Duration `json:"read-timeout"`
		WriteTimeout time.Duration `json:"write-timeout"`
	} `json:"server"`
	DB databaseConfig `json:"db"`
}

type databaseConfig struct {
	DSN string `json:"dsn"`
}

func TestLoadWithOptions(t *testing.T) {
	assert.NoError(t, os.Setenv("MY_APP_SERVER_READ-TIMEOUT", "2m"))
	assert.NoError(t, os.Setenv("MY_APP_SERVER_WRITE-TIMEOUT", "2m"))
	var conf config
	_, err := LoadWithOptions(&conf, map[string]interface{}{
		"server.port":          8080,
		"server.read-timeout":  "1m",
		"server.write-timeout": "1m",
	},
		WithEnv("MY_APP_"),
		WithConfigFile("./test-config.yml"),
		WithConfigMap(map[string]interface{}{
			"db.dsn": "root:password@tcp(127.0.0.1:3306)",
		}),
	)

	assert.NoError(t, err)
	assert.Equal(t, 8080, conf.Server.Port)
	assert.Equal(t, 2*time.Minute, conf.Server.ReadTimeout)
	assert.Equal(t, 3*time.Minute, conf.Server.WriteTimeout)
	assert.Equal(t, "root:password@tcp(127.0.0.1:3306)", conf.DB.DSN)
}

func TestLoadWithEnv(t *testing.T) {
	assert.NoError(t, os.Setenv("MY_APP_SERVER_PORT", "8080"))
	assert.NoError(t, os.Setenv("MY_APP_SERVER_READ-TIMEOUT", "5m"))
	var conf config
	_, err := LoadWithOptions(&conf, nil, WithEnv("MY_APP_"))

	assert.NoError(t, err)
	assert.Equal(t, 8080, conf.Server.Port)
	assert.Equal(t, 5*time.Minute, conf.Server.ReadTimeout)
}

func TestLoadWithConfigMap(t *testing.T) {
	var conf config
	_, err := LoadWithOptions(&conf, nil, WithConfigMap(map[string]interface{}{
		"server.port":         8080,
		"server.read-timeout": "5m",
	}))

	assert.NoError(t, err)
	assert.Equal(t, 8080, conf.Server.Port)
	assert.Equal(t, 5*time.Minute, conf.Server.ReadTimeout)
}

func TestLoadWithConfigFile(t *testing.T) {
	var conf config
	_, err := LoadWithOptions(&conf, nil, WithConfigFile("./test-config.yml"))

	assert.NoError(t, err)
	assert.Equal(t, 3*time.Minute, conf.Server.WriteTimeout)
}
