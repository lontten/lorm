```js



db.update(user).byId(id)
db.update(user).byModel(model)
db.update(user).byWhere(where)

db.update("t_user").
    setModel(user).
    setMap(map).
    setNull("name").
    setCurrentDate("create_time").
    setGreeterSelf("num",-1).
    setCurrentTime("create_time").
    setCurrentDateTime("create_time").
    byModel(user)
    byMap(map)
    byPrimaryKeys(interface{})
    filterPrimaryKeys(interface{})
    byWhere(*whereBuider)
    .exec() //返回受影响的行数
    .sql()  //获取sql



type K interface {
    Num() (int64, error)
    Err() error
}

func aaa() {
    num, err := h().Num()
    fmt.Println(num, err)
}

```