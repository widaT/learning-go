# go使用es

elasticsearch有官方的golang驱动[go-elasticsearch](https://github.com/elastic/go-elasticsearch)这个项目比较新，
另外一个常用的是 [elastic](https://github.com/olivere/elastic)，这两个驱动文档和demo都比较少。es的查询语法也相对复杂，很多查询方式去翻翻它们的test文件才能发现方式。本小节使用`elastic`做演示，注意不同elasticsearch版本对于不同的client版本，例如elasticsearch 5.5.3对应的client版本为`gopkg.in/olivere/elastic.v5`。如果这个对应关系错误，很可能程序会出错，这个在`https://github.com/olivere/elastic`的readme文档也有介绍。
本小节的demo主要基于 `elasticsearch 5.5.3`，client为`gopkg.in/olivere/elastic.v5`。

## go链接es
```golang
var client *elastic.Client
func init()  {
	var err error
	client, err = elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Fatal(err)
	}
}
```

## CURD

测试数据结构体
```golang
type Item struct {
	Id               int64  `json:"id"`
	Appid            string `json:"appid"`
	AppBAutoId       string `json:"app_b_auto_id"`
}
```

### 添加
```golang
func add()  {
	//add one
	item := Item{Id:int64(21),Appid:fmt.Sprintf("app_%d",21),AppBAutoId:fmt.Sprintf("app_%d",21+200)}
	put, err := client.Index().
		Index("es_test").
		Type("test").
		Id("1").       //这个id也可以指定,不指定的话 es自动生成一个
		BodyJson(item).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println(put)
	//add many
	bulkRequest := client.Bulk()
	for i:=0;i<20;i++ {
		item := Item{Id:int64(i),Appid:fmt.Sprintf("app_%d",i),AppBAutoId:fmt.Sprintf("app_%d",i+200)}
		bulkRequest.Add(elastic.NewBulkIndexRequest().Index("es_test").Type("test").Doc(item))
	}
	bulkRequest.Do(context.TODO())
}

```
### 查找
```golang
func find()  {
	//find one
	get1, err := client.Get().
		Index("es_test").
		Type("test").
		Id("1").
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}
	var ttyp Item
	json.Unmarshal(*get1.Source,&ttyp)
	fmt.Println("item",ttyp)

	//find many
	searchResult, err := client.Search().
		Index("es_test"). Sort("id", true).
		Type("test").From(0).Size(100).
		Do(context.TODO())
	if err != nil {
		panic(err)
	}
	if searchResult.Hits.TotalHits >0 {
		var ttyp Item
		for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
			t := item.(Item)
			fmt.Println("item ",t)
		}

	}
}
```
### 更新
```golang
func update() {
	fmt.Println(client.Update().Index("es_test").Type("test").Id("1").
		Doc(map[string]interface{}{"appid": "app_23"}).Do(context.TODO()))
}
```
### 删除
```golang
func delete()  {
	fmt.Println(client.Delete().Index("es_test").Type("test").Id("1").Do(context.TODO()))
}

```
## 统计查询
```golang
func agg() {
	//获取最大的id
	searchResult, err := client.Search().
		Index("es_test").Type("test").
		Aggregation("max_id", elastic.NewMaxAggregation().Field("id")).Size(0).Do(context.TODO())
	if err != nil {
		panic(err)
	}
	var a map[string]float32
	if searchResult != nil {
		if v, found := searchResult.Aggregations["max_id"]; found {
			json.Unmarshal([]byte(*v), &a)
			fmt.Println(a)
		}
	}
	//统计id相同的文档数
	searchResult, err = client.Search().
		Index("es_test").Type("test").
		Aggregation("count", elastic.NewTermsAggregation().Field("id")).Size(0).Do(context.TODO())
	if err != nil {
		panic(err)
	}
	if searchResult != nil {
		if v, found := searchResult.Aggregations["count"]; found {
			var ar elastic.AggregationBucketKeyItems
			err := json.Unmarshal(*v, &ar)
			if err != nil {
				fmt.Printf("Unmarshal failed: %v\n", err)
				return
			}

			for _, item := range ar.Buckets {
				fmt.Printf("id ：%v: count ：%v\n", item.Key, item.DocCount)

			}
		}
	}
}
```
# 参考资料
- [olivere/elastic](https://github.com/olivere/elastic)
- [elastic/wiki](https://github.com/olivere/elastic/wiki)
- [QueryDSL](https://github.com/olivere/elastic/wiki/QueryDSL)