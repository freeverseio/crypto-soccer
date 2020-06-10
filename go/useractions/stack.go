package useractions

type Stack []interface{}

func (b *Stack) Push(el interface{}) {
	*b = append(*b, el)
}
