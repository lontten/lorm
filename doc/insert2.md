```js

&引用，返回插入的数据的所有null值得字段


//只插入一个数据
num,err:=Insert(db,user)
//插入一个数据，并返回数据的所有null值得字段
num,err:=Insert(db,&user)

//通过遍历插入多个数据，并返回插入的数据的所有null值得字段
num,err:=Insert(db,users,new(lorm.Extra).
        .shwoSql()  // 打印sql
        .skipLgDel() //跳过逻辑删除字段
        .returnLevel(Field.All,Field.Nil,Field.Pk,Field.None) //返回所有字段，；只返回nil字段；只返回主键字段
        .setTable("")  //覆盖表名

        .setNUll("")
        .setkv("","")
    
    
        .whenDuplicateKey(name ...string,)
        .DoNothing()  //INSERT IGNORE 无法返回插入数据
        .DoUpdate(new(lorm.Set)
            .set("name","tom")
            .setNull("age")
            .field("id","created_at")
            .setModel(user)  // model 会排除零值
            .setMap(map[string]any{"name":"tom","age":18})  // map类型数据不会排除 零值
        )
        .DoReplace()


        .updateWhen(byField("").byModel().byWhere())  // 独立出来比较好，放在insert，没有复用的优势 InsertOrUpdate
   )


//先用条件查询是否有记录，有则对其更新，没有则插入
has,num,err:=InsertOrUpdate(db,&user,byField("name") // 先检查 name 是否有已存在数据，如果有返回true，否则insert



//应用场景，
添加数据时，要求，名字不能重复
// 先根据条件查询是否存在，有则返回true，否则insert
has,num,err:=InsertOrHas(db,&user,new(lorm.Extra).
                .shwoSql()  // 打印sql
                .skipLgDel() //跳过逻辑删除字段
                .returnLevel(Field.All,Field.Nil,Field.Pk,Field.None) //返回所有字段，；只返回nil字段；只返回主键字段
                .setTable("")  //覆盖表名
            
                .setNUll("")
                .setkv("","")
            
            
                .whenDuplicateKey(name ...string,)
            .DoNothing()  //INSERT IGNORE 无法返回插入数据
                .DoUpdate(new(lorm.Set)
                    .set("name","tom")
                    .setNull("age")
                    .field("id","created_at")
                    .setModel(user)  // model 会排除零值
                    .setMap(map[string]any{"name":"tom","age":18})  // map类型数据不会排除 零值
            )
            .DoReplace()
).byField("name") // 先检查 name 是否有已存在数据，如果有返回true，否则insert


has,num,err:=UpdateOrHas(db,&user,).byField("name","age") // 先检查 name,age 是否有已存在数据，如果有返回true，否则update
```