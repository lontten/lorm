# lorm

[![Build Status](https://travis-ci.org/jmoiron/sqlx.svg?branch=master)](https://travis-ci.org/jmoiron/sqlx) [![Coverage Status](https://coveralls.io/repos/github/jmoiron/sqlx/badge.svg?branch=master)](https://coveralls.io/github/jmoiron/sqlx?branch=master) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/jmoiron/sqlx) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/jmoiron/sqlx/master/LICENSE)

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
db.update("t_user").byModel(user)
map=Map{}
db.update("t_user").byModel(map)

db.update("t_user").byWhere(*whereBuider)

db.update(user).setNull("student_name",user.Name==nil)


db.update(user)
db.update().set("name",age)



```
