package wallet

import (
	"github.com/ToirovSadi/wallet/pkg/types"
	"sync"
)

func (s *Service) SumPaymentsWithProgress() <-chan types.Progress {
	parts := 100_000

	l := 0
	channels := make([]<-chan types.Progress, (len(s.payments)+parts-1)/parts)
	for index, r := range DelN(len(s.payments), parts) { // i know that's wrong :) to use 'range' like this
		ch := make(chan types.Progress)
		go func(ch chan<- types.Progress, data []*types.Payment) {
			defer close(ch)
			result := types.Progress{}
			result.Part = len(data)
			for _, payment := range data {
				result.Result += payment.Amount
			}
			ch <- result
		}(ch, s.payments[l:r])
		channels[index] = ch
		l = r
	}
	return merge(channels)
}
func merge(channels []<-chan types.Progress) <-chan types.Progress {
	merged := make(chan types.Progress)
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	for _, ch := range channels {
		go func(ch <-chan types.Progress) {
			defer wg.Done()
			for val := range ch {
				merged <- val
			}
		}(ch)
	}
	go func() {
		defer close(merged)
		wg.Wait()
	}()
	return merged
}
