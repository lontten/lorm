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


通过遍历
//如果存在则更新，不存在则插入
db.set(user).byId(id)
db.set(user).byUni(...string)

db.set(user).byModel(model)
db.set(user).byWhere(where)

//如果存在则更新，不存在则插入，
db.set(&user)



通过遍历，list
//根据id查找，若没有就添加user
db.getOrInsert(user).byId(id)

db.getOrInsert(user).byModel(model)

db.getOrInsert(user).byWhere(where)





```