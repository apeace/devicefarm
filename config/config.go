/*

Package config provides data structures and functions to configure Device Farm testing.

YAML Config

Repositories using the devicefarm tool will provide a devicefarm.yml file in the
repository root. The config.New() function parses this config from a file into a
config.Config struct:

	config, err := config.New("path/to/config.yml")

Here is an annotated example of what a config file should look like:

	# Project ARN. This property is REQUIRED.
	project_arn: arn:aws:devicefarm:us-west-2:026109802893:project:1124416c-bfb2-4334-817c-e211ecef7dc0

	# Device Pool definitions. this block defines three Device Pools:
	# a_few_devices, samsung_s4_s5, and everything. The everything pool
	# simply includes both the other pools.
	#
	# This property is REQUIRED, must have at least one pool defined,
	# and each pool must have at least one device.
	devicepools:
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

		# The device pool names that tests should be run on.
		devicepools:
		  - a_few_devices

	# Branches defines overrides for particular branches. For each branch,
	# it accepts the same properties as `defaults`. Branch configs will be
	# merged with `defaults` so that the specified properties override the
	# same properties from `defaults`. In this example, only the `devicepools`
	# property will be overridden for the `master` branch.
	#
	# This property is OPTIONAL, but building will fail on a particular branch
	# unless a full definition is available for that branch.
	branches:
	  master:
		devicepools:
		  - everything

*/
package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// An AndroidConfig specifies the location of APKs after running the build steps.
type AndroidConfig struct {
	Apk                string `yaml:"apk"`
	ApkInstrumentation string `yaml:"apk_instrumentation"`
}

// A BuildManifest specifies the whole configuration for a build: the steps to
// perform the build, the location of Android APKs, and the DevicePool names
// to run on.
type BuildManifest struct {
	Steps           []string      `yaml:"build"`
	Android         AndroidConfig `yaml:"android"`
	DevicePoolNames []string      `yaml:"devicepools"`
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
	if len(m2.DevicePoolNames) > 0 {
		merged.DevicePoolNames = m2.DevicePoolNames[:]
	} else {
		merged.DevicePoolNames = m1.DevicePoolNames[:]
	}
	return merged
}

// IsRunnable returns true and nil if the BuildManifest is properly configured
// to run, and returns false and an error otherwise. For example, if a BuildManifest
// has no DevicePoolNames, it cannot be run.
func (manifest *BuildManifest) IsRunnable() (bool, error) {
	if len(manifest.Android.Apk) == 0 || len(manifest.Android.ApkInstrumentation) == 0 {
		return false, fmt.Errorf("Missing Android apk or apk_instrumentation")
	}
	if len(manifest.DevicePoolNames) == 0 {
		return false, fmt.Errorf("Missing devicepools")
	}
	return true, nil
}

// A Config specifies configuration for a particular repo: the names of DevicePools,
// the default BuildManifest, and override BuildManifests for particular branches.
type Config struct {
	Arn                   string                   `yaml:"project_arn"`
	DevicePoolDefinitions map[string][]string      `yaml:"devicepool_definitions"`
	Defaults              BuildManifest            `yaml:"defaults"`
	Branches              map[string]BuildManifest `yaml:"branches"`
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
	// TODO: Use regex valid instead?
	if len(config.Arn) == 0 {
		return false, fmt.Errorf("project_arn is required")
	}
	if len(config.DevicePoolDefinitions) == 0 {
		return false, fmt.Errorf("devicepools must have at least one pool")
	}
	for name, devicepool := range config.DevicePoolDefinitions {
		if len(devicepool) == 0 {
			return false, fmt.Errorf("devicepool %s has no devices", name)
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
