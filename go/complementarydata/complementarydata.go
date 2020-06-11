package complementarydata

type ComplementaryData []interface{}

func (b *ComplementaryData) Push(el interface{}) {
	*b = append(*b, el)
}
