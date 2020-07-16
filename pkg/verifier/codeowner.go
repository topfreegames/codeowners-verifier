package verifier

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/topfreegames/codeowners-verifier/pkg/providers"
)

func CheckCodeowner(p providers.Provider, filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("Couldn't open file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	valid := true
	line := 1
	for scanner.Scan() {
		elements := strings.Split(scanner.Text(), " ")[1:]
		for _, element := range elements {
			owner := strings.Replace(element, "@", "", 1)
			exists, err := p.UserExists(owner)
			if err != nil || !exists {
				exists, err = p.GroupExists(owner)
				if err != nil || !exists {
					valid = false
					log.Errorf("Error parsing line %d: %s invalid", line, element)
				}
			}
		}
		line++
	}
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("Error parsing CODEOWNERS file: %s", err)
	}
	return valid, nil
}
