```js

db.delete("t_user").                      OrmDel  == by ;Exec(返回行数+err);  Err(返回err)
   
    byModel(user)
    byMap(map)
    byPrimaryKeys(interface{})
    filterPrimaryKeys(interface{})
    byWhere(*whereBuider)
    .err() //返回受影响的行数
    .num()  //获取sql


```