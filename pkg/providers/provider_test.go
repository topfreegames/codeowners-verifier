package providers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getProviders() []string {
	return []string{"gitlab"}
}

func TestInitProviderSuccess(t *testing.T) {
	token := "xyz"
	baseURL := ""
	for _, p := range getProviders() {
		t.Logf("Validating provider %s", p)
		provider, err := InitProvider(p, token, baseURL)
		assert.Equal(t, nil, err)
		_, ok := interface{}(provider).(Provider)
		assert.Equal(t, true, ok)
	}
}
func TestInitProviderError(t *testing.T) {
	token := ""
	baseURL := ""
	for _, p := range getProviders() {
		t.Logf("Validating provider %s", p)
		provider, err := InitProvider(p, token, baseURL)
		assert.Equal(t, fmt.Errorf("Token can't be empty"), err)
		_, ok := interface{}(provider).(Provider)
		assert.Equal(t, false, ok)
	}
}
