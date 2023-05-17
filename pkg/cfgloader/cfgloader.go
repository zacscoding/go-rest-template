package cfgloader

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf"
	kjson "github.com/knadh/koanf/parsers/json"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type Option func(k *koanf.Koanf) error

// LoadWithOptions loads config with given options.
func LoadWithOptions(conf interface{}, defaultConf map[string]interface{}, opts ...Option) (*koanf.Koanf, error) {
	k := koanf.New(".")
	// load default configs
	if len(defaultConf) != 0 {
		if err := k.Load(confmap.Provider(defaultConf, "."), nil); err != nil {
			return nil, err
		}
	}

	// load from options
	for _, opt := range opts {
		if err := opt(k); err != nil {
			return nil, err
		}
	}
	if err := k.UnmarshalWithConf("", conf, koanf.UnmarshalConf{Tag: "json", FlatPaths: false}); err != nil {
		return nil, err
	}
	return k, nil
}

// WithConfigMap loads config from the given configMap which "." separated keys.
func WithConfigMap(configMap map[string]interface{}) Option {
	return func(k *koanf.Koanf) error {
		return k.Load(confmap.Provider(configMap, "."), nil)
	}
}

// WithConfigFile loads config from the given file.
// Currently, supports "yaml" and "json" files.
func WithConfigFile(configFile string) Option {
	return func(k *koanf.Koanf) error {
		path, err := filepath.Abs(configFile)
		if err != nil {
			return err
		}
		var (
			parser koanf.Parser
			ext    = filepath.Ext(path)
		)
		switch ext {
		case ".yaml", ".yml":
			parser = kyaml.Parser()
		case ".json":
			parser = kjson.Parser()
		default:
			return fmt.Errorf("not supported config file extension: %s. full path: %s", ext, configFile)
		}
		return k.Load(file.Provider(path), parser)
	}
}

// WithEnv loads config from environments.
// Use Kebab case for proper parsing.
// For example, The env value "{PREFIX}SERVER_PORT=8080" will be parsed to "server.port=8080".
// "{PREFIX}SERVER_READ-TIMEOUT=5m" to "server.read-timeout=5m".
func WithEnv(prefix string) Option {
	return func(k *koanf.Koanf) error {
		return k.Load(env.ProviderWithValue(prefix, ".", func(key string, value string) (string, interface{}) {
			// trim prefix and to lowercase
			key = strings.ToLower(strings.TrimPrefix(key, prefix))
			// replace "_" to "."
			key = strings.ReplaceAll(key, "_", ".")
			// if value is array type, then split with "," separator.
			switch k.Get(key).(type) {
			case []interface{}, []string:
				return key, strings.Split(value, ",")
			}
			return key, value
		}), nil)
	}
}

// Override overrides the configuration value based on the given newValues.
func Override(k *koanf.Koanf, conf interface{}, newValues map[string]interface{}) (*koanf.Koanf, error) {
	if err := k.Load(confmap.Provider(newValues, "."), nil); err != nil {
		return nil, err
	}
	if err := k.UnmarshalWithConf("", conf, koanf.UnmarshalConf{Tag: "json", FlatPaths: false}); err != nil {
		return nil, err
	}
	return k, nil
}
