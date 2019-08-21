package module

import (
	"plugin"
	"reflect"
)

type IModule interface {
	Boot()
	Info() Info
}

type Info struct {
	// Name          string `json:"name",desc:"名称"`
	Description   string `json:"description",desc:"描述"`
	Version       string `json:"version",desc:"版本"`
	Configuration bool   `json:"configuration",desc:"是否包含配置文件选项"`
	Dependencies  []string
}

type Module struct {
	Func map[string]interface{}
	// ViewFunc      template.FuncMap
	Apidoc        map[string]interface{}
	Install       map[string]func(string)
	Update        map[string]func(string)
	Uninstall     map[string]func(string)
}

var Value *Module

func New() {
	Value = &Module{
		Func: make(map[string]interface{}),
		// ViewFunc:      make(template.FuncMap),
		Apidoc:        make(map[string]interface{}),
		Install:       make(map[string]func(string)),
		Update:        make(map[string]func(string)),
		Uninstall:     make(map[string]func(string)),
	}
}

func (self *Module) AddFunc(name string, v interface{}) {
	self.Func[name] = v
}

func (self *Module) AddViewFunc(name string, v interface{}) {
	// self.ViewFunc[name] = v
}

func (self *Module) AddApiHandler(name string, v interface{}) {
	self.Apidoc[name] = v
}

func LoadModule(path string) (IModule, Info) {
	so, err := plugin.Open(path)

	if err != nil {
		panic(err)
	}

	symbol, err := so.Lookup("ExternalModule")

	if err != nil {
		panic(err)
	}

	currentModule, ok := symbol.(IModule)

	if !ok {
		panic("unexpected type from module symbol")
	}

	// currentModule.Boot()

	// moduleInfo, _ := currentModule.Info().(Info)
	moduleInfo := currentModule.Info()

	// fmt.Sprintf("[%s]%s %s", moduleInfo.Version, moduleInfo.Name, moduleInfo.Description)

	return currentModule, moduleInfo

	// Demo 导出函数
	// if moduleInfo.Name == "Demo" {
	// 	fmt.Printf("%+v \n", vars.Kernel.Func["Demo"])
	// 	fmt.Printf("%+v \n", vars.Kernel.Func["Demo"].(function.IExternalFunc).GetFullName("matt", "ma"))
	// }

	// for _, event := range moduleInfo.Events {
	// 	observer.Dispatcher.On(event.Name, osEvent.Listener{Callback: event.Callback, Priority: event.Priority})
	// }
}

func Install(m IModule) {
	method := reflect.ValueOf(m).MethodByName("Install")

	if method.IsValid() {
		install := method.Interface().(func())
		install()
	}
}

func Remove(m IModule) {
	method := reflect.ValueOf(m).MethodByName("Remove")

	if method.IsValid() {
		remove := method.Interface().(func())
		remove()
	}
}

func Update(m IModule, version string) {
	method := reflect.ValueOf(m).MethodByName("Update")

	if method.IsValid() {
		update := method.Interface().(func(string))
		update(version)
	}
}
