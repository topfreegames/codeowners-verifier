package gitlab_client

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

// GitlabClient interface implements the Gitlab Client
type ClientInterface interface {
	NewClient(Token string, BaseURL string)
	ListUsers(name string) ([]*gitlab.User, error)
	ListGroups(name string) ([]*gitlab.Group, error)
}

// Gitlab represents a Gitlab Client configuration
type Gitlab struct {
	Token   string
	BaseURL string
	api     ClientInterface
}

// GitlabClient implements a wrapper for calling the gitlab library
type GitlabClient struct {
	client *gitlab.Client
}

// NewClient returns a new Gitlab client
func (c *GitlabClient) NewClient(Token string, BaseURL string) {
	c.client, _ = gitlab.NewClient(Token, gitlab.WithBaseURL(BaseURL))
}

// ListUsers returns a list of Gitlab users matching the name
func (c *GitlabClient) ListUsers(name string) ([]*gitlab.User, error) {
	users, _, err := c.client.Users.ListUsers(&gitlab.ListUsersOptions{Search: gitlab.String(name)})
	if err != nil {
		return nil, fmt.Errorf("Error searching for user %s: %s", name, err)
	}
	return users, nil
}

// ListGroups returns a list of Gitlab users matching the name
func (c *GitlabClient) ListGroups(name string) ([]*gitlab.Group, error) {
	groups, _, err := c.client.Groups.ListGroups(&gitlab.ListGroupsOptions{Search: gitlab.String(name)})
	if err != nil {
		return nil, fmt.Errorf("Error searching for group %s: %s", name, err)
	}
	return groups, nil
}

// Init initializes the Gitlab Client
func (g *Gitlab) Init() error {
	if g.Token == "" {
		return fmt.Errorf("Token can't be empty")
	}
	if g.BaseURL == "" {
		g.BaseURL = "https://gitlab.com/api/v4"
	}
	if g.api == nil {
		g.api = &GitlabClient{}
	}
	g.api.NewClient(g.Token, g.BaseURL)
	return nil
}

// SearchUser searches a user by name
func (g *Gitlab) SearchUser(name string) (*gitlab.User, error) {
	users, err := g.api.ListUsers(name)
	if err != nil {
		return nil, err
	}
	if len(users) > 1 {
		return nil, fmt.Errorf("Multiple users found")
	}
	return users[0], nil
}

// SearchGroup searches a group by name
func (g *Gitlab) SearchGroup(name string) (*gitlab.Group, error) {
	groups, err := g.api.ListGroups(name)
	if err != nil {
		return nil, err
	}
	if len(groups) > 1 {
		return nil, fmt.Errorf("Multiple groups found")
	}
	return groups[0], nil
}
