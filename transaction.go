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
	"fmt"
	"log"
)

type RollbackWithResult[T any] struct {
	Result T
}

func Rollback[T any](result T) {
	panic(RollbackWithResult[T]{Result: result})
}

func TransactionErr[T any](db Engine, fn func(tx Engine) (T, error)) (res T, err error) {
	tx, err := db.Begin()
	if err != nil {
		return res, err
	}

	defer func() {
		r := recover()
		if r != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				log.Printf("rollback failed: %v", rbErr)
			}

			// 检查是否是自定义的回滚类型
			switch v := r.(type) {
			case RollbackWithResult[T]:
				// 如果是自定义回滚，返回指定的结果，没有错误
				res = v.Result
				err = nil
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", err)
			}
			return
		}
	}()

	res, err = fn(tx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			log.Printf("rollback failed: %v", rbErr)
		}
		return res, err
	}

	if err = tx.Commit(); err != nil {
		return res, fmt.Errorf("commit failed: %v", err)
	}

	return res, nil
}

func Transaction[T any](db Engine, fn func(tx Engine) T) (res T) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	defer func() {
		r := recover()
		if r != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				log.Printf("rollback failed: %v", rbErr)
			}

			// 检查是否是自定义的回滚类型
			switch v := r.(type) {
			case RollbackWithResult[T]:
				// 如果是自定义回滚，返回指定的结果，没有错误
				res = v.Result
			case error:
				panic(v)
			default:
				panic(fmt.Errorf("%v", err))
			}
			return
		}
	}()

	res = fn(tx)

	if err = tx.Commit(); err != nil {
		panic(fmt.Errorf("commit failed: %v", err))
	}

	return res
}
