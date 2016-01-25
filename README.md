#DLite

The simplest way to use Docker on OSX.

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

##Usage

Just use Docker. DLite creates a `/var/run/docker.sock` in your host operating system.

When opening ports in your docker containers, connect to `local.docker` instead of `localhost`. Everything else should just work™

If you need to SSH to the VM for whatever reason, `ssh docker@local.docker` should do the trick.

##Caveats

DLite depends on [xhyve](https://github.com/mist64/xhyve) which only works on OSX versions 10.10 (Yosemite) or newer.

DLite is *not* secured via TLS. If that's important to you for local development, look elsewhere.

DLite is most definitely *not* recommended for any kind of production use.
