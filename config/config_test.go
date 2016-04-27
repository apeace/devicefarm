package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeManifests(t *testing.T) {
	assert := assert.New(t)

	// a complete manifest
	m1 := BuildManifest{
		Steps:      []string{"foo", "bar"},
		Android:    AndroidConfig{"foo", "bar"},
		DevicePool: "foo",
	}

	// m2 should override only Steps
	m2 := BuildManifest{
		Steps: []string{"bar", "foo"},
	}
	merged := MergeManifests(&m1, &m2)
	assert.Equal(BuildManifest{
		Steps:      []string{"bar", "foo"},
		Android:    AndroidConfig{"foo", "bar"},
		DevicePool: "foo",
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
		Steps:      []string{"foo", "bar"},
		Android:    AndroidConfig{"bar", "bar"},
		DevicePool: "foo",
	}, *merged)
	// m1 should override everything in m3
	merged = MergeManifests(&m3, &m1)
	assert.Equal(m1, *merged)

	// m4 should override only DevicePool
	m4 := BuildManifest{
		DevicePool: "bar",
	}
	merged = MergeManifests(&m1, &m4)
	assert.Equal(BuildManifest{
		Steps:      []string{"foo", "bar"},
		Android:    AndroidConfig{"foo", "bar"},
		DevicePool: "bar",
	}, *merged)
	// m1 should override everything in m4
	merged = MergeManifests(&m4, &m1)
	assert.Equal(m1, *merged)
}

func TestBuildManifestIsRunnable(t *testing.T) {
	assert := assert.New(t)

	// a complete manifest
	m1 := BuildManifest{
		Steps:      []string{"foo", "bar"},
		Android:    AndroidConfig{"foo", "bar"},
		DevicePool: "foo",
	}
	runnable, err := m1.IsRunnable()
	assert.True(runnable)
	assert.Nil(err)

	// missing Android.Apk
	m2 := BuildManifest{
		Steps:      []string{"foo", "bar"},
		Android:    AndroidConfig{ApkInstrumentation: "bar"},
		DevicePool: "foo",
	}
	runnable, err = m2.IsRunnable()
	assert.False(runnable)
	assert.NotNil(err)

	// missing a device pool
	m3 := BuildManifest{
		Steps:      []string{"foo", "bar"},
		Android:    AndroidConfig{"foo", "bar"},
		DevicePool: "",
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
	devicePoolDefs := map[string][]string{
		"a_few_devices": {"Samsung S3", "Blah fone"},
		"samsung_s4_s5": {"Samsung S4 TMobile", "Samsung S4 AT&T", "Samsung S5 TMobile", "Samsung S5 AT&T"},
		"everything":    {"+a_few_devices", "+samsung_s4_s5"},
	}
	defaultBuild := BuildManifest{
		Steps: []string{"echo \"Foo\"", "echo \"Bar\""},
		Android: AndroidConfig{
			Apk:                "./path/to/build.apk",
			ApkInstrumentation: "./path/to/instrumentation.apk",
		},
		DevicePool: "",
	}
	masterBuild := BuildManifest{
		DevicePool: "everything",
	}
	expected := Config{
		Arn: "arn:aws:devicefarm:us-west-2:026109802893:project:1124416c-bfb2-4334-817c-e211ecef7dc0",
		DevicePoolDefinitions: devicePoolDefs,
		Defaults:              defaultBuild,
		Branches:              map[string]BuildManifest{"master": masterBuild},
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

	// missing devicepools
	config, err = New("testdata/config_incomplete.yml")
	assert.NotNil(err)
	assert.Nil(config)
}

func TestConfigIsValid(t *testing.T) {
	assert := assert.New(t)

	// a valid config
	c1 := Config{Arn: "foo", DevicePoolDefinitions: map[string][]string{"foo": {"bar"}}}
	ok, err := c1.IsValid()
	assert.True(ok)
	assert.Nil(err)

	// invalid due to blank device pool item
	c2 := Config{Arn: "foo", DevicePoolDefinitions: map[string][]string{"foo": {""}}}
	ok, err = c2.IsValid()
	assert.False(ok)
	assert.NotNil(err)

	// invalid due to no device pool defs
	c3 := Config{Arn: "foo"}
	ok, err = c3.IsValid()
	assert.False(ok)
	assert.NotNil(err)

	// invalid due to missing Arn
	c4 := Config{DevicePoolDefinitions: map[string][]string{"foo": {"bar"}}}
	ok, err = c4.IsValid()
	assert.False(ok)
	assert.NotNil(err)

	// invalid due to empty device pool def
	c5 := Config{Arn: "foo", DevicePoolDefinitions: map[string][]string{"foo": {}}}
	ok, err = c5.IsValid()
	assert.False(ok)
	assert.NotNil(err)
}

func TestConfigBranchManifest(t *testing.T) {
	assert := assert.New(t)

	config, err := New("testdata/config.yml")
	assert.Nil(err)

	// build the expected manifest
	masterManifest := BuildManifest{
		Steps: []string{"echo \"Foo\"", "echo \"Bar\""},
		Android: AndroidConfig{
			Apk:                "./path/to/build.apk",
			ApkInstrumentation: "./path/to/instrumentation.apk",
		},
		DevicePool: "everything",
	}
	assert.Equal(masterManifest, *config.BranchManifest("master"))
}

func TestFlatDevicePoolDefinitions(t *testing.T) {
	assert := assert.New(t)

	// should successfully get flattened result
	defs := map[string][]string{
		"pool1": {"foo", "bar"},
		"pool2": {"bar", "+pool1"},
	}
	config := Config{DevicePoolDefinitions: defs}
	flat, err := config.FlatDevicePoolDefinitions()
	assert.Nil(err)
	assert.Equal(map[string][]string{
		"pool1": {"bar", "foo"},
		"pool2": {"bar", "foo"},
	}, flat)

	// should fail because of circular dependency
	defs = map[string][]string{
		"pool1": {"foo", "bar", "+pool2"},
		"pool2": {"bar", "+pool1"},
	}
	config = Config{DevicePoolDefinitions: defs}
	flat, err = config.FlatDevicePoolDefinitions()
	assert.NotNil(err)

	// should fail because "pool3" does not exist
	defs = map[string][]string{
		"pool1": {"foo", "bar", "+pool3"},
		"pool2": {"bar", "+pool1"},
	}
	config = Config{DevicePoolDefinitions: defs}
	flat, err = config.FlatDevicePoolDefinitions()
	assert.NotNil(err)
}
