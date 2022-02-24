```js

user,err:=db.select(user).byId(id).getOne<User>()
user,err:=db.select(user).byId(id).getFirst<User>()
users,err:=db.select(user).byId(...id).getList<User>()




db.select("t_user").                                                    sele-in
 
    byModel(user)
    byMap(map)
        .isNull("name")
    byPrimaryKeys(interface{})
    filterPrimaryKeys(interface{})
    byWhere(*whereBuider)
    .scanOne() //返回受影响的行数  num,dto,err
    .scanFirst() //返回受影响的行数  num,dto,err
    .scanList() //返回受影响的行数  num,dto,err
    .sql()  //获取sql







```