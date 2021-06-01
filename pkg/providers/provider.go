package providers

import "fmt"

type Provider interface {
	UserExists(username string) (bool, error)
	GroupExists(username string) (bool, error)
}

func ListProviders() []string {
	return []string{"gitlab"}
}

func InitProvider(provider string, token string, baseURL string) (Provider, error) {
	var client Provider
	switch provider {
	case "gitlab":
		var err error

		if client, err = NewGitlabProviderClient(token, baseURL); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid provider")
	}
	return client, nil
}
