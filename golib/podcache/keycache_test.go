package podcache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

type mockKeyData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Data string `json:"keyData"`
}

func mockKeyDataLoader(_ context.Context, keys []int) (map[int]mockKeyData, error) {
	data := make(map[int]mockKeyData)
	for _, key := range keys {
		data[key] = mockKeyData{
			Id:   key,
			Name: fmt.Sprintf("name-%d", key),
			Data: fmt.Sprintf("data-%d-%d", key, time.Now().UnixNano()),
		}
	}
	return data, nil
}

func Test_keyCache_BatchGet(t *testing.T) {
	cache := NewKeyCache[int, mockKeyData](&sync.Mutex{}, mockKeyDataLoader, WithExpires[int, mockKeyData](3), WithMaxSize[int, mockKeyData](100))

	printKeyData(cache.BatchGet(context.Background(), []int{1, 2, 3}))
	printKeyData(cache.BatchGet(context.Background(), []int{3}))
	printKeyData(cache.BatchGet(context.Background(), []int{4}))
	printKeyData(cache.BatchGet(context.Background(), []int{2, 5, 6}))
	time.Sleep(3 * time.Second)
	printKeyData(cache.BatchGet(context.Background(), []int{1, 2, 3}))
	printKeyData(cache.BatchGet(context.Background(), []int{3}))
	printKeyData(cache.BatchGet(context.Background(), []int{4}))
	printKeyData(cache.BatchGet(context.Background(), []int{2, 5, 6}))
}

func printKeyData(data map[int]mockKeyData, err error) {
	log.Println(data, err)
}
