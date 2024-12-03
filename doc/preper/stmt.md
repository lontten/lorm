```js
db.query("")

db.exec("")



sql,err:=db.insert(user)

sql,err:=db.del(user).byId()

sql,err:=db.del(user).byWhere(where)



sql,err:=db.update(user).by



db.query("").where(where.byId().byMode().byWhere)
UserField{
    files []string
}
UserField{}.Name().Age().Id().Fileds()

fun (u *UserField)Name() *UserField{
    u.fields = append(u.fields, "name")
}
```