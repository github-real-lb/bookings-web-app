package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDBConfigFilename = "./../../db.config.json"

func TestLoadDBConfig(t *testing.T) {
	dbConfig, err := LoadDBConfig(testDBConfigFilename)
	require.NoError(t, err)
	require.NotNil(t, dbConfig)
	assert.NotEmpty(t, dbConfig.ConnectionString)
	assert.NotEmpty(t, dbConfig.ServerAddress)

	_, err = LoadDBConfig("")
	require.Error(t, err)
}
