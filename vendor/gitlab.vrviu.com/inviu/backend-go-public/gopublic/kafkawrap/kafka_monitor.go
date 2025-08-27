package kafkawrap

import (
	"context"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
)

type KafkaMonitor interface {
	// 启动监听
	Start()

	// 停止监听
	Stop()

	// 启动监听：返回channel用于接收消息
	StartChan() <-chan interface{}

	// 启动监听：指定topic
	StartWithTopic(topic string, group ...string)

	// 启动监听：指定多个topic
	StartWithMultiTopic(topic []string, group ...string)

	// 启动监听：返回channel用于接收消息
	StartChanWithMultiTopic(topic []string, group ...string) <-chan interface{}
	StartChanWithMultiTopicWithCtx(topic []string, group ...string) <-chan interface{}
}

// 消息处理函数签名，仅抛出消息内容
type MessageLiteHandlerWithCtx func(ctx context.Context, value []byte)
type MessageLiteHandler func(value []byte)

// 消息处理函数签名，抛出消息内容及topic信息
type MessageHandlerWithCtx func(ctx context.Context, key, value []byte, topic string, partition int32, offset int64)
type MessageHandler func(key, value []byte, topic string, partition int32, offset int64)

// CreateKafkaMonitor 创建kafka监听器
// 使用默认配置文件中kafka的相关配置；例：配置文件路径为'./conf/app.conf'
// 配置项：
// [kafka]
// hosts=cg-kafka-ip-1:29092,cg-kafka-ip-2:29092
// -
// @param   topic: 监听的topic
// @param   group: 组
// @param  config: consumer配置，可使用DefaultConfig
// @param handler: 消息处理回调，函数签名为MessageLiteHandler或MessageHandler
// ---
// @return: 监听器对象
var CreateKafkaMonitor FuncCreateKafkaMonitor = func(topic, group string, config *sarama.Config, handler interface{}) KafkaMonitor {
	return CreateKafkaMonitorWithCtx(otelwrap.NewSkipTraceCtx("CreateKafkaMonitor"), topic, group, config, handler)
}

var CreateKafkaMonitorWithCtx FuncCreateKafkaMonitorWithCtx = func(ctx context.Context, topic, group string, config *sarama.Config, handler interface{}) KafkaMonitor {
	return CreateExplicitKafkaMonitorWithCtx(ctx, strings.Split(beego.AppConfig.String("kafka::hosts"), ","), topic, group, config, handler)
}

// CreateExplicitKafkaMonitor 创建kafka监听器
// @param   hosts: kafka集群broker列表
// @param   topic: 监听的topic
// @param   group: 组
// @param  config: consumer配置，可使用DefaultConfig
// @param handler: 消息处理回调，函数签名为MessageLiteHandler或MessageHandler
// ---
// @return: 监听器对象
var CreateExplicitKafkaMonitor FuncCreateExplicitKafkaMonitor = func(hosts []string, topic string, group string, config *sarama.Config, handler interface{}) KafkaMonitor {
	return CreateExplicitKafkaMonitorWithCtx(otelwrap.NewSkipTraceCtx("CreateExplicitKafkaMonitor"), hosts, topic, group, config, handler)
}

var CreateExplicitKafkaMonitorWithCtx FuncCreateExplicitKafkaMonitorWithCtx = func(ctx context.Context, hosts []string, topic string, group string, config *sarama.Config, handler interface{}) KafkaMonitor {
	monitor := &kafkaMonitor{
		ctx:    ctx,
		hosts:  hosts,
		topic:  []string{topic},
		group:  group,
		config: config,
	}

	switch handler := handler.(type) {
	case func(value []byte):
		monitor.liteHandler = MessageLiteHandlerWithCtx(func(ctx context.Context, value []byte) { handler(value) })
	case func([]byte, []byte, string, int32, int64):
		monitor.handler = MessageHandlerWithCtx(func(ctx context.Context, key, value []byte, topic string, partition int32, offset int64) {
			handler(key, value, topic, partition, offset)
		})
	case MessageLiteHandler:
		monitor.liteHandler = MessageLiteHandlerWithCtx(func(ctx context.Context, value []byte) { handler(value) })
	case MessageHandler:
		monitor.handler = MessageHandlerWithCtx(func(ctx context.Context, key, value []byte, topic string, partition int32, offset int64) {
			handler(key, value, topic, partition, offset)
		})

	case func(context.Context, []byte):
		monitor.liteHandler = MessageLiteHandlerWithCtx(handler)
	case func(context.Context, []byte, []byte, string, int32, int64):
		monitor.handler = MessageHandlerWithCtx(handler)
	case MessageLiteHandlerWithCtx:
		monitor.liteHandler = handler
	case MessageHandlerWithCtx:
		monitor.handler = handler
	default:
		panic("invalid message handler")
	}

	return monitor
}

type FuncCreateKafkaMonitor func(topic, group string, config *sarama.Config, handler interface{}) KafkaMonitor
type FuncCreateKafkaMonitorWithCtx func(ctx context.Context, topic, group string, config *sarama.Config, handler interface{}) KafkaMonitor
type FuncCreateExplicitKafkaMonitor func(hosts []string, topic string, group string, config *sarama.Config, handler interface{}) KafkaMonitor
type FuncCreateExplicitKafkaMonitorWithCtx func(ctx context.Context, hosts []string, topic string, group string, config *sarama.Config, handler interface{}) KafkaMonitor
type FuncCreateSingletonKafkaProducer func(config *sarama.Config) error
type FuncCreateSingletonKafkaProducerWithCtx func(ctx context.Context, config *sarama.Config) error
type FuncWritePartitionKafkaJSON func(topic string, object interface{}, partitionKey string) (err error)
type FuncWritePartitionKafkaJSONWithCtx func(ctx context.Context, topic string, object interface{}, partitionKey string) (err error)
type FuncWriteKafkaJSON func(topic string, object interface{}) (err error)
type FuncWriteKafkaJSONWithCtx func(ctx context.Context, topic string, object interface{}) (err error)
type FuncCreateExplicitKafkaProducer func(name string, hosts []string, config *sarama.Config) error
type FuncCreateExplicitKafkaProducerWithCtx func(ctx context.Context, name string, hosts []string, config *sarama.Config) error
type FuncWriteExplicitKafkaJSON func(name string, topic string, object interface{}) (err error)
type FuncWriteExplicitKafkaJSONWithCtx func(ctx context.Context, name string, topic string, object interface{}) (err error)
type FuncIsExistProducer func(name string) bool
type FuncIsExistProducerWithCtxWithCtx func(ctx context.Context, name string) bool

type ValueWithCtx struct {
	Ctx context.Context
	Val []byte
}
