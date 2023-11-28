package main

import "github.com/lontten/lorm"

func main2() {
	kk := Kk{KkName: "a"}

	db, err := lorm.Connect(nil, nil)
	if err != nil {
		panic(err)
	}
	query := db.Query("select * from t_kk")

	num, err := query.ScanList(&kk)
	println(num, err)

}
