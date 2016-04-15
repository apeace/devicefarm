// Package config provides data structures used to configure devicefarm testing,
// as well as a function to read a config from a YAML file.
package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// A Device is just a string: the name of the device.
type Device string

// A DeviceGroup is just a list of Devices.
type DeviceGroup []Device

func (dg *DeviceGroup) Equals(dg2 DeviceGroup) bool {
	dg1 := *dg
	if len(dg1) != len(dg2) {
		return false
	}
	// TODO should not be order dependent
	for i := range dg1 {
		if dg2[i] != dg1[i] {
			return false
		}
	}
	return true
}

// A BuildSteps is just a list of strings: the bash commands to execute in order
// to perform the build.
type BuildSteps []string

func (bs *BuildSteps) Equals(bs2 BuildSteps) bool {
	bs1 := *bs
	if len(bs1) != len(bs2) {
		return false
	}
	for i := range bs1 {
		if bs2[i] != bs1[i] {
			return false
		}
	}
	return true
}

// An AndroidConfig specifies the location of APKs after running the build steps.
type AndroidConfig struct {
	Apk                string `yaml:"apk"`
	ApkInstrumentation string `yaml:"apk_instrumentation"`
}

func (c1 *AndroidConfig) Equals(c2 *AndroidConfig) bool {
	return c1.Apk == c2.Apk && c1.ApkInstrumentation == c2.ApkInstrumentation
}

// A BuildConfig specifies the whole configuration for a build: the steps to
// perform the build, the location of Android APKs, and the DeviceGroup names
// to run on.
type BuildConfig struct {
	Steps        BuildSteps    `yaml:"build"`
	Android      AndroidConfig `yaml:"android"`
	DeviceGroups []string      `yaml:"devicegroups"`
}

func (c1 *BuildConfig) Equals(c2 *BuildConfig) bool {
	if len(c1.DeviceGroups) != len(c2.DeviceGroups) {
		return false
	}
	for i := range c1.DeviceGroups {
		if c2.DeviceGroups[i] != c1.DeviceGroups[i] {
			return false
		}
	}
	return c1.Steps.Equals(c2.Steps) && c1.Android.Equals(&c2.Android)
}

// A Config specifies configuration for a particular repo: the names of DeviceGroups,
// the default BuildConfig, and override BuildConfigs for particular branches.
type Config struct {
	DeviceGroupDefinitions map[string]DeviceGroup `yaml:"devicegroups"`
	Defaults               BuildConfig            `yaml:"defaults"`
	Branches               map[string]BuildConfig `yaml:"branches"`
}

func (config *Config) IsValid() (bool, error) {
	if len(config.DeviceGroupDefinitions) == 0 {
		return false, fmt.Errorf("devicegroups must have at least one group")
	}
	for name, devicegroup := range config.DeviceGroupDefinitions {
		if len(devicegroup) == 0 {
			return false, fmt.Errorf("devicegroup %s has no devices", name)
		}
	}
	return true, nil
}

func (c1 *Config) Equals(c2 *Config) bool {
	if len(c1.DeviceGroupDefinitions) != len(c2.DeviceGroupDefinitions) {
		return false
	}
	for k, v1 := range c1.DeviceGroupDefinitions {
		if v2, ok := c2.DeviceGroupDefinitions[k]; !ok || !v1.Equals(v2) {
			return false
		}
	}
	if !c1.Defaults.Equals(&c2.Defaults) {
		return false
	}
	if len(c1.Branches) != len(c2.Branches) {
		return false
	}
	for k, v1 := range c1.Branches {
		if v2, ok := c1.Branches[k]; !ok || !v1.Equals(&v2) {
			return false
		}
	}
	return true
}

// Creates a new Config from a YAML file.
func New(filename string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	_, err = config.IsValid()
	if err != nil {
		return nil, err
	}
	return &config, nil
}
