package main

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ScriptLoader struct {
	basePath string
	vm       *otto.Otto
	scripts  map[string]otto.Value
}

func New(basePath string) *ScriptLoader {
	self := &ScriptLoader{
		basePath: basePath,
		vm:       otto.New(),
		scripts:  map[string]otto.Value{},
	}

	_ = self.vm.Set("define", func(call otto.FunctionCall) otto.Value {
		__FILE__ := call.Argument(0).String()
		value := call.Argument(1)
		//if value.IsFunction() {
		//	fmt.Printf("define(%s,fun)", __FILE__)
		//} else {
		//	fmt.Println("define(%s,%s)", __FILE__, value.Class())
		//}
		fmt.Println("define(", __FILE__, ",", value.Class(), ")")
		self.scripts[__FILE__] = value
		return value
	})

	_ = self.vm.Set("require", func(call otto.FunctionCall) otto.Value {
		requireName := call.Argument(0).String()
		__FILE__, _ := filepath.Abs(self.basePath + requireName)

		value, ok := self.scripts[__FILE__]
		if !ok {
			v, _ := self.Load(requireName)
			return v
		}
		return value
	})
	return self
}

func (self ScriptLoader) Init() {
	//
}

func (self ScriptLoader) Load(name string) (otto.Value, error) {
	__FILE__, _ := filepath.Abs(self.basePath + name)

	//从缓存读取脚本
	value, ok := self.scripts[__FILE__]
	if ok {
		return value, nil
	}

	readBytes, readError := ReadAll(__FILE__)
	if readError != nil {
		return otto.Value{}, readError
	}

	setError := self.vm.Set("__FILE__", __FILE__)
	if setError != nil {
		return otto.Value{}, setError
	}
	return self.vm.Run(string(readBytes))
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}
