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

func TestCodeOwnerReadFile(t *testing.T) {
	directory := []string{
		"folder1",
		"folder1/subfolder1",
		"folder2/file2",
		"folder2/subfolder2/file3",
	}
	tests := []TestCase{
		{
			Name: "invalid file",
			Sample: map[string]interface{}{
				"Filename": "non-existent-file",
				"Contents": "",
				"Members":  []string{},
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
				"Contents": `
* @user1 @user2
folder1 @group1
folder2/ @group1
folder2/* @group2
!file1 @user3
`,
				"Members": []string{},
			},
			Expected: ReturnWithError{
				Value: []*CodeOwner{
					{
						Path:   "*",
						Regex:  nil,
						Negate: false,
						Owners: []string{
							"@user1",
							"@user2",
						},
						Line: 1,
					},
					{
						Path:   "folder1",
						Regex:  nil,
						Negate: false,
						Owners: []string{
							"@group1",
						},
						Line: 2,
					},
					{
						Path:   "folder2/",
						Regex:  nil,
						Negate: false,
						Owners: []string{
							"@group1",
						},
						Line: 3,
					},
					{
						Path:   "folder2/*",
						Regex:  nil,
						Negate: false,
						Owners: []string{
							"@group2",
						},
						Line: 4,
					},
					{
						Path:   "*",
						Regex:  nil,
						Negate: true,
						Owners: []string{
							"@user3",
						},
						Line: 5,
					},
				},
				Error: false,
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
