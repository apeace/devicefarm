# DeviceFarm

## Install

Go to [Releases](https://github.com/ride/devicefarm/releases/) and download the latest binary for your platform. Then
move it to `/usr/local/bin/devicefarm` and run:

```
chmod +x /usr/local/bin/devicefarm
cp ~/Dropbox\ \(Ride\)/Engineering/DeviceFarm/devicefarm.json ~/.devicefarm.json
```

## Setup

You should have a `devicefarm.yml` config in your repo, [like this one](./config/testdata/config.yml).

## Features

Only Android instrumentation tests are supported at the moment.

### Run UI tests on DeviceFarm

You will need [access to AWS](https://github.com/ride/devops/blob/master/docs/aws-access.md) to view the results of the tests, view screenshots, etc.

```bash
# go into your android project
$ cd /path/to/ride-app-android/

# run tests. the output is a URL you can visit to view your test results.
$ devicefarm run
>> Dir: /Users/apeace/code/MarkdownPreview, Config: devicefarm.yml, Branch: devicefarm
>> Running build... (silencing output)
$ ./gradlew assembleDebug
$ ./gradlew assembleAndroidTest
>> Build complete
>> Device Pool: everything (9 devices)
>> Uploading files...
/Users/apeace/code/MarkdownPreview/app/build/outputs/apk/app-debug.apk
/Users/apeace/code/MarkdownPreview/app/build/outputs/apk/app-debug-androidTest-unaligned.apk
>> Waiting for files to be processed...
>> Creating test run...
https://us-west-2.console.aws.amazon.com/devicefarm/home?region=us-west-2#/projects/1124416c-bfb2-4334-817c-e211ecef7dc0/runs/a07ca17f-d8ec-4adf-8e36-dc776b847705
```

### Update device pools

If you update your device pools in `devicefarm.yml` you can run tests on different sets of devices. See "List devices" below for how to find device names.

```bash
# go into your android project
$ cd /path/to/ride-app-android/

# if you checkout a branch, device pools will be tied only to that branch
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
(arn=device:306ABA42C96044ED9AC3EE8684B56C54) Apple iPad Mini 1st Gen
(arn=device:5F9CEB47606A4709879003E11BEAFB08) Samsung Galaxy Tab 2 10.1 (WiFi)
...

# list all devices matching "mini"
$ devicefarm devices "mini"
(arn=device:306ABA42C96044ED9AC3EE8684B56C54) Apple iPad Mini 1st Gen
(arn=device:5C748437DC1C409EA595B98B1D7A8EDD) Samsung Galaxy S3 Mini (AT&T)
...

# list all android devices matching "mini"
# there is also an --ios flag :)
$ devicefarm devices --android "mini"
(arn=device:5C748437DC1C409EA595B98B1D7A8EDD) Samsung Galaxy S3 Mini (AT&T)
(arn=device:20766AF83D3A4FEF977643BFCDC2CE3A) Samsung Galaxy S4 mini (Verizon)
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
   development

COMMANDS:
    run		Create test run based on YAML config
    build	Run local build based on YAML config
    devices	Search device farm devices

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

## Limitations, bugs & future work

See our [issue tracker](https://github.com/ride/devicefarm/issues) for known bugs, improvements, and maintenance work.

Right now only Android instrumentation tests are supported, and only command-line usage is supported.

Future work:

 * v1.x - Integrate with Github.
 * v2.0 - Rework `devicefarm.yml` format, introduce support for iOS and additional test types.

## Docs

 * [Development](./docs/development.md)
