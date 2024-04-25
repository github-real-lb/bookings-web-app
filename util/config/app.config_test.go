package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAppConfig(t *testing.T) {
	config, err := LoadAppConfig(testAppConfigFilename, ProductionMode)
	require.NoError(t, err)
	assert.Equal(t, ProductionMode, config.Mode)
	assert.Equal(t, config.StartingPathProduction+config.TemplateDirectoryName, config.TemplatePath)
	assert.Equal(t, config.StartingPathProduction+config.StaticDirectoryName, config.StaticPath)

	config, err = LoadAppConfig(testAppConfigFilename, DevelopmentMode)
	require.NoError(t, err)
	assert.Equal(t, DevelopmentMode, config.Mode)
	assert.Equal(t, config.StartingPathProduction+config.TemplateDirectoryName, config.TemplatePath)
	assert.Equal(t, config.StartingPathProduction+config.StaticDirectoryName, config.StaticPath)

	config, err = LoadAppConfig(testAppConfigFilename, TestingMode)
	require.NoError(t, err)
	assert.Equal(t, TestingMode, config.Mode)
	assert.Equal(t, config.StartingPathTesting+config.TemplateDirectoryName, config.TemplatePath)
	assert.Equal(t, config.StartingPathTesting+config.StaticDirectoryName, config.StaticPath)

	config, err = LoadAppConfig(testAppConfigFilename, DebuggingMode)
	require.NoError(t, err)
	assert.Equal(t, DebuggingMode, config.Mode)
	assert.Equal(t, config.StartingPathTesting+config.TemplateDirectoryName, config.TemplatePath)
	assert.Equal(t, config.StartingPathTesting+config.StaticDirectoryName, config.StaticPath)

	_, err = LoadAppConfig("", ProductionMode)
	require.Error(t, err)

	_, err = LoadAppConfig(testAppConfigFilename, -1)
	require.Error(t, err)
}

func TestAppConfig_SetProductionMode(t *testing.T) {
	app.SetProductionMode()
	assert.Equal(t, ProductionMode, app.Mode)
	assert.Equal(t, app.StartingPathProduction+app.TemplateDirectoryName, app.TemplatePath)
	assert.Equal(t, app.StartingPathProduction+app.StaticDirectoryName, app.StaticPath)
}

func TestAppConfig_SetDevelopementMode(t *testing.T) {
	app.SetDevelopementMode()
	assert.Equal(t, DevelopmentMode, app.Mode)
	assert.Equal(t, app.StartingPathProduction+app.TemplateDirectoryName, app.TemplatePath)
	assert.Equal(t, app.StartingPathProduction+app.StaticDirectoryName, app.StaticPath)
}

func TestAppConfig_SetTestingMode(t *testing.T) {
	app.SetTestingMode()
	assert.Equal(t, TestingMode, app.Mode)
	assert.Equal(t, app.StartingPathTesting+app.TemplateDirectoryName, app.TemplatePath)
	assert.Equal(t, app.StartingPathTesting+app.StaticDirectoryName, app.StaticPath)
}

func TestAppConfig_SetDebuggingMode(t *testing.T) {
	app.SetDebuggingMode()
	assert.Equal(t, DebuggingMode, app.Mode)
	assert.Equal(t, app.StartingPathTesting+app.TemplateDirectoryName, app.TemplatePath)
	assert.Equal(t, app.StartingPathTesting+app.StaticDirectoryName, app.StaticPath)
}

func TestAppConfig_InProductionMode(t *testing.T) {
	app.Mode = ProductionMode
	assert.True(t, app.InProductionMode())
}

func TestAppConfig_InDevelopmentMode(t *testing.T) {
	app.Mode = DevelopmentMode
	assert.True(t, app.InDevelopmentMode())
}

func TestAppConfig_InTestingMode(t *testing.T) {
	app.Mode = TestingMode
	assert.True(t, app.InTestingMode())
}

func TestAppConfig_InDebuggingMode(t *testing.T) {
	app.Mode = DebuggingMode
	assert.True(t, app.InDebuggingMode())
}
