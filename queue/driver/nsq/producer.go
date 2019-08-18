package nsq

import (
	"encoding/json"
	"errors"
	"github.com/nsqio/go-nsq"
	"strconv"
	"time"
)

// producer 生产者
type producer struct {
	isInit   bool
	nsqdAddr string
	config   *nsq.Config
	producer *nsq.Producer
	Debug    bool
}

var (
	Producer = producer{
		config: nsq.NewConfig(),
	}
)

// Init 初始化
func (p *producer) Init(debug bool) ( error) {
	// 直接连接上 nsqd,支持多个,用半角逗号隔开
	p.nsqdAddr = "127.0.0.1:4150"

	p.isInit = true
	p.Debug = debug

	return nil
}

// Run 启动 producer
func (p *producer) Run() ( error) {
	var err error

	if !p.isInit {
		return errors.New("producer not init")
	}
	if p.producer, err = nsq.NewProducer(p.nsqdAddr, p.config); err != nil {
		return err
	}

	if p.Debug {
		// p.producer.SetLogger(log.GetLogger(), nsq.LogLevelDebug)
	} else {
		// p.producer.SetLogger(log.GetLogger(), nsq.LogLevelWarning)
	}
	return nil
}

// marshalMsg 将消息解析成[]byte,如果出错,第二个参数返回 error
func (p *producer) marshalMsg(msg interface{}) (m []byte, err error) {
	switch t := msg.(type) {
	case []byte:
		m = t
	case float64:
		m = []byte(strconv.FormatFloat(t, 'f', -1, 64))
	case int64:
		m = []byte(strconv.FormatInt(t, 10))
	case string:
		m = []byte(t)
	default:
		m, err = json.Marshal(msg)
	}

	return
}

// Publish 投递消息,如果失败,返回 error
func (p *producer) Publish(topic string, msg interface{}) (err error) {
	if !p.isInit {
		err = errors.New("producer not init")
		return
	}
	var (
		m []byte
	)
	if m, err = p.marshalMsg(msg); err != nil {
		return
	}
	err = p.producer.Publish(topic, m)

	return
}

// MultiPublish 批量发布消息,如果失败,返回 error
func (p *producer) MultiPublish(topic string, msg [][]interface{}) (err error) {
	if !p.isInit {
		err = errors.New("producer not init")
		return
	}
	var (
		m   = make([][]byte, 0)
		tmp []byte
	)
	for _, v := range msg {
		if tmp, err = p.marshalMsg(v); err != nil {
			return
		}
		m = append(m, tmp)
	}
	err = p.producer.MultiPublish(topic, m)

	return
}

func (p *producer) DeferPublish(topic string, msg interface{}, deferSecond int64) (err error) {
	if !p.isInit {
		err = errors.New("producer not init")
		return
	}
	var (
		m []byte
	)
	if m, err = p.marshalMsg(msg); err != nil {
		return
	}
	err = p.producer.DeferredPublish(topic, time.Second*time.Duration(deferSecond), m)
	return
}

// Stop 停止
func (p *producer) Stop() {
	p.producer.Stop()
}
