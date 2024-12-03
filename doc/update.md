```js

db.update(user).byId(id)
db.update(user).byModel(model)
db.update(user).byWhere(where)

db.update("t_user").                            update-in
   
    setModel(user).                             
    setMap(map).
    setNull("name").
    setGreeterSelf("num",-1).
    setNowDate("create_time").
    setNowTime("create_time").
    setNowDateTime("create_time").
  
    byModel(user)                                update-by-in  
    byMap(map)
    byPrimaryKeys(any)
    filterPrimaryKeys(any)
    byWhere(*whereBuider)
   
    .exec() //返回受影响的行数 num, err


func aaa() {
    num, err := h().exec()
    fmt.Println(num, err)
}

```