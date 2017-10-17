package envi

import (
	"fmt"
	"os"
	"testing"
)

type TestEnvironments struct {
	SomeString    string            `env:"SOMESTRING"`
	DbHost        string            `env:"DB_HOST" envDefault:"postgres://localhost:5432/db"`
	Port          int               `env:"PORT"`
	CodeCountries map[string]string `env:"COUNTRIES"`
	Rate          float32           `env:"RATE"`
	Numbers       []int             `env:"NUMBERS"`
	NotNumbers    []int             `env:"NOTNUMBERS"`
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func TestParse(t *testing.T) {
	defer os.Clearenv()

	os.Setenv("SOMESTRING", "ENVI")
	os.Setenv("PORT", "8080")
	os.Setenv("RATE", "0.5")
	os.Setenv("COUNTRIES", "Chile:CL,Venezuela:VEN,Colombia:CO")
	os.Setenv("NUMBERS", "1,2,3,4,5")

	testEnv := TestEnvironments{}
	err := Parse(&testEnv)
	if err != nil {
		panic(fmt.Sprintf("%+v\n", err))
	}
	assertEqual(t, "ENVI", testEnv.SomeString, "")
	assertEqual(t, 8080, testEnv.Port, "")
	assertEqual(t, float32(0.5), testEnv.Rate, "")
	assertEqual(t, "postgres://localhost:5432/db", testEnv.DbHost, "")

	// Change data
	ChangeValue("Rate", "0.8")
	ChangeValue("Port", "2323")
	assertEqual(t, float32(0.8), testEnv.Rate, "")
	assertEqual(t, 2323, testEnv.Port, "")

	// Slice
	var numbers = [5]int{1, 2, 3, 4, 5}

	assertEqual(t, len(numbers), len(testEnv.Numbers), "")
	assertEqual(t, numbers[0], testEnv.Numbers[0], "")
	assertEqual(t, numbers[1], testEnv.Numbers[1], "")

	// Compare map
	if len(testEnv.CodeCountries) == 3 {
		assertEqual(t, "CL", testEnv.CodeCountries["Chile"], "")
		assertEqual(t, "VEN", testEnv.CodeCountries["Venezuela"], "")
		assertEqual(t, "CO", testEnv.CodeCountries["Colombia"], "")
	} else {
		t.Errorf("expected %#v", map[string]string{"Chile": "CL", "Venezuela": "VEN", "Colombia": "CO"})
	}
}
