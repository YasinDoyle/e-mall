package main

import (
	"fmt"

	conf "github.com/YasinDoyle/e-mall/config"
	// "github.com/YasinDoyle/e-mall/repository/cache"
	// "github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/es"
	// "github.com/YasinDoyle/e-mall/repository/kafka"
	// "github.com/YasinDoyle/e-mall/repository/rabbitmq"
	"github.com/YasinDoyle/e-mall/routes"
	log "github.com/YasinDoyle/e-mall/utils/log"
	"github.com/YasinDoyle/e-mall/utils/track"

	_ "github.com/apache/skywalking-go"
)

func main() {
	loading()
	r := routes.NewRouter()
	_ = r.Run(conf.Config.System.HttpPort)
	fmt.Println("启动配置成功...")
}

func loading() {
	conf.InitConfig()
	// dao.InitMysql()
	// cache.InitCache()
	// rabbitmq.InitRabbitMQ()
	es.InitES()
	// kafka.InitKafka()
	track.InitTrack()
	log.InitLogger()
	fmt.Println("加载配置完成...")
	go startScript()
}

func startScript() {
	// 这里可以放一些启动时需要执行的脚本，比如数据迁移、定时任务等
}
