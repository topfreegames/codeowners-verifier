package providers

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

// GitlabClient interface implements the Gitlab Client
//go:generate mockgen -destination=gitlab_client_mock.go -package=providers github.com/topfreegames/codeowners-verifier/pkg/providers ClientInterface
type ClientInterface interface {
	NewClient(token string, baseURL string)
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

// NewClient returns a new Gitlab Provider Client
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

// ListAllUsers returns a list of all Gitlab users
func (c *GitlabClient) ListAllUsers() ([]*gitlab.User, error) {
	var users []*gitlab.User

	opt := &gitlab.ListUsersOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		iteration_users, resp, err := c.client.Users.ListUsers(opt)
		if err != nil {
			fmt.Errorf("error retriving user list: %s", err)
		}
		users = append(users, iteration_users...)
		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	return users, nil
}

// ListAllGroups returns a list of all Gitlab groups
func (c *GitlabClient) ListAllGroups() ([]*gitlab.Group, error) {

	var groups []*gitlab.Group

	opt := &gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		iteration_groups, resp, err := c.client.Groups.ListGroups(opt)
		if err != nil {
			fmt.Errorf("error retriving group list: %s", err)
		}
		groups = append(groups, iteration_groups...)
		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
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
