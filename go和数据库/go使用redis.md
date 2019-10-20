# go使用redis

go语言比较常用的两个redis包[redis](https://github.com/go-redis/redis)和[goredis](https://github.com/gomodule/redigo)。前者封装的比较优雅，功能也比较全。后者像是个工具集合，redis command需要重新封装下才好用在项目中，流行度也很广。本文主要介绍前者[redis](https://github.com/go-redis/redis)的使用。

## redis连接池
创建redis client的代码如下，
```golang
client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
    })
```
这段代码实际上初始化了一个redis连接池而不是单一的redis链接。`NewClient`的源码如下
```golang
func NewClient(opt *Options) *Client {
	opt.init()
	c := Client{
		baseClient: baseClient{
			opt:      opt,
			connPool: newConnPool(opt),
		},
	}
	c.baseClient.init()
	c.init()
	return &c
}
```

我们了解下`redis.Options`结构体中比较关键的几个配置
```golang
type Options struct {
	Network string      //tcp 或者unix
	Addr string       // redis 地址 host:port.
	Dialer func() (net.Conn, error) //创建新的redis连接
	OnConnect func(*Conn) error //链接创建后的钩子函数
	Password string         //redis password 默认为空
	DB int                  //redis db 默认为0
	PoolSize int      //链接池最大活跃数 默认runtime.NumCPU * 10
	MinIdleConns int     //最新空闲链接数 默认是0，空闲连接的用处是加快第一次请求（少了握手过程）
	IdleTimeout time.Duration    //最大空闲时间，需要配置成比redis服务端配置的时间段，默认是5分钟
}
```
再看下 `opt.init()` 做了什么事情
```golang
func (opt *Options) init() {
	if opt.Network == "" {
		opt.Network = "tcp"
	}
	if opt.Addr == "" {
		opt.Addr = "localhost:6379"
	}
	if opt.Dialer == nil {
		opt.Dialer = func() (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   opt.DialTimeout,
				KeepAlive: 5 * time.Minute,
			}
			if opt.TLSConfig == nil {
				return netDialer.Dial(opt.Network, opt.Addr)
			} else {
				return tls.DialWithDialer(netDialer, opt.Network, opt.Addr, opt.TLSConfig)
			}
		}
	}
	if opt.PoolSize == 0 {
		opt.PoolSize = 10 * runtime.NumCPU()
	}
	...
	if opt.IdleTimeout == 0 {
		opt.IdleTimeout = 5 * time.Minute
	}
	...
}
```
从上面的代码可以看出来`opt.init()`实际上初始化了默认配置。一般来说这样子的默认配置能满足大部分项目的需求。

## qick start
```golang
var client *redis.Client

func init()  {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func ping()  {
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func main() {
    ping() 
}
```
```bash
# go run main.go
PONG <nil>
```

## kv操作
```golang
func keyv()  {
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}
```

## list操作
```golang
func list()  {
	client.LPush("wida","ddddd")
	ret,_ := client.LPop("wida").Result()
	fmt.Println(ret)
	_,err := client.LPop("wida").Result()
	if err == redis.Nil {
		fmt.Println("empty list")
	}
}
```

## set操作
```golang
func set()  {
	client.SAdd("set","wida")
	client.SAdd("set","wida1")
	client.SAdd("set","wida2")
	ret,_:= client.SMembers("set").Result()
	fmt.Println(ret)
}
```

## sortset 操作
```golang
func sortset()  {
	client.ZAdd("page_rank", redis.Z{10 ,"google.com"})
	client.ZAdd("page_rank", redis.Z{9 ,"baidu.com"},redis.Z{8 ,"bing.com"})
	ret,_:=client.ZRangeWithScores("page_rank",0,-1).Result()
	fmt.Println(ret)
}
```

## hashset 操作
```golang
func hash()  {
	client.HSet("hset","wida",1)
	ret ,_:=client.HGet("hset","wida").Result()
	fmt.Println(ret)
	ret ,err:=client.HGet("hset","wida3").Result()
	if err == redis.Nil {
		fmt.Println("key not found")
	}
	client.HSet("hset","wida2",2)
	r ,_:=client.HGetAll("hset").Result()
	fmt.Println(r)
}
```

## pipeline
```golang
func pipeline()  {
	pipe := client.Pipeline()
	pipe.HSet("hset2","wida",1)
	pipe.HSet("hset2","wida2",2)

	ret := pipe.HGetAll("hset2")
	fmt.Println(pipe.Exec())
	fmt.Println(ret.Result())
}
```