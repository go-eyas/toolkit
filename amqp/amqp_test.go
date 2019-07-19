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
	t.Logf("test msg count: %d", testCount)
	startTime := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < testCount; i++ {
		err := mq.Pub(queue, msg)
		if err != nil {
			panic(err)
		}
		wg.Add(1)
	}
	t.Logf("pub time: %d ns\n", time.Since(startTime))

	startTime1 := time.Now()
	go func() {
		msgs, err := mq.Sub(queue)
		if err != nil {
			panic(err)
		}
		//i := 0
		for range msgs {
			//i++
			//fmt.Printf("receive: %d, msg %s\n", i, string(msg.Data))
			wg.Done()
		}
	}()

	wg.Wait()
	t.Logf("sub time: %d ns\n", time.Since(startTime1))

}
