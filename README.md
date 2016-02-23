#DLite

The simplest way to use Docker on OSX. [![Build Status](https://travis-ci.org/nlf/dlite.svg?branch=master)](https://travis-ci.org/nlf/dlite)

##Thanks

DLite leverages [xhyve](https://github.com/mist64/xhyve) through the [libxhyve](https://github.com/TheNewNormal/libxhyve) Go bindings for virtualization. Without these projects and the people behind them, this project wouldn't exist.

## Installation

There are several ways to install dlite. You may install it with [Homebrew](http://brew.sh/), download it from github or compile it yourself.

### Download

1. Download the latest binary from the [releases page](https://github.com/nlf/dlite/releases) and put it somewhere in your path, or
2. Install it with homebrew: `brew install dlite`
3. If you have a working [Go development environment](https://golang.org/doc/install) you can build dlite from source by running:

    ```
    go get github.com/nlf/dlite
    ```

    After that you need to compile it (run `make dlite` in the `src/github.com/nlf/dlite` dir.)

### Initialization

To create the necessary files and a launchd agent which manages the process, simply run

```
sudo dlite install
```

See the output of `sudo dlite install --help` for additional options, like changing number of CPUs, Disk Size, et cetera.

After you've installed, you need to start the process:

```
dlite start
```

DLite will start automatically upon logging in as well.

##Updating DLite

The DLite app itself can be updated by running `dlite stop`, installing the updated binary, and then running `dlite start`.

To install the updated binary with Homebrew simply run `brew upgrade dlite`.

If you update dlite, you probably want to update your VM as well:

##Updating your VM

It's possible to update your virtual machine without having to rebuild it entirely. To do so, run the following commands

```
dlite stop
dlite update
dlite start
```

##Usage

Just use Docker. DLite creates a `/var/run/docker.sock` in your host operating system.

When opening ports in your docker containers, connect to `local.docker` instead of `localhost`. Everything else should just work™

Note that the `local.docker` hostname is configurable by passing the `-n` flag to the install command, as in `sudo dlite install -n docker.dev`

If you need to SSH to the VM for whatever reason, `ssh docker@local.docker` should do the trick.

## Seamless routing

By default, Docker creates a virtual interface named `docker0` on the host machine that forwards packets between any other network interface.

However, on OSX, this means you are not able to access the Docker network directly. To be able to do so, you need add a route and allow traffic from any host on the interface that connects to the VM.

Run the following commands on your OSX machine:

```sh
sudo route -n add 172.17.0.0/16 local.docker
DOCKER_INTERFACE=$(route get local.docker | grep interface: | cut -f 2 -d: | tr -d ' ')
DOCKER_INTERFACE_MEMBERSHIP=$(ifconfig ${DOCKER_INTERFACE} | grep member: | cut -f 2 -d: | cut -c 2-4)
sudo ifconfig "${DOCKER_INTERFACE}" -hostfilter "${DOCKER_INTERFACE_MEMBERSHIP}"
```

See if it works by pinging the IP of any running container (assuming you have at least one):

```sh
docker inspect -f '{{.NetworkSettings.IPAddress}}' $(docker ps -q | head -1)
```

For now, you may include this on a profile script if desired in case you need to repeat the same steps. Unless you reboot your OSX machine, you shouldn't need to run this often.

### Service Discovery via DNS

Now that you've got transparent routing to your Docker containers, it is time to improve their accessibility by using a DNS server.

First, let's start by configuring the default Docker service DNS server to IP where the DNS server will run (`172.17.42.1`). Currently, this requires SSH'ing into the VM and editing `/etc/default/docker`, but this likely to change in the [future](https://github.com/nlf/dlite/issues/90).

Add the DNS server static IP via `--bip` and `--dns`:

```
DOCKER_ARGS="-H unix:///var/run/docker.sock -H tcp://0.0.0.0:2375 -s btrfs --bip=172.17.42.1/24 --dns=172.17.42.1"
```

Then edit the rc script `/etc/init.d/S51docker` that starts Docker and run the DNS server on startup:

Add after `[ $? = 0 ] && echo "OK" || echo "FAIL"`:

```sh
/usr/bin/docker run -d -v /var/run/docker.sock:/var/run/docker.sock --name dnsdock -p 172.17.42.1:53:53/udp tonistiigi/dnsdock
```

Add before `start-stop-daemon -K -q -p /var/run/docker.pid`:

```sh
/usr/bin/docker rm --force dnsdock
```

After all of this is done just restart the Docker service inside the VM to get the DNS server up and running:

```sh
sudo /etc/init.d/S51docker restart
```

Lastly, configure OSX so that all `.docker` requests are forwarded to Docker's DNS server. Since routing has already been taken care of, just create a custom resolver under `/etc/resolver/docker` with the following content:

```
nameserver 172.17.42.1
```

Then restart OSX's own DNS server:

```
sudo killall -HUP mDNSResponder
```

Check if the DNS server is working as expected by querying a running image:

```sh
dig <image>.docker @172.17.42.1
```

You should see a Docker network IP resolved correctly:

```
;; ANSWER SECTION:
<image>.docker.	0	IN	A	172.17.42.3
```

#### Troubleshooting DNS

It usually takes some time to adapt to the DNS naming scheme of `dnsdock`, so if you'd like see which DNS names are being registered in real time, just follow the `dnsdock` logs:

`docker logs --follow dnsdock`

##Troubleshooting

A common cause of the virtual machine failing to start is conflicting entries in your `/etc/exports` file. Edit the file and see if any other process has an export that conflicts with the one DLite added (it will have comments before and after it, making it easy to identify). If they do, remove the conflicting entry and try starting the service again. Note that dlite adds its export when it is started, not when it is installed, so make sure to either clean your exports file or specify a shared directory that doesn't conflict with existing shares when you install.

If `docker` cli commands hang, there's a good chance that you have a stale entry in your `/etc/hosts` file. Run `dlite stop`, then use sudo to edit your `/etc/hosts` file and remove any entries that end with `# added by dlite`. Save the hosts file and run `dlite start` and try again.

Note that `launchctl` commands appear to not work correctly when run inside tmux. If you are a tmux user and are having problems, try starting the service outside of your tmux session.

##Caveats

DLite depends on [xhyve](https://github.com/mist64/xhyve) which only works on OSX versions 10.10 (Yosemite) or newer. You also need a fairly recent mac. You can tell if your computer is new enough by running `sysctl kern.hv_support` in a terminal. If you see `kern.hv_support: 1` as a response, you're good to go. If not, unfortunately your computer is too old to leverage the hypervisor framework and DLite won't work for you.

Xhyve, and therefore DLite, does not support sparse disk images. This means that when you create a virtual machine with DLite the *full size* of the image must be allocated up front. There is ongoing work to support sparse images in xhyve, and once that support lands DLite will be able to take advantage of it. See [xhyve#80](https://github.com/mist64/xhyve/pull/80), [xhyve#82](https://github.com/mist64/xhyve/pull/82), and [xhyve-xyz/xhyve#1](https://github.com/xhyve-xyz/xhyve/pull/1) for more information.

DLite is *not* secured via TLS. If that's important to you for local development, look elsewhere.

DLite is most definitely *not* recommended for any kind of production use.
