# purl

A command-line URL parser.

### Usage

`purl` reads URLs from stdin, parses them according to the provided flags, and prints the result to stdout.

```
usage: purl [<flags>]

Flags:
  --help                 Show context-sensitive help (also try --help-long and --help-man).
  --scheme               Print the URL scheme
  --opaque               Print the opaque URL part
  --user                 Print the user and password information
  --host                 Print the URL hostname
  --path                 Print the URL path
  --query                Print the URL query
  --fragment             Print the URL fragment
  --separator=SEPARATOR  Separate results by a delimeter
```

By default, URL parts are separated by a tab (\t). You can use the `--separator` flag to override this.

### Examples

Print the hostname

```
$ echo "https://google.com/test" | purl --host
google.com
```

Print the hostname and path

```
$ echo "https://google.com/test/?q=test" | purl --host --path
google.com      /test/
```

Print the user:password information:

```
$ echo "https://jordan-wright:password@google.com/test" | purl --user
jordan-wright:password
```

Use a comma separator:

```
$ echo "https://jordan-wright:password@google.com/test" | purl --user --host --path --separator ,
jordan-wright:password,google.com,/test
```

### Build From Source

To build from source, simply `go get` the package:

```
go get github.com/jordan-wright/purl
```