package math

import (
	"github.com/gansidui/geohash"
	"github.com/gansidui/nearest"
)

func EncodeGeoCode(lat, lon float64, precision int) string {
	hash, _ := geohash.Encode(lat, lon, precision)
	return hash
}

// 计算两坐标点的地球面距离, 单位为 km
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	return nearest.Distance(lat1, lon1, lat2, lon2)
}

//排序
type Sort struct {
	SortValue float64
	Info      interface{}
}
type SortList []Sort

//排序规则
func (list SortList) Len() int {
	return len(list)
}

func (list SortList) Less(i, j int) bool {
	if list[i].SortValue < list[j].SortValue {
		return true
	}
	return false
}
func (list SortList) Swap(i, j int) {
	var temp Sort = list[i]
	list[i] = list[j]
	list[j] = temp
}
