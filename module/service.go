// Copyright (c) Airy Author. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package module

import (
	"fmt"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/logger"
	"runtime"
)

var (
	modulesMap = make(map[constants.ModuleType]Module)
	modules    = []*moduleWrapper{}
)

type moduleWrapper struct {
	typ    constants.ModuleType
	module Module
}

// RegisterModule registers a module, by default it register after registered modules
func RegisterModule(module Module) error {
	return registerModuleAfter(module)
}

// startModules starts all modules in order
func StartModules(conf *config.AiryConfig) {
	for i := 0; i < len(modules); i++ {
		modules[i].module.Init(conf)
		logger.Debugf("initializing module: %s", modules[i].typ)
	}

	for i := 0; i < len(modules); i++ {
		modules[i].module.AfterInit()
		logger.Infof("module: %s successfully loaded", modules[i].typ)
	}
}

// shutdownModules starts all modules in reverse order
func ShutdownModules() {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 2<<10)
			l := runtime.Stack(buf, false)
			logger.Error("%v: %s", r, buf[:l])
		}
	}()

	for i := len(modules) - 1; i >= 0; i-- {
		logger.Debugf("before shutdown module: %s", modules[i].typ)
		modules[i].module.BeforeShutdown()
		logger.Debugf("before shutdown module: %s finish", modules[i].typ)
	}

	for i := len(modules) - 1; i >= 0; i-- {
		mod := modules[i]
		typ := mod.typ
		logger.Debugf("shutdown module: %s", modules[i].typ)
		if err := mod.module.Shutdown(); err != nil {
			logger.Warnf("error stopping module: %s", typ)
		}
		logger.Debugf("shutdown module: %s finish", modules[i].typ)
	}
}

// GetModule gets a module with a type
func GetModule(typ constants.ModuleType) (Module, error) {
	if m, ok := modulesMap[typ]; ok {
		return m, nil
	}
	return nil, fmt.Errorf("module with type %s not found", typ)
}

func alreadyRegistered(typ constants.ModuleType) error {
	if _, ok := modulesMap[typ]; ok {
		return fmt.Errorf("module with type %s already exists", typ)
	}
	return nil
}

// registerModuleAfter registers a module after all registered modules
func registerModuleAfter(module Module) error {
	typ := module.Type()
	if err := alreadyRegistered(typ); err != nil {
		return err
	}

	modulesMap[typ] = module
	modules = append(modules, &moduleWrapper{
		typ:    typ,
		module: module,
	})
	return nil
}
