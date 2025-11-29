//  Copyright 2025 lontten lontten@163.com
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package lorm

type ConvertCtx struct {
	convertFuncs ConvertFuncMap
	valBox       ConvertValBoxMap
}

type Convert struct {
	name        string
	val         any
	convertFunc ConvertFunc
}

type ConvertFunc func(o any) any
type ConvertFuncMap map[string]ConvertFunc
type ConvertValBoxMap map[string]any

func (c ConvertCtx) Init() ConvertCtx {
	c.convertFuncs = ConvertFuncMap{}
	c.valBox = ConvertValBoxMap{}
	return c
}
func (c *ConvertCtx) Add(v Convert) {
	name := v.name
	c.valBox[name] = v.val
	c.convertFuncs[name] = v.convertFunc
}

func (c ConvertCtx) Get(name string) (any, ConvertFunc) {
	vb, ok := c.valBox[name]
	if !ok {
		return nil, nil
	}
	f := c.convertFuncs[name]
	return vb, f
}

func ConvertRegister[T any](name string, f func(v T) any) Convert {
	var t = new(T)
	return Convert{
		name: name,
		val:  t,
		convertFunc: func(val any) any {
			return f(*val.(*T))
		},
	}
}
