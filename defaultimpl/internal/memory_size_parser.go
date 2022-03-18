package internal

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var SizeType = map[uint64]string{
	1 << 10: "KB",
	1 << 20: "MB",
	1 << 30: "GB",
	1 << 40: "TB",
}

func isValidSuffix(sizeStr string) (bool, uint64) {
	for k, v := range SizeType {
		if strings.HasSuffix(strings.ToUpper(sizeStr), v) {
			return true, k
		}
	}
	return false, 0
}

// ParseMemorySize 解析内存大小值sizeStr，形如1KB。支持单位K、M、G、T。返回等价的字节数值。
func ParseMemorySize(sizeStr string) (uint64, error) {
	isValid, power := isValidSuffix(sizeStr)
	if !isValid {
		log.Printf("memory size [%s] doesn't support\n", sizeStr)
		return 0, fmt.Errorf("memory size [%s] doesn't support\n", sizeStr)
	}
	size := sizeStr[:len(sizeStr)-2]
	if len(size) == 0 {
		log.Printf("memory size [%s] doesn't contain any valid number\n", sizeStr)
		return 0, fmt.Errorf("memory size [%s] doesn't contain any valid number\n", sizeStr)
	}
	sizeNum, err := strconv.ParseUint(size, 10, 64)
	if err != nil {
		log.Printf("memory size [%s] parsing failed; err: %s", size, err)
		return 0, fmt.Errorf("memory size [%s] parsing failed; err is: %s", size, err)
	}

	log.Printf("memory size parsing value is [%d]\n", sizeNum)
	return power * sizeNum, nil
}
