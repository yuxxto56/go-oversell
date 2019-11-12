package main

import (
	"fmt"
  	redisgo "go-oversell/go"
	"log"
	"net/http"

)

func HanderOverSell(w http.ResponseWriter, r *http.Request){
	red := &redisgo.Redis{
		Host:     "192.168.1.132",
		Port:     "6379",
		Protocol: "tcp",
	}
	err := red.GetConn()
	if err != nil{
		fmt.Println("connect error,err is ",err.Error())
		return
	}
	defer red.Close()
	//设置初始库存
	storeNum := 95
	//设置总库存
	limitStoreNum := 100
	//设置key
	redisKey := "huawei_p30_num_100"
	if !red.Exists(redisKey){
		//实现分布式锁
		red.SetNx(redisKey,storeNum)
	}
	//递增完再判断
	num := red.Incr(redisKey)
	if num >limitStoreNum{
		log.Println("writeDb Error,storeNum is ",num)
	}else{
		log.Println("writeDb SUCCESS,storeNum is ",num)
	}
}

func main() {
	http.HandleFunc("/", HanderOverSell)
	err := http.ListenAndServe("0.0.0.0:8000", nil)
	if(err != nil){
		fmt.Println("start Http Error,err is ",err)
	}
	fmt.Println("start Http,Success.0.0.0:8000")
}
