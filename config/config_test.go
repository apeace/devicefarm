package config

import "testing"

func TestDeviceGroupEquals(t *testing.T) {
	d1 := DeviceGroup{"foo", "bar"}
	d2 := DeviceGroup{"bar", "foo"}
	d3 := DeviceGroup{"foo"}
	d4 := DeviceGroup{"foo", "baz"}
	if !d1.Equals(d2) {
		t.Errorf("Expected %v to equal %v", d1, d2)
	}
	if d1.Equals(d3) {
		t.Errorf("Expected %v to NOT equal %v", d1, d3)
	}
	if d1.Equals(d4) {
		t.Errorf("Expected %v to NOT equal %v", d1, d4)
	}
}

func TestBuildStepEquals(t *testing.T) {
	b1 := BuildSteps{"foo", "bar"}
	b2 := BuildSteps{"foo", "bar"}
	b3 := BuildSteps{"bar", "foo"}
	b4 := BuildSteps{"foo"}
	b5 := BuildSteps{"foo", "baz"}
	if !b1.Equals(b2) {
		t.Errorf("Expected %v to equal %v", b1, b2)
	}
	if b1.Equals(b3) {
		t.Errorf("Expected %v to NOT equal %v", b1, b3)
	}
	if b1.Equals(b4) {
		t.Errorf("Expected %v to NOT equal %v", b1, b5)
	}
	if b1.Equals(b4) {
		t.Errorf("Expected %v to NOT equal %v", b1, b5)
	}
}

func TestBuildConfigEquals(t *testing.T) {
	b1 := BuildConfig{DeviceGroupNames: []string{"foo", "bar"}}
	b2 := BuildConfig{DeviceGroupNames: []string{"bar", "foo"}}
	b3 := BuildConfig{DeviceGroupNames: []string{"foo"}}
	b4 := BuildConfig{DeviceGroupNames: []string{"foo", "baz"}}
	if !b1.Equals(&b2) {
		t.Errorf("Expected %v to equal %v", b1, b2)
	}
	if b1.Equals(&b3) {
		t.Errorf("Expected %v to NOT equal %v", b1, b3)
	}
	if b1.Equals(&b4) {
		t.Errorf("Expected %v to NOT equal %v", b1, b4)
	}
}

func TestConfigIsValid(t *testing.T) {
	c1 := Config{DeviceGroupDefinitions: map[string]DeviceGroup{"foo": []Device{"bar"}}}
	c2 := Config{DeviceGroupDefinitions: map[string]DeviceGroup{"foo": []Device{}}}
	c3 := Config{}
	if ok, err := c1.IsValid(); !ok {
		t.Errorf("Expected %v to be valid but got %v", c1, err)
	}
	if ok, _ := c2.IsValid(); ok {
		t.Errorf("Expected %v to be invalid", c2)
	}
	if ok, _ := c3.IsValid(); ok {
		t.Errorf("Expected %v to be invalid", c3)
	}
}

func TestConfigEquals(t *testing.T) {
	c1 := Config{DeviceGroupDefinitions: map[string]DeviceGroup{"foo": []Device{"bar"}}}
	c2 := Config{DeviceGroupDefinitions: map[string]DeviceGroup{"foo": []Device{"bar"}}}
	c3 := Config{DeviceGroupDefinitions: map[string]DeviceGroup{"foo": []Device{"baz"}}}
	c4 := Config{DeviceGroupDefinitions: map[string]DeviceGroup{"foo": []Device{"baz"}, "bar": []Device{"blah"}}}
	if !c1.Equals(&c2) {
		t.Errorf("Expected %v to equal %v", c1, c2)
	}
	if c1.Equals(&c3) {
		t.Errorf("Expected %v to NOT equal %v", c1, c3)
	}
	if c1.Equals(&c4) {
		t.Errorf("Expected %v to NOT equal %v", c1, c4)
	}
	c5 := Config{Defaults: BuildConfig{Steps: BuildSteps{"foo", "bar"}}}
	c6 := Config{Defaults: BuildConfig{Steps: BuildSteps{"bar", "foo"}}}
	if c5.Equals(&c6) {
		t.Errorf("Expected %v to NOT equal %v", c5, c6)
	}
	c7 := Config{Branches: map[string]BuildConfig{"master": BuildConfig{}}}
	c8 := Config{Branches: map[string]BuildConfig{"staging": BuildConfig{}}}
	c9 := Config{Branches: map[string]BuildConfig{"staging": BuildConfig{}, "master": BuildConfig{}}}
	if c7.Equals(&c8) {
		t.Errorf("Expected %v to NOT equal %v", c7, c8)
	}
	if c7.Equals(&c9) {
		t.Errorf("Expected %v to NOT equal %v", c7, c9)
	}
}

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
		DeviceGroupNames: []string{"a_few_devices"},
	}
	masterBuild := BuildConfig{
		DeviceGroupNames: []string{"everything"},
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
	config, err = New("testdata/config_bad.yml")
	if err == nil {
		t.Error("Expected non-nil error")
	}
	if config != nil {
		t.Errorf("Expected nil result, got %v", config)
	}
	config, err = New("testdata/config_incomplete.yml")
	if err == nil {
		t.Error("Expected non-nil error")
	}
	if config != nil {
		t.Errorf("Expected nil result, got %v", config)
	}
}
