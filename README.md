# Gosteno

## Overview

Gosteno is a golang implementation of the
[steno log tool](https://github.com/cloudfoundry/steno).  The feature set of Gosteno is very similar with that of ruby steno.
## Get Gosteno

    go get -u github.com/cloudfoundry/gosteno

## Getting started

    package main

    import (
        steno "github.com/cloudfoundry/gosteno"
        "os"
    )

    func main() {
        c := &steno.Config{
            Sinks: []steno.Sink{
                steno.NewFileSink("./a.log"),
                steno.NewIOSink(os.Stdout),
                steno.NewSyslogSink("foobar"),
            },
            Level:     steno.LOG_INFO,
            Codec:     steno.NewJsonCodec(),
            Port:      8080,
            EnableLOC: true,
        }
        steno.Init(c)
        logger := steno.NewLogger("test")
        t := steno.NewTaggedLogger(logger, map[string]string{"foo": "bar", "hello": "world"})
        t.Info("Hello")
    }

## Change logger properties on the fly

Changing logger properties such as log level without restarting system is allowed in Gosteno. It is achieved through a http interface by some APIs and data is exchanged as JSON:

  1. GET /regexp : get something like {"RexExp": "test$", "Level": "fatal"}
  2. PUT /regexp : put with data like {"RegExp": "test$", "Level":"fatal"}
  3. GET /loggers/{name} : get information about the logger by name
  4. PUT /loggers/{name} : put with data like {"Level" : "fatal" }
  5. GET /loggers : get information about all loggers

## Supported platforms

Currently targeting modern flavors of darwin and linux.

## License

Apache 2.0

## File a Bug

To file a bug against Cloud Foundry Open Source and its components, sign up and use our
bug tracking system: [http://cloudfoundry.atlassian.net](http://cloudfoundry.atlassian.net)

