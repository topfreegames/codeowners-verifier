package verifier

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"
	"github.com/topfreegames/codeowners-verifier/pkg/providers"
)

type CodeOwner struct {
	Path   string
	Regex  *regexp.Regexp
	Line   int
	Owners []string
	Negate bool
}

func reverseSlice(a []*CodeOwner) []*CodeOwner {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
	return a
}

const commentChars = "[#]"

func stripComment(source string) string {
	if cut := strings.IndexAny(source, commentChars); cut >= 0 {
		return strings.TrimRightFunc(source[:cut], unicode.IsSpace)
	}
	return source
}

func difference(slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}

func ReadCodeownersFile(filename string) ([]*CodeOwner, error) {
	var codeowners []*CodeOwner
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := stripComment(scanner.Text())
		if line != "" {
			elements := strings.Split(line, " ")
			if len(elements) < 2 {
				return nil, fmt.Errorf("Invalid CODEOWNERS entry: %d", lineNumber)
			}
			regex, negateRegex := getPatternFromLine(elements[0])
			if regex != nil {
				c := &CodeOwner{
					Path:   elements[0],
					Regex:  regex,
					Line:   lineNumber,
					Owners: elements[1:],
					Negate: negateRegex,
				}
				codeowners = append(codeowners, c)
			}
		}
		lineNumber++
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error parsing CODEOWNERS file: %s", err)
	}
	return codeowners, nil
}

func ValidateCodeownerFile(p providers.Provider, filename string) (bool, error) {
	codeowners, err := ReadCodeownersFile(filename)
	if err != nil {
		return false, err
	}
	valid := true
	for _, c := range codeowners {
		if _, err := os.Stat(c.Path); os.IsNotExist(err) {
			log.Errorf("Error parsing line %d, path %s does not exist", c.Line, c.Path)
		}
		for _, element := range c.Owners {
			owner := strings.Replace(element, "@", "", 1)
			exists, err := p.UserExists(owner)
			if err != nil || !exists {
				exists, err = p.GroupExists(owner)
				if err != nil || !exists {
					valid = false
					log.Errorf("Error parsing line %d: user/group %s is invalid", c.Line, element)
				}
			}
		}
	}
	return valid, nil
}

// getPatternFromLine converts a line to a CODEOWNERS entry
// This is roughly addapted from https://github.com/sabhiram/go-gitignore
func getPatternFromLine(lineNumber string) (*regexp.Regexp, bool) {
	// Trim OS-specific carriage returns.
	lineNumber = strings.TrimRight(lineNumber, "\r")

	// TODO: Handle [Rule 4] which negates the match for patterns leading with "!"
	negatePattern := false
	if lineNumber[0] == '!' {
		negatePattern = true
		lineNumber = lineNumber[1:]
	}

	// Handle [Rule 2, 4], when # or ! is escaped with a \
	// Handle [Rule 4] once we tag negatePattern, strip the leading ! char
	if regexp.MustCompile(`^(\#|\!)`).MatchString(lineNumber) {
		lineNumber = lineNumber[1:]
	}

	// If we encounter a foo/*.blah in a folder, prepend the / char
	if regexp.MustCompile(`([^\/+])/.*\*\.`).MatchString(lineNumber) && lineNumber[0] != '/' {
		lineNumber = "/" + lineNumber
	}

	// Handle escaping the "." char
	lineNumber = regexp.MustCompile(`\.`).ReplaceAllString(lineNumber, `\.`)

	magicStar := "#$~"

	// Handle "/**/" usage
	if strings.HasPrefix(lineNumber, "/**/") {
		lineNumber = lineNumber[1:]
	}
	lineNumber = regexp.MustCompile(`/\*\*/`).ReplaceAllString(lineNumber, `(/|/.+/)`)
	lineNumber = regexp.MustCompile(`\*\*/`).ReplaceAllString(lineNumber, `(|.`+magicStar+`/)`)
	lineNumber = regexp.MustCompile(`/\*\*`).ReplaceAllString(lineNumber, `(|/.`+magicStar+`)`)

	// Handle escaping the "*" char
	lineNumber = regexp.MustCompile(`\\\*`).ReplaceAllString(lineNumber, `\`+magicStar)
	lineNumber = regexp.MustCompile(`\*`).ReplaceAllString(lineNumber, `([^/]*)`)

	// Handle escaping the "?" char
	lineNumber = strings.Replace(lineNumber, "?", `\?`, -1)

	lineNumber = strings.Replace(lineNumber, magicStar, "*", -1)

	// Temporary regex
	var expr = ""
	if strings.HasSuffix(lineNumber, "/") {
		expr = lineNumber + "(|.*)$"
	} else {
		expr = lineNumber + "(|/.*)$"
	}
	if strings.HasPrefix(expr, "/") {
		expr = "^(|/)" + expr[1:]
	} else {
		expr = "^(|.*/)" + expr
	}
	pattern, _ := regexp.Compile(expr)

	return pattern, negatePattern
}

// MatchesPath returns true if the given GitIgnore structure would target
// a given path string `f`.
func (co *CodeOwner) MatchesPath(f string) bool {
	// Replace OS-specific path separator.
	f = strings.Replace(f, string(os.PathSeparator), "/", -1)

	matchesPath := false
	if co.Regex.MatchString(f) {
		// If this is a regular target (not negated with a gitignore exclude "!" etc)
		if !co.Negate {
			matchesPath = true
		} else if matchesPath {
			// Negated pattern, and matchesPath is already set
			matchesPath = false
		}
	}
	return matchesPath
}

func CheckCodeowner(codeowners []*CodeOwner, filename string, ignore []string) (*CodeOwner, bool) {
	for _, c := range reverseSlice(codeowners) {
		match := c.MatchesPath(filename)
		if match {
			owners := difference(c.Owners, ignore)
			return c, len(owners) > 0
		}
	}
	return nil, false
}
