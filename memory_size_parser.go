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

import (
	"fmt"
	"strconv"
	"strings"
)

var sizeType = map[uint]string{
	1 << 10: "KB",
	1 << 20: "MB",
	1 << 30: "GB",
	1 << 40: "TB",
}

func parseMemorySize(sizeStr string, log Logger) (uint, error) {
	sizeStr = strings.TrimSpace(sizeStr)
	isValid, power := isValidSuffix(sizeStr)
	if !isValid {
		log.Error("memory size [%s] doesn't support", sizeStr)
		return 0, fmt.Errorf("memory size [%s] doesn't support", sizeStr)
	}
	size := sizeStr[:len(sizeStr)-2]
	if len(size) == 0 {
		log.Error("memory size [%s] doesn't contain any valid number", sizeStr)
		return 0, fmt.Errorf("memory size [%s] doesn't contain any valid number", sizeStr)
	}
	sizeNum, err := strconv.ParseUint(size, 10, 64)
	if err != nil {
		log.Error("memory size [%s] parsing failed; err: %s", size, err)
		return 0, fmt.Errorf("memory size [%s] parsing failed; err is: %s", size, err)
	}

	log.Info("memory size parsing value is [%d]", sizeNum)
	return power * uint(sizeNum), nil
}

func isValidSuffix(sizeStr string) (bool, uint) {
	for k, v := range sizeType {
		if strings.HasSuffix(strings.ToUpper(sizeStr), v) {
			return true, k
		}
	}
	return false, 0
}
