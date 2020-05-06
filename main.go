package main

import (
	"context"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	lineCh := make(chan string, buflen)
	cancelCh := make(chan struct{}, 1)
	countUrlCh := make(chan string, buflen)
	countMap := make(map[string]int, buflen/5)
	//fmt.Println(buflen/5)
	UrlForSortSlice := make(UrlSortSlice,  buflen/5)

	countCore := &CountAndExport{LogFilePath, cancel, lineCh, cancelCh, countUrlCh,countMap,UrlForSortSlice}
	countCore.AnalysisStart()

	waitFlag := false
	for !waitFlag {
		select {
		case <-ctx.Done():
			waitFlag = true
		default:
		}
	}

}
