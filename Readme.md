# DeviceFarm

A command-line tool for using [AWS Device Farm](https://aws.amazon.com/device-farm/).
Key features:

 * Easily manage device pools through a config file in your repo.
 * Run your build and upload artifacts & tests to Device Farm.
 * Gives you a URL to jump straight to your test results in the AWS Console.

## Install

`devicefarm` is distributed as a binary. Simply go to
[Releases](https://github.com/ride/devicefarm/releases/) and download the
latest binary for your platform.

Example for OS X:

```
# first, download the latest release. then...
mv ~/Downloads/devicefarm_darwin_amd64 /usr/local/bin/devicefarm
chmod +x /usr/local/bin/devicefarm

# that's it! you now have devicefarm
devicefarm --version
```

## Setup

**First,** you will need an AWS user who has permission to access Device Farm.
It is recommended to setup a separate user for this purpose. Once you have that,
create a file `~/.devicefarm.json` with contents like this:

```
{
  "AWS_ACCESS_KEY_ID": "... your key ...",
  "AWS_SECRET_ACCESS_KEY": "... your secret ..."
}
```

**Second,** you should setup a `devicefarm.yml` config in your repo,
[like this one](./config/testdata/config.yml).

## Features

Only Android instrumentation tests are supported at the moment. See
[future work](#limitations-bugs--future-work).

### Run instrumentation tests on Device Farm

```bash
# go into your android project
$ cd /path/to/android-app/

# run tests. the output is a URL you can visit to view your test results.
$ devicefarm run
...
https://us-west-2.console.aws.amazon.com/devicefarm/home?region=us-west-2#/projects/1124416c-bfb2-4334-817c-e211ecef7dc0/runs/a07ca17f-d8ec-4adf-8e36-dc776b847705
```

### Update device pools

If you update your device pools in `devicefarm.yml` you can run tests on
different sets of devices. See "List devices" below for how to find device
identifiers.

```bash
# go into your android project
$ cd /path/to/android-app/

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

See our [issue tracker](https://github.com/ride/devicefarm/issues) for known
bugs, improvements, and maintenance work.

Right now only Android instrumentation tests are supported. As part of our
[2.0 Milestone](https://github.com/ride/devicefarm/milestones/2.0) we'll be
adding:

 * Support for all test types (including iOS and web).
 * Commands to help setup AWS credentials and `devicefarm.yml` config.
 * A Homebrew tap so OS X users can install and update using `brew`.
 * Polishing existing commands and config to make them easier to use.

## Docs

 * [Development](./docs/development.md)

## License

Copyright © 2016 Ride Group, Inc. [MIT](./LICENSE)
