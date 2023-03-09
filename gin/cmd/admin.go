package cmd

import (
	"github.com/MQEnergy/gin-framework/bootstrap"
	"github.com/MQEnergy/gin-framework/config"
	"github.com/urfave/cli/v2"
)

var (
	account  string
	password string
)

// AccountCmd 管理者账号创建命令
func AccountCmd() *cli.Command {
	return &cli.Command{
		Name:  "account",
		Usage: "Create a new manager account",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "env",
				Aliases:     []string{"e"},
				Value:       "dev",
				Usage:       "请选择配置文件 [dev | test | prod]",
				Destination: &config.ConfEnv,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "account",
				Aliases:     []string{"a"},
				Value:       "",
				Usage:       "请输入账号名称 如：admin",
				Destination: &account,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"p"},
				Value:       "",
				Usage:       "请输入账号密码 如：admin888",
				Destination: &password,
				Required:    true,
			},
		},
		Action: func(ctx *cli.Context) error {
			bootstrap.BootService(bootstrap.MysqlService)
			return nil
		},
	}
}
