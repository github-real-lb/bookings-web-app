package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	app := LoadConfig()
	require.NotNil(t, app)
	assert.Equal(t, ProductionMode, app.AppMode)
	assert.True(t, app.UseTemplateCache)
	assert.NotEmpty(t, app.ServerAddress)
	assert.NotEmpty(t, app.TemplatePath)
}

func TestSetDevelopementMode(t *testing.T) {
	app := AppConfig{}
	app.SetDevelopementMode()
	assert.Equal(t, DevelopmentMode, app.AppMode)
	assert.False(t, app.UseTemplateCache)
}

func TestSetTestingMode(t *testing.T) {
	app := AppConfig{}
	app.SetTestingMode()
	assert.Equal(t, TestingMode, app.AppMode)
	assert.False(t, app.UseTemplateCache)
	assert.NotEmpty(t, app.TemplatePath)
}

func TestInProductionMode(t *testing.T) {
	app := AppConfig{}
	app.AppMode = ProductionMode
	assert.True(t, app.InProductionMode())
}

func TestInDevelopmentMode(t *testing.T) {
	app := AppConfig{}
	app.AppMode = DevelopmentMode
	assert.True(t, app.InDevelopmentMode())
}

func TestInTestingMode(t *testing.T) {
	app := AppConfig{}
	app.AppMode = TestingMode
	assert.True(t, app.InTestingMode())
}
