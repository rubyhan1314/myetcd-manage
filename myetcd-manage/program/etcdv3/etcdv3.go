package etcdv3

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"myetcd-manage/program/config"
	"errors"
	"go.etcd.io/etcd/pkg/transport"
	"sync"
	"fmt"
	"strings"
	"strconv"
)

// Etcd3Client etcd v3客户端
type Etcd3Client struct {
	*clientv3.Client
}

var (
	// EtcdClis etcd连接对象
	etcdClis *sync.Map
)

func init() {
	etcdClis = new(sync.Map)

}

// NewEtcdCli 创建一个etcd客户端
func NewEtcdCli(etcdCfg *config.EtcdServer) (*Etcd3Client, error) {
	if etcdCfg == nil {
		return nil, errors.New("etcdCfg is nil")
	}
	if etcdCfg.TLSEnable == true && etcdCfg.TLSConfig == nil {
		return nil, errors.New("TLSConfig is nil")
	}
	if len(etcdCfg.Address) == 0 {
		return nil, errors.New("Etcd connection address cannot be empty")
	}

	var cli *clientv3.Client
	var err error

	if etcdCfg.TLSEnable == true {
		// tls 配置
		tlsInfo := transport.TLSInfo{
			CertFile:      etcdCfg.TLSConfig.CertFile,
			KeyFile:       etcdCfg.TLSConfig.KeyFile,
			TrustedCAFile: etcdCfg.TLSConfig.CAFile,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}

		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   etcdCfg.Address,
			DialTimeout: 10 * time.Second,
			TLS:         tlsConfig,
			Username:    etcdCfg.Username,
			Password:    etcdCfg.Password,
		})
	} else {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   etcdCfg.Address,
			DialTimeout: 10 * time.Second,
			Username:    etcdCfg.Username,
			Password:    etcdCfg.Password,
		})
	}

	if err != nil {
		return nil, err
	}
	etcdClis.Store(etcdCfg.Name, cli)
	return &Etcd3Client{
		Client: cli,
	}, nil
}



// GetEtcdCli 获取一个etcd cli对象
func GetEtcdCli(etcdCfg *config.EtcdServer) (*Etcd3Client, error) {
	if etcdCfg == nil {
		return nil, errors.New("etcdCfg is nil")
	}
	val, ok := etcdClis.Load(etcdCfg.Name)
	fmt.Println("---->etcdCfg.Name",etcdCfg.Name)
	if ok == false {
		if len(etcdCfg.Address) > 0 {
			cli, err := NewEtcdCli(etcdCfg)
			if err != nil {
				return nil, err
			}
			return cli, nil
		}
		return nil, errors.New("Getting etcd client error")
	}
	return &Etcd3Client{
		Client: val.(*clientv3.Client),
	}, nil
}


const (
	DEFAULT_DIR_VALUE = "etcdv3_dir_$2H#%gRe3*t"
)



// node 列表格式化成json
func NodeJsonFormat(prefix string, list []*Node) (interface{}, error) {
	resp := make(map[string]interface{}, 0)
	if len(list) == 0 {
		return resp, nil
	}
	for _, v := range list {
		key := strings.TrimPrefix(v.FullDir, prefix)
		key = strings.TrimLeft(key, "/")
		strs := strings.Split(key, "/")
		// if len(strs) > 0 {
		// 	for _, val := range strs {
		// 		log.Println(val)
		// 		js, _ := json.Marshal(resp)
		// 		log.Println(string(js))
		// 		if _, ok := resp[val]; ok == false {
		// 			if v.Value == DEFAULT_DIR_VALUE {
		// 				resp[val] = make(map[string]interface{}, 0)
		// 			} else {
		// 				resp[val] = formatValue(v.Value) // 这里应该做个类型预判
		// 			}
		// 		}
		// 		_, ok := resp[val].(map[string]interface{})
		// 		if ok == false {
		// 			break
		// 		}
		// 	}
		// }
		// log.Println("---------------------")
		// log.Println(v.FullDir)
		// log.Println(strs)
		// log.Println(v.Value)

		recursiveJsonMap(strs, v, resp)
		// jjj, _ := json.Marshal(resp)
		// log.Println(string(jjj))

	}
	// jjj, _ := json.Marshal(resp)
	// log.Println(string(jjj))
	return resp, nil
}


// 递归的将一个值赋值到map中
func recursiveJsonMap(strs []string, node *Node, parent map[string]interface{}) interface{} {
	if len(strs) == 0 || strs[0] == "" || node == nil || parent == nil {
		return nil
	}
	if _, ok := parent[strs[0]]; ok == false {
		if node.Value == DEFAULT_DIR_VALUE {
			parent[strs[0]] = make(map[string]interface{}, 0)
		} else {
			parent[strs[0]] = formatValue(node.Value)
		}
	}
	val, ok := parent[strs[0]].(map[string]interface{})
	if ok == false {
		return val
	}
	return recursiveJsonMap(strs[1:], node, val)
}



// Format 时获取值，转为指定类型
func formatValue(v string) interface{} {
	if v == "true" {
		return true
	} else if v == "false" {
		return false
	}
	// 尝试转浮点数
	vf, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return vf
	}
	// 尝试转整数
	vi, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		return vi
	}
	return v
}

