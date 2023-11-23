```js

user,err:=db.select(user).byId(id).getOne<User>()
user,err:=db.select(user).byId(id).getFirst<User>()
users,err:=db.select(user).byId(...id).getList<User>()




db.select("t_user").                             OrmSelect: by;scanOne,scanFirst,scanList
 
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





//根据id查找，若没有就添加user
db.getOrInsert(user).byId(id)

db.getOrInsert(user).byModel(model)

db.getOrInsert(user).byWhere(where)






```