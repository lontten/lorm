```js

user, err
:
= db.select(user).byId(id).getOne(&user)
user, err
:
= db.select(user).byId(id).getFirst(&user)
users, err
:
= db.select(user).byId(...id).getList(&users)


db.select("t_user").OrmSelect
:
by;
scanOne, scanFirst, scanList

byModel(user)
byMap(map)
    .isNull("name")
byPrimaryKeys(interface
{
}
)
filterPrimaryKeys(interface
{
}
)
byWhere( * whereBuider
)

.
scanOne[User]() //返回受影响的行数  dto,err
    .scanFirst[User]() //返回受影响的行数  dto,err
    .scanList[User]() //返回受影响的行数  []dto,err

    .sql()  //获取sql


//根据id查找，若没有就添加user
db.getOrInsert(user)
    .byId(id)
    .byModel(model)
    .byWhere(where)

    .scanOne() //返回受影响的行数  dto,err


db.has("t_user").OrmSelect
:
by;
scanOne, scanFirst, scanList
byModel(user)
byMap(map)
    .isNull("name")
byPrimaryKeys(interface
{
}
)
filterPrimaryKeys(interface
{
}
)
byWhere( * whereBuider
)

.
exec() //返回   has true,err


//根据id查找，若没有就添加user
db.count(user)
    .byId(id)
    .byModel(model)
    .byWhere(where)

    .exec() //返回受影响的行数  num,err


```