package main

var (
	LogFilePath = `nginx.access.log`
	countUrlCh   = make(chan string, 50000)
)

type UrlSortSlice []UrlForSort

type UrlForSort struct {
	Often int // Url出现的次数
	Url   string // Nginx日志中筛选出的url
}

func (s UrlSortSlice) Len() int { return len(s) }
func (s UrlSortSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s UrlSortSlice) Less(i, j int) bool { return s[i].Often < s[j].Often }