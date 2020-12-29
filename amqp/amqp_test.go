package amqp

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAmqp(t *testing.T) {
	queue := &Queue{Name: "toolkit.queue1.test"}
	// exchange := &Exchange{Name: "toolkit.exchange.test"}

	msg := &Message{
		Data: []byte("{\"seqno\":\"1563541319\",\"cmd\":\"44\",\"data\":{\"mid\":1070869}}"),
	}

	mq, err := New(&Config{
		Addr:         "amqp://guest:guest@mq.aam.test:5672/",
		ExchangeName: "toolkit.exchange1.test",
	})
	if err != nil {
		panic(err)
	}

	testCount := 1000

	startTime := time.Now()

	wg := sync.WaitGroup{}
	si := 0
	for ; si < testCount; si++ {
		err := mq.Pub(queue, msg)
		if err != nil {
			panic(err)
		}
	}
	t.Logf("发送 %d 条数据, 耗时 %d 纳秒 \n", si, time.Since(startTime))

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
	queue := &Queue{Name: "toolkit.queue.test2", Key: "toolkit.queue2.*"}
	mq, err := New(&Config{
		Addr:         "amqp://guest:guest@10.0.2.252:5672/",
		ExchangeName: "toolkit.exchange.test2", // 直连交换机名称
	})
	if err != nil {
		panic(err)
	}

	count := 100

	wg2 := sync.WaitGroup{}
	wg2.Add(count)
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
			wg2.Done()
			// fmt.Printf("msg: %s \n", v)
		}
	}()

	<-time.After(100 * time.Millisecond)

	msg := &Message{
		Data: []byte(`{"seqno":"1563541319","cmd":"44","data":{"uid":1070869}}`),
	}
	ex := &Exchange{Name: "toolkit.ex.test.fanout", Kind: ExchangeFanout, AutoDelete: true}

	for i := 0; i < count; i++ {
		err := mq.Pub(queue, msg, ex)
		if err != nil {
			panic(err)
		}
	}

	wg2.Wait()

}

func TestAmqpApp(t *testing.T) {
	testQueue := &Queue{Name: "toolkit.queue.test3", Key: "toolkit.queue.test3"}
	testReplyQueue := &Queue{Name: "ttoolkit.queue.reply.test3", Key: "toolkit.queue.reply.test3"}
	mq, err := NewApp(&Config{
		Addr:         "amqp://guest:guest@10.0.2.252:5672/",
		ExchangeName: "toolkit.exchange.test3", // 直连交换机名称
	})

	if err != nil {
		t.Errorf("amqp error: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	mq.On(testQueue, func(c *MQContext) {
		t.Log("mq listener here")
		wg.Done()
	})
	mq.Route(map[*Queue]MQHandler{
		testQueue: func(c *MQContext) {
			body := map[string]interface{}{}
			if err := c.BindJSON(&body); err != nil {
				t.Errorf("bind error")
				return
			}
			t.Logf("mq context here, data: %+v", body)
			c.Pub(testReplyQueue, &Message{Data: []byte(`{"hello":"world"}`)})
			wg.Done()
		},
	})

	mq.Pub(testQueue, &Message{Data: []byte(`{"hello":"world"}`)})
	wg.Wait()
}

func ExampleSimple() {
	queue := &Queue{Name: "toolkit.queue.test", Key: "toolkit.queue.*"}
	mq, _ := New(&Config{
		Addr:         "amqp://guest:guest@10.0.2.252:5672/",
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

func TestAmqpR(t *testing.T) {
	queue := &Queue{
		Name:    "debug.desktop.v1.server.rpc.req_mall_card_detail_list",
		Durable: true,
	}
	// exchange := &Exchange{Name: "toolkit.exchange.test"}

	mq, err := New(&Config{
		Addr:         "amqp://guest:guest@10.0.2.252:5672/",
		ExchangeName: "desktop.exchange.v1",
	})
	if err != nil {
		panic(err)
	}

	testCount := 500000

	startTime := time.Now()

	// wg := sync.WaitGroup{}
	si := 0
	for ; si < testCount; si++ {
		msg := &Message{
			Data: []byte(fmt.Sprintf(`{"cmd":"req_mall_card_detail_list","data":{"card_id":1,"mid":9100251},"seqno":"%d"}`, si)),
		}
		err := mq.Pub(queue, msg)
		if err != nil {
			panic(err)
		}
	}
	t.Logf("发送 %d 条数据, 耗时 %d 纳秒 \n", si, time.Since(startTime))

	// startTime1 := time.Now()
	// wg.Add(testCount)
	// go func() {
	// 	msgs, err := mq.Sub(queue)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for range msgs {
	// 		wg.Done()
	// 	}
	// }()

	// wg.Wait()
	// t.Logf("消费 %d 条数据, 耗时 %d 纳秒 \n", testCount, time.Since(startTime1))
}
