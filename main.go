package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

// args is a simple struct that holds our arguments. It's passed to parseURL
// to indicate which fields to print and how.
type args struct {
	scheme    bool
	opaque    bool
	user      bool
	host      bool
	path      bool
	query     bool
	fragment  bool
	separator string
}

// parseURL receives the commmand line arguments and URL to parse, parses the
// URL, and returns a string with the requested parts.
// If args.separator is specified, the parts will be split using the separator.
// Otherwise, the URL will be returned using the standard format.
func parseURL(rawurl string, cli args) (string, error) {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	built := URL{}
	// Normally, we could use a switch here, but Go's switch statements only
	// execute the selected case, so we have to resort to if statements
	if cli.scheme {
		built.Scheme = parsed.Scheme
	}
	if cli.opaque {
		built.Opaque = parsed.Opaque
	}
	if cli.user {
		built.User = parsed.User
	}
	if cli.host {
		built.Host = parsed.Host
	}
	if cli.path {
		built.Path = parsed.Path
	}
	if cli.query {
		built.RawQuery = parsed.RawQuery
	}
	if cli.fragment {
		built.Fragment = parsed.Fragment
	}
	if cli.separator != "" {
		return built.Join(cli.separator), nil
	}
	return built.String(), nil
}

func main() {
	cli := args{}
	kingpin.Flag("scheme", "Print the URL scheme").BoolVar(&cli.scheme)
	kingpin.Flag("opaque", "Print the opaque URL part").BoolVar(&cli.opaque)
	kingpin.Flag("user", "Print the user and password information").BoolVar(&cli.user)
	kingpin.Flag("host", "Print the URL hostname").BoolVar(&cli.host)
	kingpin.Flag("path", "Print the URL path").BoolVar(&cli.path)
	kingpin.Flag("query", "Print the URL query").BoolVar(&cli.query)
	kingpin.Flag("fragment", "Print the URL fragment").BoolVar(&cli.fragment)
	kingpin.Flag("separator", "Separate results by a delimeter").StringVar(&cli.separator)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Enumerate through stdin, parsing and printing the results
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parsed, err := parseURL(scanner.Text(), cli)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
		fmt.Println(parsed)
	}
}
