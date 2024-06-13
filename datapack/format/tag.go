package format

import "fmt"

var tagFormatStr = "`json:\"%s\" yaml:\"%s\"`"

func Tag(tag string) string {
	return fmt.Sprintf(tagFormatStr, tag, tag)
}
