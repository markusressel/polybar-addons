# polybar-addons

This repo contains a number of utility programs I use in my [polybar] setup.

## Programs

| Name      | Example                  | Description                                                          |
|-----------|--------------------------|----------------------------------------------------------------------|
| `battery` | `00:45`                  | Prints the ETR of usage on battery/charging to full/∞                |
| `disk`    | `↑ 0.0 B/s ↓ 16.3 B/s`   | Prints the average disk read/write activity since the last call      |
| `zfs`     | `5% (3.54G), 21% (725G)` | Prints ZFS pool statistics                                           |
| `network` | `↓ 12.6MB/s ↑ 45.2 B/s`  | Prints the average network send/receive activity since the last call |

### Templating

Although all programs have a default output, they also accept an optional "template" string. You can use this template
to modify how the data gathered by the programs is printed. A description of available template placeholders can be
found below.

Example:

```shell
> disk "Read: %reads% Write: %writes%"
Read:   27.2KB/s Write:    2.7MB/s
```

Since the output of these programs is passed to polybar, you can also include color codes within the template:

```shell
> disk "%{F00C853}Read: %reads% Write: %writes%%{FDFDFDF}"
%{F00C853}Read:   27.2KB/s Write:    2.7MB/s{FDFDFDF}
```

Color codes are simply passed on to polybar, which will parse and use them to color the output accordingly.


### battery

| Name        | Example | Description                       |
|-------------|---------|-----------------------------------|
| `%hours%`   | `01`    | 2 digit padded remaining hours.   |
| `%minutes%` | `01`    | 2 digit padded remaining minutes. |

### network

| Name            | Example     | Description                                        |
|-----------------|-------------|----------------------------------------------------|
| `%received%`    | ` 12.6MB/s` | Monospaced data rate for incoming network traffic. |
| `%transmitted%` | ` 16.3 B/s` | Monospaced data rate for outgoing network traffic. |

### disk

| Name       | Example     | Description                                     |
|------------|-------------|-------------------------------------------------|
| `%reads%`  | ` 12.6MB/s` | Monospaced data rate for reading disk activity. |
| `%writes%` | ` 16.3 B/s` | Monospaced data rate for writing disk activity. |

### zfs

Placeholders must be prefixed with the name of the target pool.

| Name            | Example | Description                    |
|-----------------|---------|--------------------------------|
| `%rpool.free%`  | `750GB` | Free pool space.               |
| `%rpool.used%`  | `250GB` | Used pool space.               |
| `%rpool.cap%`   | `25%`   | Used pool capacity in percent. |
| `%rpool.total%` | `1TB`   | Total pool size.               |



## How to use

### Build and "Install" (copy to home)

To build and copy all executables to `~/.config/polybar/scripts`

```shell
git clone https://github.com/markusressel/polybar-addons.git
cd polybar-addons
make deploy
```

### Polybar Config

Then in your polybar config you can use them like this:

```

modules-right = your_name_of_choice

[...]

[module/your_name_of_choice]
type = custom/script
exec = ~/.config/polybar/scripts/battery "%{F00C853}%hours%:%minutes%%{FDFDFDF}"
interval = 2

[...]

```

[polybar]: https://github.com/polybar/polybar
