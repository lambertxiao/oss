package objects

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"oss/apiServer/heartbeat"
	"oss/apiServer/locate"
	"oss/lib/es"
	"oss/lib/rs"
	"oss/lib/utils"
)

func post(w http.ResponseWriter, r *http.Request) {
	// 获取对象名
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 获取对象的hash
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 查询hash是否已在系统中
	if locate.Exist(url.PathEscape(hash)) {
		// 即使hash已经存在，因为是post，更新一个版本号
		e = es.AddVersion(name, hash, size)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}

	// 随机选择分片数个数据节点来分片保存数据
	ds := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(ds) != rs.ALL_SHARDS {
		log.Println("cannot find enough dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// 底下会在数据服务创建分片个数个临时对象
	stream, e := rs.NewRSResumablePutStream(ds, name, url.PathEscape(hash), size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 返回一个临时token给客户端，客户端用来断点上传数据
	w.Header().Set("location", "/temp/"+url.PathEscape(stream.ToToken()))
	w.WriteHeader(http.StatusCreated)
}
