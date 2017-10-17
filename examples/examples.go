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
	err = envi.ChangeValue("Intent", "33")
	if err != nil {
		//You can handle the errors as follows
		if e, ok := err.(*envi.EnvError); ok {
			switch e.Err {
			case envi.ValueIsEmpty:
				fmt.Println(err.Error())
			case envi.FieldNotExists:
				fmt.Println(err.Error())
			}
		}
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", env)
}
