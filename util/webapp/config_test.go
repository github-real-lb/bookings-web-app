package webapp

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

func TestAppConfig_SetDevelopementMode(t *testing.T) {
	app := AppConfig{}
	app.SetDevelopementMode()
	assert.Equal(t, DevelopmentMode, app.AppMode)
	assert.False(t, app.UseTemplateCache)
}

func TestAppConfig_SetTestingMode(t *testing.T) {
	app := AppConfig{}
	app.SetTestingMode()
	assert.Equal(t, TestingMode, app.AppMode)
	assert.False(t, app.UseTemplateCache)
	assert.NotEmpty(t, app.TemplatePath)
}

func TestAppConfig_InProductionMode(t *testing.T) {
	app := AppConfig{}
	app.AppMode = ProductionMode
	assert.True(t, app.InProductionMode())
}

func TestAppConfig_InDevelopmentMode(t *testing.T) {
	app := AppConfig{}
	app.AppMode = DevelopmentMode
	assert.True(t, app.InDevelopmentMode())
}

func TestAppConfig_InTestingMode(t *testing.T) {
	app := AppConfig{}
	app.AppMode = TestingMode
	assert.True(t, app.InTestingMode())
}
