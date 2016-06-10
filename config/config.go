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
	# samsung_s4, samsung_s5, and everything. The everything pool
	# simply includes both the other pools.
	#
	# This property is REQUIRED, must have at least one pool defined,
	# and each pool must have at least one device.
	devicepool_definitions:
	samsung_s4:
	  - (arn=device:D1C28D6B913C479399C0F594E1EBCAE4) Samsung Galaxy S4 (AT&T)
	  - (arn=device:33F66BE404B543669978079E905F8637) Samsung Galaxy S4 (Sprint)
	  - (arn=device:D45C750161314335924CE0B9B7D2558E) Samsung Galaxy S4 (T-Mobile)

	samsung_s5:
	  - (arn=device:5CC0164714304CBF81BB7B7C03DFC1A1) Samsung Galaxy S5 (AT&T)
	  - (arn=device:18E28478F1D54525A15C2A821B6132FA) Samsung Galaxy S5 (Sprint)
	  - (arn=device:5931A012CB1C4E68BD3434DF722ADBC8) Samsung Galaxy S5 (T-Mobile)

	everything:
	  - +samsung_s4
	  - +samsung_s5

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

	  # The tests you want to run. Each key is a name (e.g. my_instrumentation)
	  # and the value is a map with type, plus whatever other fields that test
	  # type needs. TODO: doc for test types.
	  tests:
	    my_instrumentation:
	      type: android_instrumentation
	      app_apk: ./path/to/build.apk
	      instrumentation_apk: ./path/to/instrumentation.apk

	  # The device pool name that tests should be run on.
	  devicepool: samsung_s4

	  # DEPRECATED. Being removed in 2.0.
	  # This used to be how you would specify Android instrumentation tests.
	  # Now you do it with the tests field.
	  android:
	    apk: ./path/to/build.apk
	    apk_instrumentation: ./path/to/instrumentation.apk

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
	"fmt"
	"github.com/ride/devicefarm/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"sort"
)

// An AndroidConfig specifies the location of APKs after running the build steps.
// DEPRECATED. See ConvertDeprecatedAndroid()
type AndroidConfig struct {
	Apk                string `yaml:"apk"`
	ApkInstrumentation string `yaml:"apk_instrumentation"`
}

// A BuildManifest specifies the whole configuration for a build: the steps to
// perform the build, the location of Android APKs, and the DevicePool names
// to run on.
type BuildManifest struct {
	Steps      []string                     `yaml:"build"`
	Tests      map[string]map[string]string `yaml:"tests"`
	DevicePool string                       `yaml:"devicepool"`
	// DEPRECATED, see ConvertDeprecatedAndroid()
	Android AndroidConfig `yaml:"android"`
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
	if len(m2.Tests) > 0 {
		// TODO clone?
		merged.Tests = m2.Tests
	} else {
		merged.Tests = m1.Tests
	}
	return merged
}

// ConvertDeprecatedAndroid converts the Android field into tests in the Tests
// field. The Android field is deprecated in favor of Tests.
func (manifest *BuildManifest) ConvertDeprecatedAndroid() {
	if len(manifest.Android.Apk) == 0 {
		return
	}
	test := map[string]string{
		"type":                "android_instrumentation",
		"app_apk":             manifest.Android.Apk,
		"instrumentation_apk": manifest.Android.ApkInstrumentation,
	}
	if len(manifest.Tests) == 0 {
		manifest.Tests = map[string]map[string]string{}
	}
	manifest.Tests["deprecated_auto_instrumentation"] = test
	manifest.Android.Apk = ""
	manifest.Android.ApkInstrumentation = ""
}

// IsRunnable returns true and nil if the BuildManifest is properly configured
// to run, and returns false and an error otherwise. For example, if a BuildManifest
// has no DevicePool, it cannot be run.
func (manifest *BuildManifest) IsRunnable() (bool, error) {
	if len(manifest.Tests) == 0 {
		return false, fmt.Errorf("Missing tests")
	}
	if len(manifest.DevicePool) == 0 {
		return false, fmt.Errorf("Missing devicepool")
	}
	return true, nil
}

// A Config specifies configuration for a particular repo: the names of DevicePools,
// the default BuildManifest, and override BuildManifests for particular branches.
type Config struct {
	ProjectArn            string                   `yaml:"project_arn"`
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
	// build steps are optional
	if config.Defaults.Steps == nil {
		config.Defaults.Steps = []string{}
	}
	for _, manifest := range config.Branches {
		if manifest.Steps == nil {
			manifest.Steps = []string{}
		}
	}
	// support deprecated android config
	config.Defaults.ConvertDeprecatedAndroid()
	for _, config := range config.Branches {
		config.ConvertDeprecatedAndroid()
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
	if !util.ArnRegexp.MatchString(config.ProjectArn) {
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
			return nil, fmt.Errorf("DevicePool has no items: %v", name)
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
				return nil, fmt.Errorf("Blank DevicePool item in: %v", name)
			}
			if item == "+"+name {
				return nil, fmt.Errorf("DevicePool circular dependency: %v", name)
			}
			if seen[item] {
				continue
			}
			seen[item] = true
			if string(item[0]) == "+" {
				poolRef := item[1:]
				refItems, ok := defs[poolRef]
				if !ok {
					return nil, fmt.Errorf("DevicePool definition does not exist: %v", poolRef)
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

// DeviceArns takes a list of devices from a config file, and returns a list
// of full ARNs.
func DeviceArns(devices []string) ([]string, error) {
	itemRegexp := regexp.MustCompile("\\(arn=([^\\)]+)\\)\\s*(.+)\\s*")
	parsed := []string{}
	for _, item := range devices {
		match := itemRegexp.FindStringSubmatch(item)
		if len(match) < 3 {
			return nil, fmt.Errorf("Invalid device: %v", item)
		}
		arn := util.Arn{
			Partition: "aws",
			Service:   "devicefarm",
			Region:    "us-west-2",
			AccountId: "",
			Resource:  match[1],
		}
		parsed = append(parsed, arn.String())
	}
	return parsed, nil
}
