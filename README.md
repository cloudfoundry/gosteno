# Gosteno

Gosteno is a golang implementation of the
[steno log tool](https://github.com/cloudfoundry/steno).  The feature set of
Gosteno is very similar with that of ruby steno.

## Overview

Core concepts behind Gosteno includes codec, sink, level, tag.

### codec

A codec encodes log entries to structural data, more specifically, JSON format
data. Besides JSON codecs, Gosteno provides prettified codec which generates
more human-readable data.

### sink

Roughly speaking, a sink is the destination where you store your log data. It's
an abstraction of the underlying data storage systems. Currently Gosteno
supports two kinds of sinks, namely IOSink and SyslogSink. IOSink includes files
and standard output while SyslogSink streams your log data to syslog daemons
such as rsyslogd. You can register as many sinks as you want. Everytime you log
information, it will be written to all the sinks you have registered.

### level

Gosteno supports 9 levels(from low to high): all, debug2, debug1, debug, info,
warn, error, fatal, off. You can change the level on the fly without respawning
the process.

### tag

In gosteno, tags are extended information that will be encoded together with
other normal log information. You can add as many tags as you want. Tag makes
the log information extensive.

## Get Gosteno

    go get -u github.com/cloudfoundry/gosteno

## Getting started

Here is a short but complete program showing how to registering sinks, chosing
codec, tagging the information.

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

Changing logger properties such as log level without restarting system is
allowed in Gosteno. It is achieved through a http interface by some APIs and
data is exchanged as JSON:

  1. GET /regexp : get something like {"RexExp": "test$", "Level": "fatal"}
  2. PUT /regexp : put with data like {"RegExp": "test$", "Level":"fatal"}
  3. GET /loggers/{name} : get information about the logger by name
  4. PUT /loggers/{name} : put with data like {"Level" : "fatal" }
  5. GET /loggers : get information about all loggers

## Supported platforms

Currently targeting modern flavors of darwin and linux.

## License

Apache 2.0

