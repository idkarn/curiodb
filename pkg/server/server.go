package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/idkarn/curiodb/pkg/api"
	"github.com/idkarn/curiodb/pkg/common"
	mw "github.com/idkarn/curiodb/pkg/middleware"
)

type DBConfig struct {
	PORT uint32
}

func NewConfig(port uint32) DBConfig {
	if port < 1024 || port > 49151 {
		panic(fmt.Sprintf("Port %d is not allowed", port))
	}
	return DBConfig{port}
}

func loadData(port uint32) {
	ok, data := common.Load()
	if ok {
		log.Println("Data was successsfully loaded")
		common.Config(common.DatabaseStore{
			Tables:         data.Tables,
			TablesMetaData: data.TablesMetaData,
		})
	} else {
		log.Println("Load data failed")
		common.Config(common.DatabaseStore{
			Tables: []common.Table{
				{},
			},
			TablesMetaData: []common.TableMetaData{
				{},
			},
		})
	}
}

func initRouter() {
	api.SetupRouting([]mw.Route{
		mw.NewRouteInfo("GET", "/health", api.HealthHandler),
		mw.NewRouteInfo("POST", "/row/new", api.NewRowHandler),
		mw.NewRouteInfo("POST", "/column/new", api.NewColumnHandler),
		mw.NewRouteInfo("POST", "/row/get", api.GetRowHandler),
		mw.NewRouteInfo("POST", "/row/update", api.UpdateRowHandler),
		mw.NewRouteInfo("POST", "/row/delete", api.DeleteRowHandler),
	})

	mw.SetupMiddlewares([]mw.MiddlewareFn{
		mw.CheckRouteMethod,
	})
}

func serve(port uint32) {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

func Launch(config DBConfig) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		Terminate()
	}()

	loadData(config.PORT)
	initRouter()

	log.Printf("curiodb is running on port %d\n", config.PORT)

	serve(config.PORT)
}

func Terminate() {
	common.Dump()
	log.Println("curiodb is stopped")
	os.Exit(0)
}
