# TV Awards

TV Awards is a series of packages written in Go, used to retrieve awards data from various sources such as: [The Emmys](emmys.com) and [The BAFTAs](bafta.org).

## Installation

To install the package, you need to install Go and set your Go workspace first.

1. You first need Go installed, then you can use the below Go command to install.

```
$ go get -u github.com/camc314/tvawards
```

2. Import it in your code:

```go
import "github.com/camc314/tvawards
```

## Quick start

The below program gets all BAFTA film awards and stores it to `./data/bafta-film-awards.json`

```go
package main

import "github.com/camc314/tvawards"

func main() {
    response, err := tvawards.GetBaftaAwards("film")

    if err != nil {
        panic(err)
    }

    file, _ := json.MarshalIndent(*response, "", "  ")

    ioutil.WriteFile("./data/bafta-film-awards.json", file, 0644)
}
```

```
$ go run main.go
```

## Supported Awards

- BAFTA's

## Planned Support

- Emmy's
