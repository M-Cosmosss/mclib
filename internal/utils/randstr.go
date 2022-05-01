package utils

import "github.com/pkg/errors"

func In(nums []int32, i int32) bool {
	for _, num := range nums {
		if num == i {
			return true
		}
	}
	return false
}

func CleanNumInSlice(nums []int32, i int32) ([]int32, error) {
	var re []int32
	for _, num := range nums {
		if num == i {
			continue
		} else {
			re = append(re, num)
		}
	}
	if len(re) == len(nums) {
		return nil, errors.New("数据库异常，请联系管理员")
	}
	return re, nil
}
