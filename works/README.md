# 工作池

使用 参考test文件
```
    inqueue := make(chan interface{}, 30)
    NewDispatcher(2, inqueue, nil).Run()

    for i := 0; ; i++ {
        wi := WI(i)
        inqueue <- wi
        time.Sleep(100 * time.Millisecond)
    }
```
