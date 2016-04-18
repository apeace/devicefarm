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

// A BuildSteps is just a list of strings: the bash commands to execute in order
// to perform the build.
type BuildSteps []string

// An AndroidConfig specifies the location of APKs after running the build steps.
type AndroidConfig struct {
	Apk                string `yaml:"apk"`
	ApkInstrumentation string `yaml:"apk_instrumentation"`
}

// A BuildManifest specifies the whole configuration for a build: the steps to
// perform the build, the location of Android APKs, and the DeviceGroup names
// to run on.
type BuildManifest struct {
	Steps            BuildSteps    `yaml:"build"`
	Android          AndroidConfig `yaml:"android"`
	DeviceGroupNames []string      `yaml:"devicegroups"`
}

// A Config specifies configuration for a particular repo: the names of DeviceGroups,
// the default BuildManifest, and override BuildManifests for particular branches.
type Config struct {
	DeviceGroupDefinitions map[string]DeviceGroup   `yaml:"devicegroups"`
	Defaults               BuildManifest            `yaml:"defaults"`
	Branches               map[string]BuildManifest `yaml:"branches"`
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
