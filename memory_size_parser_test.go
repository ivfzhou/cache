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

import "testing"

func TestSizeParser(t *testing.T) {
	log := &defaultLog{}
	size, err := parseMemorySize("1KB", log)
	if err != nil {
		t.Fatal(err)
	}
	if size != 1024 {
		t.Fatal("size != 1024")
	}

	size, err = parseMemorySize("1mB", log)
	if err != nil {
		t.Fatal(err)
	}
	if size != 1024*1024 {
		t.Fatal("size != 1024*124")
	}

	size, err = parseMemorySize("0mb", log)
	if err != nil {
		t.Fatal(err)
	}
	if size != 0 {
		t.Fatal("size != 0")
	}

	size, err = parseMemorySize("-1mb", log)
	if err == nil {
		t.Fatal("err is nil")
	}

	size, err = parseMemorySize("mb", log)
	if err == nil {
		t.Fatal("err is nil")
	}
}
