package pubsub

import (
	"sync"
	"time"
)

type (
	subscriber chan interface{}
	topicFunc func(v interface{}) bool
)

type Publisher struct {
	m sync.RWMutex
	buffer int
	timeout time.Duration
	subscribers map[subscriber]topicFunc
}

func NewPublisher(publishTimeOut time.Duration, buffer int) *Publisher {
	return &Publisher{
		buffer: buffer,
		timeout: publishTimeOut,
		subscribers: make(map[subscriber]topicFunc),
	}
}

//订阅
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, p.buffer)
	p.m.Lock()
	p.subscribers[ch] = topic
	p.m.Unlock()
	return  ch
}

//订阅全部主题
func (p *Publisher) Subscribe() chan interface{} {
	return p.SubscribeTopic(nil)
}

//退出订阅
func (p *Publisher) Evict(sub chan interface{}) {
	p.m.Lock()
	defer p.m.Unlock()
	delete(p.subscribers, sub)
	close(sub)
}

//发送主题
func (p *Publisher) sendTopic(sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}
	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}

//发布一个主题
func (p *Publisher) Publish(v interface{}) {
	p.m.RLock()
	defer p.m.RUnlock()
	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

//关闭发布对象，关闭订阅者通道
func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()
	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}
