package etcdv3

import (
	"go.etcd.io/etcd/etcdserver/etcdserverpb"
	"strings"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

// Member 节点信息
type Member struct {
	*etcdserverpb.Member
	Role   string `json:"role"`
	Status string `json:"status"`
	DbSize int64  `json:"db_size"`
}

const (
	ROLE_LEADER   = "leader"
	ROLE_FOLLOWER = "follower"

	STATUS_HEALTHY   = "healthy"
	STATUS_UNHEALTHY = "unhealthy"
)

// Node 需要使用到的模型
type Node struct {
	IsDir   bool   `json:"is_dir"`
	Version int64  `json:"version,string"`
	Value   string `json:"value"`
	FullDir string `json:"full_dir"`
}


// NewNode 创建节点
func NewNode(dir string, kv *mvccpb.KeyValue) *Node {
	return &Node{
		IsDir:   string(kv.Value) == DEFAULT_DIR_VALUE,
		Value:   strings.TrimPrefix(string(kv.Key), dir),
		FullDir: string(kv.Key),
		Version: kv.Version,
	}
}