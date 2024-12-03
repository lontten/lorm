```js



db.update(user).byId(id)
db.update(user).byModel(model)
db.update(user).byWhere(where)

db.builer().                            update-in
    setModel(user).                             
    setMap(map).
    setNull("name").
    setCurrentDate("create_time").
    setGreeterSelf("num",-1).
    setCurrentTime("create_time").
    setCurrentDateTime("create_time").
    byModel(user)                                update-by-in  
    byMap(map)
    byPrimaryKeys(any)
    filterPrimaryKeys(any)
    byWhere(*whereBuider)
    .err() //返回受影响的行数 err
    .num()  //获取 num,err



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