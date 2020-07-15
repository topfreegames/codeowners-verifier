package providers

type Provider interface {
	Init() error
	UserExists(username string) (bool, error)
	GroupExists(username string) (bool, error)
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
	}
	return client, nil
}
