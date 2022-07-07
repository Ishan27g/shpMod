package pkg

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/heimdalr/dag"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func Rewrite(fileName string, enable ...string) (bool, []string) {
	var toEnable = map[string]string{}

	for _, s := range enable {
		toEnable[s] = s
	}
	modules, err := ReadModulesFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return false, nil
	}

	var vertexes = map[string]string{}
	var allEnableDeps []string

	d := dag.NewDAG()
	for _, m := range modules.m {
		v, err := d.AddVertex(m.Name)
		if err == nil {
			vertexes[m.Name] = v
		}
		for _, do := range m.DependsOn {
			if vertexes[do] == "" {
				v2, _ := d.AddVertex(do)
				vertexes[do] = v2
			}
		}
	}
	for _, m := range modules.m {
		if vertexes[m.Name] == "" {
			continue
		}
		for _, do := range m.DependsOn {
			_ = d.AddEdge(vertexes[m.Name], vertexes[do])
		}
	}

	for name, v := range vertexes {
		children, _ := d.GetChildren(v)
		for _, c := range children {
			if toEnable[c.(string)] == c.(string) {
				allEnableDeps = append(allEnableDeps, c.(string))
			}
		}
		if toEnable[name] == name {
			allEnableDeps = append(allEnableDeps, name)
		}
	}
	modules.disable(modules.Names()...)
	modules.enable(allEnableDeps...)

	return WriteToFile(fileName, modules.m)
}

func WriteToFile(fileName string, modules map[string]*module) (bool, []string) {
	_ = os.Remove(fileName)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return false, nil
	}

	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	var format = func(field reflect.StructField) string {
		tag := strings.Replace(strings.Trim(string(field.Tag), "hcl"), "\"", "", 2)[1:]
		return strings.TrimSuffix(tag, ",optional")
	}
	var modsEnabled []module
	var modsDisabled []module
	for _, m := range modules {
		if m.Disabled {
			modsDisabled = append(modsDisabled, *m)
		} else {
			modsEnabled = append(modsEnabled, *m)
		}
	}
	var enabled []string
	for _, m := range append(modsEnabled, modsDisabled...) {
		rootBody.AppendNewline()
		barBlock := rootBody.AppendNewBlock("module", []string{m.Name})
		body := barBlock.Body()
		field, ok := reflect.TypeOf(&m).Elem().FieldByName("Disabled")
		if ok {
			if m.Disabled {
				body.SetAttributeValue(format(field), cty.BoolVal(m.Disabled))
			}
		}
		field, ok = reflect.TypeOf(&m).Elem().FieldByName("DependsOn")
		if ok {
			if len(m.DependsOn) > 0 {
				var c []cty.Value
				for _, d := range m.DependsOn {
					c = append(c, cty.StringVal(d))
				}
				body.SetAttributeValue(format(field), cty.ListVal(c))
			}
		}
		field, ok = reflect.TypeOf(&m).Elem().FieldByName("Source")
		if ok {
			body.SetAttributeValue(format(field), cty.StringVal(m.Source))
		}
		if m.Disabled {
			continue
		}
		enabled = append(enabled, m.Name)
	}

	fmt.Println(fmt.Sprintf("modules enabled - %s", enabled))

	_, err = file.Write(f.Bytes())
	if err != nil {
		return false, nil
	}
	return err == nil, enabled
}
