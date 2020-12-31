# 流水线型工作

- 从再一个工作线，导入另一个工作线， 再从这个工作线拿到处理后的对象

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
