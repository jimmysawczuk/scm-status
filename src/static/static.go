package static

import (
	"bytes"
	"strings"
)

func Version() string {
	ver, err := getVersion()
	if err == nil {
		tags_list := bytes.NewBuffer(ver).String()
		tags_list = strings.Replace(tags_list, `\r\n`, `\n`, 0)
		tags := strings.Split(tags_list, `\n`)

		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if strings.HasPrefix(tag, "v") {
				return tag
			}
		}
	}

	return ""
}
