package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

/**
*  logPath  日志绝对路径
*  cancelCh 单向通道 日志读取完毕的讯号
*  lineCh   单向通道 存入单行处理过的日志
 */
func AnalysisLog(logPath string, cancelCh chan<- struct{}, lineCh chan<- string) {
	logResource, err := os.Open(LogFilePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		logResource.Close()
		close(lineCh)
		cancelCh <- struct{}{}
		close(cancelCh)
	}()
	bufIoReaderForLog := bufio.NewReader(logResource)
	for {
		lineByte, err := bufIoReaderForLog.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		lineCh <- strings.Trim(lineByte, " ")
	}

}

/**
*  lineCh   单向通道 接收单行处理过的日志
*  cancelCh 单向通道 接收日志存取完毕的讯号
 */
func AnalysisLine(lineCh <-chan string, cancelCh <-chan struct{}) {
	rxpFulUrl := regexp.MustCompile(`[a-zA-z]+://[^\s]*`)

WAITLOOP:
	for {
		select {
		case linestr := <-lineCh:

			rxpRes := rxpFulUrl.Find([]byte(linestr))
			if rxpResLen := len(rxpRes); rxpResLen > 0 {
				countUrlCh <- string(rxpRes[0 : rxpResLen-1])
			}

		case <-cancelCh:

			close(countUrlCh)
			break WAITLOOP

		default:

		}
	}
}

/**
*  cancelFunc 上下文控制
 */
func CountTask(cancelFunc context.CancelFunc) {
	var countMap = make(map[string]int, 10000)
WAITLOOP:
	for {
		select {
		case urlStr, ok := <-countUrlCh:
			if !ok {
				showData(countMap, cancelFunc)
				break WAITLOOP
			}
			if _, keyExists := countMap[urlStr]; keyExists {
				countMap[urlStr]++
			} else {
				countMap[urlStr] = 1
			}
		default:
		}
	}

}

/**
*  countMap 过渡值排序的载体
*  cancelFunc 上下文控制
 */
func showData(countMap map[string]int, cancelFunc context.CancelFunc) {

	defer cancelFunc()
	UrlSortSlice := make(UrlSortSlice, len(countMap))
	num := 0
	for k, v := range countMap {
		UrlSortSlice[num] = UrlForSort{
			Often: v,
			Url:   k,
		}
		num++
	}
	sort.Sort(UrlSortSlice)
	for _, v := range UrlSortSlice {
		fmt.Printf("%d : %s\n", v.Often, v.Url)
	}
	fmt.Printf("共有: %d条数据\n", len(UrlSortSlice))
}

func main() {
	lineCh := make(chan string, 50000)
	cancelCh := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go AnalysisLog(LogFilePath, cancelCh, lineCh)
	go AnalysisLine(lineCh, cancelCh)
	go CountTask(cancel)

	waitFlag := false
	for !waitFlag {
		select {
		case <-ctx.Done():
			waitFlag = true
		default:
		}
	}

}
