## WebsocketV2网络库基本用法

封装的Websocket库，与旧版本最大的区别为，采用的是Zero-Copy的网络层，且运行性能极大提升。
旧版本可能由于底层库的问题，导致GC不及时，存在一定的内存不释放情况，新版有所改善。

注意：需要稳定的，暂时不要使用此库，还有问题。

需要的可以自己跑一下：`samples/WSNetV2`

### 特点

* 自带心跳检测
* 可定义加密解密
* 轻便、简单、易于使用
* 自带消息广播
* 自带连接池管理
* ORIGIN检测体系（防止非法连接）
* 自定义的HEADER检测

* Zero-Copy Upgrader
* Accept直接使用tcp接口

### 测试方式

* 初步测试 1000连接并发 共发送 1178551558 字节 发送6496213次 时长1小时10分钟左右
* 二次测试 当数据包大小达到2048后，多客户端会导致报错，暂时无法使用

我电脑性能有限，客户端和服务端共同运行无法测试出更准确的数据。

### 监听服务端

```
    // 创建一个WSNet
    wsNet := wsv2.NewWebsocket(":47892", "chk_origin", "/ws")

    // 绑定事件
    wsNet.OnHandler = RecvEchoMsg
    wsNet.OnClose = OnClosed

    // 监听端口
    wsNet.Runnable()
```

非常简单，内建一个Connetion的管理器，支持broardcast行为，由于Websocket是frame体系，所以不必自己处理粘黏包，主要是心跳处理较为麻烦，已处理心跳行为和origin检测。

### Callback列表

```
	/// 所有的callBack函数
	// 创建用户DATA
	CreateUserData func() interface{}

	// 通知连接
	OnAccept func(conn *WSConn)
	// 数据包进入
	OnHandler func(conn *WSConn, ownerPak []byte)
	// 连接关闭
	OnClose func(conn *WSConn)
	// 连接成功
	OnConnected func(conn *WSConn)

	// 连接安全性检测 server only
	OnHeader func(key,value string) error
    OnRequest func(url string) error

	// 打包以及加密行为
	Package   func(val interface{}) (data []byte, err error)
	Unpackage func(conn *WSConn, spak *stream.BufferIO) (data [][]byte, err error)

	// 输出panic数据
	Panic func(conn *WSConn, src string)

    // 加密解密函数
	Encrypt func(data []byte) []byte
	Decrypt func(data []byte) []byte
```

所有Callback均提供默认的行为，所以可以选择自己需要的来处理。

### 建立一个WS连接

```
    // 建立一个空的WSNET
    wsNet := ws.NetEmptyWS("test.server.me", "/ws")

    // 绑定事件
    wsNet.OnClose = OnClose
    wsNet.OnHandler = RecvEchoMsg

    // 连接服务器
    _, err := wsNet.Dial("127.0.0.1:47892")
```

连接客户端是不需要WSNET来做Runnable的，因为自身在Goroutine里处理所有的请求行为，所以不需要额外处理。

### 发送数据包

发送数据包一共有2种体系，一种是广播例如：

```
    // 广播字节码
    wsNet.NetCM.Broadcast([]byte(fmt.Sprintf("this is echo msg:%v", i)))
```

还有一种就是单独针对CONN发送数据，例如：

```
    // 发送字节码，不简易但提供该函数
    conn.Send([]byte("test"))

    // 发送一个对象，会调用Callback Package函数
    // 目前会转换成BSON对象的[]byte，可以自定义如何转换。
    conn.SendPak(object{1,2,"test"})
```

### 接收数据

只需要绑定WSNet中的Callback，OnHandler函数即可，会传入收到数据的Conn以及全部包内容，自行处理即可。