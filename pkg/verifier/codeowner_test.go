// package verifier

// import (
// 	"path/filepath"
// 	"regexp"
// 	"testing"

// 	filet "github.com/Flaque/filet"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/topfreegames/codeowners-verifier/pkg/providers"
// 	"github.com/xanzy/go-gitlab"
// )

// type TestCase struct {
// 	Expected interface{}
// 	Sample   interface{}
// 	Name     string
// }

// type ReturnWithError struct {
// 	Value interface{}
// 	Error bool
// }

// func TestReverseCodeOwners(t *testing.T) {
// 	codeOwnerEntries := []*CodeOwner{
// 		{
// 			Path:   "*",
// 			Regex:  nil,
// 			Line:   1,
// 			Owners: []string{"@defaultGroup"},
// 			Negate: false,
// 		},
// 		{
// 			Path:   "path/to/file",
// 			Regex:  nil,
// 			Line:   2,
// 			Owners: []string{"@group1 @group2"},
// 			Negate: false,
// 		},
// 	}
// 	reversedCodeOwnerEntries := []*CodeOwner{
// 		{
// 			Path:   "path/to/file",
// 			Regex:  nil,
// 			Line:   2,
// 			Owners: []string{"@group1 @group2"},
// 			Negate: false,
// 		},
// 		{
// 			Path:   "*",
// 			Regex:  nil,
// 			Line:   1,
// 			Owners: []string{"@defaultGroup"},
// 			Negate: false,
// 		},
// 	}
// 	emptyCodeOwner := []*CodeOwner{}
// 	tests := []TestCase{
// 		{
// 			Name:     "Reversing CODEOWNERS file",
// 			Sample:   codeOwnerEntries,
// 			Expected: reversedCodeOwnerEntries,
// 		},
// 		{
// 			Name:     "Checking empty CODEOWNERS file",
// 			Sample:   emptyCodeOwner,
// 			Expected: emptyCodeOwner,
// 		},
// 	}

// 	for i, test := range tests {
// 		t.Logf("Test case %d: %s", i, test.Name)
// 		result := reverseCodeOwners(test.Sample.([]*CodeOwner))
// 		assert.Equal(t, test.Expected.([]*CodeOwner), result)
// 	}
// }

// func TestStripComment(t *testing.T) {
// 	tests := []TestCase{
// 		{
// 			Name:     "Checking line with #",
// 			Sample:   "#testing",
// 			Expected: "",
// 		},
// 		{
// 			Name:     "Checking line section []",
// 			Sample:   "[SectionThatShouldBeSanitized]",
// 			Expected: "",
// 		},
// 		{
// 			Name:     "Checking valid line",
// 			Sample:   "* @test",
// 			Expected: "* @test",
// 		},
// 	}

// 	for i, test := range tests {
// 		t.Logf("Test case %d: %s", i, test.Name)
// 		result := stripComment(test.Sample.(string))
// 		assert.Equal(t, test.Expected.(string), result)
// 	}
// }

// func TestDifference(t *testing.T) {
// 	tests := []TestCase{
// 		{
// 			Name: "Checking Difference with 2 identical slices",
// 			Sample: map[string]interface{}{
// 				"array1": []string{"a", "b", "c", "d"},
// 				"array2": []string{"a", "b", "c", "d"},
// 			},
// 			Expected: false,
// 		},
// 		{
// 			Name: "Checking Difference with 2 slightly different slices",
// 			Sample: map[string]interface{}{
// 				"array1": []string{"a", "b", "c", "d"},
// 				"array2": []string{"c", "d", "e", "f"},
// 			},
// 			Expected: true,
// 		},
// 		{
// 			Name: "Checking Difference with 2 completely different slices",
// 			Sample: map[string]interface{}{
// 				"array1": []string{"a", "b", "c", "d"},
// 				"array2": []string{"e", "f", "g", "h"},
// 			},
// 			Expected: true,
// 		},
// 	}
// 	for i, test := range tests {
// 		t.Logf("Test case %d: %s", i, test.Name)
// 		result := hasDifference(test.Sample.(map[string]interface{})["array1"].([]string), test.Sample.(map[string]interface{})["array2"].([]string))
// 		assert.Equal(t, test.Expected.(bool), result)
// 	}
// }

// func TestCodeOwnerReadFile(t *testing.T) {
// 	tests := []TestCase{
// 		{
// 			Name: "invalid file",
// 			Sample: map[string]interface{}{
// 				"Filename": "non-existent-file",
// 				"Contents": "",
// 			},
// 			Expected: ReturnWithError{
// 				Value: nil,
// 				Error: true,
// 			},
// 		},
// 		{
// 			Name: "valid codeowners",
// 			Sample: map[string]interface{}{
// 				"Filename": "valid-codeowners",
// 				"Contents": `* @user1 @user2
// folder1 @group1
// folder2/ @group1
// folder2/* @group2
// !file1 @user3
// folder1/*.tf @user4
// /**/ @group1`,
// 			},
// 			Expected: ReturnWithError{
// 				Value: []*CodeOwner{
// 					{
// 						Path:   "*",
// 						Regex:  regexp.MustCompile("^(|.*/)([^/]*)(|/.*)$"),
// 						Negate: false,
// 						Owners: []string{
// 							"@user1",
// 							"@user2",
// 						},
// 						Line: 1,
// 					},
// 					{
// 						Path:   "folder1",
// 						Regex:  regexp.MustCompile("^(|.*/)folder1(|/.*)$"),
// 						Negate: false,
// 						Owners: []string{
// 							"@group1",
// 						},
// 						Line: 2,
// 					},
// 					{
// 						Path:   "folder2/",
// 						Regex:  regexp.MustCompile("^(|.*/)folder2/(|.*)$"),
// 						Negate: false,
// 						Owners: []string{
// 							"@group1",
// 						},
// 						Line: 3,
// 					},
// 					{
// 						Path:   "folder2/*",
// 						Regex:  regexp.MustCompile("^(|.*/)folder2/([^/]*)(|/.*)$"),
// 						Negate: false,
// 						Owners: []string{
// 							"@group2",
// 						},
// 						Line: 4,
// 					},
// 					{
// 						Path:   "!file1",
// 						Regex:  regexp.MustCompile("^(|.*/)file1(|/.*)$"),
// 						Negate: true,
// 						Owners: []string{
// 							"@user3",
// 						},
// 						Line: 5,
// 					},
// 					{
// 						Path:   "folder1/*.tf",
// 						Regex:  regexp.MustCompile(`^(|/)folder1/([^/]*)\.tf(|/.*)$`),
// 						Negate: false,
// 						Owners: []string{
// 							"@user4",
// 						},
// 						Line: 6,
// 					},
// 					{
// 						Path:   "/**/",
// 						Regex:  regexp.MustCompile("^(|.*/)(|.*/)(|/.*)$"),
// 						Negate: false,
// 						Owners: []string{
// 							"@group1",
// 						},
// 						Line: 7,
// 					},
// 				},
// 				Error: false,
// 			},
// 		},
// 		{
// 			Name: "invalid codeowners entry",
// 			Sample: map[string]interface{}{
// 				"Filename": "valid-codeowners",
// 				"Contents": `*
// folder1
// folder2/
// folder2/ @user2
// !file1 @user3`,
// 			},
// 			Expected: ReturnWithError{
// 				Value: nil,
// 				Error: true,
// 			},
// 		},
// 	}

// 	for i, test := range tests {
// 		t.Logf("Test case %d: %s", i, test.Name)
// 		defer filet.CleanUp(t)
// 		expected := test.Expected.(ReturnWithError)
// 		sample := test.Sample.(map[string]interface{})

// 		if sample["Contents"].(string) != "" {
// 			filet.File(t, sample["Filename"].(string), sample["Contents"].(string))
// 		}
// 		val, err := ReadCodeownersFile(sample["Filename"].(string))
// 		if expected.Error {
// 			assert.Error(t, err, "should return an error")
// 			assert.Nil(t, val, "return should be nil on error")
// 		} else {
// 			assert.Nil(t, err, "should not return error")
// 			assert.Equal(t, expected.Value.([]*CodeOwner), val, "decoded value should match expected")
// 		}
// 	}
// }

// func TestMatchesPath(t *testing.T) {
// 	codeowners := []*CodeOwner{
// 		{
// 			Path:   "*",
// 			Regex:  regexp.MustCompile("^(|.*/)([^/]*)(|/.*)$"),
// 			Negate: false,
// 			Owners: []string{
// 				"@user1",
// 				"@user2",
// 			},
// 			Line: 1,
// 		},
// 		{
// 			Path:   "folder1",
// 			Regex:  regexp.MustCompile("^(|.*/)folder1(|/.*)$"),
// 			Negate: false,
// 			Owners: []string{
// 				"@group1",
// 			},
// 			Line: 2,
// 		},
// 		{
// 			Path:   "!file1",
// 			Regex:  regexp.MustCompile("^(|.*/)file1(|/.*)$"),
// 			Negate: true,
// 			Owners: []string{
// 				"@user3",
// 			},
// 			Line: 3,
// 		},
// 	}
// 	tests := []TestCase{
// 		{
// 			Name: "Check covered path for given file",
// 			Sample: map[string]interface{}{
// 				"CodeOwnerEntry": codeowners[0],
// 				"File":           "file2",
// 			},
// 			Expected: true,
// 		},
// 		{
// 			Name: "Check non-covered path for given file",
// 			Sample: map[string]interface{}{
// 				"CodeOwnerEntry": codeowners[1],
// 				"File":           "file2",
// 			},
// 			Expected: false,
// 		},
// 		{
// 			Name: "Check negated path for given file",
// 			Sample: map[string]interface{}{
// 				"CodeOwnerEntry": codeowners[2],
// 				"File":           "file1",
// 			},
// 			Expected: false,
// 		},
// 	}

// 	for i, test := range tests {
// 		t.Logf("Test case %d: %s", i, test.Name)
// 		sample := test.Sample.(map[string]interface{})
// 		result := sample["CodeOwnerEntry"].(*CodeOwner).MatchesPath(sample["File"].(string))
// 		assert.Equal(t, test.Expected, result)
// 	}
// }

// func TestVerifyCodeowner(t *testing.T) {
// 	codeowners := []*CodeOwner{
// 		{
// 			Path:   "*",
// 			Regex:  regexp.MustCompile("^(|.*/)([^/]*)(|/.*)$"),
// 			Negate: false,
// 			Owners: []string{
// 				"@user1",
// 				"@user2",
// 			},
// 			Line: 1,
// 		},
// 		{
// 			Path:   "folder1",
// 			Regex:  regexp.MustCompile("^(|.*/)folder1(|/.*)$"),
// 			Negate: false,
// 			Owners: []string{
// 				"@group1",
// 			},
// 			Line: 2,
// 		},
// 		{
// 			Path:   "folder2/",
// 			Regex:  regexp.MustCompile("^(|.*/)folder2/(|.*)$"),
// 			Negate: false,
// 			Owners: []string{
// 				"@group1",
// 			},
// 			Line: 3,
// 		},
// 		{
// 			Path:   "folder2/*",
// 			Regex:  regexp.MustCompile("^(|.*/)folder2/([^/]*)(|/.*)$"),
// 			Negate: false,
// 			Owners: []string{
// 				"@group2",
// 			},
// 			Line: 4,
// 		},
// 		{
// 			Path:   "!file1",
// 			Regex:  regexp.MustCompile("^(|.*/)file1(|/.*)$"),
// 			Negate: true,
// 			Owners: []string{
// 				"@user3",
// 			},
// 			Line: 5,
// 		},
// 		{
// 			Path:   "folder1/*.tf",
// 			Regex:  regexp.MustCompile(`^(|/)folder1/([^/]*)\.tf(|/.*)$`),
// 			Negate: false,
// 			Owners: []string{
// 				"@user4",
// 			},
// 			Line: 6,
// 		},
// 		{
// 			Path:   "/**/",
// 			Regex:  regexp.MustCompile("^(|.*/)(|.*/)(|/.*)$"),
// 			Negate: false,
// 			Owners: []string{
// 				"@group1",
// 			},
// 			Line: 7,
// 		},
// 	}
// 	tests := []TestCase{
// 		{
// 			Name: "Checking file on empty CODEOWNERS",
// 			Sample: map[string]interface{}{
// 				"CodeOwners": []*CodeOwner{},
// 				"File":       "",
// 				"Ignore":     []string{},
// 			},
// 			Expected: map[string]interface{}{
// 				"Codeowners": &CodeOwner{},
// 				"Valid":      false,
// 			},
// 		},
// 		{
// 			Name: "Checking covered file on CODEOWNERS - no ignores",
// 			Sample: map[string]interface{}{
// 				"CodeOwners": codeowners,
// 				"File":       "file2",
// 				"Ignore":     []string{},
// 			},
// 			Expected: map[string]interface{}{
// 				"Codeowners": codeowners[0],
// 				"Valid":      true,
// 			},
// 		},
// 		{
// 			Name: "Checking covered file on CODEOWNERS - with ignores",
// 			Sample: map[string]interface{}{
// 				"CodeOwners": codeowners,
// 				"File":       "file2",
// 				"Ignore": []string{
// 					"@user1",
// 					"@user2",
// 				},
// 			},
// 			Expected: map[string]interface{}{
// 				"Codeowners": codeowners[0],
// 				"Valid":      false,
// 			},
// 		},
// 	}
// 	for i, test := range tests {
// 		t.Logf("Test case %d: %s", i, test.Name)
// 		sample := test.Sample.(map[string]interface{})
// 		expected := test.Expected.(map[string]interface{})
// 		entry, valid := VerifyCodeowner(sample["CodeOwners"].([]*CodeOwner), sample["File"].(string), sample["Ignore"].([]string))

// 		assert.Equal(t, expected["Codeowners"].(*CodeOwner), entry)
// 		assert.Equal(t, expected["Valid"].(bool), valid)
// 	}
// }

// func TestValidateCodeownerFileGitlab(t *testing.T) {
// 	defer filet.CleanUp(t)
// 	/*
// 		Directory Tree created for tests
// 		./
// 		./folder1
// 		./folder1/file1
// 		./folder2
// 		./folder2/file2
// 		./folder2/folder3/
// 		./folder2/folder3/file3
// 		./file4
// 	*/
// 	folder1 := filet.TmpDir(t, "./")
// 	folder2 := filet.TmpDir(t, "./")
// 	folder3 := filet.TmpDir(t, filepath.Join("./", folder2))
// 	filet.TmpFile(t, folder1, "")
// 	filet.TmpFile(t, folder2, "")
// 	file1 := filet.TmpFile(t, folder3, "").Name()
// 	filet.TmpFile(t, "./", "")
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
// 	MockGitlabClient := providers.NewMockClientInterface(mockCtrl)
// 	tests := []TestCase{
// 		{
// 			Name: "invalid file",
// 			Sample: map[string]interface{}{
// 				"CodeOwners": "non-existent-file",
// 				"Provider": &providers.Gitlab{
// 					Token: "xxx",
// 					Api:   MockGitlabClient,
// 				},
// 			},
// 			Expected: ReturnWithError{
// 				Value: false,
// 				Error: true,
// 			},
// 		},
// 		{
// 			Name: "Correct CodeOwners File",
// 			Sample: map[string]interface{}{
// 				"CodeOwners": filet.TmpFile(t, "", `* @user1
// /*.* @user1
// /**/* @user1
// /`+folder1+`/** @user1
// `+folder1+` @user2 @group1
// `+folder1+`/* @user3
// `+folder2+`/** @group1
// `+folder3+` @user1
// `+file1+` @user1`).Name(),
// 				"Provider": &providers.Gitlab{
// 					Token: "xxx",
// 					Api:   MockGitlabClient,
// 				},
// 			},
// 			Expected: ReturnWithError{
// 				Value: true,
// 				Error: false,
// 			},
// 		},
// 		{
// 			Name: "invalid codeowner entry",
// 			Sample: map[string]interface{}{
// 				"CodeOwners": filet.TmpFile(t, "", `missing-owner`).Name(),
// 				"Provider": &providers.Gitlab{
// 					Token: "xxx",
// 					Api:   MockGitlabClient,
// 				},
// 			},
// 			Expected: ReturnWithError{
// 				Value: false,
// 				Error: true,
// 			},
// 		},
// 		{
// 			Name: "non-existent path",
// 			Sample: map[string]interface{}{
// 				"CodeOwners": filet.TmpFile(t, "", `invalid-path @user1`).Name(),
// 				"Provider": &providers.Gitlab{
// 					Token: "xxx",
// 					Api:   MockGitlabClient,
// 				},
// 			},
// 			Expected: ReturnWithError{
// 				Value: false,
// 				Error: false,
// 			},
// 		},
// 		{
// 			Name: "non-existent owner",
// 			Sample: map[string]interface{}{
// 				"CodeOwners": filet.TmpFile(t, "", folder1+` @user100`).Name(),
// 				"Provider": &providers.Gitlab{
// 					Token: "xxx",
// 					Api:   MockGitlabClient,
// 				},
// 			},
// 			Expected: ReturnWithError{
// 				Value: false,
// 				Error: false,
// 			},
// 		},
// 	}
// 	for i, test := range tests {
// 		t.Logf("Test case %d: %s", i, test.Name)
// 		defer filet.CleanUp(t)
// 		expected := test.Expected.(ReturnWithError)
// 		sample := test.Sample.(map[string]interface{})
// 		// We could improve this logic
// 		MockGitlabClient.EXPECT().ListUsers("user1").Return([]*gitlab.User{{Username: "user1"}}, nil).AnyTimes()
// 		MockGitlabClient.EXPECT().ListUsers("user2").Return([]*gitlab.User{{Username: "user2"}}, nil).AnyTimes()
// 		MockGitlabClient.EXPECT().ListUsers("user3").Return([]*gitlab.User{{Username: "user3"}}, nil).AnyTimes()
// 		MockGitlabClient.EXPECT().ListUsers("user100").Return([]*gitlab.User{}, nil).AnyTimes()
// 		MockGitlabClient.EXPECT().ListUsers("group1").Return([]*gitlab.User{}, nil).AnyTimes()
// 		MockGitlabClient.EXPECT().ListGroups("group1").Return([]*gitlab.Group{{Name: "group1"}}, nil).AnyTimes()
// 		MockGitlabClient.EXPECT().ListGroups("user100").Return([]*gitlab.Group{}, nil).AnyTimes()
// 		val, err := ValidateCodeownerFile(sample["Provider"].(providers.Provider), sample["CodeOwners"].(string))
// 		if expected.Error {
// 			assert.Error(t, err, "should return an error")
// 			assert.Equal(t, false, val, "return should be false on error")
// 		} else {
// 			assert.Nil(t, err, "should not return error")
// 			assert.Equal(t, expected.Value.(bool), val, "decoded value should match expected")
// 		}
// 	}
// }
