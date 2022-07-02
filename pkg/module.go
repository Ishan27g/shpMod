package pkg

import (
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Modules struct{ m map[string]*module }

type module struct {
	Name      string   `hcl:",label"`
	Disabled  bool     `hcl:"disabled,optional"`
	DependsOn []string `hcl:"depends_on,optional"`
	Source    string   `hcl:"source"`
}

func ReadModulesFile(filename string) (Modules, error) {
	var mod struct {
		Module []*module `hcl:"module,block"`
	}
	err := hclsimple.DecodeFile(filename, nil, &mod)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
		return Modules{}, err
	}
	modules := Modules{map[string]*module{}}
	for _, m := range mod.Module {
		modules.m[m.Name] = m
	}
	return modules, nil
}

func (m *Modules) Names() (names []string) {
	for _, m2 := range m.m {
		names = append(names, m2.Name)
	}
	return
}
func (m *Modules) enable(module ...string) *Modules {
	for _, name := range module {
		if m.m[name] == nil || m.m[name].Name != name {
			continue
		}
		if strings.HasPrefix(m.m[name].Name, name) {
			m.m[name].Disabled = false
			m.enable(m.m[name].DependsOn...)
		}
	}
	return m
}

func (m *Modules) disable(module ...string) {
	for _, name := range module {
		if m.m[name] == nil && m.m[name].Name != name {
			continue
		}
		m.m[name].Disabled = true
	}
}
