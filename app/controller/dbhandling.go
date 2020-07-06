package controller

import (
	"strings"
)

func ParseCategory(catSubcat string) (string, string) {
	s := strings.Split(catSubcat, "-")
	return s[0], s[1]
}
