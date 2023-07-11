package config

const RedisErrKeyDoesNotExist = "key does not exist in the redis"
const RedisTTLOneHour = 1 * 3600 * 1000 * 1000 * 1000

type Redis struct {
	Addr string `yaml:"addr"`
	Db   int    `yaml:"db"`
}

type Mysql struct {
	Dialect string `yaml:"dialect"`
	DSN     string `yaml:"dsn"`
}
