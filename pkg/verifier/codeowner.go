package verifier

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"
	"github.com/topfreegames/codeowners-verifier/pkg/providers"
)

// CodeOwner represents a line in a CODEOWNERS file
type CodeOwner struct {
	Path   string
	Regex  *regexp.Regexp
	Line   int
	Owners []string
	Negate bool
}

// reverseCodeOwners returns an inverted slice
func reverseCodeOwners(a []*CodeOwner) []*CodeOwner {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
	return a
}

// Used to remove comment lines
const commentChars = "^[#]"

// stripComment uses the commentChars to remove comments from lines
func stripComment(source string) string {
	if cut := strings.IndexAny(source, commentChars); cut >= 0 {
		return strings.TrimRightFunc(source[:cut], unicode.IsSpace)
	}
	return source
}

// hasdifference returns true if there is an element on slice1 that isn't on slice2
func hasDifference(slice1 []string, slice2 []string) bool {
	for _, s1Val := range slice1 {
		found := false
		for _, s2Val := range slice2 {
			if s1Val == s2Val {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}
	return false
}

// ReadCodeownersFile reads the file specified by filename
// and returns a list of CodeOwners strucs, as well as an error
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
		line := strings.Fields(stripComment(scanner.Text()))
		if len(line) == 1 {
			return nil, fmt.Errorf("Invalid CODEOWNERS entry: %d", lineNumber)
		} else if len(line) >= 2 {
			regex, negateRegex := getPatternFromLine(line[0])
			if regex != nil {
				c := &CodeOwner{
					Path:   line[0],
					Regex:  regex,
					Line:   lineNumber,
					Owners: line[1:],
					Negate: negateRegex,
				}
				codeowners = append(codeowners, c)
			}
		}
		lineNumber++
	}
	return codeowners, nil
}

// ValidateCodeownerFile check if every entry:
// 1. Has a valid file/path
// 2. Check if every owner is an user or a group.
func ValidateCodeownerFile(p providers.Provider, filename string) (bool, error) {
        var validEntriesCache map[string]int = make(map[string]int)
	codeowners, err := ReadCodeownersFile(filename)
	if err != nil {
		return false, err
	}
	valid := true
        currentDir, _ := os.Getwd()
        files, _ := FilePathWalkDir(currentDir)
	for _, c := range codeowners {
		fileMatches := false
		for idx := 0; idx < len(files) && !fileMatches; idx++ {
                        file := files[idx]
			fileMatches = c.MatchesPath(file)
		}
		if !fileMatches {
			log.Errorf("Error parsing line %d, path %s does not exist", c.Line, c.Path)
			valid = false
		}
		for _, element := range c.Owners {
			owner := strings.Replace(element, "@", "", 1)
			_, ok := validEntriesCache[owner]
			if ok {
				continue
			}
			exists, err := p.UserExists(owner)
			if err != nil {
				return false, err
			} else {
				if !exists {
					exists, err = p.GroupExists(owner)
					if err != nil {
						return false, err
					} else {
						if !exists {
							valid = false
							log.Errorf("Error parsing line %d: user/group %s is invalid", c.Line, element)
						} else {
							validEntriesCache[owner] = 2
						}
					}
				} else {
					validEntriesCache[owner] = 1
				}
			}
		}
	}
	return valid, nil
}

// getPatternFromLine converts a line to a CODEOWNERS entry
// This is roughly adapted from https://github.com/sabhiram/go-gitignore
func getPatternFromLine(line string) (*regexp.Regexp, bool) {
	// Trim OS-specific carriage returns.
	line = strings.TrimRight(line, "\r")

	// TODO: Handle [Rule 4] which negates the match for patterns leading with "!"
	negatePattern := false
	if line[0] == '!' {
		negatePattern = true
		line = line[1:]
	}

	// If we encounter a foo/*.blah in a folder, prepend the / char
	if regexp.MustCompile(`([^\/+])/.*\*\.`).MatchString(line) && line[0] != '/' {
		line = "/" + line
	}

	// Handle escaping the "." char
	line = regexp.MustCompile(`\.`).ReplaceAllString(line, `\.`)

	magicStar := "#$~"

	// Handle "/**/" usage
	if strings.HasPrefix(line, "/**/") {
		line = line[1:]
	}
	line = regexp.MustCompile(`/\*\*/`).ReplaceAllString(line, `(/|/.+/)`)
	line = regexp.MustCompile(`\*\*/`).ReplaceAllString(line, `(|.`+magicStar+`/)`)
	line = regexp.MustCompile(`/\*\*`).ReplaceAllString(line, `(|/.`+magicStar+`)`)

	// Handle escaping the "*" char
	line = regexp.MustCompile(`\\\*`).ReplaceAllString(line, `\`+magicStar)
	line = regexp.MustCompile(`\*`).ReplaceAllString(line, `([^/]*)`)

	// Handle escaping the "?" char
	line = strings.Replace(line, "?", `\?`, -1)

	line = strings.Replace(line, magicStar, "*", -1)

	// Temporary regex
	var expr = ""
	if strings.HasSuffix(line, "/") {
		expr = line + "(|.*)$"
	} else {
		expr = line + "(|/.*)$"
	}
	if strings.HasPrefix(expr, "/") {
		expr = "^(|/)" + expr[1:]
	} else {
		expr = "^(|.*/)" + expr
	}
	pattern, _ := regexp.Compile(expr)

	return pattern, negatePattern
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			// remove current dir from filepath to use regex properly.
			path := strings.Replace(path, root, "", 1)
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, err
}

// MatchesPath returns true if the given GitIgnore structure would target
// a given path string `f`.
func (co *CodeOwner) MatchesPath(f string) bool {
	// Replace OS-specific path separator if it is not "/".
	if string(os.PathSeparator) != "/" {
		f = strings.Replace(f, string(os.PathSeparator), "/", -1)
	}

	matchesPath := false
	if co.Regex.MatchString(f) {
		// If this is a regular target (not negated with a gitignore exclude "!" etc)
		if !co.Negate {
			matchesPath = true
		}
	}
	return matchesPath
}

// VerifyCodeowner check if a line matches any entry on the reversed list of CodeOwners
func VerifyCodeowner(codeowners []*CodeOwner, filename string, ignore []string) (*CodeOwner, bool) {
	for _, c := range reverseCodeOwners(codeowners) {
		match := c.MatchesPath(filename)
		if match {
			return c, hasDifference(c.Owners, ignore)
		}
	}
	return &CodeOwner{}, false
}
