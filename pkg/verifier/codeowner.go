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
	Line   int
	Owners []string
}

const commentChars = "[#]"

func stripComment(source string) string {
	if cut := strings.IndexAny(source, commentChars); cut >= 0 {
		return strings.TrimRightFunc(source[:cut], unicode.IsSpace)
	}
	return source
}

func reverseSlice(a []*CodeOwner) []*CodeOwner {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
	return a
}
func removeSlice(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func ReadCodeownersFile(filename string) ([]*CodeOwner, error) {
	var codeowners []*CodeOwner
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	line := 1
	for scanner.Scan() {
		codeowner := strings.Split(stripComment(scanner.Text()), " ")
		if len(codeowner) > 1 {
			c := &CodeOwner{
				Path:   codeowner[0],
				Line:   line,
				Owners: codeowner[1:],
			}
			codeowners = append(codeowners, c)
		}
		line++
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
func CheckCodeowner(codeowners []*CodeOwner, filename string, ignore []string) (*CodeOwner, bool) {
	for _, c := range reverseSlice(codeowners) {
		pathPattern := c.Path
		// Patterns ending with / should be recursive to all folders
		// Patterns like * reference the current folder/path
		pathPattern = strings.ReplaceAll(pathPattern, "*", `[^\/]*`)
		if strings.HasSuffix(pathPattern, "/") {
			pathPattern = pathPattern + ".*"
		}
		var pathRegex = regexp.MustCompile("^" + pathPattern + "$")
		match := pathRegex.MatchString(filename)
		if match {
			owners := c.Owners
			for _, i := range ignore {
				owners = removeSlice(owners, i)
			}
			if len(owners) > 0 {
				return c, true
			}
		}
	}
	return nil, false
}
