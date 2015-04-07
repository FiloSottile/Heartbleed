Heartbleed
==========

A checker (site and tool) for CVE-2014-0160.

Public site at https://filippo.io/Heartbleed/

Tool usage:

```
    Heartbleed [-service="service_name"] example.com[:443]
    Heartbleed service_name://example.com[:443]
```

Exit codes: `0` - SAFE; `1` - VULNERABLE; `2` - ERROR. (*recently changed*)

See the [online FAQ](http://filippo.io/Heartbleed/faq.html) for an explanation of error messages including `TIMEOUT` and `BROKEN PIPE`.

If a service name is specified besides `https`, the tool checks the specified service using STARTTLS.
**You do still need to specify the correct port.**

## Install

You will need Go >= 1.2, otherwise you'll get `undefined: cipher.AEAD` and other errors

```
go get github.com/FiloSottile/Heartbleed
```

You can also use Docker to get a ready to run virtual machine with Heartbleed: https://github.com/kasimon/docker-heartbleed
