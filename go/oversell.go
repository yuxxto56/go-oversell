package redisgo

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)
//声明结构体
type Redis struct {
	Host string
	Port string
	Protocol string
	redisConn redis.Conn
}

//连接redis
func (r *Redis) GetConn() (error){
	conn,err := redis.Dial(r.Protocol,r.Host+":"+string(r.Port))

	if err != nil {
		panic(fmt.Sprintf("redis Dial error,err is %s",err.Error()))
	}
	r.redisConn = conn
	return nil
}

//关闭redis
func (r *Redis) Close(){
    err :=r.redisConn.Close()
	if err != nil{
		fmt.Println("redis.Close error,err is",err.Error())
	}
}


//检测key是否存在
func (r *Redis) Exists(key string) bool{
	 isExists,err := redis.Bool(r.redisConn.Do("EXISTS",key))
	 if err != nil{
		 fmt.Println("redis.Bool error,err is ",err.Error())
		 return false
	 }
	 return isExists
}

//redis Set方法
func (r *Redis) Set(key string,value interface{}) error{
	_,err := r.redisConn.Do("SET", key, value)
	if err != nil{
		fmt.Println("redis.SET error,err is ",err.Error())
		return err
	}
	return nil
}

//redis Set方法
func (r *Redis) Get(key string) string{
	value, err := redis.String(r.redisConn.Do("GET", key))
	if err != nil {
		fmt.Println("redis.GET error,err is ", err.Error())
		return ""
	}
	return value
}

//redis SetNx 实现分布式锁
func (r *Redis) SetNx(key string,value interface{}) int{
	n,err := r.redisConn.Do("SETNX", key, value)
	if err != nil{
		fmt.Println("redis.SET error,err is ",err.Error())
		return 0
	}
	if n == int(1){
		return 1
	}
	return 0
}

//redis Incr 实现递增
func (r *Redis) Incr(key string) int{
	n,err := redis.Int(r.redisConn.Do("INCR", key))
	if err != nil{
		fmt.Println("redis.INCR error,err is ",err.Error())
		return 0
	}
	return n;
}