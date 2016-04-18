package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigIsValid(t *testing.T) {
	assert := assert.New(t)

	c1 := Config{DeviceGroupDefinitions: map[string]DeviceGroup{"foo": []Device{"bar"}}}
	ok, err := c1.IsValid()
	assert.True(ok)
	assert.Nil(err)

	// invalid due to empty device group def
	c2 := Config{DeviceGroupDefinitions: map[string]DeviceGroup{"foo": []Device{}}}
	ok, err = c2.IsValid()
	assert.False(ok)
	assert.NotNil(err)

	// invalid due to no device group defs
	c3 := Config{}
	ok, err = c3.IsValid()
	assert.False(ok)
	assert.NotNil(err)
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	config, err := New("testdata/config.yml")
	assert.Nil(err)

	deviceGroupDefs := map[string]DeviceGroup{
		"a_few_devices": {"Samsung S3", "Blah fone"},
		"samsung_s4_s5": {"Samsung S4 TMobile", "Samsung S4 AT&T", "Samsung S5 TMobile", "Samsung S5 AT&T"},
		"everything":    {"+a_few_devices", "+samsung_s4_s5"},
	}
	defaultBuild := BuildManifest{
		Steps: BuildSteps{"echo \"Foo\"", "echo \"Bar\""},
		Android: AndroidConfig{
			Apk:                "./path/to/build.apk",
			ApkInstrumentation: "./path/to/instrumentation.apk",
		},
		DeviceGroupNames: []string{"a_few_devices"},
	}
	masterBuild := BuildManifest{
		DeviceGroupNames: []string{"everything"},
	}
	expected := Config{
		DeviceGroupDefinitions: deviceGroupDefs,
		Defaults:               defaultBuild,
		Branches:               map[string]BuildManifest{"master": masterBuild},
	}
	assert.Equal(expected, *config)
}

func TestNewInvalid(t *testing.T) {
	assert := assert.New(t)

	config, err := New("testdata/config_invalid.yml")
	assert.NotNil(err)
	assert.Nil(config)

	config, err = New("testdata/non_existant.yml")
	assert.NotNil(err)
	assert.Nil(config)

	config, err = New("testdata/config_bad.yml")
	assert.NotNil(err)
	assert.Nil(config)

	config, err = New("testdata/config_incomplete.yml")
	assert.NotNil(err)
	assert.Nil(config)
}
