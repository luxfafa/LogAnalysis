package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func (ce *CountAndExport) AnalysisStart() {
	go ce.AnalysisLog()
	go ce.AnalysisLine()
	go ce.CountTask()
}

/**
*  缓冲区 按行读日志文件
 */
func (ce *CountAndExport) AnalysisLog() {
	logResource, err := os.Open(ce.LogPath)
	if err != nil {
		panic(err)
	}

	defer func() {
		logResource.Close()
		close(ce.lineCh)
		ce.cancelCh <- struct{}{}
		close(ce.cancelCh)
	}()

	bufIoReaderForLog := bufio.NewReader(logResource)
	for {
		lineByte, err := bufIoReaderForLog.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		ce.lineCh <- strings.Trim(lineByte, " ")
	}

}

/**
* 每行string的处理 可拓展
 */
func (ce *CountAndExport) AnalysisLine() {
	rxpFulUrl := regexp.MustCompile(`[a-zA-z]+://[^\s]*`)

WAITLOOP:
	for {
		select {
		case linestr := <-ce.lineCh:

			rxpRes := rxpFulUrl.Find([]byte(linestr))
			if rxpResLen := len(rxpRes); rxpResLen > 0 {
				ce.countUrlCh <- string(rxpRes[0 : rxpResLen-1])
			}

		case <-ce.cancelCh:

			close(ce.countUrlCh)
			break WAITLOOP

		default:

		}
	}
}

/**
*  统计
 */
func (ce *CountAndExport) CountTask() {
WAITLOOP:
	for {
		select {
		case urlStr, ok := <-ce.countUrlCh:
			if !ok {
				ce.SortData()
				break WAITLOOP
			}
			if _, keyExists := ce.countMap[urlStr]; keyExists {
				ce.countMap[urlStr]++
			} else {
				ce.countMap[urlStr] = 1
			}
		default:
		}
	}

}

/**
*  排序
 */
func (ce *CountAndExport) SortData() {
	num := 0
tmpUrlForSortSlice := make(UrlSortSlice,len(ce.countMap))
	for k, v := range ce.countMap {
		tmpUrlForSortSlice[num] = UrlForSort{
			Often: v,
			Url:   k,
		}
		num++
	}

	sort.Sort(tmpUrlForSortSlice)
	ce.urlForSortSlice = tmpUrlForSortSlice[:]
	ce.ExportCsv()
}

/**
*  导出整理完毕的切片为csv
 */
func (ce *CountAndExport) ExportCsv() {
	defer ce.cancel()
	csvFile, err := os.Create("aka.xls")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvFile.WriteString("\xEF\xBB\xBF")
	theCsv := csv.NewWriter(csvFile)
	theCsv.Write([]string{"次数", "网址"})
	for _, v := range ce.urlForSortSlice {

		oftenStr := strconv.Itoa(v.Often)
		theCsv.Write([]string{oftenStr, v.Url})

	}
	theCsv.Flush()
	fmt.Printf("共有: %d条数据\n", len(ce.urlForSortSlice))
}
