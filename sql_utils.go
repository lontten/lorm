package lorm

import "bytes"

func genWhere(c []string) []byte {
	if len(c) == 0 {
        return []byte("")
    }
    var buf bytes.Buffer
    buf.WriteString(" WHERE ")
    buf.WriteString(c[0])
    for i := 1; i < len(c); i++ {
        buf.WriteString(" AND ")
        buf.WriteString(c[i])
    }
    return buf.Bytes()

}
