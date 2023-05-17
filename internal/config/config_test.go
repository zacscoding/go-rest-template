package config

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Default(t *testing.T) {
	conf, err := Load("", nil)
	assert.NoError(t, err)

	cases := []struct {
		key      string
		expected interface{}
		values   []interface{}
	}{
		{key: "stage", expected: "local", values: []interface{}{conf.Stage}},
		{key: "logging.level", expected: 1, values: []interface{}{conf.Logging.Level}},
		{key: "logging.encoding", expected: "console", values: []interface{}{conf.Logging.Encoding}},
		{key: "logging.development", expected: false, values: []interface{}{conf.Logging.Development}},
		{key: "logging.disable-stacktrace", expected: true, values: []interface{}{conf.Logging.DisableStacktrace}},

		{key: "server.port", expected: 8080, values: []interface{}{conf.Server.Port}},
		{key: "server.read-timeout", expected: 5 * time.Second, values: []interface{}{conf.Server.ReadTimeout}},
		{key: "server.write-timeout", expected: 10 * time.Second, values: []interface{}{conf.Server.WriteTimeout}},
		{key: "server.graceful-shutdown", expected: 30 * time.Second, values: []interface{}{conf.Server.GracefulShutdown}},
		{key: "server.cors.allow-all", expected: true, values: []interface{}{conf.Server.Cors.AllowAll}},
		{key: "server.cors.browser-ext", expected: true, values: []interface{}{conf.Server.Cors.BrowserExt}},
		{key: "server.docs.enabled", expected: false, values: []interface{}{conf.Server.Docs.Enabled}},
		{key: "server.auth.jwt.realm", expected: "sample app", values: []interface{}{conf.Server.Auth.JWT.Realm}},
		{key: "server.auth.jwt.key", expected: "c2FtcGxlIGFwcAo=", values: []interface{}{conf.Server.Auth.JWT.Key}},
		{key: "server.auth.jwt.timeout", expected: time.Hour, values: []interface{}{conf.Server.Auth.JWT.Timeout}},
		{key: "server.auth.jwt.max-refresh", expected: 5 * time.Hour, values: []interface{}{conf.Server.Auth.JWT.MaxRefresh}},

		{key: "db.driver", expected: "mysql", values: []interface{}{conf.DB.Driver}},
		{key: "db.data-source-name", expected: "root:dbpassword@tcp(127.0.0.1:3306)/mydb?charset=utf8&parseTime=True&multiStatements=true", values: []interface{}{conf.DB.DataSourceName}},
		{key: "db.logging-level", expected: 1, values: []interface{}{conf.DB.LoggingLevel}},
		{key: "db.batch-size", expected: 500, values: []interface{}{conf.DB.BatchSize}},
		{key: "db.migrate.enabled", expected: false, values: []interface{}{conf.DB.Migrate.Enabled}},
		{key: "db.migrate.dir", expected: "", values: []interface{}{conf.DB.Migrate.Dir}},
		{key: "db.pool.max-open", expected: 10, values: []interface{}{conf.DB.Pool.MaxOpen}},
		{key: "db.pool.max-idle", expected: 10, values: []interface{}{conf.DB.Pool.MaxIdle}},
		{key: "db.pool.max-lifetime", expected: 30 * time.Minute, values: []interface{}{conf.DB.Pool.MaxLifeTime}},

		{key: "cache.enabled", expected: false, values: []interface{}{conf.Cache.Enabled}},
		{key: "cache.prefix", expected: "myapp-", values: []interface{}{conf.Cache.Prefix}},
		{key: "cache.type", expected: "redis", values: []interface{}{conf.Cache.Type}},
		{key: "cache.ttl", expected: 1 * time.Minute, values: []interface{}{conf.Cache.TTL}},
		{key: "cache.redis.read-timeout", expected: 3 * time.Second, values: []interface{}{conf.Cache.Redis.ReadTimeout}},
		{key: "cache.redis.write-timeout", expected: 3 * time.Second, values: []interface{}{conf.Cache.Redis.WriteTimeout}},
		{key: "cache.redis.dial-timeout", expected: 5 * time.Second, values: []interface{}{conf.Cache.Redis.DialTimeout}},
		{key: "cache.redis.pool-size", expected: 10, values: []interface{}{conf.Cache.Redis.PoolSize}},
		{key: "cache.redis.pool-timeout", expected: 4 * time.Second, values: []interface{}{conf.Cache.Redis.PoolTimeout}},
		{key: "cache.redis.max-conn-age", expected: 0, values: []interface{}{conf.Cache.Redis.MaxConnAge}},
		{key: "cache.redis.idle-timeout", expected: 60 * time.Second, values: []interface{}{conf.Cache.Redis.IdleTimeout}},

		{key: "metric.enabled", expected: true, values: []interface{}{conf.Metric.Enabled}},
		{key: "metric.port", expected: 8089, values: []interface{}{conf.Metric.Port}},
		{key: "metric.namespace", expected: "myapp", values: []interface{}{conf.Metric.Namespace}},
		{key: "metric.subsystem", expected: "server", values: []interface{}{conf.Metric.Subsystem}},
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			values := append(tc.values, defaultConfig[tc.key])
			if d, ok := tc.expected.(time.Duration); ok {
				equalDuration(t, d, values...)
				return
			}
			equal(t, tc.expected, values...)
		})
	}
	assert.Equal(t, len(defaultConfig), len(cases))
}

func TestMarshalJSON(t *testing.T) {
	conf, err := Load("", nil)
	assert.NoError(t, err)
	b, err := json.Marshal(conf)
	assert.NoError(t, err)
	var m map[string]interface{}
	assert.NoError(t, json.Unmarshal(b, &m))

	assert.Equal(t, "root:****@tcp(127.0.0.1:3306)/mydb?charset=utf8&parseTime=True&multiStatements=true", m["db.data-source-name"])
	assert.Equal(t, "****", m["server.auth.jwt.key"])
}

func equal(t *testing.T, expected interface{}, values ...interface{}) {
	for _, v := range values {
		assert.EqualValues(t, expected, v)
	}
}

func equalDuration(t *testing.T, expected time.Duration, values ...interface{}) {
	for _, v := range values {
		if str, ok := v.(string); ok {
			d, err := time.ParseDuration(str)
			assert.NoError(t, err)
			assert.EqualValues(t, expected, d)
			continue
		}
		assert.EqualValues(t, expected, v)
	}
}
