package main

import (
	"fmt"
	"strings"
	"strconv"
	"flag"
	"os"
	"path"
	"io/ioutil"
	
	"gopkg.in/ini.v1"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

var (
	addr string
	port string
	seekable string
	
	//paths
	pkeys []string
	pvals []*ini.Key
	//types
	tkeys []string
	tvals []*ini.Key
)

func main() {

	var ifp string	//ini file path
	flag.StringVar(&ifp, "ini", "lanv.ini", "ini file")
	flag.Parse()

	cfg, err := ini.Load(ifp)
	if err != nil {
		fmt.Println("Fail to read file: %v", err)
		os.Exit(1)
	}
	
	addr = cfg.Section("").Key("addr").String()
	port = cfg.Section("").Key("port").String()
	seekable = cfg.Section("").Key("seekable").String()
	
	pkeys = cfg.Section("paths").KeyStrings()
	pvals = cfg.Section("paths").Keys()
	tkeys = cfg.Section("types").KeyStrings()
	tvals = cfg.Section("types").Keys()
	
	lanv_svc()

}

func lanv_svc(){

	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	//app.Logger().SetLevel("debug")

	app.RegisterView(iris.HTML("./static/", ".html").Reload(true))
	app.Favicon("./static/favicon.ico")

	//app.HandleDir("assets", "./static/assets")
	//app.HandleDir("file", DataDir)

	app.Get("/", func(ctx iris.Context) { ctx.View("index.html") })
	app.Get("/list", list)
	app.Get("/file", file);
	//app.Post("/add", add)

	app.Listen(addr+":"+port)
	
}

func list(ctx iris.Context) {

	var result []iris.Map

	pixstr := ctx.URLParamDefault("pix", "-1")
	pnm := ctx.URLParamDefault("pnm", "")
	psn := ctx.URLParamDefault("psn", "")
	
	pix, err := strconv.Atoi(pixstr)
	if err!=nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	
	if pix<0 {
		rni := -1 //Root nood index
		for i, it := range pkeys {
			if "Root" == it {
				rni = i
			}else{
				result = append(result, iris.Map{"pix":i, "pnm":it, "ptt":it, "psn":"", "typ":0})
			}
		}
		
		if rni > -1 {
			result = list_pix(rni, pkeys[rni], "", result)
		}
	}else{
		result = list_pix(pix, pnm, psn, result)
	}
	
	ctx.JSON(result)
	
}

func list_pix(pix int, pnm string, psn string, result []iris.Map) []iris.Map{

	pval := pvals[pix]
	dir := path.Join(pval.MustString(""), psn)
	
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		return result
	}

	for _, fi := range rd {
		if fi.IsDir() {
			result = append(result, iris.Map{"pix":pix, "pnm":pnm, "ptt":fi.Name(), "psn":path.Join(psn, fi.Name()), "typ":0})
		} else {
			ext := strings.ToLower(path.Ext(fi.Name()))
			for i, t := range tkeys {
				it := tvals[i].MustString("")
				if "1"==it && t==ext[1:] {
					result = append(result, iris.Map{"pix":pix, "pnm":pnm, "ptt":fi.Name(), "psn":path.Join(psn, fi.Name()), "typ":1})
					break
				}
			}
		}
	}

	return result
}

func file(ctx iris.Context) {
	pixstr := ctx.URLParamDefault("pix", "-1")
	psn := ctx.URLParamDefault("psn", "")
	
	pix, err := strconv.Atoi(pixstr)
	if err!=nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	pval := pvals[pix]
	fpath := path.Join(pval.MustString(""), psn)
	
	fl, err := os.Open(fpath)
	if err!=nil {
		return
	}
	defer fl.Close()
	
	fi, err := fl.Stat()
	if err!=nil {
		return
	}
	
	rge := ctx.Request().Header.Get("Range")
	var rgs []string
	if rge!="" {
		rgs1 := strings.Split(rge, "bytes=")
		if rgs1[1]!="" {
			rgs = strings.Split(rgs1[1], "-")
		}
	}

	if rgs!=nil {
		fsz := fi.Size()
		flsz := strconv.FormatInt(fsz, 10)
		
		var posb int64
		if rgs[0]!="" {
			bpos, err := strconv.ParseInt(rgs[0], 10, 64)
			if err!=nil {bpos=0}
			posb = bpos
		}else{
			posb = 0
		}
		
		var deflen int64
		deflen = 10485760
		if fsz<deflen {deflen=fsz}
		var buflen int64
		if fsz-posb < deflen {
			buflen = fsz-posb
		}else{
			buflen = deflen
		}
		
		ctx.StatusCode(206)
		ctx.Header("status", "206");
		ctx.Header("Content-Length", strconv.FormatInt(buflen, 10));
		ctx.Header("Content-Type", "application/octet-stream");
		ctx.Header("Accept-Ranges", "bytes");
		ctx.Header("Content-Range", "bytes "+strconv.FormatInt(posb, 10)+"-"+strconv.FormatInt(posb+buflen-1, 10)+"/"+flsz);
		buf := make([]byte, buflen)
		fl.ReadAt(buf, posb)
		ctx.Write(buf);

	}else{
		ctx.SendFile(fpath, fi.Name());
	}
	
}
