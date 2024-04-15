## 链码编写记录

### 链码部署操作流程

进入对应文件夹(前面因人而异)

```
cd test-network
```

测试网络启动

```
./network.sh up
```

建立通道

```
./network.sh createChannel
```

链码部署

```
./network.sh deployCC -ccn basic -ccp ../mywork -ccl go
```

这里我新建了一个文件夹mywork,里面有链码文件did.go,文件内容与博客里面测试网络使用链码相同,后期我们编写链码只需要更改did.go的内容就行。

应用程序测试(进入对应文件夹后,但一般编写的时候会放在同一文件夹方便操作)

```
go run app.go
```

