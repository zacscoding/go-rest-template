package config

var defaultConfig = map[string]interface{}{
	"stage":                      "local",
	"logging.level":              1,
	"logging.encoding":           "console",
	"logging.development":        false,
	"logging.disable-stacktrace": true,

	"server.port":                 8080,
	"server.read-timeout":         "5s",
	"server.write-timeout":        "10s",
	"server.graceful-shutdown":    "30s",
	"server.cors.allow-all":       true,
	"server.cors.browser-ext":     true,
	"server.docs.enabled":         false,
	"server.auth.jwt.realm":       "sample app",
	"server.auth.jwt.key":         "c2FtcGxlIGFwcAo=", // echo 'sample app' | base64
	"server.auth.jwt.timeout":     "1h",
	"server.auth.jwt.max-refresh": "5h",

	"db.driver":            "mysql",
	"db.data-source-name":  "root:dbpassword@tcp(127.0.0.1:3306)/mydb?charset=utf8&parseTime=True&multiStatements=true",
	"db.logging-level":     1,
	"db.batch-size":        500,
	"db.migrate.enabled":   false,
	"db.migrate.dir":       "",
	"db.pool.max-open":     10,
	"db.pool.max-idle":     10,
	"db.pool.max-lifetime": "30m",

	"cache.enabled":             false,
	"cache.prefix":              "myapp-",
	"cache.type":                "redis",
	"cache.ttl":                 "1m",
	"cache.redis.read-timeout":  "3s",
	"cache.redis.write-timeout": "3s",
	"cache.redis.dial-timeout":  "5s",
	"cache.redis.pool-size":     10,
	"cache.redis.pool-timeout":  "4s",
	"cache.redis.max-conn-age":  0,
	"cache.redis.idle-timeout":  "60s",

	"metric.enabled":   true,
	"metric.port":      8089,
	"metric.namespace": "myapp",
	"metric.subsystem": "server",
}
