package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"code.google.com/p/go.net/websocket"

	"github.com/julienschmidt/httprouter"
)

func web404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprint(w, "404 Not Found")
	logger.Info(fmt.Sprintf("url=%s client=%s status=404", r.URL.String(), r.RemoteAddr))
}

func webIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger.Info(fmt.Sprintf("url=/ client=%s status=200", r.RemoteAddr))
	fmt.Fprint(w, "ready")
}

func webUpToDate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger.Info(fmt.Sprintf("url=/ client=%s status=200", r.RemoteAddr))
	if gsb.IsUpToDate() {
		fmt.Fprintf(w, "true")
	} else {
		fmt.Fprintf(w, "false")
	}
}

func webSafeBatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var err error
	var num = -1
	var results []string
	var req []string
	var start = time.Now()
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Critical(fmt.Sprintf("batch=%d err=%s in=%f", num, err.Error(), time.Since(start).Seconds()))
		} else {
			if verbosity >= 2 {
				logger.Info(fmt.Sprintf("batch=%d err=nil, in=%f", num, time.Since(start).Seconds()))
			}
		}
		rsp, _ := json.Marshal(results)
		fmt.Fprint(w, string(rsp))
	}()
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&req); err != nil {
		return
	}
	num = len(req)
	results = make([]string, num)
	for k, url := range req {
		results[k], _, err = gsb.MightBeListed(url)
		if err != nil {
			return
		}

		if results[k] == "" {
			continue
		}

		results[k], err = gsb.IsListed(url)
		if err != nil {
			return
		}
	}
}

func webSafe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var result string
	var err error
	var url = ps.ByName("url")[1:]
	var start = time.Now()

	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Critical(fmt.Sprintf("req=%s err=%s res=null in=%f", url, err.Error(), time.Since(start).Seconds()))
		} else {
			if verbosity >= 2 {
				logger.Info(fmt.Sprintf("req=%s err=nil, res=%s in=%f", url, result, time.Since(start).Seconds()))
			}
		}
		fmt.Fprint(w, result)
	}()

	result, _, err = gsb.MightBeListed(url)
	if err != nil || result == "" {
		return
	}

	result, err = gsb.IsListed(url)
}

func webSocket(ws *websocket.Conn) {
	defer ws.Close()
	logger.Info(ws.Request().RemoteAddr)
	for {
		var url string
		if err := websocket.Message.Receive(ws, &url); err != nil {
			logger.Warn(fmt.Sprintf("req= err=%s res=null in=nil", url))
			return
		}
		start := time.Now()
		result, _, _ := gsb.MightBeListed(url)
		if result != "" {
			result, _ = gsb.IsListed(url)
		}
		if err := websocket.Message.Send(ws, result); err != nil {
			logger.Warn(fmt.Sprintf("req= err=%s res=null in=%f", url, err.Error(), time.Since(start).Seconds()))
			return
		}
		logger.Info(fmt.Sprintf("req=%s err=nil, res=%s in=%f", url, result, time.Since(start).Seconds()))
	}
}

func handleHTTP() {
	router := httprouter.New()
	router.GET("/", webIndex)
	router.GET("/safe/*url", webSafe)
	router.POST("/batch", webSafeBatch)
	router.GET("/uptodate/", webUpToDate)
	router.Handler("GET", "/sock", websocket.Handler(webSocket))
	router.NotFound = web404
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", httpHost, httpPort), router))
}
