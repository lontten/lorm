# lorm


### init lorm
```go

	conf := lorm.MysqlConf{
		Host:     "127.0.0.1",
		Port:     "3306",
		DbName:   "test",
		User:     "root",
		Password: "123456",
		Other:    "sslmode=disable TimeZone=Asia/Shanghai",
	}
	DB := lorm.MustConnect(&conf,nil)

```


```go
type User struct {
	ID   int64
	Name string
	Age  int
}

func (u User) TableConf() *lorm.TableConfContext {
    return lorm.TableConf("t_user").
        PrimaryKeys("id").
        AutoColumn("id")
}

```
### Insert
```go
	user := User{
		Name: "tom",
		Age:  12,
	}
	num, err := lorm.Insert(DB,&user)
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	//return id
	fmt.Println(user.ID)

```

### update
```go
	user := User{
		ID:   1,
		Name: "tom",
		Age:  12,
	}
	
	//根据主键更新
	num, err := lorm.UpdateByPrimaryKey(DB,&user)
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	
	user := User{
		Name: "tom",
		Age:  12,
	}
	
	// lorm.W() 条件构造器
	num, err := lorm.Update(DB, &user, lorm.W().
            Eq("id", 1).
            In("id", 1, 2).
            Gt("id", 1).
            IsNull("name").
            Like("name", "abc"),
        )
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	
	user := User{
		Name: "tom",
		Age:  12,
	}
	
	// 特殊配置
	num, err := lorm.Update(DB, &user, lorm.W(), lorm.E().
            SetNull("age").  // 设置age字段为null
            TableName("user2"). // 临时自定义表名 为 user2
            ShowSql(). // 打印执行 sql
		    NoRun(),  // 不具体执行sql，配合 ShowSql() 用来调试sql
        )
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)

```
### delete
```go
 
	
	//根据主键删除
	num, err := db.Delete(User{}).ByPrimaryKey(id)
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	
	//----------------
	
 
	//根据条件删除
	num, err := db.Delete(User{}).ByModel(NullUser{
		Name: types.NewString("tom"),
	})
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	//-------------------
	
	 
	
	//使用条件构造器
	num, err := db.Delete(User{}).ByWhere(new(lorm.WhereBuilder).
		Eq("id", user.ID,true).
		NoLike("age", *user.Name, user.Name != nil).
		Ne("age", user.Age,false))
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)

```

###select
```go
	user := User{}
	num, err := db.Select(User{}).ByPrimaryKey(id).ScanOne(&user)
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	
	fmt.Println(user)
	//-----------------
	
	users := make([]User,0)
	num, err := db.Select(User{}).ByPrimaryKey(id1,id2,id3).ScanList(&users)
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	
	fmt.Println(user)
	//-----------------
	
	
	users := make([]User, 0)
	num, err := db.Select(User{}).ByModel(NullUser{
		Name: types.NewString("tom"),
		Age:  types.NewInt(12),
	}).ScanList(&users)
	if err != nil {
		return err
	}
	// num 查询的数据个数
	fmt.Println(num)
	
	fmt.Println(users)
	//----------------
	
	user := User{}
	//随机获取一个
	num, err := db.Select(User{}).ByModel(NullUser{
		Name: types.NewString("tom"),
		Age:  types.NewInt(12),
	}).ScanFirst(&user)
	if err != nil {
		return err
	}
	// num 查询的数据个数
	fmt.Println(num)
	
	fmt.Println(user)
	//-----------------------
	
	
	has, err := db.Select(User{}).ByModel(NullUser{
		Name: types.NewString("tom"),
		Age:  types.NewInt(12),
	})
	if err != nil {
		return err
	}
	// has 查询是否存在数据
	fmt.Println(num)
	
	
	
	//----------------------------
	has, err := db.Has(User{}).ByWhere(new(lorm.WhereBuilder).
		Eq("id", user.ID, true).
		NoLike("age", *user.Name, user.Name != nil).
		Ne("age", user.Age, false))
	if err != nil {
		return err
	}
	// has 查询是否存在数据
	fmt.Println(has)

	
	
	
	has, err := db.Has(User{}).ByWhere(new(lorm.WhereBuilder).
		Eq("id", user.ID, true).
		NoLike("age", *user.Name, user.Name != nil).
		Ne("age", user.Age, false))
	if err != nil {
		return err
	}
	// num 查询是否存在数据
	fmt.Println(has)

	
	
	
```


###tx
```go
	tx := Db.Begin()
    err := tx.Commit()
    err := tx.Rollback()
```





### init lorm pool log
```go

	path := "./log/go.log"
	writer, _ := rotatelogs.New(
		path+".%Y-%m-%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(365*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	newLogger := log.New(writer, "\r\n", log.LstdFlags)

	var dbName = pg.DbName

	pgConf := lorm.PgConf{
		Host:     pg.Ip,
		Port:     pg.Port,
		DbName:   pg.dbName,
		User:     pg.User,
		Password: pg.Pwd,
		Other:    "sslmode=disable TimeZone=Asia/Shanghai",
	}
	poolConf := lorm.PoolConf{
		MaxIdleCount: 10,
		MaxOpen:      100,
		MaxLifetime:  time.Hour,
		Logger:       newLogger,
	}
	ormConf := lorm.OrmConf{
		TableNamePrefix: "t_",
		PrimaryKeyNames: []string{"id"},
	}

	db := lorm.MustConnect(&pgConf, &poolConf).OrmConf(&ormConf)

```
