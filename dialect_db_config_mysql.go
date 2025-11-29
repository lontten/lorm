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

import (
	"database/sql"
)

type MysqlVersion int

const (
	MysqlVersionLast MysqlVersion = iota
	MysqlVersion5    MysqlVersion = iota

	MysqlVersion8_0_19
	MysqlVersion8_0_20
	MysqlVersion8Last
)

type MysqlConf struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	Other    string
	Version  MysqlVersion
}

func (c MysqlConf) dialect(ctx *ormContext) Dialecter {
	ctx.ormConf.insertCanReturn = false
	if c.Version == MysqlVersionLast {
		c.Version = MysqlVersion8Last
	}
	return &MysqlDialect{
		ctx:       ctx,
		dbVersion: c.Version,
	}
}

func (c MysqlConf) open() (*sql.DB, error) {
	dsn := c.User + ":" + c.Password +
		"@tcp(" + c.Host +
		":" + c.Port +
		")/" + c.DbName + "?"

	if c.Other == "" {
		dsn += "charset=utf8mb4&parseTime=True&loc=Local"
	} else {
		dsn += c.Other
	}
	return sql.Open("mysql", dsn)
}
