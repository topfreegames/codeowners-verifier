package verifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Expected interface{}
	Sample   interface{}
	Name     string
}

type ReturnWithError struct {
	Value interface{}
	Error bool
}

func TestReverseCodeOwners(t *testing.T) {
	codeOwnerEntries := []CodeOwner{
		{
			Path:   "*",
			Regex:  nil,
			Line:   1,
			Owners: []string{"@defaultGroup"},
			Negate: false,
		},
		{
			Path:   "path/to/file",
			Regex:  nil,
			Line:   2,
			Owners: []string{"@group1 @group2"},
			Negate: false,
		},
	}
	reversedCodeOwnerEntries := []CodeOwner{
		{
			Path:   "path/to/file",
			Regex:  nil,
			Line:   2,
			Owners: []string{"@group1 @group2"},
			Negate: false,
		},
		{
			Path:   "*",
			Regex:  nil,
			Line:   1,
			Owners: []string{"@defaultGroup"},
			Negate: false,
		},
	}
	emptyCodeOwner := []CodeOwner{}
	tests := []TestCase{
		{
			Name:     "Reversing CODEOWNERS file",
			Sample:   codeOwnerEntries,
			Expected: reversedCodeOwnerEntries,
		},
		{
			Name:     "Checking empty CODEOWNERS file",
			Sample:   emptyCodeOwner,
			Expected: emptyCodeOwner,
		},
	}

	for i, test := range tests {
		t.Logf("Test case %d: %s", i, test.Name)
		result := reverseCodeOwners(test.Sample.([]*CodeOwner))
		assert.Equal(t, test.Expected, result)
	}
}

func TestCodeOwnerReadFile(t *testing.T) {
	tests := []TestCase{
		{
			Name:   "invalid file",
			Sample: "invalid-file",
			Expected: ReturnWithError{
				Value: nil,
				Error: true,
			},
		},
	}

	for i, test := range tests {
		t.Logf("Test case %d: %s", i, test.Name)

		expected := test.Expected.(ReturnWithError)
		sample := test.Sample.(string)

		val, err := ReadCodeownersFile(sample)

		if expected.Error {
			assert.Error(t, err, "should return an error")
			assert.Nil(t, val, "return should be nil on error")
		} else {
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, expected.Value.([]map[string]interface{}), val, "decoded value should match expected")
		}
	}
}
