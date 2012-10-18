# Gosteno

## Overview

Gosteno is a golang implementation of the
[steno log tool](https://github.com/cloudfoundry/steno).  The feature set and concept in Gosteno are very similar with those in ruby steno.
## Getting started
    import (
        steno "gosteno"
        "os"
    )

    func main() {
        c := &steno.Config{
            Sinks: []steno.Sink{
                steno.NewFileSink("./a.log"),
                steno.NewIOSink(os.Stdout),
                steno.NewSyslogSink(),
            },
            Level:     steno.LOG_INFO,
            Codec:     steno.JSON_CODEC,
            Port:      8080,
            EnableLOC: true,
        }
        steno.Init(c)
        logger := steno.NewLogger("test")
        t := steno.NewTaggedLogger(logger, map[string]string{"foo": "bar", "hello": "world"})
        t.Info("Hello")
    }
## Using gerrit

    $ export GOPATH=~/gocode
    $ mkdir -p $GOPATH/src/github.com/cloudfoundry
    $ cd $GOPATH/src/github.com/cloudfoundry
    $ gerrit clone ssh://[<your username>@]reviews.cloudfoundry.org:29418/gosteno
    $ cd gosteno
    $ go test


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

