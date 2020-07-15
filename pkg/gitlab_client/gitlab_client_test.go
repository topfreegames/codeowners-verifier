package gitlab_client

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func TestNewClientSucessful(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockGitlabClient := NewMockClientInterface(mockCtrl)
	client := &Gitlab{
		Token:   "example_token",
		BaseURL: "example_url",
		api:     MockGitlabClient,
	}
	MockGitlabClient.EXPECT().NewClient(client.Token, client.BaseURL).Times(1)
	client.api.NewClient(client.Token, client.BaseURL)
}

func TestInitSucessful(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockGitlabClient := NewMockClientInterface(mockCtrl)
	client := &Gitlab{
		Token:   "Token",
		BaseURL: "BaseURL",
		api:     MockGitlabClient,
	}
	MockGitlabClient.EXPECT().NewClient(client.Token, client.BaseURL).Times(1)
	assert.Equal(t, nil, client.Init())
}
func TestInitMissingToken(t *testing.T) {
	client := &Gitlab{
		BaseURL: "BaseURL",
	}
	assert.Error(t, client.Init(), "Token can't be empty")
}
func TestEmptyBaseURL(t *testing.T) {
	client := &Gitlab{
		Token: "token",
	}
	assert.Equal(t, nil, client.Init())
	assert.Equal(t, "https://gitlab.com/api/v4", client.BaseURL)
}

func TestListUsersSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	username := "mock_user"
	gitlabUsers := []*gitlab.User{
		{
			Username: username,
		},
	}
	defer mockCtrl.Finish()
	MockGitlabClient := NewMockClientInterface(mockCtrl)
	client := &Gitlab{
		Token:   "example_token",
		BaseURL: "example_url",
		api:     MockGitlabClient,
	}
	MockGitlabClient.EXPECT().ListUsers(username).Return(gitlabUsers, nil).Times(1)
	user, err := client.api.ListUsers(username)
	assert.Equal(t, err, nil)
	assert.Equal(t, gitlabUsers, user)
}
func TestListUsersFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	username := "mock_user"
	defer mockCtrl.Finish()
	MockGitlabClient := NewMockClientInterface(mockCtrl)
	client := &Gitlab{
		Token:   "example_token",
		BaseURL: "example_url",
		api:     MockGitlabClient,
	}
	MockGitlabClient.EXPECT().ListUsers(username).Return([]*gitlab.User{}, fmt.Errorf("Error searching for user %s:", username)).Times(1)
	user, err := client.api.ListUsers(username)
	assert.Error(t, err, fmt.Errorf("Error searching for user %s:", username))
	assert.Equal(t, []*gitlab.User{}, user)
}

func TestSearchUserSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	username := "mock_user"
	gitlabUsers := []*gitlab.User{
		{
			Username: username,
		},
	}
	defer mockCtrl.Finish()
	MockGitlabClient := NewMockClientInterface(mockCtrl)
	client := &Gitlab{
		Token:   "example_token",
		BaseURL: "example_url",
		api:     MockGitlabClient,
	}
	MockGitlabClient.EXPECT().ListUsers(username).Return(gitlabUsers, nil).Times(1)
	user, err := client.SearchUser(username)
	assert.Equal(t, err, nil)
	assert.Equal(t, gitlabUsers[0].Username, user.Username)
}
func TestSearchUserFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	username := "mock_user"
	defer mockCtrl.Finish()
	MockGitlabClient := NewMockClientInterface(mockCtrl)
	client := &Gitlab{
		Token:   "example_token",
		BaseURL: "example_url",
		api:     MockGitlabClient,
	}
	MockGitlabClient.EXPECT().ListUsers(username).Return([]*gitlab.User{}, fmt.Errorf("Error searching for user %s:", username)).Times(1)
	user, err := client.SearchUser(username)
	assert.Equal(t, err, fmt.Errorf("Error searching for user %s:", username))
	assert.Equal(t, (*gitlab.User)(nil), user)
}
func TestSearchUserMultipleUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	username := "mock_user"
	gitlabUsers := []*gitlab.User{
		{
			Username: username,
		},
		{
			Username: username + "_2",
		},
	}
	defer mockCtrl.Finish()
	MockGitlabClient := NewMockClientInterface(mockCtrl)
	client := &Gitlab{
		Token:   "example_token",
		BaseURL: "example_url",
		api:     MockGitlabClient,
	}
	MockGitlabClient.EXPECT().ListUsers(username).Return(gitlabUsers, nil).Times(1)
	user, err := client.SearchUser(username)
	assert.Equal(t, err, fmt.Errorf("Multiple users found"))
	assert.Equal(t, (*gitlab.User)(nil), user)
}
func TestSearchGroupSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	groupName := "mock_group"
	gitlabGroups := []*gitlab.Group{
		{
			Name: groupName,
		},
	}
	defer mockCtrl.Finish()
	MockGitlabClient := NewMockClientInterface(mockCtrl)
	client := &Gitlab{
		Token:   "example_token",
		BaseURL: "example_url",
		api:     MockGitlabClient,
	}
	MockGitlabClient.EXPECT().ListGroups(groupName).Return(gitlabGroups, nil).Times(1)
	group, err := client.SearchGroup(groupName)
	assert.Equal(t, err, nil)
	assert.Equal(t, gitlabGroups[0].Name, group.Name)
}
func TestSearchGroupFailure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	groupName := "mock_group"
	defer mockCtrl.Finish()
	MockGitlabClient := NewMockClientInterface(mockCtrl)
	client := &Gitlab{
		Token:   "example_token",
		BaseURL: "example_url",
		api:     MockGitlabClient,
	}
	MockGitlabClient.EXPECT().ListGroups(groupName).Return([]*gitlab.Group{}, fmt.Errorf("Error searching for group %s", groupName)).Times(1)
	group, err := client.SearchGroup(groupName)
	assert.Equal(t, fmt.Errorf("Error searching for group %s", groupName), err)
	assert.Equal(t, (*gitlab.Group)(nil), group)
}
func TestSearchGroupMultiplegroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	groupName := "mock_group"
	gitlabGroups := []*gitlab.Group{
		{
			Name: groupName,
		},
		{
			Name: groupName + "_2",
		},
	}
	defer mockCtrl.Finish()
	MockGitlabClient := NewMockClientInterface(mockCtrl)
	client := &Gitlab{
		Token:   "example_token",
		BaseURL: "example_url",
		api:     MockGitlabClient,
	}
	MockGitlabClient.EXPECT().ListGroups(groupName).Return(gitlabGroups, nil).Times(1)
	group, err := client.SearchGroup(groupName)
	assert.Equal(t, err, fmt.Errorf("Multiple groups found"))
	assert.Equal(t, (*gitlab.Group)(nil), group)
}
