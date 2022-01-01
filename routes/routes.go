package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/informeai/cachingo/dto"
	"github.com/informeai/cachingo/repository"
)

type Router struct {
	Router *mux.Router
	Db     *repository.RedisRepository
}

func NewRouter() *Router {
	return &Router{Router: mux.NewRouter(), Db: repository.NewRedisRepository()}
}

func (ro *Router) GetRedis() error {
	var e error
	ro.Router.HandleFunc("/rdb/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]
		value, err := ro.Db.Get(key)
		if err != nil {
			e = err
			w.Write([]byte(fmt.Sprintf("{error: %v}", err.Error())))
			return
		}
		w.Write([]byte(fmt.Sprintf("{status: success, value: %v}", value)))
	}).Methods("GET")
	return e
}

func (ro *Router) SetRedis() error {
	var e error
	ro.Router.HandleFunc("/rdb", func(w http.ResponseWriter, r *http.Request) {
		dto := dto.RedisDto{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			e = err
			w.Write([]byte(fmt.Sprintf("{error: %v}", err.Error())))
			return
		}
		if err = json.Unmarshal(body, &dto); err != nil {
			e = err
			w.Write([]byte(fmt.Sprintf("{error: %v}", err.Error())))
			return
		}
		if err = ro.Db.Set(dto.Key, dto.Value); err != nil {
			e = err
			w.Write([]byte(fmt.Sprintf("{error: %v}", err.Error())))
			return
		}
		w.Write([]byte("{status: success}"))
	}).Methods("POST")
	return e
}

func (ro *Router) Start() error {
	if err := ro.GetRedis(); err != nil {
		return err
	}
	if err := ro.SetRedis(); err != nil {
		return err
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), ro.Router); err != nil {
		return err
	}
	return nil
}
