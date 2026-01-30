# singleflight泛型库

## 使用技巧

### 1. 超时控制

DoChan() 通过 channel 返回结果。因此可以使用 select 语句实现超时控制:

```go
var g singleflight.Group[User]

ch := g.DoChan(key, func() (User, error) {
    ret, err := find(context.Background(), key)
    return ret, err
})
// Create our timeout
timeout := time.After(500 * time.Millisecond)

var ret singleflight.Result[User]
select {
case <-timeout: // Timeout elapsed
        fmt.Println("Timeout")
    return
case ret = <-ch: // Received result from channel
    fmt.Printf("index: %d, val: %v, shared: %v\n", j, ret.Val, ret.Shared)
}
```


### 2. 饱和度控制

在一些对可用性要求极高的场景下，往往需要一定的请求饱和度来保证业务的最终成功率。一次请求还是多次请求，对于下游服务而言并没有太大区别，此时使用 singleflight 只是为了降低请求的数量级，那么使用 Forget() 提高下游请求的并发:

```go
var g singleflight.Group[User]

v, _, shared := g.Do(key, func() (User, error) {
    go func() {
        time.Sleep(10 * time.Millisecond)
        fmt.Printf("Deleting key: %v\n", key)
        g.Forget(key)
    }()
    ret, err := find(context.Background(), key)
    return ret, err
})
```

当有一个并发请求超过 10ms，那么将会有第二个请求发起，此时只有 10ms 内的请求最多发起一次请求，即最大并发：100 QPS。单次请求失败的影响大大降低。



PS https://www.cyningsun.com/01-11-2021/golang-concurrency-singleflight.html
