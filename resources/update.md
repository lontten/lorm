```js



db.update(user).byId(id)
db.update(user).byModel(model)
db.update(user).byWhere(where)

db.update("t_user").                            update-in
    setModel(user).                             
    setMap(map).
    setNull("name").
    setCurrentDate("create_time").
    setGreeterSelf("num",-1).
    setCurrentTime("create_time").
    setCurrentDateTime("create_time").
    byModel(user)                                update-by-in  
    byMap(map)
    byPrimaryKeys(interface{})
    filterPrimaryKeys(interface{})
    byWhere(*whereBuider)
    .exec() //返回受影响的行数 err
    .execNum()  //获取 num,err



type UdateByIn interface {
    Num() (int64, error)
    Err() error
    by....
}

func aaa() {
    num, err := h().Num()
    fmt.Println(num, err)
}

```