package goclub_book_today_first_view

import (
	"context"
	xtime "github.com/goclub/time"
	"github.com/mediocregopher/radix/v4"
	"log"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

var client radix.Client
func init () {
	ctx := context.Background()
	var err error
	client, err = (radix.PoolConfig{}).New(ctx, "tcp", "127.0.0.1:6379") ; if err != nil {
		panic(err)
	}
}
func getKey(userID uint64) string {
	return strings.Join([]string{"todayFirstView", "date", xtime.FormatChinaDate(time.Now()), "userID", strconv.FormatUint(userID, 10)}, ":")
}
// GET SET 是2次操作不符合原子性，所以并发时每个 client 都认为自己的查询结果是没有值
func IncorrectRedisGetSet(userID uint64) error {
	ctx := context.TODO()
	var value string
	result := radix.Maybe{Rcv: &value}
	// todayFirstView:date:2000-11-11:userID:1
	key := getKey(userID)
	// redis: GET todayFirstView:date:2000-11-11:userID:1
	err := client.Do(ctx, radix.Cmd(&result, "GET", key)) ; if err != nil {
		return err
	}
	if result.Null {
		log.Print("User(", userID, ")", "first viewing today")
		// redis: SET todayFirstView:date:2000-11-11:userID:1 1 EX 86400
		err := client.Do(ctx, radix.Cmd(nil, "SET", key, "1", "EX", strconv.FormatInt(60*60*24, 10))) ; if err != nil {
			return err
		}
	} else {
		log.Print("User(", userID, ")", "viewed today")
	}
	return nil
}
func DelKey(userID uint64) {
	key := getKey(userID)
	err := client.Do(context.TODO(), radix.Cmd(nil,"DEL", key)) ; if err != nil {
		panic(err)
	}
}
func TestIncorrectRedisGetSet(t *testing.T) {
	var id uint64 = 1
	DelKey(id)
	var wg sync.WaitGroup
	for i:=0;i<10;i++ {
		wg.Add(1)
		// time.Sleep(time.Second) // 加上 sleep 可以查看不是并发时的结果
		// 并发是必须借助代码测试，自己认为很难测试出并发问题
		go func() {
			defer wg.Done()
			err := IncorrectRedisGetSet(id) ; if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
}
func RedisSetNX(userID uint64) error {
	ctx := context.TODO()
	key := getKey(userID)
	var value string
	result := radix.Maybe{Rcv: &value}
	err := client.Do(ctx, radix.Cmd(&result, "SET", key, "1", "EX", strconv.FormatInt(60*60*24, 10), "NX")) ; if err != nil {
		return err
	}
	// SET NX 结果为 nil 时表示 key 已存在
	if result.Null {
		log.Print("User(", userID, ")", "viewed today")
	} else {
		log.Print("User(", userID, ")", "first viewing today")
	}
	return nil
}


// SETNX 只运行了一次操作，符合原子性操作
func TestRedisSetNX(t *testing.T) {
	var id uint64 = 2
	DelKey(id)
	var wg sync.WaitGroup
	for i:=0;i<10;i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := RedisSetNX(id) ; if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
}