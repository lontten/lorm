```js

user,err:=db.select(user).byId(id).getOne<User>()
user,err:=db.select(user).byId(id).getFirst<User>()
users,err:=db.select(user).byId(...id).getList<User>()




pages,err:=db.select(user).byMode(model).getPage<User>(page,size)






```