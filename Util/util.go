package Util

import "strings"

func TrimArticleText(text string) string {
	art := strings.TrimSpace(text)
	txt := strings.Split(art, "\n")

	return string(txt[0]) + string(" - ") + string(txt[len(txt)-1])
}
