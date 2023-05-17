package config

import (
	"encoding/json"
	"time"

	"github.com/jeremywohl/flatten"
	"github.com/knadh/koanf"
	"github.com/zacscoding/go-rest-template/pkg/cache"
	"github.com/zacscoding/go-rest-template/pkg/cfgloader"
	"github.com/zacscoding/go-rest-template/pkg/database"
	"github.com/zacscoding/go-rest-template/pkg/utils/maskingutil"
)

const EnvPrefix = "APP_SERVER_"

type Config struct {
	K       *koanf.Koanf    `json:"-"`
	Stage   string          `json:"stage" yaml:"stage"`
	Logging LoggingConfig   `json:"logging" yaml:"logging"`
	Server  ServerConfig    `json:"server" yaml:"server"`
	DB      database.Config `json:"db" yaml:"db"`
	Cache   cache.Config    `json:"cache" yaml:"cache"`
	Metric  MetricsConfig   `json:"metric" yaml:"metric"`
}

type LoggingConfig struct {
	Level             int    `json:"level" yaml:"level"`
	Encoding          string `json:"encoding" yaml:"encoding"`
	Development       bool   `json:"development" yaml:"development"`
	DisableStacktrace bool   `json:"disable-stacktrace" yaml:"disable-stacktrace"`
}

type ServerConfig struct {
	Port             int           `json:"port" yaml:"port"`
	ReadTimeout      time.Duration `json:"read-timeout" yaml:"read-timeout"`
	WriteTimeout     time.Duration `json:"write-timeout" yaml:"write-timeout"`
	GracefulShutdown time.Duration `json:"graceful-shutdown" yaml:"graceful-shutdown"`
	Cors             struct {
		AllowAll   bool     `json:"allow-all" yaml:"allow-all"`
		Origin     []string `json:"origin" yaml:"origin"`
		BrowserExt bool     `json:"browser-ext" yaml:"browser-ext"`
	} `json:"cors" yaml:"cors"`
	Docs struct {
		Enabled bool   `json:"enabled" yaml:"enabled"`
		Path    string `json:"path" yaml:"path"`
	} `json:"docs" yaml:"docs"`
	Auth struct {
		JWT struct {
			Realm      string        `json:"realm" yaml:"realm"`
			Key        string        `json:"key" yaml:"key"`
			Timeout    time.Duration `json:"timeout" yaml:"timeout"`
			MaxRefresh time.Duration `json:"max-refresh" yaml:"max-refresh"`
		} `json:"jwt" yaml:"jwt"`
	} `json:"auth" yaml:"auth"`
}

type MetricsConfig struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	Port      int    `json:"port" yaml:"port"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Subsystem string `json:"subsystem" yaml:"subsystem"`
}

// Load loads config with below orders.
// 1. defaultConfig constants.
// 2. environment variables having "APP_SERVER_" prefix.
// 3. config file if provided
// 4. configMap if not empty.
func Load(configPath string, configMap map[string]interface{}) (*Config, error) {
	var opts []cfgloader.Option
	opts = append(opts, cfgloader.WithEnv(EnvPrefix))
	if configPath != "" {
		opts = append(opts, cfgloader.WithConfigFile(configPath))
	}
	if len(configMap) != 0 {
		opts = append(opts, cfgloader.WithConfigMap(configMap))
	}

	var conf Config
	k, err := cfgloader.LoadWithOptions(&conf, defaultConfig, opts...)
	if err != nil {
		return nil, err
	}
	conf.K = k
	return &conf, nil
}

func (c *Config) MarshalJSON() ([]byte, error) {
	type conf Config
	alias := conf(*c)

	data, err := json.Marshal(&alias)
	if err != nil {
		return nil, err
	}

	flat, err := flatten.FlattenString(string(data), "", flatten.DotStyle)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal([]byte(flat), &m)
	if err != nil {
		return nil, err
	}

	// add keys if u want to mask some properties.
	maskKeys := map[string]struct{}{
		"server.auth.jwt.key": {},
	}

	for key, val := range m {
		if v, ok := val.(string); ok {
			m[key] = maskingutil.MaskPassword(v)
		}
		if _, ok := maskKeys[key]; ok {
			m[key] = "****"
		}
	}
	return json.Marshal(&m)
}
