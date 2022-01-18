package lorm

type Page struct {
	size         int64
	current      int64
	selectTokens []string
}

func (db DB) Page(size, current int64) Page {
	return Page{size: size, current: current}
}

func (p Page) Select(name string, condition ...bool) OrmPageSelect {
	for _, b := range condition {
		if !b {
			return OrmPageSelect{base: p}
		}
	}
	p.selectTokens = append(p.selectTokens, name)
	return OrmPageSelect{base: p}
}

func (p OrmPageSelect) Form(query string, args ...interface{}) {

}

type OrmPageSelect struct {
	base Page
}
