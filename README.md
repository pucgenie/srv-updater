dynv6 SRV Updater
=================

This program sends signed UDP packets to dynv6 to keep the port of SRV records up to date.
It only works with ed25519 keys.

## Installation

    go get github.com/dynv6/srv-updater

If you don't have a ed25519 key, generate one:

    openssh-keygen -t ed25519

## Usage

The SRV updater uses raw sockets. Therefore It needs root privilege or `CAP_NET_RAW` capability.

### Command line arguments

```
  -dst-host string
    	the destination host (default "dynv6.com")
  -dst-port int
    	the destination port (default 55)
  -fqdn string
    	the hostname of the record you want to update (default "_service._tcp.example.com")
  -interval duration
    	sending interval (default 1m)
  -key string
    	private key (default "~/.ssh/id_ed25519")
  -priority uint
    	the priority of the record you want to update (default 10)
  -src-ip string
    	the local ip address (default should be determined)
  -src-port int
    	the local port (default 10000)
  -weight uint
    	the weight of the record you want to update (default 1)
```

## Debugging

If the server does not understand the packet, it responds with a UDP packet containing an error message.
You can read it with network protocol analyzers like tcpdump or wireshark.
