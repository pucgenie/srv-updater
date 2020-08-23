dynv6 SRV Updater
=================

This program sends signed UDP datagrams to dynv6 to keep the port of SRV records up to date.
It only works with ed25519 keys.

## Installation

    go get github.com/dynv6/srv-updater

If you don't have an ed25519 key, generate one:

    openssh-keygen -t ed25519

## Usage

The SRV updater uses raw sockets. Therefore it needs root privilege or `CAP_NET_RAW` capability.

### Command line arguments

```
  -dst-host string
    	the destination host (default "dynv6.com")
  -dst-port int
    	the destination port (default 55)
  -fqdn string
    	the hostname of the record you want to update (default "_service._tcp.example.com")
  -interval duration
    	sending interval (default 28s)
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

If the server does not understand the datagram, it responds with a UDP datagram containing an error message.
You can read it using network protocol analyzers like tcpdump or wireshark.

## Technical limitations

It won't work for network links where the MTU is tiny (say, less than 500 bytes) - the exact minimum limit depends on your SRV resource record name, though.
