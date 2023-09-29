package learn_redis_lock

import (
	"context"
	"sync"
	"testing"
)

func Test_blockingLock(t *testing.T){
	addr := "127.0.0.1:6379"
	passwd := "123456"

	client := NewClient("tcp", addr, passwd)
	lock1 := NewRedisLock("test_key", client, WithExpireSeconds(1))
	//lock2 := NewRedisLock("test_key", client, WithBlock(), WithBlockWaitingSeconds(2))

	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := lock1.Lock(ctx); err != nil {
			t.Error(err)
			return
		}
	}()
	//
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	if err := lock2.Lock(ctx); err != nil {
	//		t.Error(err)
	//		return
	//	}
	//}()

	wg.Wait()

	t.Log("success")
}

func Test_nonblockingLock(t *testing.T) {
	// 请输入 redis 节点的地址和密码
	addr := "127.0.0.1:6379"
	passwd := "123456"

	client := NewClient("tcp", addr, passwd)
	lock1 := NewRedisLock("test_key", client, WithExpireSeconds(1))


	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := lock1.Lock(ctx); err != nil {
			t.Error(err)
			return
		}
	}()


	wg.Wait()
	t.Log("success")
}