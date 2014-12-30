package config

/**
 * Load config and prepare environment
 * for using it.
 */
import (
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"time"
	"encoding/json"
	"errors"
)

type cache struct {
	Protocol string
	Port     string
	Driver   string
}

type postal struct {
	Username string
	Password string
}

type config struct {
	Listen      string
	Cache       cache
	Postcode    postal
}

var (
	Pref  config
	Redis *redis.Pool
)

func Init(path string) error {
	f, e := ioutil.ReadFile(path)
	if e != nil {
		return e
	}
	if e := json.Unmarshal(f, &Pref); e != nil {
		return e
	}

	return loadRedis()
}

func RedisPool(protocol string, server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(protocol, server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func loadRedis() error {
	if Pref.Cache.Driver != "redis" {
		return errors.New("Unsupporte driver: " + Pref.Cache.Driver)
	}
	Redis = RedisPool(Pref.Cache.Protocol, Pref.Cache.Port)
	return nil
}

func Close() error {
	if e := Redis.Close(); e != nil {
		return e
	}
	return nil
}
