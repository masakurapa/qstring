package qstringer

type urlValueMap map[string]urlValue

func (vm urlValueMap) firstValue() []string {
	var val []string
	for _, v := range vm {
		val = v.values
	}
	return val
}

type urlValue struct {
	key      string
	values   []string
	isString bool
	child    urlValueMap
}

func (uv urlValue) hasChild() bool {
	return uv.child != nil && len(uv.child) > 0
}

func (uv urlValue) hasSingleValue() bool {
	return len(uv.values) == 1 && !uv.hasChild()
}
