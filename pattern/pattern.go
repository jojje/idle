package pattern

import (
	"regexp"
	"strings"
)

type Matcher func(s string) bool

func NewMatcher(s string, caseInsensitive bool) (Matcher, error) {
	start, end := 0, len(s)
	useRE := false

	if len(s) > 2 && s[0] == '/' && s[len(s)-1] == '/' {
		start, end = 1, len(s)-1
		useRE = true
	}

	expr := s[start:end]

	// exact string matching
	if !useRE {
		if caseInsensitive {
			expr = strings.ToLower(expr)
			return func(s string) bool {
				return expr == strings.ToLower(s)
			}, nil
		}
		return func(s string) bool {
			return expr == s
		}, nil
	}

	// regular expressions
	if caseInsensitive {
		expr = "(?i)" + expr
	}

	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return func(s string) bool {
		return re.MatchString(s)
	}, nil
}
