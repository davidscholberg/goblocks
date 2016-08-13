## Goblocks

Goblocks is an [i3status](https://i3wm.org/i3status/) replacement written in [Go](https://golang.org/), using the [Go-i3barjson](https://github.com/davidscholberg/go-i3barjson) library to communicate with [i3bar](https://i3wm.org/i3bar/).

The main goal of this project is to match the features of [i3blocks](https://github.com/vivien/i3blocks) while keeping all of the modules written in pure Go. This will keep Goblocks fast and lightweight, allowing the user to configure Goblocks with a very high update frequency without fear of taking up excessive system resource and battery.

**WARNING:** Goblocks is still very rough around the edges. There's a lot of hardcoded stuff and no configuration yet. See the [TODO](#todo) list.

### Get

Fetch and build Goblocks:

```
go get github.com/davidscholberg/goblocks
```

### TODO

* Add configuration support (probably with [viper](https://github.com/spf13/viper)).
* Add cli arg support (probably with [cobra](https://github.com/spf13/cobra)).
* Add color support.
