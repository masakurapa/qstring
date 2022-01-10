package qstring

import "sort"

type urlValueMap map[string]urlValue

func (vm urlValueMap) firstValue() []string {
	var val []string
	for _, v := range vm {
		val = v.values
		break
	}
	return val
}

func (vm urlValueMap) sortedChild() []urlValue {
	uvs := make([]urlValue, 0, len(vm))
	for _, uv := range vm {
		uvs = append(uvs, uv)
	}
	sort.Slice(uvs, func(i, j int) bool {
		return uvs[i].key < uvs[j].key
	})
	return uvs
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
