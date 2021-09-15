# lorm
```javascript
db.update("t_user").byId(1)
db.update("t_user").byId([1,2])

user=User{}
db.update("t_user").byModel(user)
db.update("t_user").byWhere(*whereBuider)

db.update(user).setNull("student_name",user.Name==nil)






```
