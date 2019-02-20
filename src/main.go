package main

import (
	//	"dao"
	"facade/handler"
	"filecache"
	"fmt"
	"redisclusteradapter"
	"sysconf"
	"time"
	"web/restadapter"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gitstliu/log4go"
)

func main() {

	defer panicHandler()

	log4go.LoadConfiguration("sysconf/log.xml")
	defer log4go.Close()

	sysconf.LoadConfigure("sysconf/config.toml")

	log4go.Info("Config Load Success!!!")

	//	initFileSysErr := datafile.GetFileSys().InitFileSys()

	//	if initFileSysErr != nil {
	//		log4go.Error(initFileSysErr)
	//	}
	//	dao.InitDB(sysconf.GetConfigure().DBType, sysconf.GetConfigure().ConnectionString, sysconf.GetConfigure().LogMode)

	restAdapter := restadapter.RestAdapter{initUrl(), sysconf.GetConfigure().ServicePort}

	redisclusteradapter.CreadRedisCluster(sysconf.GetConfigure().RedisClusterAddresses, sysconf.GetConfigure().RedisClusterConnTimeout, sysconf.GetConfigure().RedisClusterReadTimeout, sysconf.GetConfigure().RedisClusterWriteTimeout, sysconf.GetConfigure().RedisClusterKeepAlive, sysconf.GetConfigure().RedisClusterAliveTime)
	filecache.InitFileCache(sysconf.GetConfigure().LocalCacheMemSize)
	//	idadapter.InitIdWorker(sysconf.GetConfigure().ServiceId, sysconf.GetConfigure().DatacenterId)
	//	initUrlFilter()

	go flushLocalCache()

	restAdapter.Start()

}

func flushLocalCache() {
	for true {
		filecache.CurrFileCache.Flush(sysconf.GetConfigure().LocalCacheFlushThreshold)
		time.Sleep(time.Duration(sysconf.GetConfigure().LocalCacheTimeOut * int64(time.Second)))
	}
}

func initUrl() []*restadapter.UrlMap {

	//	usersfacade := handler.UsersFacade{}

	filefacade := handler.FileFacade{}

	urls := make([]*restadapter.UrlMap, 0, 100)

	urls = append(urls, &restadapter.UrlMap{
		//		Url: "/files/:group_name",
		Url: "/files",
		MethodMap: map[string]rest.HandlerFunc{
			"POST": filefacade.SaveFile}})

	urls = append(urls, &restadapter.UrlMap{
		Url: "/files/:file_name",
		MethodMap: map[string]rest.HandlerFunc{
			"GET": filefacade.ReadFile}})

	urls = append(urls, &restadapter.UrlMap{
		Url: "/files",
		MethodMap: map[string]rest.HandlerFunc{
			"PUT":    filefacade.ReplaceFile,
			"DELETE": filefacade.DeleteFile}})

	//	urls = append(urls, &restadapter.UrlMap{
	//		Url: "/files/:group_name/:file_name",
	//		MethodMap: map[string]rest.HandlerFunc{
	//			"DELETE": filefacade.DeleteFile}})

	urls = append(urls, &restadapter.UrlMap{
		Url: "/cache/:group_name/:file_name",
		MethodMap: map[string]rest.HandlerFunc{
			"DELETE": filefacade.CleanCache}})

	return urls
}

//func initUrlFilter() {
//	http.HttpBeforeFilters["/users"] = []http.FilterHandler{filter.CheckToken}
//	//http.HttpAfterFilters[""] = []http.FilterHandler{filter.CheckToken}

//	printFiltersInfo()
//}

//func printFiltersInfo() {
//	log4go.Debug("#####HttpBeforeFilters")

//	for key, _ := range http.HttpBeforeFilters {
//		log4go.Debug(key)
//	}

//	log4go.Debug("#####HttpAfterFilters")

//	for key, _ := range http.HttpAfterFilters {
//		log4go.Debug(key)
//	}
//	log4go.Debug("#####FilterInfoFinished")
//}

func panicHandler() {
	if r := recover(); r != nil {
		fmt.Println(r)
		fmt.Printf("%T", r)
		panic(r)
	}
}
