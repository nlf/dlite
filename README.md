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
dlite start
```

as your user to start the process. DLite will start automatically upon logging in as well.

##Usage

Just use Docker. DLite creates a `/var/run/docker.sock` in your host operating system.

When opening ports in your docker containers, connect to `local.docker` instead of `localhost`. Everything else should just work™

Note that the `local.docker` hostname is configurable by passing the `-n` flag to the install command, as in `sudo dlite install -n docker.dev`

If you need to SSH to the VM for whatever reason, `ssh docker@local.docker` should do the trick.

##Troubleshooting

A common cause of the virtual machine failing to start is conflicting entries in your `/etc/exports` file. Edit the file and see if any other process has an export that conflicts with the one DLite added (it will have comments before and after it, making it easy to identify). If they do, remove the conflicting entry and try starting the service again.

If `docker` cli commands hang, there's a good chance that you have a stale entry in your `/etc/hosts` file. Run `dlite stop`, then use sudo to edit your `/etc/hosts` file and remove any entries that end with `# added by dlite`. Save the hosts file and run `dlite start` and try again.

Note that `launchctl` commands appear to not work correctly when run inside tmux. If you are a tmux user and are having problems, try starting the service outside of your tmux session.

##Caveats

DLite depends on [xhyve](https://github.com/mist64/xhyve) which only works on OSX versions 10.10 (Yosemite) or newer.

Xhyve, and therefor DLite, does not support sparse disk images. This means that when you create a virtual machine with DLite the *full size* of the image must be allocated up front. There is ongoing work to support sparse images in xhyve, and once that support lands DLite will be able to take advantage of it. See [xhyve#80](https://github.com/mist64/xhyve/pull/80), [xhyve#82](https://github.com/mist64/xhyve/pull/82), and [xhyve-xyz/xhyve#1](https://github.com/xhyve-xyz/xhyve/pull/1) for more information.

DLite is *not* secured via TLS. If that's important to you for local development, look elsewhere.

DLite is most definitely *not* recommended for any kind of production use.
