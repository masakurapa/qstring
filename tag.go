package qstringer

import (
	"reflect"
	"strings"
)

const (
	tagName      = "qstringer"
	optSeparator = ","
	omitempty    = "omitempty"
)

type tagOption struct {
	omitempty bool
}

func parseTag(tag reflect.StructTag) (string, tagOption) {
	s := tag.Get(tagName)
	if idx := strings.Index(s, optSeparator); idx != -1 {
		return s[:idx], tagOption{
			omitempty: strings.Contains(s[idx:], omitempty),
		}
	}
	return s, tagOption{}
}
