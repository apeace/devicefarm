/*

Package config provides data structures and functions to configure Device Farm testing.

YAML Config

Repositories using the devicefarm tool will provide a devicefarm.yml file in the
repository root. The config.New() function parses this config from a file into a
config.Config struct:

	config, err := config.New("path/to/config.yml")

Here is an annotated example of what a config file should look like:

	# Device Group definitions. this block defines three Device Groups:
	# a_few_devices, samsung_s4_st, and everything. The everything group
	# simply includes both the other groups.
	#
	# This property is REQUIRED, must have at least one group defined,
	# and each group must have at least one device.
	devicegroups:
	  a_few_devices:
		- Samsung S3
		- Blah fone

	  samsung_s4_s5:
		- Samsung S4 TMobile
		- Samsung S4 AT&T
		- Samsung S5 TMobile
		- Samsung S5 AT&T

	  everything:
		- +a_few_devices
		- +samsung_s4_s5

	# Defaults defines the build config that will be used for all branches,
	# unless overrides are specified in the branches section.
	#
	# This property is OPTIONAL, but building will fail on a particular branch
	# unless a full definition is available for that branch.
	defaults:
		# The bash commands to run for this build.
		build:
		  - echo "Foo"
		  - echo "Bar"

		# The location of APK files, after build commands have been run.
		android:
		  apk: ./path/to/build.apk
		  apk_instrumentation: ./path/to/instrumentation.apk

		# The device group names that tests should be run on.
		devicegroups:
		  - a_few_devices

	# Branches defines overrides for particular branches. For each branch,
	# it accepts the same properties as `defaults`. Branch configs will be
	# merged with `defaults` so that the specified properties override the
	# same properties from `defaults`. In this example, only the `devicegroups`
	# property will be overridden for the `master` branch.
	#
	# This property is OPTIONAL, but building will fail on a particular branch
	# unless a full definition is available for that branch.
	branches:
	  master:
		devicegroups:
		  - everything

*/
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
	existsInDg1 := map[Device]bool{}
	for _, device := range dg1 {
		existsInDg1[device] = true
	}
	for _, device := range dg2 {
		if !existsInDg1[device] {
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
	Steps            BuildSteps    `yaml:"build"`
	Android          AndroidConfig `yaml:"android"`
	DeviceGroupNames []string      `yaml:"devicegroups"`
}

func (c1 *BuildConfig) Equals(c2 *BuildConfig) bool {
	if len(c1.DeviceGroupNames) != len(c2.DeviceGroupNames) {
		return false
	}
	existsInDg1 := map[string]bool{}
	for _, deviceGroup := range c1.DeviceGroupNames {
		existsInDg1[deviceGroup] = true
	}
	for _, deviceGroup := range c2.DeviceGroupNames {
		if !existsInDg1[deviceGroup] {
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
	valid, err := config.IsValid()
	if !valid {
		return nil, err
	}
	return &config, nil
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
		if v2, ok := c2.Branches[k]; !ok || !v1.Equals(&v2) {
			return false
		}
	}
	return true
}
