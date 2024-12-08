## 泛型设计
泛型方法？想都不要想，你永远别想在golang里玩你那花里胡哨的链式调用！
在理解了golang的设计哲学后，对lorm进行了重构，使其更加完美。
```js

user, err:= One[User](db,byId(id))
user, err:= First[User](db,byId(id...))
users, err:= List[User](db,byModel(model)
    .byMap(map)
    .byPrimaryKeys(id...)
    .byWhere(where)

    new(Extra)
        .showsSql()
        .skipSoftDelete()
        .select("id","name"..)
        


).orderBy("id desc")
    .limit(10)
.offset(10)


First
Has
Count
list,err := List[User](db,model)
    .byModel(model)
    .byMap(map)
    .byWhere(where)
    .orderby("id desc")
    .limit(10)
    .offset(10)
.showsSql()
.skipSoftDelete()
.select("id","name"..)
.exec()


user, err:= GetOrInsert[User](db,&user,set().setE)
    byModel(model)
    .byMap(map)
    .byPrimaryKeys(id...)
    .byWhere(where)

    .showsSql()
    .skipSoftDelete()
    .select("id","name"..)
.query()


```