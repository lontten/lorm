

nuller 的作用，判断 字段是非为null，insert update where 时，只使用非null字段。
nuller 无法对 sql.nullXX 类型 进行检查，所以干脆不检查 nuller，
现在不用nuller，如何判断 字段 是否为null？
直接对 value进行检查，值为nil的，则为null，忽略字段。

sql.NullXX ,只能判断 zero，不能判断 nil，

确定方案。
如果字段是 指针类型，当字段为nil， 为 null状态，忽略字段。
如果不是指针类型，用零值进行判断，零值字段为 null状态，忽略字段。
