package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAppConfig(t *testing.T) {
	config, err := LoadAppConfig(testAppConfigFilename, ProductionMode)
	require.NoError(t, err)
	assert.Equal(t, ProductionMode, config.AppMode)
	assert.Equal(t, config.TemplatePathProduction, config.TemplatePath)
	assert.True(t, config.UseTemplateCache)

	config, err = LoadAppConfig(testAppConfigFilename, DevelopmentMode)
	require.NoError(t, err)
	assert.Equal(t, DevelopmentMode, config.AppMode)
	assert.Equal(t, config.TemplatePathProduction, config.TemplatePath)
	assert.False(t, config.UseTemplateCache)

	config, err = LoadAppConfig(testAppConfigFilename, TestingMode)
	require.NoError(t, err)
	assert.Equal(t, TestingMode, config.AppMode)
	assert.Equal(t, config.TemplatePathTesting, config.TemplatePath)
	assert.False(t, config.UseTemplateCache)

	_, err = LoadAppConfig("", ProductionMode)
	require.Error(t, err)

	_, err = LoadAppConfig(testAppConfigFilename, -1)
	require.Error(t, err)
}

func TestAppConfig_SetProductionMode(t *testing.T) {
	app.SetProductionMode()
	assert.Equal(t, ProductionMode, app.AppMode)
	assert.Equal(t, app.TemplatePathProduction, app.TemplatePath)
	assert.True(t, app.UseTemplateCache)
}

func TestAppConfig_SetDevelopementMode(t *testing.T) {
	app.SetDevelopementMode()
	assert.Equal(t, DevelopmentMode, app.AppMode)
	assert.Equal(t, app.TemplatePathProduction, app.TemplatePath)
	assert.False(t, app.UseTemplateCache)
}

func TestAppConfig_SetTestingMode(t *testing.T) {
	app.SetTestingMode()
	assert.Equal(t, TestingMode, app.AppMode)
	assert.Equal(t, app.TemplatePathTesting, app.TemplatePath)
	assert.False(t, app.UseTemplateCache)
}

func TestAppConfig_InProductionMode(t *testing.T) {
	app.AppMode = ProductionMode
	assert.True(t, app.InProductionMode())
}

func TestAppConfig_InDevelopmentMode(t *testing.T) {
	app.AppMode = DevelopmentMode
	assert.True(t, app.InDevelopmentMode())
}

func TestAppConfig_InTestingMode(t *testing.T) {
	app.AppMode = TestingMode
	assert.True(t, app.InTestingMode())
}
