package providers

import "fmt"

type Provider interface {
	UserExists(username string) (bool, error)
	GroupExists(username string) (bool, error)
}

func ListProviders() []string {
	return []string{"gitlab"}
}

func InitProvider(providerName string, token string, baseURL string) (Provider, error) {
	var provider Provider
	switch providerName {
	case "gitlab":
		var err error

		if provider, err = NewGitlabProvider(token, baseURL); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid provider")
	}
	return provider, nil
}
