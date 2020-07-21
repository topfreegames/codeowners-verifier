package providers

import "fmt"

type Provider interface {
	Init() error
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
		client = &Gitlab{
			Token:   token,
			BaseURL: baseURL,
		}
		if err := client.Init(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Invalid provider")
	}
	return client, nil
}
