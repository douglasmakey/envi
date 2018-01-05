## ENVI
Golang library for managing configuration from environment variables

## Usage

Basic example (in `examples` folder):

```go
package main

import (
    "fmt"
    "github.com/douglasmakey/envi"
)

type environments struct {
    Intent int            `env:"INTENT"`
    Ports  []int          `env:"PORTS" envDefault:"3000"`
    IsProd bool           `env:"PROD,required"`
    IsDev  bool           `env:"DEV"`
    Hosts  []string       `env:"HOSTS" envSeparator:":"`
    Sector map[string]int `env:"SECTOR"`
}

func main() {
    env := environments{}
    err := envi.Parse(&env)
    if err != nil {
        // You can handle the errors as follows
        if e, ok := err.(*envi.EnvError); ok {
            switch e.Err {
            case envi.IsRequired:
                // You can get info details using
                // e.KeyName --> Name of key
                // e.Value --> Value
                panic(e.Error())
            }
        }
        
        fmt.Println(err)
	}

    fmt.Printf("%+v\n", env)
}

```

You can run it like this:

```sh
$ INTENT=5 PROD=true HOSTS="127.0.0.1:localhost" SECTOR="a:1,b:2,c:4"  go run examples/examples.go
Intent:5 Ports:[3 0 0 0] IsProd:true IsDev:false Hosts:[127.0.0.1 localhost] Sector:map[a:1 b:2 c:4]}
```

## Supported types and defaults

This library has support for the following types:

* `string`
* `int`
* `uint`
* `int64`
* `bool`
* `float32`
* `float64`
* `Map[string]int`
* `Map[string]string`
* `[]string`
* `[]int`
* `[]bool`
* `[]float32`
* `[]float64`


You can set the `envDefault` tag for something, this value will be used in the
case of absence of it in the environment. If you don't do that AND the
environment variable is also not set, the zero-value
of the type will be used: empty for `string`s, `false` for `bool`s
and `0` for `int`s.

By default, slice types will split the environment value on `,`; you can change this behavior by setting the `envSeparator` tag.

The `env` tag option `required` for example `env:"MyKey,required"` can be added
to ensure that some environment variable is set.

## TODO
- Implement errors handler
- Implement httpHandler for list env and change value.
