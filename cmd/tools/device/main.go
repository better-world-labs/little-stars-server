package main

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/tencent"
	"flag"
	"time"

	"aed-api-server/internal/pkg/config"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/utils"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const pageSize = 20

func readFile_V2(filepath string) []entities.AddDevice {
	var list []entities.AddDevice
	filepath = "device_v2.csv"
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return list
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		arr := strings.Split(line, ",")
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return list
			}
		}
		d := entities.AddDevice{
			Longitude: utils.ToFloat(arr[3]),
			Latitude:  utils.ToFloat(arr[2]),
			Address:   arr[0],
			// Contract:  arr[2],
			Title: arr[1],
			State: 1,
		}
		if d.Title == "" {
			d.Title = d.Address
		}

		if d.Address == "成都市红十字会" && d.Title == "天府大道1480号拉德方斯大厦西楼5层" {
			d.DeviceImage = "https://openview-oss.oss-cn-chengdu.aliyuncs.com/star-static/image1.png.compress.jpeg"
		} else if d.Address == "天堂岛海洋乐园" && d.Title == "4号门内" {
			d.DeviceImage = "https://openview-oss.oss-cn-chengdu.aliyuncs.com/star-static/image3.png.compress.jpeg"
		} else if d.Address == "新世纪环球购物中心c" && d.Title == "一楼服务台" {
			d.DeviceImage = "https://openview-oss.oss-cn-chengdu.aliyuncs.com/star-static/image5.png.compress.jpeg"
		}

		list = append(list, d)

		fmt.Printf("%+v\n", d)
	}

	return list
}

func main() {
	configPath := "../../../config-testing.yaml"
	c, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	tencent.Init(&c.MapConfig)
	db.InitEngine(c.Database)

	var cmd = flag.String("cmd", "cmd", "cmd to run")
	var file = flag.String("file", "file", "file name")
	// parse
	flag.Parse()
	fmt.Print(*cmd)
	fmt.Println(*file)
	sess := db.GetSession()
	defer sess.Close()

	var del []string
	switch *cmd {
	case "del":
		var index = 1
		for {
			list, err := tencent.ListDevice(index, pageSize)
			if err != nil {
				panic(err)
			}
			if len(list) == 0 {
				break
			}
			index += 1
			for _, v := range list {
				del = append(del, v.UdID)
			}
		}

		k := 20
		run := true
		for run && len(del) > 0 {
			var dels []string
			time.Sleep(time.Second * 1)

			if len(del) >= k {
				dels = del[:k]
				del = del[k:]
			} else {
				dels = del
				run = false
			}
			fmt.Println("dels", dels)
			err = tencent.DelDevice(dels)
			if err != nil {
				panic(err)
			}
		}

		_, err = sess.Exec("truncate table device")
		if err != nil {
			panic(err)
		}

	case "import":
		list := readFile_V2(*file)
		for _, req := range list {
			de := new(entities.Device)
			de.Address = req.Address
			de.Title = req.Title
			de.Latitude = req.Latitude
			de.Longitude = req.Longitude
			de.Tel = req.Contract
			de.DeviceImage = req.DeviceImage

			id, err := tencent.AddDevice(req.Longitude, req.Latitude, req.Title)
			if err != nil {
				panic(err)
			}

			de.Id = id
			fmt.Println("import ", de)
			_, err = sess.Insert(de)
			if err != nil {
				panic(err)
			}
		}

	case "list":
		var index = 1
		for {
			fmt.Println("pageIndex:", index)
			list, err := tencent.ListDevice(index, pageSize)
			if err != nil {
				panic(err)
			}

			if len(list) == 0 {
				break
			}
			index += 1
			for k, v := range list {
				fmt.Println(k, v)
			}
		}
	}
}
