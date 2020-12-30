package strutil

func ValueFindinInt(value int64, values ...int64) bool {
	for i := 0; i < len(values); i++ {
		if value == values[i] {
			return true
		}
	}
	return false
}

func ValueFindinString(value string, values ...string) bool {
	for i := 0; i < len(values); i++ {
		if value == values[i] {
			return true
		}
	}
	return false
}

func ValueFindinUint(value uint64, values ...uint64) bool {
	for i := 0; i < len(values); i++ {
		if value == values[i] {
			return true
		}
	}
	return false
}

//移除切片元素
func RemoveElementInt64(nums []int64, val int64) []int64 {
	if len(nums) == 0 {
		return nums
	}
	index := 0
	for ; index < len(nums); {
		if nums[index] == val {
			nums = append(nums[:index], nums[index+1:]...)
			continue
		}
		index++
	}
	return nums
}
func RemoveElementUint64(nums []uint64, val uint64) []uint64 {
	if len(nums) == 0 {
		return nums
	}
	index := 0
	for ; index < len(nums); {
		if nums[index] == val {
			nums = append(nums[:index], nums[index+1:]...)
			continue
		}
		index++
	}
	return nums
}
func RemoveElementString(nums []string, val string) []string {
	if len(nums) == 0 {
		return nums
	}
	index := 0
	for ; index < len(nums); {
		if nums[index] == val {
			nums = append(nums[:index], nums[index+1:]...)
			continue
		}
		index++
	}
	return nums
}

//slice去重
func RemoveRepeatedElement(slc []uint64) []uint64 {
	result := []uint64{}         //存放返回的不重复切片
	tempMap := map[uint64]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0 //当e存在于tempMap中时，再次添加是添加不进去的，，因为key不允许重复
		//如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e) //当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

// 过滤数组 去除src中item，在dst中存在的item
// src[1,2,3,4,5]   dst[2,4,6,8]	result[1,3,5]
func FilterSlice(src []int, dst []int) (result []int) {
	aMap := make(map[int]struct{})
	for _, v := range dst {
		aMap[v] = struct{}{}
	}
	for _, v := range src {
		if _, has := aMap[v]; !has {
			result = append(result, v)
		}
	}
	return result
}


