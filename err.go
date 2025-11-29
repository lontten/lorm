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

import "github.com/pkg/errors"

var (
	ErrNil          = errors.New("nil")
	ErrContainEmpty = errors.New("slice empty")
	ErrNoPkOrUnique = errors.New(" ERROR: there is no unique or exclusion constraint matching the ON CONFLICT specification (SQLSTATE 42P10) ")
	ErrNoPk         = errors.New("no set primary key")
	ErrTypePkArgs   = errors.New("type of args is err")
	ErrNeedMultiPk  = errors.New("need multi primary key")
	ErrNoTableName  = errors.New("no set table name")
)
