package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MQEnergy/gin-framework/load_balance"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type Server struct {
	Metadata map[string]string `json:"metadata"`
}

func parseServer(value []byte) (string, error) {
	var sv Server
	err := json.Unmarshal(value, &sv)
	if err != nil {
		return "", err
	}
	return sv.Metadata["grpcHost"] + ":" + sv.Metadata["acceptorPort"], nil
}

var (
	loadBalance load_balance.LoadBalance
)

func init() {
	loadBalance = load_balance.LoadBalanceFactory(load_balance.LbRandom)
}

type ETCDConfig struct {
	Endpoint    string
	DialTimeout time.Duration
	UserName    string
	PassWord    string
	Prefix      string
}

// NewETCD ETCD连接
func NewETCD(config ETCDConfig) (*clientv3.Client, error) {
	var cli *clientv3.Client
	var err error
	c := clientv3.Config{
		Endpoints:   []string{config.Endpoint},
		DialTimeout: config.DialTimeout,
	}
	if config.UserName != "" && config.PassWord != "" {
		c.Username = config.UserName
		c.Password = config.PassWord
	}
	cli, err = clientv3.New(c)
	if err != nil {
		fmt.Printf("error initializing etcd client: %s", err.Error())
		return nil, err
	}
	initService(cli, config.Prefix)
	return cli, nil
}

func initService(client *clientv3.Client, prefix string) {
	resp, err := client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return
	}
	extractAddrs(resp)
	go watcher(client, prefix)
}

func watcher(client *clientv3.Client, prefix string) {
	rch := client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				setServiceList(string(ev.Kv.Key), ev.Kv.Value)
			case mvccpb.DELETE:
				delServiceList(string(ev.Kv.Key))
			}
		}
	}
}

func extractAddrs(resp *clientv3.GetResponse) {
	if resp == nil || resp.Kvs == nil {
		return
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			setServiceList(string(resp.Kvs[i].Key), resp.Kvs[i].Value)
		}
	}
}

func setServiceList(key string, val []byte) {
	ip, err := parseServer(val)
	if err != nil {
		fmt.Println("解析IP错误，err:", err.Error())
	}
	loadBalance.Add(ip)
	fmt.Println("发现服务：", key, " 地址:", ip)
}

func delServiceList(key string) {
	loadBalance.Del(key)
	log.Println("服务下线:", key)
}

func GetURL() string {
	url, err := loadBalance.Get("")
	if err != nil {
		fmt.Println("not found server")
		return ""
	}
	return url
}
