# 尝试用golang 1.18泛型实现orm

这几天golang社区对泛型的讨论非常多的，一片热火朝天的景象。对我们广大gopher来说总归是好事。

泛型很有可能会颠覆我们之前的很多设计，带着这种疑问和冲动，我准备尝试用golang泛型实现几个orm的常见功能。

本文并没完全实现通用的orm，只是探讨其实现的一种方式提供各位读者做借鉴。

## 创建Table

虽然golang有了泛型，但是目前在标准库sql底层还没有改造，目前还有很多地方需要用到`reflect`。

```go
func CreateTable[T any](db *sql.DB) {
	var a T
	t := reflect.TypeOf(a)
	tableName := strings.ToLower(t.Name())

	var desc string
	for i := 0; i < t.NumField(); i++ {
		columnsName := strings.ToLower(t.Field(i).Name)
		var columnType string
		switch t.Field(i).Type.Kind() {
		case reflect.Int:
			columnType = "integer"
		case reflect.String:
			columnType = "text"
		}
		desc += columnsName + " " + columnType
		if i < t.NumField()-1 {
			desc += ","
		}
	}
	sqlStmt := fmt.Sprintf(`create table if not exists %s (%s);`, tableName, desc)
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
```

调用方式

```go
type Person struct {
	ID   int
	Name string
	Age  int
}

type Student struct {
	ID   int
	Name string
	No   string
}

var db sql.DB
//init db
//...
CreateTable[Person](db)
CreateTable[Student](db)
```

这个部分跟传统的orm使用上没有太大区别，没办法不使用反射的情况下，泛型的方式可能变得有点繁琐。

## 写入数据

```go
func Create[T any](db *sql.DB, a T) {
	//没有办法这边还是得使用反射
	t := reflect.TypeOf(a)
	tableName := strings.ToLower(t.Name())

	var columns []string
	var spacehold []string
	for i := 0; i < t.NumField(); i++ {
		columns = append(columns, strings.ToLower(t.Field(i).Name))
		spacehold = append(spacehold, "?")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(
		fmt.Sprintf("insert into %s(%s) values(%s)",
			tableName,
			strings.Join(columns, ","),
			strings.Join(spacehold, ",")))

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	v := reflect.ValueOf(a)

	var values []any

	for i := 0; i < t.NumField(); i++ {
		if v.FieldByName(t.Field(i).Name).CanInt() {
			values = append(values, v.FieldByName(t.Field(i).Name).Int())
		} else {
			values = append(values, v.FieldByName(t.Field(i).Name).String())
		}
	}
	_, err = stmt.Exec(values...)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}
```

调用方式

```go
var p1 = Person{
	ID:   1,
	Name: "wida",
}
Create[Person](db, p1)
var s1 = Student{
	ID:   1,
	Name: "wida",
	No:   "1111",
}
Create[Person](db, p1)
Create[Student](db, s1)
```
和创建table类似，写入数据好像比没有之前的orm有优势。


## 读取数据

读取数据是非常高频的操作，所以我们稍作封装。

```go
type Client struct {
	db *sql.DB
}

type Query[T any] struct {
	client *Client
}

func NewQuery[T any](c *Client) *Query[T] {
	return &Query[T]{
		client: c,
	}
}
//反射到struct
func ToStruct[T any](rows *sql.Rows, to T) error {
	v := reflect.ValueOf(to)
	if v.Elem().Type().Kind() != reflect.Struct {
		return errors.New("Expect a struct")
	}

	scanDest := []any{}
	columnNames, _ := rows.Columns()

	addrByColumnName := map[string]any{}

	for i := 0; i < v.Elem().NumField(); i++ {
		oneValue := v.Elem().Field(i)
		columnName := strings.ToLower(v.Elem().Type().Field(i).Name)
		addrByColumnName[columnName] = oneValue.Addr().Interface()
	}

	for _, columnName := range columnNames {
		scanDest = append(scanDest, addrByColumnName[columnName])
	}
	return rows.Scan(scanDest...)
}

func (q *Query[T]) FetchAll(ctx context.Context) ([]T, error) {
	var items []T

	var a T
	t := reflect.TypeOf(a)

	tableName := strings.ToLower(t.Name())
	rows, err := q.client.db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var c T
		ToStruct(rows, &c)
		items = append(items, c)
	}
	return items, nil
}
```

调用方式

```go
var client = &Client{
	db: db,
}

{
	query := NewQuery[Person](client)
	all, err := query.FetchAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, person := range all {
		log.Println(person)
	}
}
{
	query := NewQuery[Student](client)
	all, err := query.FetchAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, person := range all {
		log.Println(person)
	}
}
```

稍微比原先的orm方式有了多一点想象空间，比如 在`[T any]`做更明确的约束，比如要求实现`Filter`定制方法。

# 总结

鉴于本人能力还认证有限，目前还没有发现泛型对orm剧烈的改进和突破的可能。未来如果go对底层sql做出改动，或者实现诸如`Rust`那种`Enum`方式，可能会带来更多的惊喜。