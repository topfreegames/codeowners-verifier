package verifier

import (
	"fmt"
	"regexp"
	"testing"

	filet "github.com/Flaque/filet"
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
	codeOwnerEntries := []*CodeOwner{
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
	reversedCodeOwnerEntries := []*CodeOwner{
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
	emptyCodeOwner := []*CodeOwner{}
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
		assert.Equal(t, test.Expected.([]*CodeOwner), result)
	}
}

func TestStripComment(t *testing.T) {
	tests := []TestCase{
		{
			Name:     "Checking line with #",
			Sample:   "#testing",
			Expected: "",
		},
		{
			Name:     "Checking line section []",
			Sample:   "[SectionThatShouldBeSanitized]",
			Expected: "",
		},
		{
			Name:     "Checking valid line",
			Sample:   "* @test",
			Expected: "* @test",
		},
	}

	for i, test := range tests {
		t.Logf("Test case %d: %s", i, test.Name)
		result := stripComment(test.Sample.(string))
		assert.Equal(t, test.Expected.(string), result)
	}
}

func TestDifference(t *testing.T) {
	tests := []TestCase{
		{
			Name:     "Checking 2 identical slices",
			Sample:   []string{"a", "b", "c", "d"},
			Expected: "",
		},
	}
	fmt.Print(tests)
}

func TestCodeOwnerReadFile(t *testing.T) {
	tests := []TestCase{
		{
			Name: "invalid file",
			Sample: map[string]interface{}{
				"Filename": "non-existent-file",
				"Contents": "",
			},
			Expected: ReturnWithError{
				Value: nil,
				Error: true,
			},
		},
		{
			Name: "valid codeowners",
			Sample: map[string]interface{}{
				"Filename": "valid-codeowners",
				"Contents": `* @user1 @user2
folder1 @group1
folder2/ @group1
folder2/* @group2
!file1 @user3
folder1/*.tf @user4
/**/ @group1`,
			},
			Expected: ReturnWithError{
				Value: []*CodeOwner{
					{
						Path:   "*",
						Regex:  regexp.MustCompile("^(|.*/)([^/]*)(|/.*)$"),
						Negate: false,
						Owners: []string{
							"@user1",
							"@user2",
						},
						Line: 1,
					},
					{
						Path:   "folder1",
						Regex:  regexp.MustCompile("^(|.*/)folder1(|/.*)$"),
						Negate: false,
						Owners: []string{
							"@group1",
						},
						Line: 2,
					},
					{
						Path:   "folder2/",
						Regex:  regexp.MustCompile("^(|.*/)folder2/(|.*)$"),
						Negate: false,
						Owners: []string{
							"@group1",
						},
						Line: 3,
					},
					{
						Path:   "folder2/*",
						Regex:  regexp.MustCompile("^(|.*/)folder2/([^/]*)(|/.*)$"),
						Negate: false,
						Owners: []string{
							"@group2",
						},
						Line: 4,
					},
					{
						Path:   "!file1",
						Regex:  regexp.MustCompile("^(|.*/)file1(|/.*)$"),
						Negate: true,
						Owners: []string{
							"@user3",
						},
						Line: 5,
					},
					{
						Path:   "folder1/*.tf",
						Regex:  regexp.MustCompile(`^(|/)folder1/([^/]*)\.tf(|/.*)$`),
						Negate: false,
						Owners: []string{
							"@user4",
						},
						Line: 6,
					},
					{
						Path:   "/**/",
						Regex:  regexp.MustCompile("^(|.*/)(|.*/)(|/.*)$"),
						Negate: false,
						Owners: []string{
							"@group1",
						},
						Line: 7,
					},
				},
				Error: false,
			},
		},
		{
			Name: "invalid codeowners entry",
			Sample: map[string]interface{}{
				"Filename": "valid-codeowners",
				"Contents": `*
folder1
folder2/
folder2/ @user2
!file1 @user3`,
			},
			Expected: ReturnWithError{
				Value: nil,
				Error: true,
			},
		},
	}

	for i, test := range tests {
		t.Logf("Test case %d: %s", i, test.Name)
		defer filet.CleanUp(t)
		expected := test.Expected.(ReturnWithError)
		sample := test.Sample.(map[string]interface{})

		if sample["Contents"].(string) != "" {
			filet.File(t, sample["Filename"].(string), sample["Contents"].(string))
		}
		val, err := ReadCodeownersFile(sample["Filename"].(string))
		if expected.Error {
			assert.Error(t, err, "should return an error")
			assert.Nil(t, val, "return should be nil on error")
		} else {
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, expected.Value.([]*CodeOwner), val, "decoded value should match expected")
		}
	}
}
