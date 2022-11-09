package main

import (
	"bytes"
	"fmt"
	"github.com/linjifu/videoHandle/models"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func zh(str string) (i int, time float32) {
	ptsTimeReg2 := regexp.MustCompile("n:\\s*\\d*")
	IndexStr := ptsTimeReg2.FindString(str)
	IndexArr := strings.Split(IndexStr, ":")
	index := strings.Trim(IndexArr[1], " ")
	tempIndex, _ := strconv.Atoi(index)

	ptsTimeReg3 := regexp.MustCompile("pts_time:\\s*(([1-9]\\d*\\.?\\d*)|(0\\.\\d*[1-9]))")
	startTimeStr := ptsTimeReg3.FindString(str)
	startTimeArr := strings.Split(startTimeStr, ":")
	startTime := strings.Trim(startTimeArr[1], " ")
	tempStartTime, _ := strconv.ParseFloat(startTime, 32)

	return tempIndex + 1, float32(tempStartTime)
}

func processor(video chan *models.VideoSource) {
	num := 0
	ffmpegExePath, _ := os.Getwd()
	ffmpegPath := ffmpegExePath + "\\ffmpeg\\bin\\ffmpeg.exe"
	for {
		select {
		case v := <-video:

			startTime := fmt.Sprintf("%v", v.StartTime())
			endTime := fmt.Sprintf("%v", v.EndTime())

			var cmdArguments2 = []string{}
			if startTime == "0" {
				cmdArguments2 = []string{"-i", "f:\\videos\\2.mp4", "-hide_banner", "-ss", "0.00", "-to", endTime, "-c", "copy", "-y", fmt.Sprintf("F:\\videos\\test\\%v.mp4", v.Index())}
			} else if endTime == "0" {
				cmdArguments2 = []string{"-i", "f:\\videos\\2.mp4", "-hide_banner", "-ss", startTime, "-c", "copy", "-y", fmt.Sprintf("F:\\videos\\test\\%v.mp4", v.Index())}
			} else {
				cmdArguments2 = []string{"-i", "f:\\videos\\2.mp4", "-hide_banner", "-ss", startTime, "-to", endTime, "-c", "copy", "-y", fmt.Sprintf("F:\\videos\\test\\%v.mp4", v.Index())}
			}
			if len(cmdArguments2) > 0 {
				cmd2 := exec.Command(ffmpegPath, cmdArguments2...)
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				var err2 error

				cmd2.Stdout = &stdout
				cmd2.Stderr = &stderr

				if err2 = cmd2.Run(); err2 != nil {
					fmt.Println("失败", stderr.String())

				} else {
					fmt.Println("成功")
				}
			}

		default:
			if num > 5 {
				fmt.Println("退出了一个协程")
				return
			} else {
				num += 1
				time.Sleep(time.Second * 1)
			}
		}
	}
}

func main() {
	ffmpegExePath, _ := os.Getwd()
	ffmpegPath := ffmpegExePath + "\\ffmpeg\\bin\\ffmpeg.exe"
	cmdArguments := []string{"-i", "f:\\videos\\2.mp4", "-hide_banner", "-vf", "select='gt(scene,0.1)',showinfo", "-f", "null", "-"}
	cmd := exec.Command(ffmpegPath, cmdArguments...)
	out, err := cmd.CombinedOutput()
	if err == nil {
		outStr := string(out)
		ptsTimeReg := regexp.MustCompile("n:\\s*\\d* pts:\\s*\\d* pts_time:\\s*(([1-9]\\d*\\.?\\d*)|(0\\.\\d*[1-9]))")
		ptsTimeArr := ptsTimeReg.FindAllString(outStr, -1)

		dataArr := models.VideoSourceSlice{}
		for i, v := range ptsTimeArr {
			if i == 0 {
				_, startTime := zh(v)
				videoSource := models.NewVideoSource(0, 0.00, startTime, "")
				dataArr = append(dataArr, videoSource)
			}

			index, startTime2 := zh(v)
			videoSource := models.NewVideoSource(index, startTime2, 0.00, v)
			dataArr = append(dataArr, videoSource)
		}
		//排序
		sort.Sort(dataArr)
		dataArrLen := len(dataArr)
		for i, v := range dataArr {
			if i+1 != dataArrLen {
				_, startTime := zh(dataArr[i+1].DataStr())
				v.SetEndTime(startTime)
			} else {
				v.SetEndTime(0.00)
			}
		}

		imageChan := make(chan *models.VideoSource, 1000)
		for i := 0; i < 8; i++ {
			fmt.Println("开启了一个协程")
			go processor(imageChan)
		}

		for _, v := range dataArr {
			imageChan <- v
			//startTime := fmt.Sprintf("%v", v.StartTime())
			//endTime := fmt.Sprintf("%v", v.EndTime())
			//
			//fmt.Printf("%v", startTime)
			//fmt.Printf("\t")
			//fmt.Printf("%v", endTime)
			//fmt.Println()
			//var cmdArguments2 = []string{}
			//if startTime == "0" {
			//	cmdArguments2 = []string{"-i", "f:\\videos\\2.mp4", "-hide_banner", "-ss", "0.00", "-to", endTime, "-c", "copy", "-y", fmt.Sprintf("F:\\videos\\test\\%v.mp4", v.Index())}
			//} else if endTime == "0" {
			//	cmdArguments2 = []string{"-i", "f:\\videos\\2.mp4", "-hide_banner", "-ss", startTime, "-c", "copy", "-y", fmt.Sprintf("F:\\videos\\test\\%v.mp4", v.Index())}
			//} else {
			//	cmdArguments2 = []string{"-i", "f:\\videos\\2.mp4", "-hide_banner", "-ss", startTime, "-to", endTime, "-c", "copy", "-y", fmt.Sprintf("F:\\videos\\test\\%v.mp4", v.Index())}
			//}
			//if len(cmdArguments2) > 0 {
			//	cmd2 := exec.Command(ffmpegPath, cmdArguments2...)
			//	var stdout bytes.Buffer
			//	var stderr bytes.Buffer
			//	var err2 error
			//
			//	cmd2.Stdout = &stdout
			//	cmd2.Stderr = &stderr
			//
			//	if err2 = cmd2.Run(); err2 != nil {
			//		fmt.Println("失败", stderr.String())
			//
			//	} else {
			//		fmt.Println("成功")
			//	}
			//}
		}

		for {
			time.Sleep(time.Second * 30)
		}
	}
}
