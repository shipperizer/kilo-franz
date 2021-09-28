package publisher

import (
	"context"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/monitoring"
	"github.com/shipperizer/kilo-franz/refresh"
)

// Producer is the implementation of a kafka publisher
// it will autorefresh tls secrets automatically
// can be configured with multiple writers, each of those will need a nickname
// and will be known by that when calling the Publish method
type Producer struct {
	// embed AutoRefresh object to simplify mutex setup
	af     map[string]refresh.AutoRefreshXInterface
	topics map[string]string

	// mutex
	mutexes map[string]*sync.RWMutex
}

// Close flushes all the kafka.Writers
func (p *Producer) Close() {
	for _, w := range p.af {
		w.Stop()
	}
}

// ListTopics returns a list of topics with their nicknames
func (p *Producer) ListTopics() map[string]string {
	return p.topics
}

// Publihs is the main method of the class, used to write messages on kafka
func (p *Producer) Publish(topicNickname string, messages ...MessageInterface) error {
	writer, err := p.unwrapWriter(topicNickname)

	if err != nil {
		return err
	}

	w, err := p.af[topicNickname].Object(context.TODO())

	if err != nil {
		return err
	}

	cfg, ok := w.Config().(config.WriterConfigInterface)

	if !ok {
		panic("cannot parse refreshable config into config.WriterConfigInterface")
	}

	kMessages := make([]kafka.Message, 0)

	for _, msg := range messages {
		payload, err := cfg.GetEncoder().Encode(msg.GetValue())

		if err != nil {
			return err
		}

		kMessages = append(kMessages, kafka.Message{Key: msg.GetKey(), Value: payload})
	}

	return writer.WriteMessages(
		context.TODO(),
		kMessages...,
	)

}

// Stats returns kafka.WriterStats
func (p *Producer) Stats(topicNickname string) kafka.WriterStats {
	mutex := p.mutexes[topicNickname]
	mutex.RLock()
	defer mutex.RUnlock()

	if stats, ok := p.af[topicNickname].Stats().(kafka.WriterStats); ok {
		return stats
	}

	panic("producer writer is setup incorrectly, not of type *kafka.Witer")
}

// unwrapWriter is a helper methods to remove the different interfaces and reach the final kafka.Writer
func (p *Producer) unwrapWriter(topicNickname string) (*kafka.Writer, error) {
	w, ok := p.af[topicNickname]

	if !ok {
		return nil, fmt.Errorf("no writer found for the topic nickname %v", topicNickname)
	}

	ww, err := w.Object(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("failed fetching kafka writer: %s", err)

	}

	www, ok := ww.(*core.Writer)

	if !ok {
		return nil, fmt.Errorf("kafka writer is not of type core.Writer %v", ww)
	}

	wwI, err := www.Get(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("failed fetching kafka writer interface: %s", err)

	}

	writer, ok := wwI.(*kafka.Writer)

	if !ok {
		return nil, fmt.Errorf("writer is not of type *kafka.Writer, casting failed")
	}

	return writer, nil
}

// NewProducer creates a new object implementing ProducerInterface
func NewProducer(monitor monitoring.MonitorInterface, refreshables ...config.RefreshableInterface) ProducerInterface {
	p := new(Producer)

	p.topics = make(map[string]string)
	p.af = make(map[string]refresh.AutoRefreshXInterface)
	p.mutexes = make(map[string]*sync.RWMutex)

	for _, ref := range refreshables {
		cfg, ok := ref.Config().(config.WriterConfigInterface)

		if !ok {
			panic("cannot parse refreshable config into config.WriterConfigInterface")
		}

		nick := cfg.GetNickname()

		p.topics[nick] = cfg.GetTopic()

		p.mutexes[nick] = &sync.RWMutex{}
		m := p.mutexes[nick]

		p.af[nick] = refresh.NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(m, cfg), ref, cfg.GetLogger(), monitor)
	}

	return p
}
