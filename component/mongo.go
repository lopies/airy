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

package component

import (
	"context"
	"fmt"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Mongo struct {
	client *mongo.Client
	BaseComponent
}

func NewMongo() *Mongo {
	m := new(Mongo)
	m.SetType(constants.MongoComponent)
	return m
}

func (m *Mongo) Init(config *config.AiryConfig) {
	param := fmt.Sprintf("mongodb://%s:%s@%s",
		config.Mongo.Username,
		config.Mongo.Password,
		config.Mongo.Path,
	)
	clientOptions := options.Client().ApplyURI(param).
		SetMinPoolSize(config.Mongo.MinPoolSize).
		SetMaxPoolSize(config.Mongo.MaxPoolSize)

	clientOptions.SetWriteConcern(writeconcern.New(writeconcern.WMajority(), writeconcern.J(true), writeconcern.WTimeout(1000)))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Errorf("mongo connect failed", err)
		panic(fmt.Sprintf("mongo connect failed,err=%s", err.Error()))
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger.Errorf("mongo ping failed", err)
		panic(fmt.Sprintf("mongo ping failed,err=%s", err.Error()))
	}
	m.client = client
	logger.Infof("mongo component init success")
}

func (m *Mongo) Shutdown() error {
	return m.client.Disconnect(context.TODO())
}

//NEXT CRUD
