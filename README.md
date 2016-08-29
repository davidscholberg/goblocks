## Goblocks

Goblocks is an [i3status](https://i3wm.org/i3status/) replacement written in [Go](https://golang.org/), using the [Go-i3barjson](https://github.com/davidscholberg/go-i3barjson) library to communicate with [i3bar](https://i3wm.org/i3bar/).

The main goal of this project is to match the features of [i3blocks](https://github.com/vivien/i3blocks) while keeping all of the modules written in pure Go. This will keep Goblocks fast and lightweight, allowing the user to configure Goblocks with a very high update frequency without fear of taking up excessive system resource and battery.

**WARNING:** Goblocks is still somewhat rough around the edges. See the [TODO](#todo) list.

### Get

Gobocks requires Go version 1.7+.

Fetch and build Goblocks:

```
go get github.com/davidscholberg/goblocks
```

### Configure

Goblocks configuration is specified in [YAML](http://yaml.org/). The configuration file path is `$HOME/.config/goblocks/goblocks.yml`. Here is a simple example configuration:

```yaml
global:
    debug: False

blocks:
    load:
        block_index: 1
        update_interval: 1
        label: "L: "
        crit_load: 4

    interfaces:
        - block_index: 2
          update_interval: 1
          label: "E: "
          interface_name: enp3s0

        - block_index: 3
          update_interval: 1
          label: "W: "
          interface_name: wlp4s2

    volume:
        block_index: 4
        update_interval: 60
        label: "V: "
        update_signal: 8

    time:
        block_index: 5
        update_interval: 1
        time_format: 2006-01-02 15:04
```

A full configuration example with all available block types and options can be found at [config/goblocks.yml](config/goblocks.yml).

### Contributing

If you would like to see a new feature or enhancement in Goblocks, please feel free to submit an [issue](/../../issues) or [pull request](/../../pulls).

### TODO

* Update ticker handling to allow for times less than a second.
* Add battery block.
* Add wifi block.
* Allow disk mounts to be loaded via config for disk block.
* Add cli arg support.
* Add color support.
