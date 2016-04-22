package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeManifests(t *testing.T) {
	assert := assert.New(t)

	// a complete manifest
	m1 := BuildManifest{
		Steps:            BuildSteps{"foo", "bar"},
		Android:          AndroidConfig{"foo", "bar"},
		DeviceGroupNames: []string{"foo", "bar"},
	}

	// m2 should override only Steps
	m2 := BuildManifest{
		Steps: BuildSteps{"bar", "foo"},
	}
	merged := MergeManifests(&m1, &m2)
	assert.Equal(BuildManifest{
		Steps:            BuildSteps{"bar", "foo"},
		Android:          AndroidConfig{"foo", "bar"},
		DeviceGroupNames: []string{"foo", "bar"},
	}, *merged)
	// m1 should override everything in m2
	merged = MergeManifests(&m2, &m1)
	assert.Equal(m1, *merged)

	// m3 should override only Android.Apk
	m3 := BuildManifest{
		Android: AndroidConfig{Apk: "bar"},
	}
	merged = MergeManifests(&m1, &m3)
	assert.Equal(BuildManifest{
		Steps:            BuildSteps{"foo", "bar"},
		Android:          AndroidConfig{"bar", "bar"},
		DeviceGroupNames: []string{"foo", "bar"},
	}, *merged)
	// m1 should override everything in m3
	merged = MergeManifests(&m3, &m1)
	assert.Equal(m1, *merged)

	// m4 should override only DeviceGroupNames
	m4 := BuildManifest{
		DeviceGroupNames: []string{"bar", "foo"},
	}
	merged = MergeManifests(&m1, &m4)
	assert.Equal(BuildManifest{
		Steps:            BuildSteps{"foo", "bar"},
		Android:          AndroidConfig{"foo", "bar"},
		DeviceGroupNames: []string{"bar", "foo"},
	}, *merged)
	// m1 should override everything in m4
	merged = MergeManifests(&m4, &m1)
	assert.Equal(m1, *merged)
}

func TestBuildManifestIsRunnable(t *testing.T) {
	assert := assert.New(t)

	// a complete manifest
	m1 := BuildManifest{
		Steps:            BuildSteps{"foo", "bar"},
		Android:          AndroidConfig{"foo", "bar"},
		DeviceGroupNames: []string{"foo", "bar"},
	}
	runnable, err := m1.IsRunnable()
	assert.True(runnable)
	assert.Nil(err)

	// missing Android.Apk
	m2 := BuildManifest{
		Steps:            BuildSteps{"foo", "bar"},
		Android:          AndroidConfig{ApkInstrumentation: "bar"},
		DeviceGroupNames: []string{"foo", "bar"},
	}
	runnable, err = m2.IsRunnable()
	assert.False(runnable)
	assert.NotNil(err)

	// missing a device group
	m3 := BuildManifest{
		Steps:            BuildSteps{"foo", "bar"},
		Android:          AndroidConfig{"foo", "bar"},
		DeviceGroupNames: []string{},
	}
	runnable, err = m3.IsRunnable()
	assert.False(runnable)
	assert.NotNil(err)
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	// a valid, complete config
	config, err := New("testdata/config.yml")
	assert.Nil(err)

	// build the expected config
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
		DeviceGroupNames: nil,
	}
	masterBuild := BuildManifest{
		DeviceGroupNames: []string{"everything"},
	}
	expected := Config{
		DeviceGroupDefinitions: deviceGroupDefs,
		Defaults:               defaultBuild,
		Branches:               map[string]BuildManifest{"master": masterBuild},
	}

	// config from the file should match the expected config
	assert.Equal(expected, *config)
}

func TestNewInvalid(t *testing.T) {
	assert := assert.New(t)

	// invalid because it is not valid YAML
	config, err := New("testdata/config_invalid.yml")
	assert.NotNil(err)
	assert.Nil(config)

	// file does not exist
	config, err = New("testdata/non_existant.yml")
	assert.NotNil(err)
	assert.Nil(config)

	// contains valid YAML but invalid properties
	config, err = New("testdata/config_bad.yml")
	assert.NotNil(err)
	assert.Nil(config)

	// missing devicegroups
	config, err = New("testdata/config_incomplete.yml")
	assert.NotNil(err)
	assert.Nil(config)
}

func TestConfigIsValid(t *testing.T) {
	assert := assert.New(t)

	// a valid config
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

func TestConfigBranchManifest(t *testing.T) {
	assert := assert.New(t)

	config, err := New("testdata/config.yml")
	assert.Nil(err)

	// build the expected manifest
	masterManifest := BuildManifest{
		Steps: BuildSteps{"echo \"Foo\"", "echo \"Bar\""},
		Android: AndroidConfig{
			Apk:                "./path/to/build.apk",
			ApkInstrumentation: "./path/to/instrumentation.apk",
		},
		DeviceGroupNames: []string{"everything"},
	}
	assert.Equal(masterManifest, *config.BranchManifest("master"))
}
