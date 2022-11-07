package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/linjifu/videoHandle2/models"
	"github.com/linjifu/videoHandle2/utils"
	"github.com/xuri/excelize/v2"
	"image"
	color2 "image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 界面
func Showinput() (dirPath string, savePath string) {
	d := color.New(color.FgBlue, color.Bold)
	d.Printf("请输入视频文件夹地址（默认为：F:\\videos）:")
	fmt.Scan(&dirPath)
	if dirPath == "1" {
		dirPath = "F:\\videos\\分析"
	} else {
		dirPath = strings.Replace(dirPath, "\n", "", -1)
		dirPath = strings.Replace(dirPath, "\r", "", -1)
	}
	_, err := os.Stat(dirPath)
	if err != nil {
		color.Red("视频文件夹不存在")
		time.Sleep(time.Second * 10)
		os.Exit(0)
	}
	color.Green("视频地址设置成功：" + dirPath)

	nowTimeFormat := time.Now().Format("2006-01-02-150405")
	d.Printf("请输入文件保存地址（默认为：F:\\videos\\%v）:", nowTimeFormat)
	fmt.Scan(&savePath)
	if savePath == "1" {
		savePath = "F:\\videos\\" + nowTimeFormat
	} else {
		savePath = strings.Replace(savePath, "\n", "", -1)
		savePath = strings.Replace(savePath, "\r", "", -1)
		savePath = strings.Trim(savePath, "\\")
		savePath = strings.Trim(savePath, "/")

		savePath = savePath + "\\" + nowTimeFormat
	}

	_, err = os.Stat(savePath)
	if err != nil {
		err = os.Mkdir(savePath, os.ModePerm)
		if err != nil {
			color.Red("创建文件保存地址失败" + err.Error())
			time.Sleep(time.Second * 10)
			os.Exit(0)
		}

		err = os.Mkdir(savePath+"\\images", os.ModePerm)
		if err != nil {
			color.Red("创建封面图保存地址失败" + err.Error())
			time.Sleep(time.Second * 10)
			os.Exit(0)
		}
	}
	color.Green("保存地址设置成功：" + savePath)

	color.Green("**************************************************************")
	color.Green("***开始处理视频文件***")
	color.Green("**************************************************************")

	return
}

// 读取文件夹中的文件
func putFile(fileChan chan string, dirPath string) {
	dirPath = strings.Replace(dirPath, "/", "\\", -1)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		color.Red("视频文件夹地址错误,无法打开")
		time.Sleep(time.Second * 10)
		os.Exit(0)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dirPath, file.Name())
			fileChan <- filePath
			color.Cyan("读取到文件：" + filePath)
		}
	}
}

// 载入视频
func loadVide(tempFilePath string, videoChan chan *models.Video, savePath string) {
	videoInfo := models.NewVideo(tempFilePath, savePath)
	if videoInfo != nil {
		videoChan <- videoInfo
		color.Cyan("生成视频基础信息与封面图：" + videoInfo.Path())
	}
}

// 生成背景图
func createBackgroundImage(len int) *image.RGBA {
	color.Cyan("开始生成画板背景图")
	width := 160            //单张图片宽度
	height := 90            //单张图片高度
	jg := 10                //宽间隔
	jg2 := 10               //高间隔
	imagesNum := len        //图片数量
	rowNum := 50            //一行几张
	totalWidth := 0         //总宽度
	totalHeight := 0        //总高度
	if imagesNum > rowNum { //如果图片数量大于每行数量   总宽度 = 每行数量 * (图片宽度 + 宽间隔) + 第一张间隔
		totalWidth = rowNum*(width+jg) + jg
	} else { //如果图片数量不大于每行数量   总宽度 = 图片数量 * (图片宽度 + 宽间隔) + 第一张间隔
		totalWidth = imagesNum*(width+jg) + jg
	}

	if imagesNum <= rowNum { //如果图片数量小于等于每行数量	总高度 = 图片高度 + 高间隔 * 2
		totalHeight = height + jg2*2
	} else { //如果图片数量大于每行数量
		if imagesNum%rowNum == 0 { //如果刚好铺满 总高度 = 图片数量 / 每行数量 * （图片高度 + 高间隔）+ 第一行间隔
			totalHeight = imagesNum/rowNum*(height+jg2) + jg2
		} else { //如果铺不满则代表多一行 总高度 = （图片数量 / 每行数量 + 1） * （图片高度 + 高间隔）+ 第一行间隔
			totalHeight = (imagesNum/rowNum+1)*(height+jg2) + jg2
		}
	}
	img := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))
	for x := 0; x < img.Bounds().Dx(); x++ { // 将背景图涂黑
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, color2.Black)
		}
	}

	color.Cyan("画板背景图生成完成")
	return img
}

// 生成图片
func createImage(img *image.RGBA, video *models.Video) {
	width := 160 //单张图片宽度
	height := 90 //单张图片高度
	jg := 10     //宽间隔
	jg2 := 10    //高间隔
	rowNum := 50 //一行几张

	imageIndex := video.Sort() + 1

	f, err := os.Open(video.Cover())
	if err != nil {
		fmt.Println("打开图片错误：" + err.Error())
		return
	}

	x0 := (imageIndex - 1) % rowNum
	y0 := (imageIndex - 1) / rowNum

	x3 := jg + x0*(width+jg)
	y3 := jg2 + y0*(height+jg2)

	color.Cyan(video.Path() + "开始拼接画板")

	gopherImg, err := png.Decode(f)
	draw.Draw(img, img.Bounds(), gopherImg, image.Pt(-x3, -y3), draw.Over)
	err2 := f.Close()
	if err2 != nil {
		fmt.Println("关闭图片错误：" + err2.Error())
		return
	}

	color.Cyan(video.Path() + "画板拼接完成")
}

// 保存图片
func saveImage(img *image.RGBA, NumCPU int, savePath string, imageChanExitChan chan bool, mainExitChan chan bool) {
	for i := 0; i < NumCPU; i++ {
		<-imageChanExitChan
	}
	color.Cyan("图片拼接画板完成,开始导出总图")
	outFile, err3 := os.Create(savePath + "\\cover.jpeg")
	defer outFile.Close()
	if err3 != nil {
		fmt.Println(err3.Error())
	}
	b := bufio.NewWriter(outFile)
	err3 = jpeg.Encode(b, img, &jpeg.Options{Quality: 80})
	//err3 = png.Encode(b, img)
	if err3 != nil {
		fmt.Println(err3.Error())
	}
	err3 = b.Flush()
	if err3 != nil {
		fmt.Println(err3.Error())
	}
	color.Cyan("总图导出完成")
	mainExitChan <- true
}

// 生成表格
func createTable(savePath string, totalVideoNum int, TotalTime float32, maxTimeStr string, minTimeStr string, pjTimeStr string, mainExitChan chan bool) {
	color.Cyan("开始生成excel")
	f := excelize.NewFile() // 设置单元格的值
	// 设置表头样式
	headStyleID, _ := f.NewStyle(`{
   "font":{
      "color":"#333333",
      "bold":true,
      "size":16,
      "family":"arial"
   },
   "alignment":{
      "vertical":"center",
      "horizontal":"center"
   }
}`)

	f.SetCellStyle("Sheet1", "A1", "E1", headStyleID)
	f.SetCellStyle("Sheet1", "A3", "D3", headStyleID)

	textLeftStyleID, _ := f.NewStyle(`{
   "alignment":{
      "horizontal":"left"
   }
}`)

	f.SetColWidth("Sheet1", "A", "E", 30)
	// 设置行样式
	f.SetCellStyle("Sheet1", "A2", "E2", textLeftStyleID)
	// 这里设置表头
	f.SetCellValue("Sheet1", "A1", "视频数量")
	f.SetCellValue("Sheet1", "B1", "总秒数")
	f.SetCellValue("Sheet1", "C1", "最长时间")
	f.SetCellValue("Sheet1", "D1", "最短时间")
	f.SetCellValue("Sheet1", "E1", "平均时间")

	f.SetCellValue("Sheet1", "A2", totalVideoNum)
	f.SetCellValue("Sheet1", "B2", TotalTime)
	f.SetCellValue("Sheet1", "C2", maxTimeStr)
	f.SetCellValue("Sheet1", "D2", minTimeStr)
	f.SetCellValue("Sheet1", "E2", pjTimeStr)

	f.SetCellValue("Sheet1", "A3", "视频名称")
	f.SetCellValue("Sheet1", "B3", "视频时长")
	f.SetCellValue("Sheet1", "C3", "视频秒数")
	f.SetCellValue("Sheet1", "D3", "视频帧率")

	line := 3
	for _, v := range videos {
		line++
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", line), v.Name()+v.Suffix())
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", line), v.Duration())
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", line), v.TotalSecond())
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", line), v.Tbr())

		// 设置行样式
		f.SetCellStyle("Sheet1", fmt.Sprintf("A%d", line), fmt.Sprintf("D%d", line), textLeftStyleID)
	}

	// 保存文件
	if err := f.SaveAs(savePath + "\\detail.xlsx"); err != nil {
		fmt.Println(err.Error())
	}

	color.Cyan("excel生成完成")
	mainExitChan <- true
}

// 视频协程处理
func videoProcessor(fileChan chan string, videoChan chan *models.Video, savePath string, exitChan chan bool) {

	num := 0

	for {
		select {
		case filePath := <-fileChan:
			tempFilePath := filePath
			loadVide(tempFilePath, videoChan, savePath)
		case video := <-videoChan:
			tempVideo := video
			videos = append(videos, tempVideo)
		default:
			if num > 5 {
				exitChan <- true
				fmt.Println("退出了一个视频处理协程")
				return
			} else {
				num += 1
				time.Sleep(time.Second * 1)
			}
		}
	}

}

// 封面图协程处理
func coverProcessor(img *image.RGBA, imageChan chan *models.Video, exitChan chan bool) {

	num := 0

	for {
		select {
		case video := <-imageChan:
			createImage(img, video)
		default:
			if num > 5 {
				exitChan <- true
				fmt.Println("退出了一个封面图协程")
				return
			} else {
				num += 1
				time.Sleep(time.Second * 1)
			}
		}
	}

}

var videos models.VideoSlice

func main() {
	//操作界面
	dirPath, savePath := Showinput()

	NumCPU := runtime.NumCPU()
	mainExitChan := make(chan bool, 2)
	exitChan := make(chan bool, NumCPU)
	fileChan := make(chan string, 1000)
	videoChan := make(chan *models.Video, 1000)

	start := time.Now().Unix()

	go putFile(fileChan, dirPath)

	time.Sleep(time.Second * 3)

	maxTime := float32(0)   //最长时间
	maxTimeStr := ""        //最长时间字符串格式
	minTime := float32(0)   //最小时间
	minTimeStr := ""        //最小时间字符串格式
	TotalTime := float32(0) //总时间
	totalVideoNum := 0      //总视频数量
	pjTimeStr := ""         //平均时间
	timeTool := utils.TimeTool{}

	//开启多携程处理
	for i := 0; i < NumCPU; i++ {
		fmt.Println("开启了一个视频处理协程")
		go videoProcessor(fileChan, videoChan, savePath, exitChan)
	}
	for i := 0; i < NumCPU; i++ {
		<-exitChan
	}

	color.Green("**************************************************************")
	color.Green("***处理视频文件完成，开始生成总图与excel表格***")
	color.Green("**************************************************************")

	//数据计算
	totalVideoNum = len(videos)
	if totalVideoNum > 0 {
		//排序
		sort.Sort(videos)

		imageChan := make(chan *models.Video, totalVideoNum)
		imageChanExitChan := make(chan bool, NumCPU)
		img := createBackgroundImage(totalVideoNum)

		color.Cyan("开始图片拼接画板")
		for i := 0; i < NumCPU; i++ {
			fmt.Println("退出了一个封面图协程")
			go coverProcessor(img, imageChan, imageChanExitChan)
		}
		go saveImage(img, NumCPU, savePath, imageChanExitChan, mainExitChan)

		for index, video := range videos {
			video.SetSort(index)
			imageChan <- video

			TotalTime += video.TotalSecond()
			if video.TotalSecond() > maxTime {
				maxTime = video.TotalSecond()
				//_, hour, minute, second, frame := timeTool.ResolveTime(maxTime)
				//maxTimeStr = fmt.Sprintf("%v:%v:%v:%v", hour, minute, second, frame)
				maxTimeStr = video.Duration()
			}
			if minTime == 0 {
				minTime = video.TotalSecond()
				minTimeStr = video.Duration()
			} else if video.TotalSecond() < minTime {
				minTime = video.TotalSecond()
				minTimeStr = video.Duration()
			}
			//if minTime != 0 {
			//	_, hour, minute, second, frame := timeTool.ResolveTime(minTime)
			//	minTimeStr = fmt.Sprintf("%v:%v:%v:%v", hour, minute, second, frame)
			//}
		}

		pjTime := TotalTime / float32(totalVideoNum)
		a, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", pjTime), 32)
		pjTime = float32(a)
		_, hour, minute, second := timeTool.ResolveTime(pjTime)
		pjTimeStr = fmt.Sprintf("%v:%v:%v", hour, minute, second)

		go createTable(savePath, totalVideoNum, TotalTime, maxTimeStr, minTimeStr, pjTimeStr, mainExitChan)

		for i := 0; i < 2; i++ {
			<-mainExitChan
		}
	}

	color.Green("**************************************************************")
	color.Green("***全部处理完成***")
	color.Green("**************************************************************")

	fmt.Printf("视频数量：%v\t总秒数：%v\t最长时间：%v\t最短时间：%v\t平均时间：%v\t", totalVideoNum, TotalTime, maxTimeStr, minTimeStr, pjTimeStr)
	fmt.Println()
	end := time.Now().Unix()
	fmt.Println("总耗时=", end-start, "秒")
	fmt.Println("请手动关闭程序")
	for {
		time.Sleep(time.Second * 30)
	}
}
