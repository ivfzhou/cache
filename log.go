/*
 * Copyright (c) 2023 ivfzhou
 * cache is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package cache

import "fmt"

type defaultLog struct {
	level Level
}

func (log *defaultLog) Info(format string, args ...any) {
	if log.level <= LevelInfo {
		fmt.Printf("cache INFO: "+format+"\n", args...)
	}
}

func (log *defaultLog) Warn(format string, args ...any) {
	if log.level <= LevelWarn {
		fmt.Printf("cache WARN: "+format+"\n", args...)
	}
}

func (log *defaultLog) Error(format string, args ...any) {
	if log.level <= LevelError {
		fmt.Printf("cache ERROR: "+format+"\n", args...)
	}
}

func (log *defaultLog) SetLevel(l Level) {
	log.level = l
}
