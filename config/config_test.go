package config

import "testing"

func TestNew(t *testing.T) {
	config, err := New("testdata/config.yml")
	if err != nil {
		t.Error(err)
	}
	deviceGroupDefs := map[string]DeviceGroup{
		"a_few_devices": {"Samsung S3", "Blah fone"},
		"samsung_s4_s5": {"Samsung S4 TMobile", "Samsung S4 AT&T", "Samsung S5 TMobile", "Samsung S5 AT&T"},
		"everything":    {"+a_few_devices", "+samsung_s4_s5"},
	}
	defaultBuild := BuildConfig{
		Steps: BuildSteps{"echo \"Foo\"", "echo \"Bar\""},
		Android: AndroidConfig{
			Apk:                "./path/to/build.apk",
			ApkInstrumentation: "./path/to/instrumentation.apk",
		},
		DeviceGroups: []string{"a_few_devices"},
	}
	masterBuild := BuildConfig{
		DeviceGroups: []string{"everything"},
	}
	expected := &Config{
		DeviceGroupDefinitions: deviceGroupDefs,
		Defaults:               defaultBuild,
		Branches:               map[string]BuildConfig{"master": masterBuild},
	}
	if !expected.Equals(config) {
		t.Errorf("Expected %v to equal %v", expected, config)
	}
}

func TestNewInvalid(t *testing.T) {
	config, err := New("testdata/config_invalid.yml")
	if err == nil {
		t.Errorf("Expected non-nil error, got %v", err)
	}
	if config != nil {
		t.Errorf("Expected nil result, got %v", config)
	}
	config, err = New("testdata/non_existant.yml")
	if err == nil {
		t.Errorf("Expected non-nil error, got %v", err)
	}
	if config != nil {
		t.Errorf("Expected nil result, got %v", config)
	}
}

func TestNewBad(t *testing.T) {
	config, err := New("testdata/config_bad.yml")
	if err == nil {
		t.Error("Expected non-nil error")
	}
	if config != nil {
		t.Errorf("Expected nil result, got %v", config)
	}
}
