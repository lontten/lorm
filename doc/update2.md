```js

num,err:=LnUpdate(db,user,byId(id))
num,err:=LnUpdate(db,user,byModel(model))
num,err:=LnUpdate(db,user,where(where))



num,err:=LnUpdate(db,tableName("")
    .set()
    .set()
    .set()
    .by()
    .by()
    .by()
)


//通过遍历插入多个数据，并返回插入的数据的所有null值得字段
num,err:=Update(db,users,new(ldb.Extra).
    .shwoSql()  // 打印sql
    .skipLgDel() //跳过逻辑删除字段
    .returnLevel(Field.All,Field.Nil,Field.Pk,Field.None) //返回所有字段，；只返回nil字段；只返回主键字段
    .setTable("")  //覆盖表名

    .setNUll("")
    .setkv("","")


)
    .by()
    .by()
    .by()
    .exec()



```