#DLite

The simplest way to use Docker on OSX.

##Installation

```
git clone git://github.com/nlf/dlite
cd dlite
./install
```

Note that installation uses [homebrew](http://brew.sh) for both xhyve and socat. If you don't already use homebrew, you should really install it first. If you don't want to, install xhyve and socat to `/usr/local/bin` yourself.

##Usage

Just use Docker. DLite creates a `/var/run/docker.sock` in your host operating system.

When opening ports in your docker containers, connect to `local.docker` instead of `localhost`. Everything else should just work™

If you need to SSH to the VM for whatever reason, `ssh docker@local.docker` should do the trick.

##Configuration

The number of CPUs and amount of memory allocated to the virtual machine is configurable. After running the install script, simply edit `/etc/dlite.conf` and you may add lines similar to the following:

```
DLITE_CPUS=2
DLITE_MEM=2G
```

The default is 1 CPU and 1GB of memory. Do *not* delete or change the `DLITE_UUID` setting in that file.

After making your changes:

```
sudo launchctl stop local.dlite
# I would recommend checking the output of `ps aux | grep xhyve` here
# and wait until the virtual machine has actually stopped
sudo launchctl start local.dlite
```

##Caveats

DLite depends on [xhyve](https://github.com/mist64/xhyve) which only works on OSX versions 10.10 (Yosemite) or newer.

DLite is *not* secured via TLS. If that's important to you for local development, look elsewhere.

DLite is most definitely *not* recommended for any kind of production use.
