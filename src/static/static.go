package static

import (
	"bytes"
	"strings"
)

func Version() string {

	tags_list := bytes.NewBuffer(getVersion()).String()
	tags_list = strings.Replace(tags_list, `\r\n`, `\n`, 0)
	tags := strings.Split(tags_list, `\n`)

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if strings.HasPrefix(tag, "v") {
			return tag
		}
	}

	return ""
}
