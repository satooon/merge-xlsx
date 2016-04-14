package main

import (
	"io/ioutil"
	"os"

	"regexp"

	"fmt"
	"strings"

	"log"

	"sync"

	"github.com/codegangsta/cli"
	"github.com/tealeg/xlsx"
)

const (
	csvMarker = "csv@"
)

func main() {
	app := NewApp()
	app.Run(os.Args)
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "marge-xlsx"
	app.Usage = "marge-xlsx"
	app.Author = "satooon"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, vv",
			Usage: "verbose mode",
		},
	}
	app.Action = Action
	return app
}

func Action(ctx *cli.Context) {
	verbose := ctx.Bool("verbose")
	debugLog = debugLogStruct{verbose: verbose}

	fileNames, err := getFiles(ctx.Args())
	if err != nil {
		panic(err)
	}

	reg := regexp.MustCompile(csvMarker)
	csvMap := map[string]string{}
	for _, fileName := range fileNames {
		file, err := xlsx.OpenFile(fileName)
		if err != nil {
			panic(err)
		}
		for _, sheet := range file.Sheets {
			debugLog.Println("sheet.Name", sheet.Name)
			if reg.MatchString(sheet.Name) == false {
				continue
			}
			sheetName := strings.Replace(sheet.Name, csvMarker, "", 1)
			if _, ok := csvMap[sheetName]; ok == false {
				csvMap[sheetName] = ""
			}
			rows := make([]string, len(sheet.Rows))
			for i, row := range sheet.Rows {
				if i <= 0 && csvMap[sheetName] != "" {
					continue
				}
				cells := make([]string, len(row.Cells))
				for j, cell := range row.Cells {
					if err != nil {
						panic(err)
					}
					if i <= 0 {
						cells[j] = cell.Value
						continue
					}
					if number, err := cell.Int(); err == nil {
						cells[j] = fmt.Sprintf("%d", number)
					} else {
						reg := regexp.MustCompile("\n")
						str := reg.ReplaceAllString(cell.Value, "\\n")
						cells[j] = fmt.Sprintf("\"%s\"", str)
					}
				}
				rows[i] = strings.Join(cells, ",")
			}
			csvMap[sheetName] += strings.Join(rows, "\n")
		}
		debugLog.Println("csvMap", csvMap)
	}

	currentDirName, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := os.RemoveAll(currentDirName + "/csv"); err != nil {
		log.Println("os.Remove(./csv)", err)
	}
	if err := os.Mkdir(currentDirName+"/csv", os.ModePerm); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for i, v := range csvMap {
		wg.Add(1)
		go func(dirName, sheetName, csv string) {
			log.Println("hogehhgoe")
			filePath := fmt.Sprintf("%s/%s.csv", dirName, sheetName)
			debugLog.Println("filePath", filePath)
			log.Println("err", ioutil.WriteFile(filePath, []byte(csv), os.ModePerm))
			wg.Done()
		}(currentDirName+"/csv", i, v)
	}
	wg.Wait()
}

func getFiles(args cli.Args) (fileNames []string, err error) {
	debugLog.Println("args", args)
	if len(args) > 0 {
		for _, v := range args {
			fileNames = append(fileNames, v)
		}
	} else {
		dirName, err := os.Getwd()
		if err != nil {
			return fileNames, err
		}
		debugLog.Println("dirName", dirName)
		fileInfos, err := ioutil.ReadDir(dirName)
		if err != nil {
			return fileNames, err
		}
		reg := regexp.MustCompile("xlsx$")
		for _, v := range fileInfos {
			if v.IsDir() {
				continue
			}
			if reg.MatchString(v.Name()) == false {
				continue
			}
			fileNames = append(fileNames, v.Name())
		}
	}
	return
}
