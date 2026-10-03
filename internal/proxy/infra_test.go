package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shopwell-shop/shopwell-cli/internal/shop"
)

func TestIsProxyProjectForDomain(t *testing.T) {
	t.Parallel()

	proxied := func(url string) *shop.Config {
		return &shop.Config{URL: url}
	}

	assert.True(t, IsProxyProjectForDomain(proxied("https://my-shop.shopwell.local"), "shopwell.local"))
	assert.True(t, IsProxyProjectForDomain(proxied("https://shopwell.local"), "shopwell.local"))
	assert.True(t, IsProxyProjectForDomain(proxied("https://my-shop.dev.internal"), "dev.internal"))

	// The local environment url overrides the top-level one.
	envProxied := &shop.Config{
		URL:          "http://127.0.0.1:8000",
		Environments: map[string]*shop.EnvironmentConfig{"local": {URL: "https://my-shop.shopwell.local"}},
	}
	assert.True(t, IsProxyProjectForDomain(envProxied, "shopwell.local"))

	// Port-based projects are not proxy projects.
	assert.False(t, IsProxyProjectForDomain(proxied("http://127.0.0.1:8000"), "shopwell.local"))
	assert.False(t, IsProxyProjectForDomain(proxied("http://localhost:8000"), "shopwell.local"))
	// Hostname under a different domain is not matched.
	assert.False(t, IsProxyProjectForDomain(proxied("https://my-shop.example.com"), "shopwell.local"))
	// A hostname that merely ends with the domain as a substring (no dot) must not match.
	assert.False(t, IsProxyProjectForDomain(proxied("https://notshopwell.local"), "shopwell.local"))
	assert.False(t, IsProxyProjectForDomain(nil, "shopwell.local"))
	assert.False(t, IsProxyProjectForDomain(&shop.Config{}, "shopwell.local"))
}
