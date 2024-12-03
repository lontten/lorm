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
)


user, err:= GetOrInsert[User](db,&user)
    byModel(model)
    .byMap(map)
    .byPrimaryKeys(id...)
    .byWhere(where)
)

has,err:=Has[User](db,User{})
    byModel(model)
    .byMap(map)
    .byPrimaryKeys(id...)
    .byWhere(where)
)

num,err:=Count[User](db,User{})
    byModel(model)
    .byMap(map)
    .byPrimaryKeys(id...)
    .byWhere(where)
)


```