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

// MergeManfiests merges together two BuildManifests, giving the second manifest
// priority. In other words, any non-blank field in the second manifest will
// override the value in the first manifest.
func MergeManifests(m1 *BuildManifest, m2 *BuildManifest) *BuildManifest {
	merged := &BuildManifest{}
	if len(m2.Steps) > 0 {
		merged.Steps = m2.Steps[:]
	} else {
		merged.Steps = m1.Steps[:]
	}
	if len(m2.Android.Apk) > 0 {
		merged.Android.Apk = m2.Android.Apk
	} else {
		merged.Android.Apk = m1.Android.Apk
	}
	if len(m2.Android.ApkInstrumentation) > 0 {
		merged.Android.ApkInstrumentation = m2.Android.ApkInstrumentation
	} else {
		merged.Android.ApkInstrumentation = m1.Android.ApkInstrumentation
	}
	if len(m2.DeviceGroupNames) > 0 {
		merged.DeviceGroupNames = m2.DeviceGroupNames[:]
	} else {
		merged.DeviceGroupNames = m1.DeviceGroupNames[:]
	}
	return merged
}

// IsRunnable returns true and nil if the BuildManifest is properly configured
// to run, and returns false and an error otherwise. For example, if a BuildManifest
// has no DeviceGroupNames, it cannot be run.
func (manifest *BuildManifest) IsRunnable() (bool, error) {
	if len(manifest.Android.Apk) == 0 || len(manifest.Android.ApkInstrumentation) == 0 {
		return false, fmt.Errorf("Missing Android apk or apk_instrumentation")
	}
	if len(manifest.DeviceGroupNames) == 0 {
		return false, fmt.Errorf("Missing devicegroups")
	}
	return true, nil
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

// IsValid returns true and nil if the Config is valid, and returns false and
// an error otherwise. Valid does not necessesarily mean that all builds will
// work, it just means that the config does not have any obvious errors.
//
// For example, a config with no `defaults` and no config for the `master`
// branch could be valid. But if you try to run a build on the `master` branch
// the build will still fail, because there is no config available.
//
// See also BuildManifest#IsRunnable().
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

// BranchManifest returns a BuildManifest for the given branch name, by starting
// from the Defaults manifest (if any) and merging branch-specific overrides on
// top of it.
func (config *Config) BranchManifest(branch string) *BuildManifest {
	manifest := MergeManifests(&BuildManifest{}, &config.Defaults)
	if branchOverrides, ok := config.Branches[branch]; ok {
		manifest = MergeManifests(manifest, &branchOverrides)
	}
	return manifest
}
