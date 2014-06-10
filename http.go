package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func webIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "ready")
}

func webUpToDate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if gsb.IsUpToDate() {
		fmt.Fprintf(w, "true")
	} else {
		fmt.Fprintf(w, "false")
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
			logger.Info(fmt.Sprintf("req=%s err=nil, res=%s in=%f", url, result, time.Since(start).Seconds()))
		}
		fmt.Fprint(w, result)
	}()

	result, _, err = gsb.MightBeListed(url)
	if err != nil || result == "" {
		return
	}

	result, err = gsb.IsListed(url)
}

func handleHTTP() {
	router := httprouter.New()
	router.GET("/", webIndex)
	router.GET("/safe/*url", webSafe)
	router.GET("/uptodate/", webUpToDate)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", httpHost, httpPort), router))
}
