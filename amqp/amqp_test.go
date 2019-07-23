package amqp

import (
	"sync"
	"testing"
	"time"
)

func TestAmqp(t *testing.T) {
	queue := &Queue{Name: "toolkit.queue.test"}
	//exchange := &Exchange{Name: "toolkit.exchange.test"}

	msg := &Message{
		Data: []byte("{\"seqno\":\"1563541319\",\"cmd\":\"44\",\"data\":{\"mid\":1070869}}"),
	}

	mq, err := New(&Config{
		Addr:         "amqp://guest:guest@10.0.3.252:5672/",
		ExchangeName: "toolkit.exchange.test",
	})
	if err != nil {
		panic(err)
	}

	testCount := 1000000

	startTime := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < testCount; i++ {
		err := mq.Pub(queue, msg)
		if err != nil {
			panic(err)
		}
	}
	t.Logf("发送 %d 条数据, 耗时 %d 纳秒 \n", testCount, time.Since(startTime))

	startTime1 := time.Now()
	wg.Add(testCount)
	go func() {
		msgs, err := mq.Sub(queue)
		if err != nil {
			panic(err)
		}
		for range msgs {
			wg.Done()
		}
	}()

	wg.Wait()
	t.Logf("消费 %d 条数据, 耗时 %d 纳秒 \n", testCount, time.Since(startTime1))

}
