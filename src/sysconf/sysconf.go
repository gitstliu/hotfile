package sysconf

import (
	"syscommon"
	"time"

	"github.com/pelletier/go-toml"
)

type Configure struct {

	//System config
	FileContentPath          string
	ServicePort              int
	ServiceId                int64
	DatacenterId             int64
	RedisCacheTimeOut        int64
	LocalCacheMemSize        int
	LocalCacheFlushThreshold int
	LocalCacheTimeOut        int64

	//DB config
	DBType           string
	ConnectionString string
	LogMode          bool

	//	//Redis
	//	RedisProtocol    string
	//	RedisAddress     string
	//	RedisMaxPoolSize int

	//Redis-cluster
	RedisClusterAddresses    []string
	RedisClusterConnTimeout  time.Duration
	RedisClusterReadTimeout  time.Duration
	RedisClusterWriteTimeout time.Duration
	RedisClusterKeepAlive    int
	RedisClusterAliveTime    time.Duration
}

var configure *Configure

func GetConfigure() *Configure {

	return configure
}

func LoadConfigure(fileName string) error {

	config, err := toml.LoadFile(fileName)

	if err != nil {
		return err
	}

	conf := Configure{FileContentPath: "ContentFile"}
	conf.FileContentPath = config.Get("sysconf.FileContentPath").(string)
	//	conf.DefaultPhysicalNodeCount = int(config.Get("sysconf.DefaultPhysicalNodeCount").(int64))
	conf.ServicePort = int(config.Get("sysconf.ServicePort").(int64))
	conf.ServiceId = config.Get("sysconf.ServiceId").(int64)
	conf.DatacenterId = config.Get("sysconf.DatacenterId").(int64)

	conf.RedisCacheTimeOut = config.Get("sysconf.RedisCacheTimeOut").(int64)
	conf.LocalCacheMemSize = int(config.Get("sysconf.LocalCacheMemSize").(int64)) * 1024 * 1024
	conf.LocalCacheFlushThreshold = int(config.Get("sysconf.LocalCacheFlushThreshold").(int64)) * 1024 * 1024
	conf.LocalCacheTimeOut = config.Get("sysconf.LocalCacheTimeOut").(int64)

	conf.DBType = config.Get("db.DBType").(string)
	conf.ConnectionString = config.Get("db.ConnectionString").(string)
	conf.LogMode = config.Get("db.LogMode").(bool)

	//	conf.RedisAddress = config.Get("redis.Address").(string)
	//	conf.RedisMaxPoolSize = int(config.Get("redis.MaxPoolSize").(int64))
	//	conf.RedisProtocol = config.Get("redis.Protocol").(string)
	conf.RedisClusterAddresses = syscommon.InterfacesToStrings(config.Get("redis-cluster.Addresses").([]interface{}))
	conf.RedisClusterConnTimeout = time.Duration(config.Get("redis-cluster.ConnTimeout").(int64)) * time.Millisecond
	conf.RedisClusterReadTimeout = time.Duration(config.Get("redis-cluster.ReadTimeout").(int64)) * time.Millisecond
	conf.RedisClusterWriteTimeout = time.Duration(config.Get("redis-cluster.WriteTimeout").(int64)) * time.Millisecond
	conf.RedisClusterKeepAlive = int(config.Get("redis-cluster.KeepAlive").(int64))
	conf.RedisClusterAliveTime = time.Duration(config.Get("redis-cluster.AliveTime").(int64)) * time.Millisecond

	configure = &conf

	return nil

}
