package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_ShouldUseDefaultAddrWhenEnvVariableIsNotSet(t *testing.T) {
	c := NewConfig(nil)
	assert.Equal(t, defaultAddr, c.Addr)
}

func TestNewConfig_ShouldUseAddrWhenEnvVariableIsSet(t *testing.T) {
	_ = os.Setenv("APP_ADDR", ":30")
	c := NewConfig(nil)
	assert.Equal(t, ":30", c.Addr)
}

func TestNewConfig_ShouldUseDefaultShortURLDomainWhenEnvVariableIsNotSet(t *testing.T) {
	c := NewConfig(nil)
	assert.Equal(t, defaultShortURLDomain, c.ShortURLDomain)
}

func TestNewConfig_ShouldUseShortURLDomainFromEnvVariableIfItIsSetted(t *testing.T) {
	_ = os.Setenv("SHORT_URL_DOMAIN", "tujix.me")
	c := NewConfig(nil)
	assert.Equal(t, "tujix.me", c.ShortURLDomain)
}
