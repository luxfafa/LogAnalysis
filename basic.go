package main

type UrlSortSlice []UrlForSort

var (
	// LogFilePath = `nginx.access.log`
	buflen      = 50000
	LogFilePath = `E:\project_all\go_test\www.csit18.com.access.log`
)

type CountAndExport struct {
	LogPath         string // 日志绝对路径
	cancel          func()
	lineCh          chan string   // 接收单行处理过的日志
	cancelCh        chan struct{} // 接收日志存取完毕的讯号
	countUrlCh      chan string
	countMap        map[string]int
	urlForSortSlice UrlSortSlice
}

type UrlForSort struct {
	Often int    // Url出现的次数
	Url   string // Nginx日志中筛选出的url
}
