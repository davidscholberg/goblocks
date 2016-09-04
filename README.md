## Goblocks

Goblocks is an [i3status](https://i3wm.org/i3status/) replacement written in [Go](https://golang.org/), using the [Go-i3barjson](https://github.com/davidscholberg/go-i3barjson) library to communicate with [i3bar](https://i3wm.org/i3bar/).

The main goal of this project is to match most of the features of [i3blocks](https://github.com/vivien/i3blocks) while keeping all of the modules written in pure Go. This will keep Goblocks fast and lightweight, allowing the user to configure Goblocks with a very high update frequency without fear of taking up excessive system resource and battery.

Features of Goblocks include:

* Status indicators for:
    * RAID status (mdraid only)
    * filesystem usage
    * system load
    * memory availability
    * CPU temperature
    * network interfaces
    * volume (ALSA only)
    * date/time
* Configuration in [YAML](http://yaml.org/) format (see [config/goblocks.yml](config/goblocks.yml)).
* Ability to configure UNIX signal handlers to refresh individual blocks.
* Ability to reload the configuration by sending the HUP signal (e.g. `pkill -HUP goblocks`).
* Debug option to [pretty-print](https://en.wikipedia.org/wiki/Prettyprint) Goblocks' JSON output.

**WARNING:** Goblocks is still somewhat rough around the edges. See the [TODO](#todo) list.

### Get

Gobocks requires Go version 1.7+.

Install Goblocks and the sample config file:

```bash
go get github.com/davidscholberg/goblocks
mkdir -p $HOME/.config/goblocks
cp $GOPATH/src/github.com/davidscholberg/goblocks/config/goblocks.yml $HOME/.config/goblocks/
```

### Configure

Goblocks configuration is specified in [YAML](http://yaml.org/). The configuration file path is `$HOME/.config/goblocks/goblocks.yml`. A full configuration example with all available block types and options can be found at [config/goblocks.yml](config/goblocks.yml).

### Run

To use Goblocks in your i3bar, add the Goblocks binary to the [bar section of your i3 config](https://i3wm.org/docs/userguide.html#_configuring_i3bar). Note that if `$GOPATH/bin/` is not in your `$PATH` variable, then you'll have to specify the full path to the Goblocks binary.

You can reload Goblocks' configuration without restarting i3 by sending the HUP signal to Goblocks:

```bash
pkill -HUP goblocks
```

You can debug Goblocks' output by running it on the command line. If you set the `debug` [config option](config/goblocks.yml) to true, then Goblocks will [pretty-print](https://en.wikipedia.org/wiki/Prettyprint) the JSON output, making it easier to read.

### Contributing

If you would like to see a new feature or enhancement in Goblocks, please feel free to submit an [issue](/../../issues) or [pull request](/../../pulls).

### TODO

* Add battery block.
* Add wifi block.
* Only send update to i3bar if a block has updated within the update time interval.
* Add cli arg support.
* Add color support.
