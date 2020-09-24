package main

import (
	"fmt")

type viewView struct {
	dataSet string
	name string
	query string
}

type View interface {
	DataSet() string
	Name() string
	Query() string
}

func (b viewView) DataSet() string {
	fmt.Println(b.dataSet)
	return b.dataSet
}

func (b viewView) Name() string {
	fmt.Println(b.name)
	return b.name
}

func (b viewView) Query() string {
	fmt.Println(b.query)
	return b.query
}

func vvv (v View) {
	v.DataSet()
	v.Name()
	v.Query()
	// fmt.Sprintf("%s, %s, %s", v.DataSet(), v.Name(), v.Query())
}

func (b viewView) Name() string {
	fmt.Println("aiueoooo")
	return b.name
}

func main() {
	v := viewView{"A","B","C"}
	var oya_view View
	oya_view = v
	vvv(oya_view)
	
}
