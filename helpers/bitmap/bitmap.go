package bitmap

import "errors"

//基于字符串做位的操作(从1开始)
func SetBit(m string, bit int) (string, error) {
	if bit < 1 {
		return m, errors.New("bit must be positive number")
	}
	var flag []byte
	length := bit - len(m)
	if length > 0 { //本身长度不足，要补位了
		for i := 0; i < length; i++ {
			if i == (length - 1) {
				flag = append(flag, byte('1'))
			} else {
				flag = append(flag, byte('0'))
			}
		}
		return m + string(flag), nil
	} else {
		for k, v := range m {
			if k+1 == bit {
				flag = append(flag, byte('1'))
			} else {
				flag = append(flag, byte(v))
			}
		}
		return string(flag), nil
	}
	return m, nil

}
func UnsetBit(m string, bit int) (string, error) {
	if bit < 1 || bit-len(m) > 0 {
		return m, errors.New("bit must be positive number or bit must be less than length of string")
	}
	var flag []byte
	for k, v := range m {
		if k+1 == bit {
			flag = append(flag, byte('0'))
		} else {
			flag = append(flag, byte(v))
		}
	}
	return string(flag), nil
}
func IsBitSet(m string, bit int) (bool, error) {
	if bit < 1 || bit-len(m) > 0 {
		return false, errors.New("bit must be positive number or bit must be less than length of string")
	}
	for k, v := range m {
		if k+1 == bit {
			if v == '1' {
				return true, nil
			}else{
				return false,nil
			}
		}
	}
	return false, nil
}
