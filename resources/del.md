```js

db.delete("t_user").
   
    byModel(user)
    byMap(map)
    byPrimaryKeys(interface{})
    filterPrimaryKeys(interface{})
    byWhere(*whereBuider)
    .exec() //返回受影响的行数
    .sql()  //获取sql


```