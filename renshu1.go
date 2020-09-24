

type View interface {
	DataSet() string
	Name() string
	Query() string
}

type Viewview interface {
	Renshu(view View)
}

func main() {
	v := View{"","",""}

}