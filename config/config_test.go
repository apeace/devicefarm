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
		"samsung_s3": {
			"(arn=device:50E24178F2274CFFA577EF130440D066) Samsung Galaxy S3 (AT&T)",
			"(arn=device:71F791A0C3CA4E9999304A1E8484339B) Samsung Galaxy S3 (Sprint)",
			"(arn=device:9E079354B7E9422CA52FF61B0BE345A1) Samsung Galaxy S3 (T-Mobile)",
			"(arn=device:E024E20134534DE1AFD87038726AB05C) Samsung Galaxy S3 (Verizon)",
			"(arn=device:BD86B8701031476BA30AF3D03F06B665) Samsung Galaxy S3 (Verizon)",
			"(arn=device:B6100FEA90BC4B21BD6C607865AD46F2) Samsung Galaxy S3 LTE (T-Mobile)",
			"(arn=device:5C748437DC1C409EA595B98B1D7A8EDD) Samsung Galaxy S3 Mini (AT&T)",
		},
		"samsung_s4": {
			"(arn=device:D1C28D6B913C479399C0F594E1EBCAE4) Samsung Galaxy S4 (AT&T)",
			"(arn=device:449870B9550C4840ACC1C1B59A7027FB) Samsung Galaxy S4 (AT&T)",
			"(arn=device:2A81F49C0CBD4AB6B1C2C58C1498F51F) Samsung Galaxy S4 (AT&T)",
			"(arn=device:33F66BE404B543669978079E905F8637) Samsung Galaxy S4 (Sprint)",
			"(arn=device:D45C750161314335924CE0B9B7D2558E) Samsung Galaxy S4 (T-Mobile)",
			"(arn=device:9E882A633A8E4ADC9C402AD22B1455E4) Samsung Galaxy S4 (US Cellular)",
			"(arn=device:47869F01A5F44B8999030BC0580703E5) Samsung Galaxy S4 (Verizon)",
			"(arn=device:6E920D51A4624ECA9EC856E0CAE733B9) Samsung Galaxy S4 (Verizon)",
			"(arn=device:577DC08D6B964346B86610CFF090CD59) Samsung Galaxy S4 Active (AT&T)",
			"(arn=device:F17F20E555C54544B722557AF43B015E) Samsung Galaxy S4 Tri-band (Sprint)",
			"(arn=device:20766AF83D3A4FEF977643BFCDC2CE3A) Samsung Galaxy S4 mini (Verizon)",
		},
		"samsung_s5": {
			"(arn=device:5CC0164714304CBF81BB7B7C03DFC1A1) Samsung Galaxy S5 (AT&T)",
			"(arn=device:53586C603C5A4FA38602D11AD917B01E) Samsung Galaxy S5 (AT&T)",
			"(arn=device:18E28478F1D54525A15C2A821B6132FA) Samsung Galaxy S5 (Sprint)",
			"(arn=device:D6F125CF316C47B09F5190C16DE979A9) Samsung Galaxy S5 (Sprint)",
			"(arn=device:5931A012CB1C4E68BD3434DF722ADBC8) Samsung Galaxy S5 (T-Mobile)",
			"(arn=device:C30737D1E582482C9D06BC4878E7F795) Samsung Galaxy S5 (Verizon)",
			"(arn=device:9710D509338C4639ADEFC5D6E99F45E6) Samsung Galaxy S5 Active (AT&T)",
		},
		"everything": {
			"+samsung_s3",
			"+samsung_s4",
			"+samsung_s5",
		},
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
		ProjectArn:            "arn:aws:devicefarm:us-west-2:026109802893:project:1124416c-bfb2-4334-817c-e211ecef7dc0",
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
