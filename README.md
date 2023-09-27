# Go-Profiles Configuration Framework
<img align="right" width="159px" src="https://raw.githubusercontent.com/go-profiles/logo/master/color.png">

[![Go](https://github.com/go-profiles/goprofiles/actions/workflows/goprofiles.yml/badge.svg)](https://github.com/go-profiles/goprofiles/actions/workflows/go.yml) [![codecov](https://codecov.io/gh/go-profiles/goprofiles/graph/badge.svg?token=LV851U823H)](https://codecov.io/gh/go-profiles/goprofiles)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-profiles/goprofiles)](https://goreportcard.com/report/github.com/go-profiles/goprofiles)

Go-Profiles is a lightweight configuration management framework written in Go. It simplifies configuration handling by allowing developers to define different profiles and easily switch between them based on their application environment.

**The key features of Go-Profiles are:**

- Absolutely no config structs needed
- Single YAML-based configuration file
- Multiple profile activation
- YAML validation
- Easy configuration value retrieval

## Getting started

### Prerequisites

- **[Go](https://go.dev/)**: any one of the **three latest major** [releases](https://go.dev/doc/devel/release) (we test it with these).

### Getting Go-Profiles

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```go
import "github.com/go-profiles/goprofiles"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `go-profiles` package:

```sh
$ go get -u github.com/go-profiles/goprofiles
```


### Running Go-Profiles

First you need to define your Go-Profiles in a yaml file such as the following `profiles.yaml`:

```yaml
goprofiles:
  dev:
    db:
      host: localhost
      port: 5432
    prod:
      db:
        host: db.prod.com
        port: 5432
```
Next you need to import the Go-Profiles package to load and access your configuration. The following simple `example.go` illustrates how simple Go-Profiles makes configuration management:

```go
package main

import "github.com/go-profiles/goprofiles"

func main() {
  // Import the Go-Profile package and load your configuration:
  profile := goprofiles.New(goprofiles.WithProfile("dev"), goprofiles.WithFile("profiles.yaml"))
  
  // Retrieve values from the profile:
  host := profile.GetString("db.host")
  port := profile.GetInt("db.port")
  
  // Print the database configuration:
  println("Database config - Host:", host, " Port:", port)
}

```

And use the Go command to run the demo:

```
# run example.go
$ go run example.go
