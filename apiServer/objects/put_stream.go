package objects

import (
	"fmt"
	"oss/apiServer/heartbeat"
	"oss/lib/rs"
)

// 在随机到的数据服务上生成分片个hash临时对象
func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	servers := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(servers) != rs.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}

	return rs.NewRSPutStream(servers, hash, size)
}
