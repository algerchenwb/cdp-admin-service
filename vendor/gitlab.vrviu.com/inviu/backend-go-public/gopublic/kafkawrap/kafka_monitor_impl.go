package kafkawrap

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	metrics "github.com/rcrowley/go-metrics"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/otelwrap/pkg/otelsarama"
	vlog "gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
	"go.opentelemetry.io/otel"
)

// DefaultConsumerConfig consumer默认配置；(DO NOT MUST MODIFY，不要修改；若要使用特殊非默认配置，请重新新建一个cluster.Config变量)
var DefaultConsumerConfig = sarama.NewConfig()

// DefaultProducerConfig producer默认配置；(DO NOT MUST MODIFY，不要修改；若要使用特殊非默认配置，请重新新建一个sarama.Config变量)
var DefaultProducerConfig = sarama.NewConfig()

// _kafkaProducer 全局生成者对象(默认)
var _kafkaProducer sarama.AsyncProducer

// _mapkafkaProducer 全局生成者对象映射表 <string: sarama.AsyncProducer>
var _mapkafkaProducer sync.Map

func init() {
	// 设置防止sarama引用的metrics内存泄漏
	metrics.UseNilMetrics = true

	// 服务端kafka版本（向下兼容，内网测试环境为2.1版本，其他环境为2.3版本）
	DefaultConsumerConfig.Version = sarama.V2_1_0_0
	// 接收消息失败的重试间隔，默认为2s
	DefaultConsumerConfig.Consumer.Retry.Backoff = time.Second * 2
	// 接收consumer处理消费期间发生的错误
	DefaultConsumerConfig.Consumer.Return.Errors = true
	// 提交offset的频率，默认为1s;(意味着如果在consumer异常退出，那么最长1s内消费的消息有可能在consumer重启后再次消费)
	DefaultConsumerConfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	// 消费的offset位置，若是一个全新的consumer(未提交过offset)，默认从最新开始(即只能接收到在开始消费之后发送到kafka的消息)
	DefaultConsumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	// group的offset信息的保留时长
	DefaultConsumerConfig.Consumer.Offsets.Retention = time.Hour * 24
	// 仅维持本身需要的topic信息，防止topic数量增加是，占用大量内存；
	DefaultConsumerConfig.Metadata.Full = false
	// 会话过期时长，默认为20s
	DefaultConsumerConfig.Consumer.Group.Session.Timeout = time.Second * 20

	// 发送超时时间，默认为3s
	DefaultProducerConfig.Producer.Timeout = time.Second * 3
	// 返回发送错误
	DefaultProducerConfig.Producer.Return.Errors = true
	// 不返回发送成功通知
	DefaultProducerConfig.Producer.Return.Successes = false
	// producer在leader已成功收到的数据并得到确认后发送下一条message
	DefaultProducerConfig.Producer.RequiredAcks = sarama.WaitForLocal
	// 设置producer使用的分区器为随机分区器
	DefaultProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
}

func NewDefaultConsumerConfig() *sarama.Config {
	config := *DefaultConsumerConfig
	return &config
}

func NewDefaultProducerConfig() *sarama.Config {
	config := *DefaultProducerConfig
	return &config
}

// kafkaMonitor kafka监听器
type kafkaMonitor struct {
	hosts  []string       // kafka broker节点列表
	topic  []string       // 监听的topic
	group  string         // consumer的group
	config *sarama.Config // consumer配置

	handler     MessageHandlerWithCtx     // 消息处理回调
	liteHandler MessageLiteHandlerWithCtx // 消息处理回调（简洁版本）
	mchan       chan interface{}          // 消息channel

	lock   sync.Mutex // 锁
	status bool       // 状态

	ctx    context.Context    // 控制信道
	cancel context.CancelFunc // 取消方法
}

func (m *kafkaMonitor) Ctx() context.Context {
	if m.ctx == nil {
		m.ctx = otelwrap.NewSkipTraceCtx("kafkaMonitor_ctx_nil")
		vlog.Errorf(m.ctx, "kafkaMonitor.Ctx(). ctx is null")
	}
	return m.ctx
}

// Start 启动监听器
func (m *kafkaMonitor) Start() {
	vlog.Infof(m.Ctx(), "[kafkawrap] start monitor. topic(%s) group(%s)", m.topic, m.group)

	m.lock.Lock()
	defer m.lock.Unlock()

	if m.status {
		return
	}

	m.status = true
	m.ctx, m.cancel = context.WithCancel(m.Ctx())

	if m.config == nil {
		m.config = DefaultConsumerConfig
	}

	go MonitorKafkaTopic(m.ctx, m.hosts, m.topic, m.group, m.config, m)
}

// StartChan 启动监听器
func (m *kafkaMonitor) StartChan() <-chan interface{} {
	vlog.Infof(m.Ctx(), "[kafkawrap] start monitor. topic(%s) group(%s)", m.topic, m.group)

	m.lock.Lock()
	defer m.lock.Unlock()

	if m.status {
		return m.mchan
	}

	m.status = true
	m.ctx, m.cancel = context.WithCancel(m.Ctx())
	m.mchan = make(chan interface{}, len(m.topic))
	m.liteHandler = func(ctx context.Context, data []byte) { m.mchan <- data }

	if m.config == nil {
		m.config = DefaultConsumerConfig
	}

	go MonitorKafkaTopic(m.ctx, m.hosts, m.topic, m.group, m.config, m)

	return m.mchan
}

// StartChan 启动监听器
func (m *kafkaMonitor) StartChanWithCtx() <-chan interface{} {
	vlog.Infof(m.Ctx(), "[kafkawrap] start monitor. topic(%s) group(%s)", m.topic, m.group)

	m.lock.Lock()
	defer m.lock.Unlock()

	if m.status {
		return m.mchan
	}

	m.status = true
	m.ctx, m.cancel = context.WithCancel(m.Ctx())
	m.mchan = make(chan interface{}, len(m.topic))
	m.liteHandler = func(ctx context.Context, data []byte) { m.mchan <- &ValueWithCtx{Ctx: ctx, Val: data} }

	if m.config == nil {
		m.config = DefaultConsumerConfig
	}

	go MonitorKafkaTopic(m.ctx, m.hosts, m.topic, m.group, m.config, m)

	return m.mchan
}

// StartWithTopic 更新topic并启动监听器
// @param topic: topic
// @param group: 使用group[0]更新group
func (m *kafkaMonitor) StartWithTopic(topic string, group ...string) {
	m.StartWithMultiTopic([]string{topic}, group...)
}

// StartWithMultiTopic 更新topic并启动监听器
// @param topic: topic列表
// @param group: 使用group[0]更新group
func (m *kafkaMonitor) StartWithMultiTopic(topic []string, group ...string) {
	vlog.Infof(m.Ctx(), "[kafkawrap] start monitor. topic(%s) group(%s)", m.topic, m.group)

	m.lock.Lock()
	defer m.lock.Unlock()

	if m.status {
		return
	}

	m.status = true
	m.topic = topic
	m.ctx, m.cancel = context.WithCancel(m.Ctx())

	if len(group) > 0 {
		m.group = group[0]
	}

	if m.config == nil {
		m.config = DefaultConsumerConfig
	}

	go MonitorKafkaTopic(m.ctx, m.hosts, m.topic, m.group, m.config, m)
}

// StartChanWithMultiTopic 更新topic并启动监听器
// @param topic: topic列表
// @param group: 使用group[0]更新group
func (m *kafkaMonitor) StartChanWithMultiTopic(topic []string, group ...string) <-chan interface{} {
	vlog.Infof(m.Ctx(), "[kafkawrap] start monitor. topic(%s) group(%s)", m.topic, m.group)

	m.lock.Lock()
	defer m.lock.Unlock()

	if m.status {
		return m.mchan
	}

	m.status = true
	m.topic = topic
	m.ctx, m.cancel = context.WithCancel(m.Ctx())
	m.mchan = make(chan interface{}, len(m.topic))
	m.liteHandler = func(ctx context.Context, data []byte) { m.mchan <- data }

	if len(group) > 0 {
		m.group = group[0]
	}

	if m.config == nil {
		m.config = DefaultConsumerConfig
	}

	go MonitorKafkaTopic(m.ctx, m.hosts, m.topic, m.group, m.config, m)

	return m.mchan
}

func (m *kafkaMonitor) StartChanWithMultiTopicWithCtx(topic []string, group ...string) <-chan interface{} {
	vlog.Infof(m.Ctx(), "[kafkawrap] start monitor. topic(%s) group(%s)", m.topic, m.group)

	m.lock.Lock()
	defer m.lock.Unlock()

	if m.status {
		return m.mchan
	}

	m.status = true
	m.topic = topic
	m.ctx, m.cancel = context.WithCancel(m.Ctx())
	m.mchan = make(chan interface{}, len(m.topic))
	m.liteHandler = func(ctx context.Context, data []byte) { m.mchan <- &ValueWithCtx{Ctx: ctx, Val: data} }

	if len(group) > 0 {
		m.group = group[0]
	}

	if m.config == nil {
		m.config = DefaultConsumerConfig
	}

	go MonitorKafkaTopic(m.ctx, m.hosts, m.topic, m.group, m.config, m)

	return m.mchan
}

// Setup -
func (m *kafkaMonitor) Setup(sess sarama.ConsumerGroupSession) error {
	vlog.Infof(m.Ctx(), "[kafkawrap] session setup. claims(%s) generation_id(%d) cluster_member_id(%s)", gopublic.ToJSON(sess.Claims()), sess.GenerationID(), sess.MemberID())
	return nil
}

// Cleanup -
func (m *kafkaMonitor) Cleanup(sess sarama.ConsumerGroupSession) error {
	vlog.Infof(m.Ctx(), "[kafkawrap] session cleanup. claims(%s) generation_id(%d) cluster_member_id(%s)", gopublic.ToJSON(sess.Claims()), sess.GenerationID(), sess.MemberID())
	return nil
}

// ConsumeClaim -
func (m *kafkaMonitor) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := otel.GetTextMapPropagator().Extract(m.Ctx(), otelsarama.NewConsumerMessageCarrier(msg))

		if m.liteHandler != nil {
			m.liteHandler(ctx, msg.Value)
		}

		if m.handler != nil {
			m.handler(ctx, msg.Key, msg.Value, msg.Topic, msg.Partition, msg.Offset)
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

// Stop 停止监听器
func (m *kafkaMonitor) Stop() {
	vlog.Infof(m.Ctx(), "[kafkawrap] stop monitor. topic(%s) group(%s)", m.topic, m.group)

	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.status {
		return
	}

	m.status = false
	m.cancel()
}

// MonitorKafkaTopic 监听Kafka topic
// @param   hosts: kafka集群broker列表
// @param   topic: 监听的topic
// @param   group: group会用来计算offset，因此同样一个消费者应该一直使用相同的group来保证所有消息都会被消费
// @param handler: 消息处理方法
func MonitorKafkaTopic(ctx context.Context, hosts, topic []string, group string, config *sarama.Config, handler sarama.ConsumerGroupHandler) {
	client, err := sarama.NewClient(hosts, config)
	if err != nil {
		vlog.Errorf(ctx, "[kafkawrap] create kafka client failed. host(%s) errmsg:%s", gopublic.ToJSON(hosts), err.Error())
		return
	}
	defer func() { _ = client.Close() }()

	consumergroup, err := sarama.NewConsumerGroupFromClient(group, client)
	if err != nil {
		vlog.Errorf(ctx, "[kafkawrap] create kafka consumer failed. host(%s) group(%s) errmsg:%s", gopublic.ToJSON(hosts), group, err.Error())
		return
	}
	defer func() { _ = consumergroup.Close() }()

	if EnableConsumerOtel(ctx) {
		vlog.Infof(ctx, "[kafkawrap] kafka consumer enable otel. host(%s) topic(%s) group(%s)", gopublic.ToJSON(hosts), gopublic.ToJSON(topic), group)
		handler = otelsarama.WrapConsumerGroupHandler(handler)
	}

	// 捕获错误
	go func() {
		for err := range consumergroup.Errors() {
			vlog.Errorf(ctx, "[kafkawrap] kafka consumer meet error. host(%s) topic(%s) group(%s) errmsg:%v",
				gopublic.ToJSON(hosts), gopublic.ToJSON(topic), group, err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// TODO check
			cctx, cancel := context.WithCancel(ctx)
			defer cancel()

			err := consumergroup.Consume(cctx, topic, handler)
			if err != nil {
				vlog.Errorf(ctx, "[kafkawrap] join consumer failed. host(%s) topic(%s) group(%s) errmsg:%s",
					gopublic.ToJSON(hosts), gopublic.ToJSON(topic), group, err.Error())
				continue
			}
		}
	}
}

// CreateSingletonKafkaProducer 创建默认kafka producer
// 使用默认配置文件中kafka的相关配置；例：配置文件路径为'./conf/app.conf'
// 配置项：
// [kafka]
// hosts=cg-kafka-ip-1:29092,cg-kafka-ip-2:29092
// -
// @param config: Producer配置，可使用DefaultProducerConfig
var CreateSingletonKafkaProducer FuncCreateSingletonKafkaProducer = func(config *sarama.Config) error {
	return CreateSingletonKafkaProducerWithCtx(otelwrap.NewSkipTraceCtx("CreateSingletonKafkaProducer"), config)
}

var CreateSingletonKafkaProducerWithCtx FuncCreateSingletonKafkaProducerWithCtx = func(ctx context.Context, config *sarama.Config) error {
	if _kafkaProducer != nil {
		return nil
	}

	hosts := strings.Split(beego.AppConfig.String("kafka::hosts"), ",")

	if EnableProducerOtel(ctx) {
		// So we can know the partition and offset of messages.
		config.Producer.Return.Successes = true
	}

	producer, err := sarama.NewAsyncProducer(hosts, config)
	if err != nil {
		vlog.Errorf(ctx, "[kafkawrap] create kafka producer fail, %s", err.Error())
		return err
	}
	vlog.Infof(ctx, "[kafkawrap] create async producer success. name(default)")

	if EnableProducerOtel(ctx) {
		vlog.Infof(ctx, "[kafkawrap] kafka producer enable otel. host(%s)", gopublic.ToJSON(hosts))
		producer = otelsarama.WrapAsyncProducer(config, producer)
	}

	go func() {
		for emsg := range producer.Errors() {
			ctx := otel.GetTextMapPropagator().Extract(ctx, otelsarama.NewProducerMessageCarrier(emsg.Msg))
			vlog.Errorf(ctx, "[kafkawrap] send error:%s", emsg.Error())
		}
	}()

	go func() {
		for msg := range producer.Successes() {
			ctx := otel.GetTextMapPropagator().Extract(ctx, otelsarama.NewProducerMessageCarrier(msg))
			if beego.AppConfig.DefaultBool("kafka::enable_succ_output", false) {
				data, err := msg.Value.Encode()
				vlog.GetPrintfLogger(err)(ctx, "[kafkawrap] send succ, value(%v) err(%v)", string(data), err)
			}
		}
	}()

	// 保存
	_kafkaProducer = producer

	return nil
}

// WritePartitionKafkaJSON 将object序列化为json字符串并写入默认的kafka集群的消息队列
// @param  topic: 写入的topic
// @param object: 消息
// @param partitionKey: 卡夫卡分区key
var WritePartitionKafkaJSON FuncWritePartitionKafkaJSON = func(topic string, object interface{}, partitionKey string) (err error) {
	return WritePartitionKafkaJSONWithCtx(otelwrap.NewSkipTraceCtx("WritePartitionKafkaJSON"), topic, object, partitionKey)
}

var WritePartitionKafkaJSONWithCtx FuncWritePartitionKafkaJSONWithCtx = func(ctx context.Context, topic string, object interface{}, partitionKey string) (err error) {
	if _kafkaProducer == nil {
		return fmt.Errorf("producer not exist")
	}

	data, err := json.Marshal(&object)
	if err != nil {
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}
	if partitionKey != "" {
		msg.Key = sarama.StringEncoder(partitionKey)
	}
	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(msg))

	_kafkaProducer.Input() <- msg

	return
}

// WriteKafkaJSON 将object序列化为json字符串并写入默认的kafka集群的消息队列
// @param  topic: 写入的topic
// @param object: 消息
var WriteKafkaJSON FuncWriteKafkaJSON = func(topic string, object interface{}) (err error) {
	return WriteKafkaJSONWithCtx(otelwrap.NewSkipTraceCtx("WriteKafkaJSON"), topic, object)
}

var WriteKafkaJSONWithCtx FuncWriteKafkaJSONWithCtx = func(ctx context.Context, topic string, object interface{}) (err error) {
	if _kafkaProducer == nil {
		return fmt.Errorf("producer not exist")
	}

	data, err := json.Marshal(&object)
	if err != nil {
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}
	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(msg))

	_kafkaProducer.Input() <- msg

	return
}

// CreateExplicitKafkaProducer 创建kafka producer
// @param   name: producer名称
// @param  hosts: kafka集群broker列表
// @param config: Producer配置，可使用DefaultProducerConfig
var CreateExplicitKafkaProducer FuncCreateExplicitKafkaProducer = func(name string, hosts []string, config *sarama.Config) error {
	return CreateExplicitKafkaProducerWithCtx(otelwrap.NewSkipTraceCtx("CreateExplicitKafkaProducer"), name, hosts, config)
}

var CreateExplicitKafkaProducerWithCtx FuncCreateExplicitKafkaProducerWithCtx = func(ctx context.Context, name string, hosts []string, config *sarama.Config) error {
	if _, ok := _mapkafkaProducer.Load(name); ok {
		return gopublic.ErrAlreadyExist
	}

	if EnableProducerOtel(ctx) {
		// So we can know the partition and offset of messages.
		config.Producer.Return.Successes = true
	}

	producer, err := sarama.NewAsyncProducer(hosts, config)
	if err != nil {
		vlog.Errorf(ctx, "[kafkawrap] create kafka producer fail, %s", err.Error())
		return err
	}
	vlog.Infof(ctx, "[kafkawrap] create async producer success. name(%s)", name)

	if EnableProducerOtel(ctx) {
		vlog.Infof(ctx, "[kafkawrap] kafka producer enable otel. name(%s) host(%s)", name, gopublic.ToJSON(hosts))
		producer = otelsarama.WrapAsyncProducer(config, producer)
	}

	go func() {
		for emsg := range producer.Errors() {
			ctx := otel.GetTextMapPropagator().Extract(ctx, otelsarama.NewProducerMessageCarrier(emsg.Msg))
			vlog.Errorf(ctx, "[kafkawrap] send error:%s", emsg.Error())
		}
	}()

	go func() {
		for msg := range producer.Successes() {
			ctx := otel.GetTextMapPropagator().Extract(ctx, otelsarama.NewProducerMessageCarrier(msg))
			if beego.AppConfig.DefaultBool("kafka::enable_succ_output", false) {
				data, err := msg.Value.Encode()
				vlog.GetPrintfLogger(err)(ctx, "[kafkawrap] send succ, value(%v) err(%v)", string(data), err)
			}
		}
	}()

	// 保存
	_mapkafkaProducer.Store(name, producer)

	return nil
}

// WriteExplicitKafkaJSON 将object序列化为json字符串并写入默认的kafka集群的消息队列
// @param   name: producer名称
// @param  topic: 写入的topic
// @param object: json消息对象
var WriteExplicitKafkaJSON FuncWriteExplicitKafkaJSON = func(name string, topic string, object interface{}) (err error) {
	return WriteExplicitKafkaJSONWithCtx(otelwrap.NewSkipTraceCtx("WriteExplicitKafkaJSON"), name, topic, object)
}

var WriteExplicitKafkaJSONWithCtx FuncWriteExplicitKafkaJSONWithCtx = func(ctx context.Context, name string, topic string, object interface{}) (err error) {
	v, ok := _mapkafkaProducer.Load(name)
	if !ok {
		return gopublic.ErrNotExist
	}

	data, err := json.Marshal(&object)
	if err != nil {
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}
	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(msg))

	v.(sarama.AsyncProducer).Input() <- msg

	return
}

// IsExistProducer 判断kafka producer是否已存在
var IsExistProducer FuncIsExistProducer = func(name string) bool {
	_, ok := _mapkafkaProducer.Load(name)
	return ok
}

func EnableConsumerOtel(ctx context.Context) bool {
	return (beego.AppConfig.DefaultBool("kafka::enable_otel", false) || beego.AppConfig.DefaultBool("kafka::enable_consumer_otel", false)) && !otelwrap.IsSkip(ctx)
}

func EnableProducerOtel(ctx context.Context) bool {
	return (beego.AppConfig.DefaultBool("kafka::enable_otel", false) || beego.AppConfig.DefaultBool("kafka::enable_producer_otel", false)) && !otelwrap.IsSkip(ctx)
}
