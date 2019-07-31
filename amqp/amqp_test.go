package amqp

import (
	"fmt"
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

	testCount := 100000

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

func TestExchangePub(t *testing.T) {
	queue := &Queue{Name: "toolkit.queue.test", Key: "toolkit.queue.*"}
	mq, _ := New(&Config{
		Addr:         "amqp://guest:guest@10.0.3.252:5672/",
		ExchangeName: "toolkit.exchange.test", // 直连交换机名称
	})

	count := 100

	var wg sync.WaitGroup
	wg.Add(count)
	go func() {
		msgs, err := mq.Sub(queue)
		if err != nil {
			panic(err)
		}
		for msg := range msgs {
			var v interface{}
			err := msg.JSON(&v)
			if err != nil {
				panic(err)
			}
			wg.Done()
			fmt.Printf("msg: %s \n", v)
		}
	}()

	<-time.After(100 * time.Millisecond)

	msg := &Message{
		Data: []byte("{\"seqno\":\"1563541319\",\"cmd\":\"44\",\"data\":{\"mid\":1070869}}"),
	}
	ex := &Exchange{Name: "toolkit.ex.test.fanout", Kind: ExchangeFanout, AutoDelete: true}

	for i := 0; i < count; i++ {
		err := mq.Pub(queue, msg, ex)
		if err != nil {
			panic(err)
		}
	}

	wg.Wait()
}

func ExampleSimple() {
	queue := &Queue{Name: "toolkit.queue.test", Key: "toolkit.queue.*"}
	mq, _ := New(&Config{
		Addr:         "amqp://guest:guest@10.0.3.252:5672/",
		ExchangeName: "toolkit.exchange.test", // 直连交换机名称
	})
	go func() {
		msgs, err := mq.Sub(queue)
		if err != nil {
			panic(err)
		}
		for msg := range msgs {
			var v interface{}
			err := msg.JSON(&v)
			if err != nil {
				panic(err)
			}
			fmt.Printf("msg: %s", v)
		}
	}()

}
