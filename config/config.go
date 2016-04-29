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

		# The device pool name that tests should be run on.
		devicepool: a_few_devices

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
		devicepool: everything

*/
package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sort"
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
	Steps      []string      `yaml:"build"`
	Android    AndroidConfig `yaml:"android"`
	DevicePool string        `yaml:"devicepool"`
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
	if len(m2.DevicePool) > 0 {
		merged.DevicePool = m2.DevicePool[:]
	} else {
		merged.DevicePool = m1.DevicePool[:]
	}
	return merged
}

// IsRunnable returns true and nil if the BuildManifest is properly configured
// to run, and returns false and an error otherwise. For example, if a BuildManifest
// has no DevicePool, it cannot be run.
func (manifest *BuildManifest) IsRunnable() (bool, error) {
	if len(manifest.Android.Apk) == 0 || len(manifest.Android.ApkInstrumentation) == 0 {
		return false, fmt.Errorf("Missing Android apk or apk_instrumentation")
	}
	if len(manifest.DevicePool) == 0 {
		return false, fmt.Errorf("Missing devicepool")
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
	// TODO: Use regex validator instead?
	if len(config.Arn) == 0 {
		return false, fmt.Errorf("project_arn is required")
	}
	if len(config.DevicePoolDefinitions) == 0 {
		return false, fmt.Errorf("devicepools must have at least one pool")
	}
	_, err := config.FlatDevicePoolDefinitions()
	if err != nil {
		return false, err
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

// FlatDevicePoolDefinitions returns a map of device pool definitions with the "+"
// references flattened, so that each device pool is just a list of device names.
// It also sorts each list of device names. It detects circular references or
// non-existant references and returns an error in those cases.
func (config *Config) FlatDevicePoolDefinitions() (map[string][]string, error) {
	defs := config.DevicePoolDefinitions
	flat := map[string][]string{}
	for name, items := range defs {
		if len(items) == 0 {
			return nil, errors.New("DevicePool has no items: " + name)
		}
		seen := map[string]bool{}
		deviceNames := []string{}
		queue := items[:]
		for {
			if len(queue) == 0 {
				break
			}
			item := queue[0]
			queue = queue[1:]
			if len(item) == 0 {
				return nil, errors.New("Blank DevicePool item in: " + name)
			}
			if item == "+"+name {
				return nil, errors.New("DevicePool circular dependency: " + name)
			}
			if seen[item] {
				continue
			}
			seen[item] = true
			if string(item[0]) == "+" {
				poolRef := item[1:]
				refItems, ok := defs[poolRef]
				if !ok {
					return nil, errors.New("DevicePool definition does not exist: " + poolRef)
				}
				queue = append(queue, refItems...)
				continue
			}
			deviceNames = append(deviceNames, item)
		}
		sort.Strings(deviceNames)
		flat[name] = deviceNames
	}
	return flat, nil
}
