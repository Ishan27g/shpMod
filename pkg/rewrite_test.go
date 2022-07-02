package pkg

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	allDeps = `module "1" {
  disabled = true
  depends_on = [
    "2",
  ]
  source = "/one"
}
module "2" {
  disabled = true
  depends_on = [
    "3"
  ]
  source = "/two"
}
module "3" {
  disabled = true
  depends_on = [
    "4",
  ]
  source = "/three"
}
module "4" {
  disabled = true
  depends_on = [
#    "5",
#    "6",
#    "7"
  ]
  source = "/four"
}`
	singleDeps = `
module "1" {
  disabled = true
  depends_on = []
  source     = "/one"
}

module "2" {
  disabled = true
  depends_on = ["3"]
  source     = "/two"
}

module "3" {
  disabled = true
  depends_on = ["4"]
  source     = "/three"
}`
)

func toFile(val string) string {
	file, err := ioutil.TempFile("", "tmp.*.hcl")
	if err != nil {
		return ""
	}
	defer file.Close()
	file.Write([]byte(val))
	return file.Name()
}
func TestModuleRewrite(t *testing.T) {

	type testS struct {
		name         string
		file         string
		enable       string
		shouldEnable []string
	}
	var tests = []testS{{
		name:         "all dependencies enabled",
		file:         toFile(allDeps),
		enable:       "1",
		shouldEnable: []string{"1", "2", "3", "4"},
	}, {
		name:         "single dependency enabled -> 1 (self)",
		file:         toFile(singleDeps),
		enable:       "1",
		shouldEnable: []string{"1"},
	}, {
		name:         "single dependency enables -> 2 (self) & 3(dependency)",
		file:         toFile(singleDeps),
		enable:       "2",
		shouldEnable: []string{"2", "3"},
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok, enabled := Rewrite(test.file, test.enable)
			if !ok {
				t.Error()
			}
			for _, se := range test.shouldEnable {
				assert.Contains(t, enabled, se)
			}
		})
	}

	for _, test := range tests {
		_ = os.Remove(test.file)
	}

}
