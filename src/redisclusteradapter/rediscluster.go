package redisclusteradapter

import (
	"fmt"
	"syscommon"
	"time"

	"github.com/gitstliu/go-redis-cluster"
	"github.com/gitstliu/log4go"
)

type RedisPipelineCommand struct {
	CommandName string
	Key         string
	Args        []interface{}
}

var adapter = RedisClusterAdapter{}

func GetAdapter() *RedisClusterAdapter {
	return &adapter
}

type RedisClusterAdapter struct {
}

var redisClusterClient *redis.Cluster

func CreadRedisCluster(hosts []string, connTimeout time.Duration, readTimeout time.Duration, writeTimeout time.Duration, keepAlive int, aliveTime time.Duration) {
	cluster, err := redis.NewCluster(
		&redis.Options{
			StartNodes:   hosts,
			ConnTimeout:  connTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			KeepAlive:    keepAlive,
			AliveTime:    aliveTime,
		})

	if err != nil {
		log4go.Error("Cluster Create Error: %v", err)
	} else {
		redisClusterClient = cluster
	}
}

func (this *RedisClusterAdapter) SET(key, value string) (string, error) {
	log4go.Debug("key is %v", key)
	client := redisClusterClient
	//defer client.Close()
	return redis.String(client.Do("SET", key, value))
}

func (this *RedisClusterAdapter) SETEX(key string, value []byte, ex int64) (string, error) {
	log4go.Debug("key is %v", key)
	client := redisClusterClient
	//defer client.Close()
	return redis.String(client.Do("SET", key, value, "ex", ex))
}

func (this *RedisClusterAdapter) EXPIRE(key string, ex int64) (string, error) {
	log4go.Debug("key is %v", key)
	client := redisClusterClient
	//defer client.Close()
	return redis.String(client.Do("EXPIRE", key, ex))
}

func (this *RedisClusterAdapter) GET(key string) (string, error) {
	log4go.Debug("key is %v", key)
	client := redisClusterClient
	//defer client.Close()
	return redis.String(client.Do("GET", key))
}

func (this *RedisClusterAdapter) GETBYTES(key string) ([]byte, error) {
	log4go.Debug("key is %v", key)
	client := redisClusterClient
	//defer client.Close()
	getBytesReuslt, getBytesErr := client.Do("GET", key)
	log4go.Debug("getBytesReuslt = %v", getBytesReuslt)
	log4go.Debug(getBytesErr)
	return redis.Bytes(client.Do("GET", key))
}

func (this *RedisClusterAdapter) DEL(key string) (int64, error) {
	log4go.Debug("key is %v", key)
	client := redisClusterClient
	//defer client.Close()
	return redis.Int64(client.Do("DEL", key))
}

func (this *RedisClusterAdapter) KEYS(key string) ([]string, error) {
	log4go.Debug("key is %v", key)
	client := redisClusterClient
	//defer client.Close()
	return redis.Strings(client.Do("KEYS", key))
}

func (this *RedisClusterAdapter) LPUSH(key string, value []interface{}) (interface{}, error) {
	log4go.Debug("key is %v, value is %v", key, value)
	client := redisClusterClient
	//defer client.Close()
	return client.Do("LPUSH", append([](interface{}){key}, value...)...)
}

func (this *RedisClusterAdapter) RPUSH(key string, value []interface{}) (interface{}, error) {
	//log4go.Debug("key is %v, value is %v", key, value)
	defer syscommon.PanicHandler()
	client := redisClusterClient
	//defer client.Close()
	log4go.Debug(append([](interface{}){key}, value...))
	result, err := client.Do("RPUSH", append([](interface{}){key}, value...)...)
	log4go.Debug("err is %v, key is %v, value is %v", err, key, value)
	return result, err
}

func (this *RedisClusterAdapter) LPOP(key string) (string, error) {
	client := redisClusterClient
	//defer client.Close()
	result, err := redis.String(client.Do("LPOP", key))
	log4go.Debug("err is %v, key is %v, value is %v", err, key, result)
	return result, err
}

func (this *RedisClusterAdapter) LLEN(key string) (int, error) {
	client := redisClusterClient
	//defer client.Close()
	result, err := redis.Int(client.Do("LLEN", key))
	log4go.Debug("err is %v, key is %v, value is %v", err, key, result)
	return result, err
}

func (this *RedisClusterAdapter) LRANGE(key string, start int, stop int) ([]string, error) {
	client := redisClusterClient
	//defer client.Close()
	result, err := redis.Strings(client.Do("LRANGE", key, start, stop))
	log4go.Debug("err is %v, key is %v, value is %v", err, key, result)
	return result, err
}

func (this *RedisClusterAdapter) LTRIM(key string, start int, stop int) error {
	client := redisClusterClient
	//defer client.Close()
	_, err := client.Do("LTRIM", key, start, stop)
	log4go.Debug("err is %v, key is %v", err, key)
	return err
}

func (this *RedisClusterAdapter) SADD(key string, value string) error {
	client := redisClusterClient
	//defer client.Close()
	_, err := client.Do("SADD", key, value)
	log4go.Debug("err is %v, key is %v", err, key)
	return err
}

func (this *RedisClusterAdapter) SMEMBERS(key string) ([]string, error) {
	client := redisClusterClient
	//defer client.Close()
	result, err := redis.Strings(client.Do("SMEMBERS", key))
	log4go.Debug("err is %v, key is %v", err, key)
	return result, err
}

func (this *RedisClusterAdapter) SendPipelineCommands(commands []RedisPipelineCommand) ([]interface{}, []error) {
	//	log4go.Debug("commands %v", commands)
	errorList := make([]error, 0, len(commands)+1)

	client := redisClusterClient
	//defer client.Close()

	batch := client.NewBatch()
	for _, value := range commands {
		//		log4go.Debug("Curr Commands index is %v value is %v", index, value)
		//params :=
		//tempParams :=

		//params := [](interface{}){"Cache:/api/brands/888/view:ListValue", "666", "999"}
		//params := []interface{}{"LPUSH", "666", "999"}
		//currErr := conn.Send(value.CommandName, params...)
		//params := append([](interface{}){value.Key}, value.Args...)
		//log4go.Debug("Params : %v", params)
		log4go.Debug("********************")
		log4go.Debug("%v", [](interface{}){value.Key})

		//		for in, v := range value.Args {
		//			log4go.Debug("===== %v %v", in, v)
		//		}

		//		log4go.Debug("%v", value.Args...)
		//		log4go.Debug("%v", append([](interface{}){value.Key}, value.Args...))
		//		log4go.Debug("%v", append([](interface{}){value.Key}, value.Args...)...)
		//client.
		//currErr := client.Send(value.CommandName, append([](interface{}){value.Key}, value.Args...)...)
		currErr := batch.Put(value.CommandName, append([](interface{}){value.Key}, value.Args...)...)

		//currErr := conn.Send(value.CommandName, value.Key, "666", "999")

		if currErr != nil {
			errorList = append(errorList, currErr)
		}
	}

	log4go.Debug("Send finished!!")

	reply, batchErr := client.RunBatch(batch)
	//fulshErr := client.Flush()

	if batchErr != nil {
		errorList = append(errorList, batchErr)

		return nil, errorList
	}

	replys := [](interface{}){}

	replysLength := len(commands)

	var resp int
	for i := 0; i < replysLength; i++ {
		reply, receiveErr := redis.Scan(reply, &resp)

		if receiveErr != nil {
			errorList = append(errorList, receiveErr)
		}

		replys = append(replys, reply)
	}

	log4go.Debug("Receive finished!!")

	if len(errorList) != 0 {
		return replys, errorList
	}

	return replys, nil
}

func Test() {

	hosts := []string{"10.0.71.114:7001", "10.0.71.115:7002", "10.0.71.116:7003"}
	connTimeout := time.Duration(3000 * 1000000)
	readTimeout := time.Duration(3000 * 1000000)
	writeTimeout := time.Duration(3000 * 1000000)
	keepAlive := 500
	aliveTime := time.Duration(50 * 1000 * 1000000)

	CreadRedisCluster(hosts, connTimeout, readTimeout, writeTimeout, keepAlive, aliveTime)

	result, err := GetAdapter().SET("test:test", "test1")
	fmt.Println(result)
	fmt.Println(err)

	result1, err1 := GetAdapter().GET("test:test")
	fmt.Println(result1)
	fmt.Println(err1)

	GetAdapter().DEL("test:test")

	result2, err2 := GetAdapter().GET("test:test")
	fmt.Println(result2)
	fmt.Println(err2)
	//fmt.
}
