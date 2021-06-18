package providers

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

// gitlabProvider represents a Gitlab Provider configuration
type gitlabProvider struct {
	token   string
	baseURL string
	client  *gitlab.Client
	users   []*gitlab.User
	groups  []*gitlab.Group
}

// NewClient returns a new Gitlab Provider
func NewGitlabProvider(token string, baseURL string) (*gitlabProvider, error) {

	gitlabProvider := &gitlabProvider{
		token:   token,
		baseURL: baseURL,
	}

	if gitlabProvider.token == "" {
		return nil, fmt.Errorf("token can't be empty")
	}
	if gitlabProvider.baseURL == "" {
		gitlabProvider.baseURL = "https://gitlab.com/api/v4"
	}

	client, err := gitlabProvider.newClient(gitlabProvider.token, gitlabProvider.baseURL)
	if err != nil {
		return nil, err
	}
	gitlabProvider.client = client

	users, err := gitlabProvider.listAllUsers()
	if err != nil {
		return nil, err
	}
	gitlabProvider.users = users

	groups, err := gitlabProvider.listAllGroups()
	if err != nil {
		return nil, err
	}
	gitlabProvider.groups = groups

	return gitlabProvider, nil
}

// newClient returns a new Gitlab client
func (c *gitlabProvider) newClient(Token string, BaseURL string) (*gitlab.Client, error) {
	gitlabClient, err := gitlab.NewClient(Token, gitlab.WithBaseURL(BaseURL))
	if err != nil {
		fmt.Errorf("error creating gitlab client: %s", err)
	}

	return gitlabClient, err
}

// listAllUsers returns a list of all Gitlab users
func (c *gitlabProvider) listAllUsers() ([]*gitlab.User, error) {
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

// listAllGroups returns a list of all Gitlab groups
func (g *gitlabProvider) listAllGroups() ([]*gitlab.Group, error) {

	var groups []*gitlab.Group

	opt := &gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		iteration_groups, resp, err := g.client.Groups.ListGroups(opt)
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

// UserExists verifies if a user exists in Gitlab by name
func (g *gitlabProvider) UserExists(name string) (bool, error) {

	for _, user := range g.users {
		if user.Username == name {
			return true, nil
		}
	}
	return false, nil
}

// GroupExists verifies if a group exists in Gitlab by name
func (g *gitlabProvider) GroupExists(name string) (bool, error) {

	for _, group := range g.groups {
		if group.Name == name {
			return true, nil
		}
	}
	return false, nil
}
