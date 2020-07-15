package verifier

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/topfreegames/codeowners-verifier/pkg/providers"
)

func CheckCodeowner(p providers.Provider, filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("Couldn't open file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		elements := strings.Split(scanner.Text(), " ")[1:]
		for _, element := range elements {
			owner := strings.Replace(element, "@", "", 1)
			valid, err := p.UserExists(owner)
			if err != nil || !valid {
				valid, err = p.GroupExists(owner)
				if err != nil || !valid {
					return false, nil
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("Error parsing CODEOWNERS file: %s", err)
	}
	return true, nil
}
