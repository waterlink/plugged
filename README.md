# plugged

Library for writing extendable CLI applications.

## Usage

```go
import "github.com/waterlink/plugged"
```

### "Gateway" application

```go
func main() {
        plugged.Gateway("appname", "My super cli application.", os.Args)
}
```

Example usage:

```bash
$ ./appname
USAGE: appname command [options]

appname - My super cli application.

Available commands:

- find     - Find some stuff.
- activate - Activate stuff.
- help     - This help info.

To get help for any of commands you can do `appname help command` or `appname
command --help`.
$ ./appname find --help
# .. here output of `appname-find --help` ..
```

### Plugin application

```go
func main() {
        plugged.Plugin("appname", "find", "Find some stuff.", os.Args, handler)
}

func handler(args []string) {
        // .. Find some stuff here ..
}
```

## Installing plugin

Make sure you have installed plugin on your `PATH` and just run:

```bash
appname --plugged-install find
```

## Plugin interface

Plugin does not necessary need to be written in `go` and/or using `plugged`
library. Instead it is required for a plugin to abide to the following
interface:

```bash
appname-find --plugged-description  # => Find some stuff.
appname-find --help                 # => .. help message ..
```

## Development

You will need to have working recent `golang` installation (`1.5+` at a time of
writing). And the repo needs to be cloned into your `GOPATH`.

- `go test` runs tests.
- Please follow TDD.

## Contributing

1. Fork it ( https://github.com/waterlink/plugged/fork )
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create a new Pull Request

## Contributors

1. [waterlink](https://github.com/waterlink) - Oleksii Fedorov, creator,
   maintainer.
