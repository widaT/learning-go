# go使用mongodb
mongodb的go语言驱动比较流行是[mongo-go-driver](https://github.com/mongodb/mongo-go-driver)和[mgo](https://github.com/go-mgo/mgo),前者是mongodb官方出的go驱动，mgo之前一直是个人在维护的，后来作者由于个原因已经放弃维护，后面出现[mgo](https://github.com/globalsign/mgo)的分支还在继续维护不过更新频率还是比较低。本文主要介绍`mongo-go-driver`这个mongodb官方维护的版本，目前`mongo-go-driver`的版本已经是1.1.2已经可以用在生产环境。

## mongodb client
目前大多数golang数据库驱动client 都会用连接池的方式来运行。 `mongo-go-driver`也不例外。
连接MongoDB代码如下：
```golang
client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
```
运行下main.go程序
```bash
# go run main.go
Connected to MongoDB!
```
这边需要注意下，采用连接池后一般很少自己close client而是继续放在链接池中举行供其他请求使用。
如果你确定已经不需要使用Client，可以使用如下代码收到关闭。
```golang
err = client.Disconnect(nil）
if err != nil {
    log.Fatal(err)
}
fmt.Println("Connection to MongoDB closed.")
```

## 获取集合handle
```golang
collection = client.Database("testing").Collection("students")
```

## 为collection创建索引
```golang
indexView := collection.Indexes()
	ret,err :=indexView.CreateOne(context.Background(), mongo.IndexModel{
		Keys:  bsonx.Doc{{"name", bsonx.Int32(-1)}},
		Options: options.Index().SetName("testname").SetUnique(true), //这边设置了唯一限定，不设定默认不是唯一的
	})
fmt.Println(ret,err)
```

## CURD操作

###　添加
```golang
func add()  {
	wida := Student{"wida", 32, "8001"}
	amy := Student{"amy", 25, "8002"}
	sunny := Student{"sunny", 35, "8003"}

	ret, err := collection.InsertOne(nil, wida)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("写入一个文档", ret.InsertedID)
	student := []interface{}{amy, sunny}

	ret2, err := collection.InsertMany(nil, student)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("写入多个文档 ", ret2.InsertedIDs)
}
```

### 查询
```golang
func find()  {
	var result Student
	err := collection.FindOne(nil,  bson.D{{"name", "amy"}, {"age", 25}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", result)
	findOptions := options.Find()
	findOptions.SetLimit(3)

	var results []*Student

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem Student
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(nil)

	for _,r := range results {
		fmt.Printf("%+v\n", r)
	}
}
```
### 更新
```golang
filter := bson.D{{"no", "8001"}}
	update := bson.D{
		{"$set", bson.D{
			{"age", 33},
			{"name", "wida2"},
		}},
	}
	ret, err := collection.UpdateOne(nil, filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", ret.MatchedCount, ret.ModifiedCount)
}
```

### 删除
```golang
func del()  {
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}}) // 这边删除全部文档
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents\n", deleteResult.DeletedCount)
}
```

## 参考文档
- [MongoDB Go Driver](https://docs.mongodb.com/ecosystem/drivers/go/)
- [MongoDB Go Driver Tutorial](https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial)