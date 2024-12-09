package main

import (
	"fmt"
	"github.com/lontten/lorm"
	"test/ldb"
)

func Build1() {
	var u User
	num, err := lorm.SelectBuilder(ldb.DB).Native(`
select count(*)
from xjwy_enquiry_order o
         join xjwy_china_area a on a.id = o.order_area_id
         join xjwy_organization org on org.organization_id = o.organization_id
where (o.order_title like ?
    or o.project_name like ?
    or org.organization_name like ?
    )
  and o.enquiry_order_id in (select distinct oi.enquiry_order_id
                             from xjwy_enquiry_inquirer oi
                             where oi.inquirer_id in (1, 2))
  and o.order_status in (1, 2)
`).
		WhereIng().Args(1, 2, 3).
		Where("id = 1").
		WhereIn(`
   o.enquiry_order_id in (select distinct oi.enquiry_order_id
                             from xjwy_enquiry_inquirer oi
                             where oi.inquirer_id in ? )
`, 1, 2).
		Where("id = 2").
		Limit(2).
		ScanOne(&u)
	fmt.Println(num, err)
	fmt.Println(*u.Id)
	fmt.Println(*u.Name)
}

func Build2() {
	var u User
	num, err := lorm.SelectBuilder(ldb.DB).
		Select("id").Select("name").
		From("t_user").Where("id = ?", 1 == 2).Arg(1, 1 == 2).
		Where("id = 2").Limit(2).
		ScanOne(&u)
	fmt.Println(num, err)
	fmt.Println(*u.Id)
	fmt.Println(*u.Name)
}
