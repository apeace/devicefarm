package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeManifests(t *testing.T) {
	assert := assert.New(t)

	// a complete manifest
	complete := BuildManifest{
		Steps:      []string{"foo", "bar"},
		Tests:      map[string]map[string]string{"foo": {"bar": "baz"}},
		DevicePool: "foo",
	}

	// should override only Steps
	override_steps := BuildManifest{
		Steps: []string{"bar", "foo"},
	}
	merged := MergeManifests(&complete, &override_steps)
	assert.Equal(BuildManifest{
		Steps:      []string{"bar", "foo"},
		Tests:      map[string]map[string]string{"foo": {"bar": "baz"}},
		DevicePool: "foo",
	}, *merged)
	// complete should override everything
	merged = MergeManifests(&override_steps, &complete)
	assert.Equal(complete, *merged)

	// should override only Tests
	override_tests := BuildManifest{
		Tests: map[string]map[string]string{"bar": {"baz": "blah"}},
	}
	merged = MergeManifests(&complete, &override_tests)
	assert.Equal(BuildManifest{
		Steps:      []string{"foo", "bar"},
		Tests:      map[string]map[string]string{"bar": {"baz": "blah"}},
		DevicePool: "foo",
	}, *merged)
	// complete should override everything
	merged = MergeManifests(&override_tests, &complete)
	assert.Equal(complete, *merged)

	// should override only DevicePool
	override_devicepool := BuildManifest{
		DevicePool: "bar",
	}
	merged = MergeManifests(&complete, &override_devicepool)
	assert.Equal(BuildManifest{
		Steps:      []string{"foo", "bar"},
		Tests:      map[string]map[string]string{"foo": {"bar": "baz"}},
		DevicePool: "bar",
	}, *merged)
	// complete should override everything
	merged = MergeManifests(&override_devicepool, &complete)
	assert.Equal(complete, *merged)

	// an empty manifest
	// steps and tests should default to empty but non-nil values
	empty := BuildManifest{}
	merged = MergeManifests(&empty, &empty)
	assert.Equal(BuildManifest{
		Steps:      []string{},
		Tests:      map[string]map[string]string{},
		DevicePool: "",
	}, *merged)
}

func TestBuildManifestIsRunnable(t *testing.T) {
	assert := assert.New(t)

	// a complete manifest, should be runnable
	m1 := BuildManifest{
		Steps: []string{"foo", "bar"},
		Tests: map[string]map[string]string{
			"instrumentation": {
				"type":                "android_instrumentation",
				"app_apk":             "./path/to/build.apk",
				"instrumentation_apk": "./path/to/instrumentation.apk",
			},
		},
		DevicePool: "foo",
	}
	runnable, err := m1.IsRunnable()
	assert.True(runnable)
	assert.Nil(err)

	// a complete manifest missing build steps, should be runnable
	m2 := BuildManifest{
		Steps: []string{},
		Tests: map[string]map[string]string{
			"instrumentation": {
				"type":                "android_instrumentation",
				"app_apk":             "./path/to/build.apk",
				"instrumentation_apk": "./path/to/instrumentation.apk",
			},
		},
		DevicePool: "foo",
	}
	runnable, err = m2.IsRunnable()
	assert.True(runnable)
	assert.Nil(err)

	// missing Tests, should NOT be runnable
	m3 := BuildManifest{
		Steps:      []string{"foo", "bar"},
		Tests:      map[string]map[string]string{},
		DevicePool: "foo",
	}
	runnable, err = m3.IsRunnable()
	assert.False(runnable)
	assert.NotNil(err)

	// missing a device pool, should NOT be runnable
	m4 := BuildManifest{
		Steps: []string{"foo", "bar"},
		Tests: map[string]map[string]string{
			"instrumentation": {
				"type":                "android_instrumentation",
				"app_apk":             "./path/to/build.apk",
				"instrumentation_apk": "./path/to/instrumentation.apk",
			},
		},
		DevicePool: "",
	}
	runnable, err = m4.IsRunnable()
	assert.False(runnable)
	assert.NotNil(err)
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	// build the expected config
	devicePoolDefs := map[string][]string{
		"samsung_s3": {
			"(arn=device:50E24178F2274CFFA577EF130440D066) Samsung Galaxy S3 (AT&T)",
			"(arn=device:71F791A0C3CA4E9999304A1E8484339B) Samsung Galaxy S3 (Sprint)",
			"(arn=device:9E079354B7E9422CA52FF61B0BE345A1) Samsung Galaxy S3 (T-Mobile)",
		},
		"samsung_s4": {
			"(arn=device:D1C28D6B913C479399C0F594E1EBCAE4) Samsung Galaxy S4 (AT&T)",
			"(arn=device:33F66BE404B543669978079E905F8637) Samsung Galaxy S4 (Sprint)",
			"(arn=device:D45C750161314335924CE0B9B7D2558E) Samsung Galaxy S4 (T-Mobile)",
		},
		"samsung_s5": {
			"(arn=device:5CC0164714304CBF81BB7B7C03DFC1A1) Samsung Galaxy S5 (AT&T)",
			"(arn=device:18E28478F1D54525A15C2A821B6132FA) Samsung Galaxy S5 (Sprint)",
			"(arn=device:5931A012CB1C4E68BD3434DF722ADBC8) Samsung Galaxy S5 (T-Mobile)",
		},
		"everything": {
			"+samsung_s3",
			"+samsung_s4",
			"+samsung_s5",
		},
	}
	defaultBuild := BuildManifest{
		Steps: []string{"echo \"Foo\"", "echo \"Bar\""},
		Tests: map[string]map[string]string{
			"instrumentation": {
				"type":                "android_instrumentation",
				"app_apk":             "./path/to/build.apk",
				"instrumentation_apk": "./path/to/instrumentation.apk",
			},
		},
		Android: AndroidConfig{
			Apk:                "",
			ApkInstrumentation: "",
		},
		DevicePool: "samsung_s5",
	}
	masterBuild := BuildManifest{
		DevicePool: "everything",
	}
	expected := Config{
		ProjectArn:            "arn:aws:devicefarm:us-west-2:026109802893:project:1124416c-bfb2-4334-817c-e211ecef7dc0",
		DevicePoolDefinitions: devicePoolDefs,
		Defaults:              defaultBuild,
		Branches:              map[string]BuildManifest{"master": masterBuild},
	}

	// a valid, complete config
	config, err := New("testdata/config.yml")
	assert.Nil(err)
	assert.Equal(expected, *config)

	// a valid config with no build
	config, err = New("testdata/config_nobuild.yml")
	assert.Nil(err)
	expected.Defaults.Steps = []string{}
	assert.Equal(expected, *config)

	// a valid config using the deprecated android field
	config, err = New("testdata/config_deprecated.yml")
	assert.Nil(err)
	expected.Defaults.Tests["deprecated_auto_instrumentation"] = expected.Defaults.Tests["instrumentation"]
	delete(expected.Defaults.Tests, "instrumentation")
	assert.Equal(expected, *config)
}

func TestNewInvalid(t *testing.T) {
	assert := assert.New(t)

	// invalid because it is not valid YAML
	config, err := New("testdata/config_invalid_yaml.yml")
	assert.NotNil(err)
	assert.Nil(config)

	// file does not exist
	config, err = New("testdata/non_existant.yml")
	assert.NotNil(err)
	assert.Nil(config)

	// contains valid YAML but invalid properties
	config, err = New("testdata/config_invalid_properties.yml")
	assert.NotNil(err)
	assert.Nil(config)

	// missing devicepools
	config, err = New("testdata/config_nodevicepooldef.yml")
	assert.NotNil(err)
	assert.Nil(config)
}

func TestConfigIsValid(t *testing.T) {
	assert := assert.New(t)

	arn := "arn:aws:devicefarm:us-west-2:026109802893:project:1124416c-bfb2-4334-817c-e211ecef7dc0"

	// a valid config
	c1 := Config{ProjectArn: arn, DevicePoolDefinitions: map[string][]string{"foo": {"bar"}}}
	ok, err := c1.IsValid()
	assert.True(ok)
	assert.Nil(err)

	// invalid due to blank device pool item
	c2 := Config{ProjectArn: arn, DevicePoolDefinitions: map[string][]string{"foo": {""}}}
	ok, err = c2.IsValid()
	assert.False(ok)
	assert.NotNil(err)

	// invalid due to no device pool defs
	c3 := Config{ProjectArn: arn}
	ok, err = c3.IsValid()
	assert.False(ok)
	assert.NotNil(err)

	// invalid due to missing Arn
	c4 := Config{DevicePoolDefinitions: map[string][]string{"foo": {"bar"}}}
	ok, err = c4.IsValid()
	assert.False(ok)
	assert.NotNil(err)

	// invalid due to empty device pool def
	c5 := Config{ProjectArn: arn, DevicePoolDefinitions: map[string][]string{"foo": {}}}
	ok, err = c5.IsValid()
	assert.False(ok)
	assert.NotNil(err)
}

func TestConfigBranchManifest(t *testing.T) {
	assert := assert.New(t)

	// build the expected manifest for branch "master"
	masterManifest := BuildManifest{
		Steps: []string{"echo \"Foo\"", "echo \"Bar\""},
		Tests: map[string]map[string]string{
			"instrumentation": {
				"type":                "android_instrumentation",
				"app_apk":             "./path/to/build.apk",
				"instrumentation_apk": "./path/to/instrumentation.apk",
			},
		},
		Android: AndroidConfig{
			Apk:                "",
			ApkInstrumentation: "",
		},
		DevicePool: "everything",
	}

	config, err := New("testdata/config.yml")
	assert.Nil(err)
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

func TestDeviceArns(t *testing.T) {
	assert := assert.New(t)

	// should fail because invalid format
	devices := []string{"foo"}
	parsed, err := DeviceArns(devices)
	assert.NotNil(err)
	assert.Nil(parsed)

	// should succeed
	devices = []string{
		"(arn=device:50E24178F2274CFFA577EF130440D066) Samsung Galaxy S3 (AT&T)",
	}
	parsed, err = DeviceArns(devices)
	assert.Nil(err)
	assert.Equal([]string{
		"arn:aws:devicefarm:us-west-2::device:50E24178F2274CFFA577EF130440D066",
	}, parsed)
}
