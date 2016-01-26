#DLite

The simplest way to use Docker on OSX.

##Thanks

DLite leverages [xhyve](https://github.com/mist64/xhyve) through the [libxhyve](https://github.com/TheNewNormal/libxhyve) Go bindings for virtualization. Without these projects and the people behind them, this project wouldn't exist.

##Installation

Download the latest binary from the [releases page](https://github.com/nlf/dlite/releases) and put it somewhere in your path, then:

```
sudo dlite install
```

See the output of `sudo dlite install --help` for additional options.

This will create the necessary files and a launchd agent to manage the process. After you've installed, run:

```
launchctl load ~/Library/LaunchAgents/local.dlite.plist
```

as your user to start the process. DLite will start automatically upon logging in as well.

Then to start the virtual machine run:

```
launchctl start local.dlite
```

##Usage

Just use Docker. DLite creates a `/var/run/docker.sock` in your host operating system.

When opening ports in your docker containers, connect to `local.docker` instead of `localhost`. Everything else should just work™

Note that the `local.docker` hostname is configurable by passing the `-n` flag to the install command, as in `sudo dlite install -n docker.dev`

If you need to SSH to the VM for whatever reason, `ssh docker@local.docker` should do the trick.

##Caveats

DLite depends on [xhyve](https://github.com/mist64/xhyve) which only works on OSX versions 10.10 (Yosemite) or newer.

Xhyve, and therefor DLite, does not support sparse disk images. This means that when you create a virtual machine with DLite the *full size* of the image must be allocated up front. There is ongoing work to support sparse images in xhyve, and once that support lands DLite will be able to take advantage of it. See [xhyve#80](https://github.com/mist64/xhyve/pull/80), [xhyve#82](https://github.com/mist64/xhyve/pull/82), and [xhyve-xyz/xhyve#1](https://github.com/xhyve-xyz/xhyve/pull/1) for more information.

DLite is *not* secured via TLS. If that's important to you for local development, look elsewhere.

DLite is most definitely *not* recommended for any kind of production use.
