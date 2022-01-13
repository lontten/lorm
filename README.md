# lorm

[![Build Status](https://travis-ci.org/lontten/lorm.svg?branch=main)](https://travis-ci.org/lontten/lorm) 
[![Coverage Status](https://coveralls.io/repos/github/lontten/lorm/badge.svg?branch=main)](https://coveralls.io/github/lontten/lorm?branch=main) 
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/lontten/lorm) 
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/lontten/lorm/main/LICENSE)


### init lorm
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

	engine := lorm.MustConnect(&pgConf, &poolConf).Db(&ormConf)

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
	num, err := engine.Table.Create(&user)
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
	num, err := engine.Table.Create(user)
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
	num, err := Engine.Table.CreateOrUpdate(&user).ByPrimaryKey()
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
	num, err := Engine.Table.CreateOrUpdate(&user).ByUnique([]string{"name","age"})
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
	num, err := Engine.Table.Update(&user).ByPrimaryKey()
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
	num, err := Engine.Table.Update(&user).ByModel(NullUser{
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
	num, err := Engine.Table.Update(&user).ByWhere(new(lorm.WhereBuilder).
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
	num, err := Engine.Table.Delete(User{}).ByPrimaryKey(id)
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	
	//----------------
	
 
	//根据条件删除
	num, err := Engine.Table.Delete(User{}).ByModel(NullUser{
		Name: types.NewString("tom"),
	})
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	//-------------------
	
	 
	
	//使用条件构造器
	num, err := Engine.Table.Delete(User{}).ByWhere(new(lorm.WhereBuilder).
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
	num, err := Engine.Table.Select(User{}).ByPrimaryKey(id).ScanOne(&user)
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	
	fmt.Println(user)
	//-----------------
	
	users := make([]User,0)
	num, err := Engine.Table.Select(User{}).ByPrimaryKey(id1,id2,id3).ScanList(&users)
	if err != nil {
		return err
	}
	// num=1
	fmt.Println(num)
	
	fmt.Println(user)
	//-----------------
	
	
	users := make([]User, 0)
	num, err := Engine.Table.Select(User{}).ByModel(NullUser{
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
	num, err := Engine.Table.Select(User{}).ByModel(NullUser{
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
	
	
	has, err := Engine.Table.Select(User{}).ByModel(NullUser{
		Name: types.NewString("tom"),
		Age:  types.NewInt(12),
	})
	if err != nil {
		return err
	}
	// has 查询是否存在数据
	fmt.Println(num)
	
	
	
	//----------------------------
	has, err := Engine.Table.Has(User{}).ByWhere(new(lorm.WhereBuilder).
		Eq("id", user.ID, true).
		NoLike("age", *user.Name, user.Name != nil).
		Ne("age", user.Age, false))
	if err != nil {
		return err
	}
	// has 查询是否存在数据
	fmt.Println(has)

	
	
	
	has, err := Engine.Table.Has(User{}).ByWhere(new(lorm.WhereBuilder).
		Eq("id", user.ID, true).
		NoLike("age", *user.Name, user.Name != nil).
		Ne("age", user.Age, false))
	if err != nil {
		return err
	}
	// num 查询是否存在数据
	fmt.Println(has)

	
	
	
```