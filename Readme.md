# DeviceFarm

Alpha release: v0.0.3

## Install

[Click here](https://github.com/ride/devicefarm/releases/download/v0.0.3/devicefarm-osx) to download the `devicefarm` binary for OS X. Then complete the install by running:

```
mv ~/Downloads/devicefarm-osx /usr/local/bin/devicefarm
chmod +x /usr/local/bin/devicefarm
cp ~/Dropbox\ \(Ride\)/Engineering/DeviceFarm/devicefarm.json ~/.devicefarm.json
```

## Setup

You should have a `devicefarm.yml` config in your repo, [like this one](https://github.com/ride/ride-app-android/pull/1297).

## Features

Only Android instrumentation tests are supported at the moment.

### Run UI tests on DeviceFarm

You will need [access to AWS](https://github.com/ride/devops/blob/master/docs/aws-access.md) to view the results of the tests, view screenshots, etc.

```bash
# go into your android project
$ cd /path/to/ride-app-android/

# run tests. after a couple of minutes, you will get a URL
$ devicefarm run
2016/04/29 09:13:54 >> Dir: .
2016/04/29 09:13:54 >> Config: devicefarm.yml
2016/04/29 09:13:54 >> Branch: devicefarm
2016/04/29 09:13:54 >> Running build... (silencing output)
2016/04/29 09:13:54 $ ./gradlew assembleDevelopment
2016/04/29 09:15:19 $ ./gradlew assembleAndroidTest
2016/04/29 09:17:54 >> Build complete
2016/04/29 09:17:55 >> Device Pool: random (2 devices)
2016/04/29 09:17:55 >> Uploading files...
2016/04/29 09:17:55 /Users/apeace/code/ride-app-android/app/build/outputs/apk/app-development-debug.apk
2016/04/29 09:18:01 /Users/apeace/code/ride-app-android/app/build/outputs/apk/app-development-debug-androidTest-unaligned.apk
2016/04/29 09:18:01 >> Waiting for files to be processed...
2016/04/29 09:19:01 >> Creating test run...
2016/04/29 09:19:02 https://us-west-2.console.aws.amazon.com/devicefarm/home?region=us-west-2#/projects/018d4112-c644-4e53-aa37-e83415f83f9f/runs/bd7c001b-dc93-4695-958f-8abbae775532
```

### Update device pools

If you update your device pools in `devicefarm.yml` you can run tests on different sets of devices. See "List devices" below for how to find device names.

```bash
# go into your android project
$ cd /path/to/ride-app-android/

# if you checkout a branch, device pools will be tied only to that branch
# you can just use master/develop if you want, but you may have conflicts
# with others if they test on the same branch at the same time
$ git checkout -b my-test-branch

# now edit devicefarm.yml...
# you can change the default config, OR specify overrides for your specific branch

# now just create another run
$ devicefarm run
```

### List devices

This way you can find devices you want to add to your device pools.

```bash
# list all devices
$ devicefarm devices
Apple iPad Mini 2
Samsung Galaxy Tab 3 10.1 (WiFi)
...

# list all devices matching "mini"
$ devicefarm devices "mini"
Apple iPad Mini 1st Gen
Samsung Galaxy S3 Mini (AT&T)
...

# list all android devices matching "mini"
# there is also an --ios flag :)
$ devicefarm devices --android "mini"
Samsung Galaxy S3 Mini (AT&T)
Samsung Galaxy S4 mini (Verizon)
```

### More

You can run `devicefarm help` or `devicefarm help COMMAND` to get help info:

```bash
$ devicefarm help
NAME:
   devicefarm - Run UI tests in AWS Device Farm

USAGE:
   devicefarm [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
    run		Create test run based on YAML config
    build	Run local build based on YAML config
    devicepools	Sync devicepools with your YAML config
    devices	Search device farm devices

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

## Limitations, bugs & future work

This is an alpha release for testing. Work is ongoing to stabilize this feature set for version 1.0.

See our [issue tracker](https://github.com/ride/devicefarm/issues) for known bugs, improvements, and maintenance work.

Right now only Android instrumentation tests are supported, and only command-line usage is supported.

Future work:

 * [v1.0](https://github.com/ride/devicefarm/milestones/v1.0.0) - Stabilize the existing feature set, higher level of unit-test coverage, automated releases (continuous delivery through CI).
 * v1.1 - Update to latest version with `devicefarm update`.
 * v1.x - Integrate with Github.
 * v2.0 - Rework `devicefarm.yml` format, introduce support for iOS and additional test types.

## Docs

 * [Development](./docs/development.md)
