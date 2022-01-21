# lsql

[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/lontten/lsql/main/LICENSE)


### init lsql
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

	pgConf := lsql.PgConf{
		Host:     pg.Ip,
		Port:     pg.Port,
		DbName:   pg.dbName,
		User:     pg.User,
		Password: pg.Pwd,
		Other:    "sslmode=disable TimeZone=Asia/Shanghai",
	}
	poolConf := lsql.PoolConf{
		MaxIdleCount: 10,
		MaxOpen:      100,
		MaxLifetime:  time.Hour,
		Logger:       newLogger,
	}
	ormConf := lsql.OrmConf{
		TableNamePrefix: "t_",
		PrimaryKeyNames: []string{"id"},
	}

	db := lsql.MustConnect(&pgConf, &poolConf).OrmConf(&ormConf)

```
```go
type User struct {
	ID   types.UUID `json:"id"  tableName:"public.t_user"`
	Name string     `json:"info"`
	Age  int        `json:"age"`
}

type NullUser struct {
	ID   *types.UUID `json:"id"  tableName:"public.t_user"`
	Name *string     `json:"info"`
	Age  *int        `json:"age"`
}

```
### create
```go
	user := NullUser{
		ID:   types.NewV4P(),
		Name: types.NewString("tom"),
		Age:  types.NewInt(12),
	}
	// create 是引用，会返回id
	num, err := db.Insert(&user)
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	//return id
	fmt.Println(user.ID)
	
	//-----------------------

	user := NullUser{
		ID:   types.NewV4P(),
		Name: types.NewString("tom"),
		Age:  types.NewInt(12),
	}
	
	// create 不是引用，不会返回id
	num, err := db.Insert(user)
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	// nil
	fmt.Println(user.ID)

```

###create or update
```go
	user := NullUser{
		ID:   types.NewV4P(),
		Name: types.NewString("tom"),
		Age:  types.NewInt(12),
	}
	
	// 创建或更新，根据主键
	num, err := db.InsertOrUpdate(&user).ByPrimaryKey()
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	//------------------
	
	user := NullUser{
		Name: types.NewString("tom"),
		Age:  types.NewInt(12),
	}
	
	// 创建或更新，根据 name,age组合的唯一索引；mysql不支持此功能
	num, err := db.InsertOrUpdate(&user).ByUnique([]string{"name","age"})
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)

```

### update
```go
	user := NullUser{
		ID:   types.NewV4P(),
		Name: types.NewString("tom"),
		Age:  types.NewInt(12),
	}
	
	//根据主键更新
	num, err := db.Update(&user).ByPrimaryKey()
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	
	//----------------
	
	user := NullUser{
		ID:   types.NewV4P(),
		Name: types.NewString("tom"),
		Age:  types.NewInt(12),
	}
	
	//根据条件更新
	num, err := db.Update(&user).ByModel(NullUser{
		Name: types.NewString("tom"),
	})
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	//-------------------
	
	
	user := NullUser{
		ID:   types.NewV4P(),
		Name: types.NewString("tom"),
		Age:  types.NewInt(12),
	}
	
	//使用条件构造器
	num, err := db.Update(&user).ByWhere(new(lorm.WhereBuilder).
		Eq("id", user.ID,true).
		NoLike("age", *user.Name, user.Name != nil).
		Ne("age", user.Age,false))
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
	num, err := db.Delete(User{}).ByWhere(new(lsql.WhereBuilder).
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
	has, err := db.Has(User{}).ByWhere(new(lsql.WhereBuilder).
		Eq("id", user.ID, true).
		NoLike("age", *user.Name, user.Name != nil).
		Ne("age", user.Age, false))
	if err != nil {
		return err
	}
	// has 查询是否存在数据
	fmt.Println(has)

	
	
	
	has, err := db.Has(User{}).ByWhere(new(lsql.WhereBuilder).
		Eq("id", user.ID, true).
		NoLike("age", *user.Name, user.Name != nil).
		Ne("age", user.Age, false))
	if err != nil {
		return err
	}
	// num 查询是否存在数据
	fmt.Println(has)

	
	
	
```