package learn_redis_lock

import (
	"context"
	"errors"
	"github.com/gomodule/redigo/redis"
	"strings"
	"time"
)

type LockClient interface {
	SetNEX(ctx context.Context, key, value string, expireSeconds int64)(int64, error)
	Eval(ctx context.Context, src string, keyCount int, keysAndArgs []interface{})(interface{}, error)
}

type Client struct {
	ClientOptions
	pool *redis.Pool
}

func (c *Client) getRedisPool() *redis.Pool{
	return &redis.Pool{
		MaxIdle: c.maxIdle,
		IdleTimeout: time.Duration(c.idleTimeoutSeconds) * time.Second,
		Dial: func() (redis.Conn, error) {
			c,err  := c.getRedisConn()
			if err != nil{
				return nil,err
			}
			return c,nil
		},
		MaxActive: c.maxActive,
		Wait: c.wait,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_,err := c.Do("PING")
			return err
		},
	}
}

func (c *Client)GetConn(ctx context.Context)(redis.Conn, error){
	return c.pool.GetContext(ctx)
}

func (c *Client) getRedisConn() (redis.Conn, error) {
	if c.address == ""{
		panic("cannot get redis address from config")
	}

	var dialOpts []redis.DialOption
	if len(c.password) > 0{
		dialOpts = append(dialOpts, redis.DialPassword(c.password))
	}
	conn,err := redis.DialContext(context.Background(), c.network,c.address, dialOpts...)
	if err != nil{
		return nil,err
	}
	return conn,err
}

func NewClient(network, address, password string, opts ...ClientOption)*Client{
	c := Client{
		ClientOptions:ClientOptions{
			network: network,
			password: password,
			address: address,
		},
	}
	for _,opt := range  opts{
		opt(&c.ClientOptions)
	}
	repairClient(&c.ClientOptions)
	pool := c.getRedisPool()
	return &Client{
		pool: pool,
	}
}


func (c *Client)Get(ctx context.Context, key string)(string, error){
	if key == ""{
		return "",errors.New("redis get cannot be empty")
	}

	conn,err := c.pool.GetContext(ctx)
	if err != nil{
		return "",err
	}
	defer conn.Close()

	return redis.String(conn.Do("GET", key))
}

func (c *Client)Set(ctx context.Context, key ,value string)(int64,error){
	if key == "" || value == ""{
		return -1,errors.New("redis set key or value cannot be empty")
	}

	conn,err := c.pool.GetContext(ctx)
	if err != nil{
		return -1,err
	}
	defer  conn.Close()

	resp,err := conn.Do("SET", key,value)
	if err != nil{
		return -1,err
	}
	if respStr,ok := resp.(string); ok && strings.ToLower(respStr) == "ok"{
		return 1,nil
	}

	return redis.Int64(resp,err)
}

func (c *Client)SetNEX(ctx context.Context, key,value string, expireSeconds int64)(int64,error){
	if key == "" || value == ""{
		return -1,errors.New("redis set keynx or value cannot be empty")
	}

	conn,err := c.pool.GetContext(ctx)
	if err != nil{
		return -1,err
	}
	defer  conn.Close()

	resp,err := conn.Do("SET", key,value,"EX", expireSeconds, "NX")
	if err != nil{
		return -1,err
	}
	if respStr,ok := resp.(string); ok && strings.ToLower(respStr) == "ok"{
		return 1,nil
	}

	return redis.Int64(resp,err)
}

func (c *Client)SetNX(ctx context.Context, key,value string) (int64,error){
	if key == "" || value == ""{
		return -1,errors.New("redis set key nx or value cannot be empty")
	}

	conn,err := c.pool.GetContext(ctx)
	if err != nil{
		return -1,err
	}
	defer  conn.Close()

	resp,err := conn.Do("SET", key,value,"NX")
	if err != nil{
		return -1,err
	}
	if respStr,ok := resp.(string); ok && strings.ToLower(respStr) == "ok"{
		return 1,nil
	}

	return redis.Int64(resp,err)
}

func (c *Client) Del(ctx context.Context, key string)error{
	if key == ""{
		return errors.New("redis del key cannot be empty")
	}

	conn,err := c.pool.GetContext(ctx)
	if err != nil{
		return err
	}
	defer conn.Close()
	_,err = conn.Do("DEL",key)
	return err
}


func (c *Client)Incr(ctx context.Context, key string)(int64,error){
	if key == ""{
		return -1,errors.New("redis incr key cannot be empty")
	}

	conn,err := c.pool.GetContext(ctx)
	if err != nil{
		return -1,err
	}
	defer conn.Close()
	return redis.Int64(conn.Do("INCR", key))
}

func (c *Client)Eval(ctx context.Context, src string, keyCount int, keysAndArgs []interface{})(interface{}, error){
	args := make([]interface{}, 2 + len(keysAndArgs))

	args[0] = src
	args[1] = keyCount
	copy(args[2:], keysAndArgs)

	conn,err := c.pool.GetContext(ctx)
	if err != nil{
		return -1,err
	}
	defer conn.Close()

	return conn.Do("EVAL", args...)
}

