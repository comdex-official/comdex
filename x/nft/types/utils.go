package types

import (
	"github.com/google/uuid"
	"regexp"
	"strings"
)

var (
	IsAlphaNumeric   = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	IsBeginWithAlpha = regexp.MustCompile(`^[a-zA-Z].*`).MatchString
	IsAlpha          = regexp.MustCompile(`^[a-zA-Z]+`).MatchString
)

func GenUniqueID(prefix string) string {
	return prefix + strings.ReplaceAll(uuid.New().String(), "-", "")
}
