# lorm

[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/lontten/lorm/main/LICENSE)

target - 必须是struct，用于tableName

model - struct-map，用于 ，生成where

scan ptr-* slice-* map


```javascript
db.update("t_user").byId(1)
db.update("t_user").byId([1,2])

//num是int类型，
db.update("t_user").byWhere(
    num.in([1,2,3])
)

//num是 []int 类型, []int 必须封装成 struct,
//要有 valuer
db.update("t_user").byWhere(
    num.in([
    [1,2] AS Array,
    [1,2] AS Array,
    [1,2] AS Array,
])
)


// map 和struct等价,filed中的slice，自动转成
//数组型字段。
user=User{}


// setModel/setMap nil会被排除，setnull，若前面没有，则添加 set null
db.update("t_user").
    setModel(user).
    setMap(map).
    setNull("name").
    setCurrentDate("create_time").
    setGreeterSelf("num",-1).
    setCurrentTime("create_time").
    setCurrentDateTime("create_time").
    byModel(user)
    byMap(map)
    byPrimaryKey(interface{})
    byWhere(*whereBuider)
    .exec() //返回受影响的行数
    .sql()  //获取sql
    .prepare() //返回*sql.Stmt

db.update(user)
db.update().set("name",age)



```
