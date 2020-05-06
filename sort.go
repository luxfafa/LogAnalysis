package main

func (s UrlSortSlice) Len() int           { return len(s) }
func (s UrlSortSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s UrlSortSlice) Less(i, j int) bool { return s[i].Often < s[j].Often }
