```js

&引用，返回插入的数据的所有null值得字段

//只插入一个数据
db.insert(user)
//插入一个数据，并返回数据的所有null值得字段
db.insert(&user)

//通过遍历插入多个数据，并返回插入的数据的所有null值得字段
db.insert(users)

//通过prepare添加，并返回插入的数据的所有null值得字段
db.insertFast(users)

# 上面后续方法： Exec Err 两种

通过遍历
//如果存在则更新，不存在则插入
db.insertOrUpdate(user).byId(id)
db.insertOrUpdate(user).byUni(...string)

db.insertOrUpdate(user).byModel(model)
db.insertOrUpdate(user).byWhere(where)

//如果存在则更新，不存在则插入，
db.insertOrUpdate(&user)

# 上面后续方法：byid byUni byModel byWhere

通过遍历，list


```