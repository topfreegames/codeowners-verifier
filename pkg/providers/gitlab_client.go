package providers

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

// GitlabClient interface implements the Gitlab Client
//go:generate mockgen -destination=gitlab_client_mock.go -package=providers github.com/topfreegames/codeowners-verifier/pkg/providers ClientInterface
type ClientInterface interface {
	NewClient(token string, baseURL string)
	ListUsers(name string) ([]*gitlab.User, error)
	ListGroups(name string) ([]*gitlab.Group, error)
	ListAllUsers() ([]*gitlab.User, error)
	ListAllGroups() ([]*gitlab.Group, error)
}

// gitlabProvider represents a Gitlab Provider configuration
type gitlabProvider struct {
	token   string
	baseURL string
	api     ClientInterface
	users   []*gitlab.User
	groups  []*gitlab.Group
}

// GitlabClient implements a wrapper for calling the gitlab library
type GitlabClient struct {
	client *gitlab.Client
}

// NewClient returns a new Gitlab client
func (c *GitlabClient) NewClient(Token string, BaseURL string) {
	c.client, _ = gitlab.NewClient(Token, gitlab.WithBaseURL(BaseURL))
}

// NewClient returns a new Gitlab
func NewGitlabProviderClient(token string, baseURL string) (*gitlabProvider, error) {
	g := &gitlabProvider{
		token:   token,
		baseURL: baseURL,
	}

	if g.token == "" {
		return nil, fmt.Errorf("token can't be empty")
	}
	if g.baseURL == "" {
		g.baseURL = "https://gitlab.com/api/v4"
	}
	if g.api == nil {
		g.api = &GitlabClient{}
	}
	g.api.NewClient(g.token, g.baseURL)

	users, err := g.api.ListAllUsers()
	if err != nil {
		return nil, err
	}
	g.users = users

	groups, err := g.api.ListAllGroups()
	if err != nil {
		return nil, err
	}
	g.groups = groups

	return g, nil
}

// ListUsers returns a list of Gitlab users matching the name
func (c *GitlabClient) ListUsers(name string) ([]*gitlab.User, error) {
	users, _, err := c.client.Users.ListUsers(&gitlab.ListUsersOptions{Search: gitlab.String(name)})
	if err != nil {
		return nil, fmt.Errorf("error searching for user %s: %s", name, err)
	}
	return users, nil
}

// ListAllUsers returns a list of all Gitlab users
func (c *GitlabClient) ListAllUsers() ([]*gitlab.User, error) {
	users, _, err := c.client.Users.ListUsers(&gitlab.ListUsersOptions{})
	if err != nil {
		return nil, fmt.Errorf("error retriving user list: %s", err)
	}
	return users, nil
}

// ListGroups returns a list of Gitlab groups matching the name
func (c *GitlabClient) ListGroups(name string) ([]*gitlab.Group, error) {
	groups, _, err := c.client.Groups.ListGroups(&gitlab.ListGroupsOptions{Search: gitlab.String(name)})
	if err != nil {
		return nil, fmt.Errorf("error searching for group %s: %s", name, err)
	}
	return groups, nil
}

// ListAllGroups returns a list of all Gitlab groups
func (c *GitlabClient) ListAllGroups() ([]*gitlab.Group, error) {
	groups, _, err := c.client.Groups.ListGroups(&gitlab.ListGroupsOptions{})
	if err != nil {
		return nil, fmt.Errorf("error searching for group %s", err)
	}
	return groups, nil
}

// SearchUser searches a user by name
func (g *gitlabProvider) UserExists(name string) (bool, error) {

	for _, user := range g.users {
		if user.Username == name {
			return true, nil
		}
	}
	return false, nil
}

// SearchGroup searches a group by name
func (g *gitlabProvider) GroupExists(name string) (bool, error) {

	for _, group := range g.groups {
		if group.Name == name {
			return true, nil
		}
	}
	return false, nil
}
