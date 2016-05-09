# DLite

## THIS BRANCH IS CURRENTLY BROKEN. DO NOT USE UNLESS YOU KNOW WHAT YOU'RE DOING

The simplest way to use Docker on OSX.

## Usage
Just use Docker. DLite creates a `/var/run/docker.sock` in your host operating system.

When opening ports in your docker containers, connect to `local.docker` instead of `localhost`. Everything else should just work™

Note that the `local.docker` hostname is configurable by passing the `-n` flag to the install command, as in `sudo dlite install -n docker.dev`

If you need to SSH to the VM for whatever reason, `ssh docker@local.docker` should do the trick.

## Seamless routing
By passing the `-r` or `--route` flag to the install command, or editing your config file to set `"route": true`, DLite will set up routing tables to allow you to directly access your containers on the 172.17.0.0/16 network.

Some events cause OSX to clear the routing table. If you find you're unable to reach your containers, run the `dlite route` command to readd the routing entries.

You can find the IP of an individual container by running `docker inspect -f '{{.NetworkSettings.IPAddress}}' <container_name>` where `<container_name>` is the name of the container you wish to connect to.

### Service Discovery via DNS
If you wish to use DNS records to improve your containers accessibility, you can easily do so by leveraging the [dnsdock](https://github.com/tonistiigi/dnsdock) container.

Note that doing so, however, will cause docker to ignore any DNS server you configured in DLite. If you use a non-standard DNS server, add `-nameserver="8.8.8.8:53"` to the very end of the command below, replacing `8.8.8.8` with your desired DNS server.

First, run the dnsdock service:

```sh
docker run -d -v /var/run/docker.sock:/var/run/docker.sock --name dnsdock --restart always -p 172.17.0.1:53:53/udp tonistiigi/dnsdock
```

Next, edit your config file for DLite via `dlite config`. Set the value of the `"extra"` option to `"--bip=172.17.0.1/24 --dns=172.17.0.1"` and exit your editor.

Lastly, configure OSX so that all `.docker` requests are forwarded to Docker's DNS server. Since routing has already been taken care of, just create a custom resolver:

```sh
sudo mkdir -p /etc/resolver
echo "nameserver 172.17.0.1" | sudo tee /etc/resolver/docker
```

Then restart OSX's own DNS server (10.10.4+):

```
sudo dscacheutil -flushcache
sudo killall -HUP mDNSResponder
```

Or when using OSX Yosemite (up until 10.10.4):

```
sudo discoveryutil mdnsflushcache
sudo discoveryutil udnsflushcaches
```

Check if the DNS server is working as expected by querying a running image:

```sh
dig <image>.docker @172.17.0.1
```

You should see a Docker network IP resolved correctly:

```
;; ANSWER SECTION:
<image>.docker.    0    IN    A    172.17.0.2
```

#### Diagnosing DNS
It usually takes some time to adapt to the DNS naming scheme of `dnsdock`, so if you'd like see which DNS names are being registered in real time, just follow the `dnsdock` logs:

`docker logs --follow dnsdock`

Sometimes the DNS server does not pickup the new resolver immediately. This may require a full `mDNSResponder` restart:

```
sudo launchctl unload -w /System/Library/LaunchDaemons/com.apple.mDNSResponder.plist
sudo launchctl load -w /System/Library/LaunchDaemons/com.apple.mDNSResponder.plist
```

Then check with `scutil --dns` if you have the following entry:

```
resolver #<n>
  domain   : docker
  nameserver[0] : 172.17.0.1
  flags    : Request A records
Reachable
```

## Troubleshooting
### Unresponsive `docker` cli
If `docker` cli commands hang, there's a good chance that you have a stale entry in your `/etc/hosts` file. Run `dlite stop`, then use sudo to edit your `/etc/hosts` file and remove any entries that end with `# added by dlite` or are surrounded by `# begin dlite` and `# end dlite` comments. Save the hosts file and run `dlite start` and try again.

### Accessing the terminal
If your virtual machine is misbehaving and you're unable to SSH to it, you can use the psuedo terminal that is allocated to the machine by running `screen /dev/ttys000`. Log in with the username `root` and the password `dhyve` and you can then perform some basic troubleshooting from there.

#### Tmux sessions
Note that `launchctl` commands appear to not work correctly when run inside tmux. If you are a tmux user and are having problems, try starting the service outside of your tmux session.

#### If you had boot2docker or docker-machine installed
Make sure to unset the environment variables so that the docker client uses the default configuration, as required by DLite. Add this lines to your `.bashrc` or corresponding shell profile file:

```
unset DOCKER_HOST
unset DOCKER_CERT_PATH
unset DOCKER_TLS_VERIFY
unset DOCKER_MACHINE_NAME
```

## Caveats
### Hypervisor framework
DLite depends on [xhyve](https://github.com/mist64/xhyve) which only works on OSX versions 10.10 (Yosemite) or newer. You also need a fairly recent mac. You can tell if your computer is new enough by running `sysctl kern.hv_support` in a terminal. If you see `kern.hv_support: 1` as a response, you're good to go. If not, unfortunately your computer is too old to leverage the Hypervisor framework and DLite won't work for you.

### Crash when waking after long sleep
There is an open issue with Xhyve ([https://github.com/mist64/xhyve/issues/86](https://github.com/mist64/xhyve/issues/86)) that causes OSX to crash when waking after a long sleep.

### TLS
DLite is _not_ secured via TLS. If that's important to you for local development, look elsewhere.

### Production usage
DLite is most definitely _not_ recommended for any kind of production use.

## Acknowledgements
DLite leverages [xhyve](https://github.com/mist64/xhyve) through the [libxhyve](https://github.com/TheNewNormal/libxhyve) Go bindings for virtualization. Without these projects and the people behind them, this project wouldn't exist.

## License
MIT

[travis-image]: https://img.shields.io/travis/nlf/dlite.svg?style=flat-square
[travis-url]: https://travis-ci.org/nlf/dlite
[gitter-image]: https://img.shields.io/gitter/room/nlf/dlite.svg?style=flat-square
[gitter-url]: https://gitter.im/nlf/dlite
