```js

num,err:=LnUpdate(db,user,byId(id))
num,err:=LnUpdate(db,user,byModel(model))
num,err:=LnUpdate(db,user,where(where))



num,err:=LnUpdate(db,tableName("")
    .set()
    .set()
    .set()
    .by()
    .by()
    .by()
)


UpdateBuilder(db,User{}).
   .setTableName("user")
    setModel(user).                             
    setMap(map).
    setNull("name").
    setGreeterSelf("num",-1).
    setNowDate("create_time").
    setNowTime("create_time").
    setNowDateTime("create_time").
  
    byModel(user)                              
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