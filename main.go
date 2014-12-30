package main
/**
 * Simple HTTP-abstraction for reading
 * from multiple public web services.
 *
 * @author mdroog <mpdroog@icloud.com>
 */
import (
	"flag"
	"fmt"
	"net/http"
	"github.com/xsnews/webutils/report"
	"github.com/xsnews/webutils/muxdoc"
	"github.com/xsnews/webutils/httpd"
	"github.com/xsnews/webutils/ratelimit"
	"github.com/xsnews/webutils/middleware"
	"confirmd/openkvk"
	"confirmd/config"
	"confirmd/vies"
	"confirmd/postcode"
)

var (
	mux muxdoc.MuxDoc
)

func doc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, mux.String())
}

func kvk(w http.ResponseWriter, r *http.Request) {
	kvk := r.URL.Query().Get("kvk")
	if len(kvk) < 4 {
		httpd.Error(w, nil, "Number too short")
		return
	}

	k, e := openkvk.Get(kvk)
	if e != nil {
		httpd.Error(w, e, "Failed resolving kvk")
		return
	}
	if e := httpd.FlushJson(w, k); e != nil {
		httpd.Error(w, e, "Failed flushing to browser")
		return
	}
}

func euvies(w http.ResponseWriter, r *http.Request) {
	vat := r.URL.Query().Get("vat")

	if len(vat) < 4 {
		httpd.Error(w, nil, "Number too short")
		return
	}
	k, e := vies.Get(vat)
	if e != nil {
		httpd.Error(w, e, "Failed calling VIES webservice")
		return
	}
	if e := httpd.FlushJson(w, k); e != nil {
		httpd.Error(w, e, "Failed flushing to browser")
		return
	}
}

func nlpostal(w http.ResponseWriter, r *http.Request) {
	postal := r.URL.Query().Get("postal")
	houseno := r.URL.Query().Get("houseno")
	additional := r.URL.Query().Get("additional")

	if len(postal) < 6 {
		httpd.Error(w, nil, "Postal too short")
		return
	}
	if len(houseno) == 0 {
		httpd.Error(w, nil, "Houseno too short")
		return
	}

	entity, e := postcode.Get(postal, houseno, additional)
	if e != nil {
		httpd.Error(w, nil, "Failed calling postal webservice")
		return
	}

	if e := httpd.FlushJson(w, entity); e != nil {
		httpd.Error(w, e, "Failed flushing to browser")
		return
	}
}

func main() {
	var (
		path string
		listen string = "localhost:8008"
		isVerbose bool
	)
	flag.StringVar(&path, "p", "./config.json", "Config path")
	flag.BoolVar(&isVerbose, "v", false, "Show all that happens")
	flag.Parse()

	if e := config.Init(path); e != nil {
		panic(e)
	}
	defer config.Close()

	if e := report.Init("confirmd", "/var/log/confirmd", isVerbose); e != nil {
		panic(e)
	}
	defer report.Close()

	mux.Add("/", doc, "This page")
	mux.Add("/doc", doc, "This page")
	mux.Add("/kvk", kvk, "GET Dutch Chamber of Commerce (CoC) info")
	mux.Add("/vies", euvies, "GET European VAT info")
	mux.Add("/postcode", nlpostal, "GET postal code info by postalcode + houseno + housenoadditional")
	ratelimit.SetRedis(config.Redis)
	http.Handle("/", middleware.Use(mux.Mux))

	report.Msg("Start listening on " + listen)
	if e := http.ListenAndServe(listen, nil); e != nil {
		panic(e)
	}
}

