// Copyright (c) Airy Author. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package config

import (
	"encoding/json"
	"github.com/google/uuid"
	"os"
	"time"
)

func (s *Server) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

type (
	Stop struct {
		SgChan  chan os.Signal
		DieChan chan struct{}
	}

	Server struct {
		ID            string            `json:"id"`
		Type          string            `json:"type"`
		StartAt       int64             `json:"startAt"`
		Metadata      map[string]string `json:"metadata"`
		Authorization bool
	}

	TLS struct {
		Cert string
		Key  string
	}

	Port struct {
		GatePort  int
		GMSPort   int
		LogicPort int
		TCPPort   int
	}

	Log struct {
		Level string
		File  string
	}

	Heartbeat struct {
		Interval time.Duration
	}

	GRPCClient struct {
		DialTimeout    time.Duration
		LazyConnection bool
		RequestTimeout time.Duration
	}

	GRPCServer struct {
		MaxRecvMessage int
	}

	Metrics struct {
		Period            time.Duration
		PrometheusEnabled bool
		StatsdEnabled     bool
	}

	Pipelines struct {
		StructValidationEnabled bool
	}

	InfoRetrieverConfig struct {
		Region string
	}

	Prometheus struct {
		Port             int
		AdditionalLabels map[string]string
		Game             string
		ConstLabels      map[string]string
	}

	Statsd struct {
		Host        string
		Prefix      string
		Rate        float64
		ConstLabels map[string]string
	}

	// EtcdServiceDiscovery Etcd service discovery config
	EtcdServiceDiscovery struct {
		Endpoints               []string
		User                    string
		Pass                    string
		DialTimeout             time.Duration
		Prefix                  string
		HeartbeatTTL            time.Duration
		HeartbeatLog            bool
		SyncServersInterval     time.Duration
		SyncServersParallelism  int
		RevokeTimeout           time.Duration
		GrantLeaseTimeout       time.Duration
		GrantLeaseMaxRetries    int
		GrantLeaseRetryInterval time.Duration
		ShutdownDelay           time.Duration
		ServerTypesBlacklist    []string
	}

	Mysql struct {
		Path                        string
		Config                      string
		Dbname                      string
		Username                    string
		Password                    string
		MaxIdleConns                int
		MaxOpenConns                int
		LogLevel                    string
		Secret                      string // 数据加密密钥
		SlowThreshold               time.Duration
		IgnoreRecordNotFoundError   bool
		Colorful                    bool
		NamingStrategySingularTable bool
	}

	Mongo struct {
		Path        string
		Database    string
		Username    string
		Password    string
		MinPoolSize uint64
		MaxPoolSize uint64
	}

	Redis struct {
		//Connection pool maximum number of connections, indeterminate can be 0 (0 indicates automatic definition), allocated on demand
		MaxIdle int
		//zero no limit
		MaxActive int
		//zero no close
		IdleTimeout    time.Duration
		Path           string
		Password       string
		Index          int
		ConnectTimeout time.Duration
		Wait           bool
	}

	Docker struct {
	}
)

// Global config
type AiryConfig struct {
	*Stop
	*Server
	*TLS
	*Port
	*Log
	*Heartbeat
	*GRPCClient
	*GRPCServer
	*Metrics
	*Pipelines
	*InfoRetrieverConfig
	*Prometheus
	*EtcdServiceDiscovery
	*Mysql
	*Mongo
	*Redis
}

// NewDefaultAiryConfig provides default configuration for Pitaya App
func NewDefaultAiryConfig(typ string) *AiryConfig {
	config := &AiryConfig{
		Stop: &Stop{
			SgChan:  make(chan os.Signal),
			DieChan: make(chan struct{}),
		},

		Server: &Server{
			ID:      uuid.New().String(),
			Type:    typ,
			StartAt: time.Now().Unix(),
		},

		TLS: &TLS{
			Cert: "",
			Key:  "",
		},

		Port: &Port{
			GatePort:  42000,
			GMSPort:   44000,
			LogicPort: 46000,
			TCPPort:   48000,
		},

		Log: &Log{
			Level: "debug",
			File:  "./tplog/test.logger",
		},

		Heartbeat: &Heartbeat{
			Interval: 6 * time.Second,
		},

		GRPCClient: &GRPCClient{
			DialTimeout:    5 * time.Second,
			LazyConnection: false,
			RequestTimeout: 5 * time.Second,
		},

		GRPCServer: &GRPCServer{
			MaxRecvMessage: 2 << 26,
		},

		Metrics: &Metrics{
			Period:            15 * time.Second,
			PrometheusEnabled: false,
			StatsdEnabled:     false,
		},

		Pipelines: &Pipelines{
			StructValidationEnabled: false,
		},

		InfoRetrieverConfig: &InfoRetrieverConfig{
			Region: "",
		},

		Prometheus: &Prometheus{
			Port:             9090,
			AdditionalLabels: map[string]string{},
			ConstLabels:      map[string]string{},
		},

		EtcdServiceDiscovery: &EtcdServiceDiscovery{
			Endpoints:               []string{"localhost:2379"},
			User:                    "",
			Pass:                    "",
			DialTimeout:             5 * time.Second,
			Prefix:                  "airy/",
			HeartbeatTTL:            60 * time.Second,
			HeartbeatLog:            false,
			SyncServersInterval:     120 * time.Second,
			SyncServersParallelism:  10,
			RevokeTimeout:           5 * time.Second,
			GrantLeaseTimeout:       60 * time.Second,
			GrantLeaseMaxRetries:    15,
			GrantLeaseRetryInterval: 5 * time.Second,
			ShutdownDelay:           300 * time.Millisecond,
			ServerTypesBlacklist:    nil,
		},

		Mysql: &Mysql{
			Path:                        "127.0.0.1:3306",
			Config:                      "parseTime=True&loc=Asia%2FShanghai&charset=utf8mb4&collation=utf8mb4_unicode_ci",
			MaxIdleConns:                10,
			MaxOpenConns:                30,
			Dbname:                      "db1",
			Username:                    "root",
			Password:                    "123456",
			LogLevel:                    "info",
			Secret:                      "",
			SlowThreshold:               30 * time.Millisecond,
			IgnoreRecordNotFoundError:   true,
			Colorful:                    true,
			NamingStrategySingularTable: true,
		},

		Mongo: &Mongo{
			Path:        "127.0.0.1:27017",
			MaxPoolSize: 20,
			MinPoolSize: 5,
			Username:    "admin",
			Password:    "123456",
			Database:    "db1",
		},

		Redis: &Redis{
			MaxIdle:        10,
			MaxActive:      0,
			IdleTimeout:    time.Minute,
			Path:           "127.0.0.1:6379",
			Password:       "",
			Index:          1,
			ConnectTimeout: time.Second * 5,
			Wait:           true,
		},
	}
	return config
}
