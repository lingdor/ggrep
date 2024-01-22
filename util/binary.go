package util

import (
	"math"
	"strconv"
)

func IntRemainBit(bit int) int {
	return strconv.IntSize - bit
}

func FullIntBinary(bit int) uint {
	return math.MaxUint >> IntRemainBit(bit)
}

func IsTrue(v uint, bit int) bool {
	if bit > 0 {
		v = v >> bit
	}
	return v&1 == 1
}
func SetTrue(v uint, bit int) uint {
	var num uint = 1
	if bit > 0 {
		num = num << bit
	}
	return v | num
}
