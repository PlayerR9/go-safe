package buffer

import (
	"context"
	"sync"
	"testing"

	"github.com/PlayerR9/go-safe/common"
)

func TestInit(t *testing.T) {
	const (
		MaxCount int = 100
	)

	ctx, cancel := NewContext[int](context.Background())

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		for i := 0; i < MaxCount; i++ {
			err := common.Run(ctx, Send(i))
			if err != nil {
				t.Errorf("could not send %d: %v", i, err)
				return
			}
		}

		cancel()
	}()

	go func() {
		defer wg.Done()

		i := 0

		for i < MaxCount {
			var x int

			err := common.Run(ctx, Receive(&x))
			if err != nil {
				t.Errorf("could not receive %d: %v", i, err)
				return
			}

			if x != i {
				t.Errorf("expected %d, got %d", i, x)
				return
			}

			i++
		}
	}()

	wg.Wait()

	// t.Fatalf("Test completed")
}

func TestTrimFrom(t *testing.T) {
	const (
		MaxCount int = 100
	)

	ctx, cancel := NewContext[int](context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	go func(max int) {
		defer wg.Done()

		for {
			var x int

			err := common.Run(ctx, Receive(&x))
			if err != nil {
				break
			}

			t.Logf("Received %d", x)
		}
	}(MaxCount)

	for i := 0; i < MaxCount; i++ {
		err := common.Run(ctx, Send(i))
		if err != nil {
			t.Errorf("could not send %d: %v", i, err)
			return
		}
	}

	err := common.Run(ctx, Send(MaxCount))
	if err != nil {
		t.Errorf("could not send %d: %v", MaxCount, err)
		return
	}

	cancel()

	wg.Wait()
}
