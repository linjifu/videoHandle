package models

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Video struct {
	//路径
	path string
	//名称
	name string
	//后缀
	suffix string
	//封面
	cover string
	//时长
	duration string
	//总秒数
	totalSecond float32
	//排序位置
	sort int
	//帧率
	tbr float32
}

func (v *Video) Tbr() float32 {
	return v.tbr
}

func (v *Video) SetTbr(tbr float32) {
	v.tbr = tbr
}

func (v *Video) Sort() int {
	return v.sort
}

func (v *Video) SetSort(sort int) {
	v.sort = sort
}

func (v *Video) TotalSecond() float32 {
	return v.totalSecond
}

func (v *Video) SetTotalSecond(totalSecond float32) {
	v.totalSecond = totalSecond
}

func (v *Video) Name() string {
	return v.name
}

func (v *Video) SetName(name string) {
	v.name = name
}

func (v *Video) Suffix() string {
	return v.suffix
}

func (v *Video) SetSuffix(suffix string) {
	v.suffix = suffix
}

func (v *Video) Cover() string {
	return v.cover
}

func (v *Video) SetCover(cover string) {
	v.cover = cover
}

//func (v *Video) SetCover(savePath string) bool {
//
//	path, _ := os.Getwd()
//	ffmpegPath := path + "\\ffmpeg\\bin\\ffmpeg.exe"
//	savePath = savePath + "\\images\\" + v.Name() + ".png"
//	cmdArguments := []string{"-i", strings.Replace(v.Path(), "\\", "/", -1), "-hide_banner", "-v", "error", "-ss", "00:00:00.1", "-s", "160*90", "-y", "-f", "image2", "-frames:v", "1", savePath}
//	cmd := exec.Command(ffmpegPath, cmdArguments...)
//	var stdout bytes.Buffer
//	var stderr bytes.Buffer
//	var err error
//
//	cmd.Stdout = &stdout
//	cmd.Stderr = &stderr
//
//	if err = cmd.Run(); err != nil {
//		color.Red(v.Path() + "封面图处理错误")
//		return false
//	}
//	v.cover = savePath
//	return true
//}

func (v *Video) Duration() string {
	return v.duration
}

func (v *Video) SetDuration(savePath string) int {

	ffmpegExePath, _ := os.Getwd()
	ffmpegPath := ffmpegExePath + "\\ffmpeg\\bin\\ffmpeg.exe"
	savePath = savePath + "\\images\\" + v.Name() + ".png"
	cmdArguments := []string{"-i", strings.Replace(v.Path(), "\\", "/", -1), "-hide_banner", "-ss", "00:00:00.1", "-s", "160*90", "-y", "-f", "image2", "-frames:v", "1", savePath}
	cmd := exec.Command(ffmpegPath, cmdArguments...)
	out, err := cmd.CombinedOutput()
	if err == nil {
		v.SetCover(savePath)
		outStr := string(out)
		durationReg := regexp.MustCompile("Duration: (.*?),")
		durationArr := durationReg.FindAllStringSubmatch(outStr, 1)
		if len(durationArr) != 0 {
			v.duration = durationArr[0][1]
		} else {
			return 1
		}

		tbrReg := regexp.MustCompile("(([1-9]\\d*\\.?\\d*)|(0\\.\\d*[1-9])) tbr")
		tbrArr := tbrReg.FindAllStringSubmatch(outStr, 1)
		if len(tbrArr) != 0 {
			tbr, err2 := strconv.ParseFloat(tbrArr[0][1], 32)
			if err2 != nil {
				return 2
			}
			v.SetTbr(float32(tbr))
		} else {
			return 3
		}

		durationSplit := strings.Split(v.Duration(), ":")
		if len(durationSplit) != 3 {
			return 4
		}

		hour, err3 := strconv.ParseFloat(durationSplit[0], 32)
		if err3 != nil {
			return 5
		} else {
			hour = hour * 3600
		}
		minute, err4 := strconv.ParseFloat(durationSplit[1], 32)
		if err4 != nil {
			return 6
		} else {
			minute = minute * 60
		}
		second, err5 := strconv.ParseFloat(durationSplit[2], 32)
		if err5 != nil {
			return 7
		} else {
			////取出帧数
			//temp := float32(second) - float32(int(second))
			////帧数转换秒
			//temp2 := temp * v.tbr
			//tempSecond := float32(int(second)) + float32(temp2)
			tempTotalSecond := float32(hour) + float32(minute) + float32(second)

			a, error6 := strconv.ParseFloat(fmt.Sprintf("%.4f", tempTotalSecond), 32)
			if error6 != nil {
				return 8
			} else {
				v.SetTotalSecond(float32(a))
			}
		}
	} else {
		return 9
	}

	return 0
}

//func (v *Video) SetDuration(duration string) bool {
//
//	path, _ := os.Getwd()
//	ffmpegPath := path + "\\ffmpeg\\bin\\ffprobe.exe"
//	cmdArguments := []string{"-i", strings.Replace(v.Path(), "\\", "/", -1), "-hide_banner", "-v", "error", "-show_entries", "format=duration", "-of", "csv=p=0"}
//	cmd := exec.Command(ffmpegPath, cmdArguments...)
//	var stdout bytes.Buffer
//	var stderr bytes.Buffer
//	var err error
//
//	cmd.Stdout = &stdout
//	cmd.Stderr = &stderr
//
//	if err = cmd.Run(); err != nil {
//		color.Red(v.Path() + "时长处理错误")
//		return false
//	} else {
//		timeStr := strings.Replace(stdout.String(), " ", "", -1)
//		timeStr = strings.Replace(timeStr, "\n", "", -1)
//		timeStr = strings.Replace(timeStr, "\r", "", -1)
//		timeFloat, _ := strconv.ParseFloat(timeStr, 32)
//
//		v.SetTotalSecond(float32(timeFloat))
//
//		timeTool := utils.TimeTool{}
//		_, hour, minute, second, frame := timeTool.ResolveTime(v.TotalSecond())
//
//		v.duration = fmt.Sprintf("%v:%v:%v:%v", hour, minute, second, frame)
//	}
//	return true
//}

func (v *Video) Path() string {
	return v.path
}

func (v *Video) SetPath(path string) {
	v.path = path
}

func (v *Video) SetBase(pathStr string, savePath string) bool {
	//设置文件路径
	v.SetPath(pathStr)
	//文件全称
	fullFileName := path.Base(v.Path())
	//后缀
	suffix := path.Ext(fullFileName)
	//文件名
	fileName := strings.TrimSuffix(filepath.Base(fullFileName), suffix)
	v.SetName(fileName)
	v.SetSuffix(suffix)
	//b1 := v.SetDuration("")
	//b2 := v.SetCover(savePath)
	//if b1 && b2 {
	//	return true
	//}

	b1 := v.SetDuration(savePath)
	if b1 == 0 {
		return true
	}
	return false
}

func NewVideo(videoPath string, savePath string) *Video {
	if path.Ext(videoPath) == ".mp4" {
		tempVideo := new(Video)
		if tempVideo.SetBase(videoPath, savePath) {
			return tempVideo
		} else {
			return nil
		}
	}
	return nil
}

// 2.声明一个Hero结构体切片类型
type VideoSlice []*Video

func (vs VideoSlice) Len() int {
	return len(vs)
}
func (vs VideoSlice) Less(i, j int) bool {
	aa := []byte(vs[i].name)
	bb := []byte(vs[j].name)
	if len(aa) < len(bb) {
		return true
	} else if len(aa) > len(bb) {
		return false
	} else {
		a := strings.Compare(vs[i].name, vs[j].name)
		if a > 0 {
			return false
		} else {
			return true
		}
	}
}
func (vs VideoSlice) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}
