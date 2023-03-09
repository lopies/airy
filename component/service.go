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

package component

import (
	"fmt"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/logger"
)

var (
	componentsMap = make(map[constants.ComponentType]Component)
	componentsArr = []componentWrapper{}
)

type componentWrapper struct {
	component Component
	typ       constants.ComponentType
}

// RegisterComponent registers a component, by default it register after registered components
func RegisterComponent(component Component) error {
	return registerComponentAfter(component)
}

// registerComponentAfter registers a component after all registered components
func registerComponentAfter(component Component) error {
	typ := component.Type()
	if err := alreadyRegistered(typ); err != nil {
		return err
	}

	componentsMap[typ] = component
	componentsArr = append(componentsArr, componentWrapper{
		component: component,
		typ:       typ,
	})
	return nil
}

// GetComponent gets a component with a type
func GetComponent(typ constants.ComponentType) (Component, error) {
	if m, ok := componentsMap[typ]; ok {
		return m, nil
	}
	return nil, fmt.Errorf("component with type %s not found", typ)
}

func alreadyRegistered(typ constants.ComponentType) error {
	if _, ok := componentsMap[typ]; ok {
		return fmt.Errorf("component with type %s already exists", typ)
	}

	return nil
}

// StartComponents starts all components in order
func StartComponents(conf *config.AiryConfig) {
	logger.Debugf("initializing all components")
	for _, modWrapper := range componentsArr {
		logger.Debugf("init component: %s", modWrapper.typ)
		modWrapper.component.Init(conf)
		logger.Infof("init component: %s finish", modWrapper.typ)
	}

	for _, modWrapper := range componentsArr {
		logger.Debugf("after init component: %s", modWrapper.typ)
		modWrapper.component.AfterInit()
		logger.Infof("after init component: %s finish", modWrapper.typ)
	}
}

// ShutdownComponents starts all components in reverse order
func ShutdownComponents() {
	for i := len(componentsArr) - 1; i >= 0; i-- {
		logger.Debugf("before shutdown component: %s", componentsArr[i].typ)
		componentsArr[i].component.BeforeShutdown()
		logger.Infof("before shutdown component: %s finish", componentsArr[i].typ)
	}

	for i := len(componentsArr) - 1; i >= 0; i-- {
		typ := componentsArr[i].typ
		comp := componentsArr[i].component

		logger.Debugf("shutdown component: %s", componentsArr[i].typ)
		if err := comp.Shutdown(); err != nil {
			logger.Warnf("error stopping component: %s", typ)
		}
		logger.Infof("shutdown component: %s finish", componentsArr[i].typ)
	}
}
