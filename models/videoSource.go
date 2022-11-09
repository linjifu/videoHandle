package models

type VideoSource struct {
	dataStr   string
	index     int
	startTime float32
	endTime   float32
}

func (v *VideoSource) DataStr() string {
	return v.dataStr
}

func (v *VideoSource) SetDataStr(dataStr string) {
	v.dataStr = dataStr
}

func (v *VideoSource) Index() int {
	return v.index
}

func (v *VideoSource) SetIndex(index int) {
	v.index = index
}

func (v *VideoSource) StartTime() float32 {
	return v.startTime
}

func (v *VideoSource) SetStartTime(startTime float32) {
	v.startTime = startTime
}

func (v *VideoSource) EndTime() float32 {
	return v.endTime
}

func (v *VideoSource) SetEndTime(endTime float32) {
	v.endTime = endTime
}

func NewVideoSource(index int, startTime float32, endTime float32, dataStr string) *VideoSource {
	return &VideoSource{
		index:     index,
		startTime: startTime,
		endTime:   endTime,
		dataStr:   dataStr,
	}
}

// 2.声明一个Hero结构体切片类型
type VideoSourceSlice []*VideoSource

func (vss VideoSourceSlice) Len() int {
	return len(vss)
}
func (vss VideoSourceSlice) Less(i, j int) bool {
	return vss[i].index < vss[j].index
}
func (vss VideoSourceSlice) Swap(i, j int) {
	vss[i], vss[j] = vss[j], vss[i]
}
