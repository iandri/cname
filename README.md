# cname

The _cname_ plugin eliminates CNAME records.

## Syntax

```
cname
```

This will transform responses like this:

```
;; ANSWER SECTION:
example.com.		3600	IN	CNAME	two.example.org.
two.example.org.	3600	IN	CNAME	one.example.net.
one.example.net.	3600	IN	A	127.0.0.1
```

into

```
;; ANSWER SECTION:
example.com.		3600	IN	A	127.0.0.1
```

## Installation

As per [CoreDNS docs](https://coredns.io/2017/07/25/compile-time-enabling-or-disabling-plugins/), there are two ways.

### Build with compile-time configuration file

```
$ git clone https://github.com/coredns/coredns
$ cd coredns
$ vim plugin.cfg
# Add the line cname:github.com/iandri/cname before the file middleware
$ go generate
$ go build
$ ./coredns -plugins | grep cname
```

### Build with external golang source code

```
$ git clone https://github.com/iandri/cname
$ cd cname/coredns
$ go build
$ ./coredns -plugins | grep cname
```
