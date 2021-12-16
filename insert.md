```js

//只添加
db.insert(user)
//添加并返回id
db.insert(&user)
//通过遍历添加
db.insert(users)
//通过prepare添加
db.insertFast(users)

//插入或更新
db.insertOrUpdate(user)

//根据id查找，若没有就添加user
db.getOrInsert(user).byId(id)


db.delete(user).byId(id)


db.update(user).getId(id)




```