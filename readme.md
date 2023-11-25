# timecounter

基于时间轮思想进行短时间的时序统计。 见 https://juejin.cn/post/7305013302416244774

用法：


```go

counter := NewCounter(10 * time.Minute, 100)

counter.PutNow(1)
counter.PutN(time.Now(), 10)

fmt.Println(counter.SumAll()) // output: 11
```