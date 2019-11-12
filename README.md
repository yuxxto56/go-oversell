# go-oversell
go与redis结合，实现并发场景中的防超卖
### 使用,在main.go文件中编写
```
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
	//目前总库存存
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

```
### 测试
* 编译
```
1、go build main.go
2、./main &
```
* 模拟访问请求
```
 ab -c 100 -n 200  http://127.0.0.1:8000/
 //模拟200次请求，100个用户并发
```
* 结果
```
//结果发现，与我们预想的结果一样，只有5个库存写入库成功，其他失败，防止了超卖现象
2019/11/12 13:22:07 writeDb SUCCESS,storeNum is  96
2019/11/12 13:22:07 writeDb SUCCESS,storeNum is  97
2019/11/12 13:22:07 writeDb SUCCESS,storeNum is  98
2019/11/12 13:22:07 writeDb SUCCESS,storeNum is  99
2019/11/12 13:22:07 writeDb SUCCESS,storeNum is  100
2019/11/12 13:22:07 writeDb Error,storeNum is  101
2019/11/12 13:22:07 writeDb Error,storeNum is  102
2019/11/12 13:22:07 writeDb Error,storeNum is  103
2019/11/12 13:22:07 writeDb Error,storeNum is  104
2019/11/12 13:22:07 writeDb Error,storeNum is  105
2019/11/12 13:22:07 writeDb Error,storeNum is  123
2019/11/12 13:22:07 writeDb Error,storeNum is  124
2019/11/12 13:22:07 writeDb Error,storeNum is  122
2019/11/12 13:22:07 writeDb Error,storeNum is  121
2019/11/12 13:22:07 writeDb Error,storeNum is  119
2019/11/12 13:22:07 writeDb Error,storeNum is  118
2019/11/12 13:22:07 writeDb Error,storeNum is  117
2019/11/12 13:22:07 writeDb Error,storeNum is  116
2019/11/12 13:22:07 writeDb Error,storeNum is  115
2019/11/12 13:22:07 writeDb Error,storeNum is  114
2019/11/12 13:22:07 writeDb Error,storeNum is  113
2019/11/12 13:22:07 writeDb Error,storeNum is  112
2019/11/12 13:22:07 writeDb Error,storeNum is  111
2019/11/12 13:22:07 writeDb Error,storeNum is  110
2019/11/12 13:22:07 writeDb Error,storeNum is  109
2019/11/12 13:22:07 writeDb Error,storeNum is  108
2019/11/12 13:22:07 writeDb Error,storeNum is  107
2019/11/12 13:22:07 writeDb Error,storeNum is  106

```