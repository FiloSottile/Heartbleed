Heartbleed
==========

A checker (site and tool) for CVE-2014-0160.

Public site at http://filippo.io/Heartbleed/

Tool usage: `Heartbleed [-service="service_name"] example.com[:443]`
        or: `Heartbleed service_name://example.com[:443]`

Exit codes: `0` - SAFE; `1` - VULNERABLE; `2` - ERROR. (*recently changed*)

Please note that the code is a bit of a mess, not exactly release-ready.

If a service is specified besides `https`, the tool checks the specified service using STARTTLS.
**You do still need to specify the correct port.**

## Install

You will need Go 1.2.x, otherwise you get `undefined: cipher.AEAD` and other errors

```
go get github.com/mozilla-services/Heartbleed
go install github.com/mozilla-services/Heartbleed
```

You can also use docker to get a ready to run virtual machine with heartbleed, see https://github.com/kasimon/docker-heartbleed.
