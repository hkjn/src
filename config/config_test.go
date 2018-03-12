// Tests for package config.
package config

import (
	"reflect"
	"testing"
)

func ExampleMustLoad() {
	c := struct {
		Env struct {
			Address, Password string
		}
	}{}
	MustLoad(&c)
	// If there's a config.yaml in the directory of the code importing
	// the package (or up to `MaxSteps` directories above) with the
	// following YAML contents,
	// env:
	//   # Comments look like this.
	//   address: 1.2.3.4
	//   password: something
	// then `c` will now have the values from the config. See also
	// testdata/ for more examples.

}

type testConfig struct {
	Twitter struct {
		Token, Secret string
		TestAccounts  []string
	}
}

func TestMustLoad(t *testing.T) {
	BasePath = "testdata"

	wantPanic := false
	defer func() {
		if r := recover(); r != nil {
			if !wantPanic {
				t.Fatalf("unexpected panic: %+v\n", r)
			}
		}
	}()

	want := testConfig{}
	want.Twitter.Token = "YEAHaTOKEN"
	want.Twitter.Secret = "NAHNOTREALLY"
	want.Twitter.TestAccounts = []string{
		"notarealaccount1",
		"notarealaccount2",
	}
	c := testConfig{}
	MustLoad(&c)

	if !reflect.DeepEqual(c, want) {
		t.Fatalf("MustLoad() got %+v, want %+v\n", c, want)
	}

	// Now we're nested beyond `MaxSteps`, and expect a panic.
	BasePath = "testdata/1/2/3/4/5/6"
	MaxSteps = 5
	c = testConfig{}
	wantPanic = true
	MustLoad(&c)
}

func TestLoad(t *testing.T) {
	BasePath = "testdata"
	c := testConfig{}
	if err := Load(&c); err != nil {
		t.Fatalf("Load() failed: %v\n", err)
	}

	want := testConfig{}
	want.Twitter.Token = "YEAHaTOKEN"
	want.Twitter.Secret = "NAHNOTREALLY"
	want.Twitter.TestAccounts = []string{
		"notarealaccount1",
		"notarealaccount2",
	}
	if !reflect.DeepEqual(c, want) {
		t.Fatalf("Load() got %+v, want %+v\n", c, want)
	}
}

func TestLoad_Custom(t *testing.T) {
	BasePath = "testdata"
	c := struct {
		Twitter struct {
			Custom string
		}
	}{}
	want := c
	want.Twitter.Custom = "yes"
	if err := Load(&c, Name("custom.yaml")); err != nil {
		t.Fatalf("Load(Name(%s)) failed: %v\n", "custom.yaml", err)
	}
	if !reflect.DeepEqual(c, want) {
		t.Fatalf("Load(Name(%s)) got %+v, want %+v\n", "custom.yaml", c, want)
	}
}

func TestLoadPath(t *testing.T) {
	c := testConfig{}
	want := testConfig{}
	want.Twitter.Token = "YEAHaTOKEN"
	want.Twitter.Secret = "NAHNOTREALLY"
	want.Twitter.TestAccounts = []string{"notarealaccount1", "notarealaccount2"}

	if err := loadPath("testdata/config.yaml", &c); err != nil {
		t.Fatalf("loadPath(testdata/config.yaml) got err: %v\n", err)
	}
	if !reflect.DeepEqual(c, want) {
		t.Fatalf("loadPath() got %+v, want %+v\n", c, want)
	}
}

func TestTryLoad(t *testing.T) {
	// Loading the config when starting in the same directory should
	// work regardless of `MaxSteps` setting.
	BasePath = "testdata/"
	MaxSteps = 0
	want := testConfig{}
	want.Twitter.Token = "YEAHaTOKEN"
	want.Twitter.Secret = "NAHNOTREALLY"
	want.Twitter.TestAccounts = []string{"notarealaccount1", "notarealaccount2"}
	c := testConfig{}
	if err := tryLoad("config.yaml", &c); err != nil {
		t.Fatalf("tryLoad(config.yaml) got err: %v\n", err)
	}
	if !reflect.DeepEqual(c, want) {
		t.Fatalf("loadPath() got %+v, want %+v\n", c, want)
	}

	// A broken config should produce an error.
	if err := tryLoad("badconfig.yaml", &c); err == nil {
		t.Fatalf("tryLoad(badconfig.yaml) want err, got none")
	}

	// Loading the config when starting from a directory `MaxSteps` down
	// should work.
	BasePath = "testdata/1/2/3/4/5"
	MaxSteps = 5
	c = testConfig{}
	if err := tryLoad("config.yaml", &c); err != nil {
		t.Fatalf("tryLoad(config.yaml) got err: %v\n", err)
	}
	if !reflect.DeepEqual(c, want) {
		t.Fatalf("loadPath() got %+v, want %+v\n", c, want)
	}

	// Now we're too deeply nested, should give up and return err.
	MaxSteps = 4
	c = testConfig{}
	if err := tryLoad("config.yaml", &c); err == nil {
		t.Fatalf("tryLoad(config.yaml) want err, got none")
	}
}
