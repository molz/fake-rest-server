package main

import (
	"flag"
	"github.com/kataras/iris"
	"github.com/patrickmn/go-cache"
	"time"
	"io/ioutil"
	"log"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"fmt"
)

var (
	bind          = flag.String("bind", ":8080", "http bind port")
	keyValueStore = cache.New(24*time.Hour, 1*time.Hour)
	data          = flag.String("data", "", "Filename to load data. Data must be named *.yaml or *.json format, see data.example.yaml")
)

func main() {
	flag.Parse()
	if *data != "" {
		loadData(*data)
	}
	app := iris.New()
	app.OnAnyErrorCode(handler)

	app.Run(iris.Addr(*bind), iris.WithoutVersionChecker)
}

func handler(ctx iris.Context) {
	addCorsHeaders(ctx)
	switch ctx.Method() {
	case iris.MethodGet:
		handlerGet(ctx)
		break
	case iris.MethodPost:
		handlerPostPut(ctx)
		break
	case iris.MethodPut:
		handlerPostPut(ctx)
		break
	case iris.MethodOptions:
		handlerOptions(ctx)
		break
	case iris.MethodDelete:
		handlerDelete(ctx)
		break
	}
}

func handlerGet(ctx iris.Context) {
	if value, ok := keyValueStore.Get(ctx.Path()); ok {
		ctx.StatusCode(iris.StatusOK)
		ctx.ContentType("application/json")
		ctx.Write(value.([]byte))
		return
	}
	ctx.StatusCode(iris.StatusNotFound)
}

func handlerPostPut(ctx iris.Context) {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err == nil {
		defer ctx.Request().Body.Close()
		ctx.StatusCode(iris.StatusOK)
		ctx.ContentType("application/json")
		keyValueStore.Set(ctx.Path(), body, -1)
		ctx.Write(body)
		return
	}
	log.Printf("fail to read body: %s", err.Error())
	ctx.StatusCode(iris.StatusBadRequest)
}

func handlerDelete(ctx iris.Context) {
	if _, ok := keyValueStore.Get(ctx.Path()); ok {
		ctx.StatusCode(iris.StatusOK)
		keyValueStore.Delete(ctx.Path())
		return
	}
	ctx.StatusCode(iris.StatusNotFound)
}

func handlerOptions(ctx iris.Context) {
	ctx.StatusCode(iris.StatusNoContent)
	ctx.Header("Access-Control-Max-Age", "1728000")
	ctx.Header("Content-Type", "text/plain; charset=utf-8")
	ctx.Header("Content-Length", "0")
}

func addCorsHeaders(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
}

func loadData(filename string) {
	if len(filename) < 6 {
		log.Printf("invalid filename")
		return
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("fail to read file %s: %s", filename, err.Error())
		return
	}
	var data *DataFile
	extension := filename[len(filename)-4:]
	switch extension {
	case "json":
		err = json.Unmarshal(content, &data)
		break
	case "yaml":
		err = yaml.Unmarshal(content, &data)
		break
	default:
		err = fmt.Errorf("invalid file extension: %s", extension)
	}
	if err != nil {
		log.Printf("fail to parse %s: %s", filename, err.Error())
		return
	}
	count := 0
	for key, value := range data.Data {
		keyValueStore.Set(key, value, -1)
		count++
	}
	log.Printf("%d record has been loaded", count)
}
