package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/idkarn/curio-db/pkg/api"
	"github.com/idkarn/curio-db/pkg/common"
	"github.com/idkarn/curio-db/pkg/middleware"
)

type DBConfig struct {
	PORT uint32
}

func NewConfig(port uint32) DBConfig {
	if port < 1024 || port > 49151 {
		panic("This port is not allowed")
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
	api.SetupRouting([]api.Route{
		{Method: "GET", Path: "/health", Handler: api.HealthHandler},
		{Method: "POST", Path: "/row/new", Handler: api.NewRowHandler},
		{Method: "POST", Path: "/column/new", Handler: api.NewColumnHandler},
		{Method: "POST", Path: "/row/get", Handler: api.GetRowHandler},
		{Method: "POST", Path: "/row/update", Handler: api.UpdateRowHandler},
		{Method: "POST", Path: "/row/delete", Handler: api.DeleteRowHandler},
		{Method: "POST", Path: "/row/all", Handler: api.GetAllRowsHandler},
	})

	middleware.SetupMiddlewares([]middleware.MiddlewareFn{
		middleware.DecodeRequestPayload,
	})
}

func serve(port uint32) {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

func Launch(config DBConfig) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		Terminate()
	}()

	loadData(config.PORT)
	initRouter()

	log.Printf("curio-db is running on port %d\n", config.PORT)

	serve(config.PORT)
}

func Terminate() {
	common.Dump()
	log.Println("curio-db is stopped")
	os.Exit(0)
}
