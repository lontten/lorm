package lorm

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
