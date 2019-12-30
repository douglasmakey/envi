package envi

import (
	"fmt"
	"os"
	"testing"
)

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
	t.Run("boolean", func(tt *testing.T) {
		defer os.Clearenv()

		type TestEnvironments struct {
			IsValid bool `env:"ISVALID"`
			IsProd  bool `env:"ISPROD"`
		}

		os.Setenv("ISVALID", "TRUE")
		os.Setenv("ISPROD", "FALSE")

		testEnv := TestEnvironments{}
		err := Parse(&testEnv)
		if err != nil {
			panic(fmt.Sprintf("%+v\n", err))
		}

		assertEqual(tt, true, testEnv.IsValid, "")
		assertEqual(tt, false, testEnv.IsProd, "")

	})

	t.Run("string", func(tt *testing.T) {
		defer os.Clearenv()

		type TestEnvironments struct {
			SomeString string `env:"SOMESTRING"`
			Name       string `env:"NAME"`
		}

		os.Setenv("SOMESTRING", "GOLANG")
		os.Setenv("NAME", "ENVI")

		testEnv := TestEnvironments{}
		err := Parse(&testEnv)
		if err != nil {
			panic(fmt.Sprintf("%+v\n", err))
		}

		assertEqual(tt, "GOLANG", testEnv.SomeString, "")
		assertEqual(tt, "ENVI", testEnv.Name, "")
	})

	t.Run("int", func(tt *testing.T) {
		defer os.Clearenv()

		type TestEnvironments struct {
			Port    int `env:"PORT"`
			Version int `env:"VERSION"`
		}

		os.Setenv("PORT", "8080")
		os.Setenv("VERSION", "1")

		testEnv := TestEnvironments{}
		err := Parse(&testEnv)
		if err != nil {
			panic(fmt.Sprintf("%+v\n", err))
		}

		assertEqual(tt, 8080, testEnv.Port, "")
		assertEqual(tt, 1, testEnv.Version, "")
	})

	t.Run("float", func(tt *testing.T) {
		defer os.Clearenv()

		type TestEnvironments struct {
			Rate    float32 `env:"RATE"`
			RateTwo float32 `env:"RATETWO"`
		}

		os.Setenv("RATE", "0.5")
		os.Setenv("RATETWO", "1.0")

		testEnv := TestEnvironments{}
		err := Parse(&testEnv)
		if err != nil {
			panic(fmt.Sprintf("%+v\n", err))
		}

		assertEqual(tt, float32(0.5), testEnv.Rate, "")
		assertEqual(tt, float32(1.0), testEnv.RateTwo, "")
	})

	t.Run("map", func(tt *testing.T) {
		defer os.Clearenv()

		type TestEnvironments struct {
			CodeCountries map[string]string `env:"COUNTRIES"`
		}

		os.Setenv("COUNTRIES", "Chile:CL,Venezuela:VEN,Colombia:CO")

		testEnv := TestEnvironments{}
		err := Parse(&testEnv)
		if err != nil {
			panic(fmt.Sprintf("%+v\n", err))
		}

		// Compare map
		if len(testEnv.CodeCountries) == 3 {
			assertEqual(tt, "CL", testEnv.CodeCountries["Chile"], "")
			assertEqual(tt, "VEN", testEnv.CodeCountries["Venezuela"], "")
			assertEqual(tt, "CO", testEnv.CodeCountries["Colombia"], "")
		} else {
			tt.Errorf("expected %#v", map[string]string{"Chile": "CL", "Venezuela": "VEN", "Colombia": "CO"})
		}
	})

	t.Run("slice", func(tt *testing.T) {
		defer os.Clearenv()

		type TestEnvironments struct {
			Numbers []int `env:"NUMBERS"`
		}

		os.Setenv("NUMBERS", "1,2,3,4,5")

		testEnv := TestEnvironments{}
		err := Parse(&testEnv)
		if err != nil {
			panic(fmt.Sprintf("%+v\n", err))
		}

		// Slice
		var numbers = [5]int{1, 2, 3, 4, 5}

		assertEqual(tt, len(numbers), len(testEnv.Numbers), "")
		assertEqual(tt, numbers[0], testEnv.Numbers[0], "")
		assertEqual(tt, numbers[1], testEnv.Numbers[1], "")
	})

	t.Run("error not a pointer", func(tt *testing.T) {
		type TestEnvironments struct {
			Numbers []int `env:"NUMBERS"`
		}

		testEnv := TestEnvironments{}
		err := Parse(testEnv)
		assertEqual(tt, errNotAPointer, err, "")

	})

	t.Run("error env is required", func(tt *testing.T) {
		type TestEnvRequired struct {
			IsProd bool `env:"PROD,required"`
		}

		testEnv := TestEnvRequired{}
		err := Parse(&testEnv)

		if err == nil {
			tt.Error("expected error IsRequired")
		}
	})
}
