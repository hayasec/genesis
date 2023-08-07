package channels_test

import (
	"context"
	"sync"
	"testing"

	"github.com/life4/genesis/channels"
	"github.com/matryer/is"
)

func TestToSlice(t *testing.T) {
	is := is.New(t)
	f := func(given []int) {
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		actual := channels.ToSlice(c)
		is.Equal(given, actual)
	}
	f([]int{})
	f([]int{1})
	f([]int{1, 2, 3, 1, 2})
}

func TestAny(t *testing.T) {
	is := is.New(t)
	f := func(given []int, expected bool) {
		even := func(t int) bool { return t%2 == 0 }
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		actual := channels.Any(c, even)
		is.Equal(expected, actual)
	}
	f([]int{}, false)
	f([]int{1}, false)
	f([]int{2}, true)
	f([]int{1, 2}, true)
	f([]int{1, 2, 3}, true)
	f([]int{1, 3, 5}, false)
	f([]int{1, 3, 5, 7, 9, 11}, false)
	f([]int{1, 3, 5, 7, 10, 11}, true)
}

func TestAll(t *testing.T) {
	is := is.New(t)
	f := func(given []int, expected bool) {
		even := func(t int) bool { return t%2 == 0 }
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		actual := channels.All(c, even)
		is.Equal(expected, actual)
	}
	f([]int{}, true)
	f([]int{1}, false)
	f([]int{2}, true)
	f([]int{1, 2}, false)
	f([]int{2, 4}, true)
	f([]int{2, 4, 6, 8, 10, 12}, true)
	f([]int{2, 4, 6, 8, 11, 12}, false)
}

func TestClose(t *testing.T) {
	is := is.New(t)
	is.True(!channels.Close[int](nil))
	c := make(chan int)
	is.True(channels.Close(c))
	is.True(!channels.Close(c))
}

func TestEach(t *testing.T) {
	is := is.New(t)
	f := func(given []int) {
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		result := make(chan int, len(given))
		mapper := func(t int) { result <- t }
		channels.Each(c, mapper)
		close(result)
		actual := channels.ToSlice(result)
		is.Equal(given, actual)
	}

	f([]int{})
	f([]int{1})
	f([]int{1, 2, 3})
	f([]int{1, 2, 3, 4, 5, 6, 7})
}

func TestChunkEvery(t *testing.T) {
	is := is.New(t)
	f := func(size int, given []int, expected [][]int) {
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		result := channels.ChunkEvery(c, size)
		actual := make([][]int, 0)
		for el := range result {
			actual = append(actual, el)
		}
		is.Equal(expected, actual)
	}
	f(2, []int{}, [][]int{})
	f(2, []int{1}, [][]int{{1}})
	f(2, []int{1, 2}, [][]int{{1, 2}})
	f(2, []int{1, 2, 3, 4}, [][]int{{1, 2}, {3, 4}})
	f(2, []int{1, 2, 3, 4, 5}, [][]int{{1, 2}, {3, 4}, {5}})
}

func TestCount(t *testing.T) {
	is := is.New(t)
	f := func(element int, given []int, expected int) {
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		actual := channels.Count(c, element)
		is.Equal(expected, actual)
	}
	f(1, []int{}, 0)
	f(1, []int{1}, 1)
	f(1, []int{2}, 0)
	f(1, []int{1, 2, 3, 1, 4}, 2)
}

func TestDrop(t *testing.T) {
	is := is.New(t)
	f := func(count int, given []int, expected []int) {
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		result := channels.Drop(c, count)
		actual := make([]int, 0)
		for el := range result {
			actual = append(actual, el)
		}
		is.Equal(expected, actual)
	}
	f(1, []int{}, []int{})
	f(1, []int{2}, []int{})
	f(1, []int{2, 3}, []int{3})
	f(1, []int{1, 2, 3}, []int{2, 3})
	f(0, []int{1, 2, 3}, []int{1, 2, 3})
	f(3, []int{1, 2, 3, 4, 5, 6}, []int{4, 5, 6})
	f(1, []int{1, 2, 3, 4, 5, 6}, []int{2, 3, 4, 5, 6})
}

func TestFilter(t *testing.T) {
	is := is.New(t)
	f := func(given []int, expected []int) {
		even := func(t int) bool { return t%2 == 0 }
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		result := channels.Filter(c, even)
		actual := channels.ToSlice(result)
		is.Equal(expected, actual)
	}
	f([]int{}, []int{})
	f([]int{1}, []int{})
	f([]int{2}, []int{2})
	f([]int{1, 2, 3, 4}, []int{2, 4})
}

func TestMap(t *testing.T) {
	is := is.New(t)
	f := func(given []int, expected []int) {
		double := func(el int) int { return (el * 2) }
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		result := channels.Map(c, double)

		// convert chan int to chan G
		c2 := make(chan int, 1)
		go func() {
			for el := range result {
				c2 <- el
			}
			close(c2)
		}()

		actual := channels.ToSlice(c2)
		is.Equal(expected, actual)
	}
	f([]int{}, []int{})
	f([]int{1}, []int{2})
	f([]int{1, 2, 3}, []int{2, 4, 6})
}

func TestMax(t *testing.T) {
	is := is.New(t)
	f := func(given []int, expected int, expectedErr error) {
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		actual, actualErr := channels.Max(c)
		is.Equal(expected, actual)
		is.Equal(expectedErr, actualErr)
	}
	f([]int{}, 0, channels.ErrEmpty)
	f([]int{1, 4, 2}, 4, nil)
	f([]int{1, 2, 4}, 4, nil)
	f([]int{4, 2, 1}, 4, nil)
}

func TestMin(t *testing.T) {
	is := is.New(t)
	f := func(given []int, expected int, expectedErr error) {
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		actual, actualErr := channels.Min(c)
		is.Equal(expected, actual)
		is.Equal(expectedErr, actualErr)
	}
	f([]int{}, 0, channels.ErrEmpty)
	f([]int{4, 1, 2}, 1, nil)
	f([]int{1, 2, 4}, 1, nil)
	f([]int{4, 2, 1}, 1, nil)
}

func TestPop(t *testing.T) {
	is := is.New(t)

	// exit on cancelled context
	c := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	v, ok := channels.Pop(ctx, c)
	is.True(!ok)
	is.Equal(v, 0)

	// succesful read
	ctx = context.Background()
	go func() { c <- 12 }()
	v, ok = channels.Pop(ctx, c)
	is.True(ok)
	is.Equal(v, 12)

	// exit on closed channel
	ctx = context.Background()
	close(c)
	v, ok = channels.Pop(ctx, c)
	is.True(!ok)
	is.Equal(v, 0)
}

func TestPush(t *testing.T) {
	// exit on cancelled context
	c := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	channels.Push(ctx, c, 11)

	// push into the channel when there is someone to read
	ctx = context.Background()
	go channels.Push(ctx, c, 12)
	is := is.New(t)
	is.Equal(<-c, 12)
}

func TestReduce(t *testing.T) {
	is := is.New(t)
	f := func(given []int, expected int) {
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		sum := func(el int, acc int) int { return (el) + acc }
		actual := channels.Reduce(c, 0, sum)
		is.Equal(expected, actual)
	}
	f([]int{}, 0)
	f([]int{1}, 1)
	f([]int{1, 2}, 3)
	f([]int{1, 2, 3, 4, 5}, 15)
}

func TestScan(t *testing.T) {
	is := is.New(t)
	f := func(given []int, expected []int) {
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		sum := func(el int, acc int) int { return (el) + acc }
		result := channels.Scan(c, 0, sum)

		// convert chan int to chan G
		c2 := make(chan int, 1)
		go func() {
			for el := range result {
				c2 <- (el)
			}
			close(c2)
		}()

		actual := channels.ToSlice(c2)
		is.Equal(expected, actual)
	}
	f([]int{}, []int{})
	f([]int{1}, []int{1})
	f([]int{1, 2}, []int{1, 3})
	f([]int{1, 2, 3, 4, 5}, []int{1, 3, 6, 10, 15})
}

func TestSum(t *testing.T) {
	is := is.New(t)
	f := func(given []int, expected int) {
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		actual := channels.Sum(c)
		is.Equal(expected, actual)
	}
	f([]int{}, 0)
	f([]int{1}, 1)
	f([]int{1, 2}, 3)
	f([]int{1, 2, 3, 4, 5}, 15)
}

func TestTake(t *testing.T) {
	is := is.New(t)
	f := func(count int, given int, expected []int) {
		ctx, cancel := context.WithCancel(context.Background())
		seq := channels.Repeat(ctx, given)
		seq2 := channels.Take(seq, count)
		actual := channels.ToSlice(seq2)
		cancel()
		is.Equal(expected, actual)
	}
	f(0, 1, []int{})
	f(1, 1, []int{1})
	f(2, 1, []int{1, 1})
}

func TestTee(t *testing.T) {
	is := is.New(t)
	f := func(count int, given []int) {
		c := make(chan int, 1)
		go func() {
			for _, el := range given {
				c <- el
			}
			close(c)
		}()
		chans := channels.Tee(c, count)
		is.Equal(count, len(chans))
		for _, ch := range chans {
			go func(ch chan int) {
				actual := channels.ToSlice(ch)
				is.Equal(given, actual)
			}(ch)
		}
	}
	f(1, []int{})
	f(1, []int{1})
	f(1, []int{1, 2})
	f(1, []int{1, 2, 3})
	f(1, []int{1, 2, 3, 1, 2})

	f(2, []int{})
	f(2, []int{1})
	f(2, []int{1, 2})
	f(2, []int{1, 2, 3})
	f(2, []int{1, 2, 3, 1, 2})

	f(10, []int{1, 2, 3, 1, 2})
}

func TestWithBuffer(t *testing.T) {
	is := is.New(t)
	c1 := make(chan int)
	c2 := channels.WithBuffer(c1, 2)
	is.Equal(cap(c2), 2)

	// two values can be pushed and then it stops
	// because we don't read from c2
	c1 <- 11
	c1 <- 12
	select {
	case c1 <- 13:
		is.Fail()
	default:
	}

	// if we read from c2, c1 unblocks
	v := <-c2
	is.Equal(v, 11)
	c1 <- 13

	close(c1)
	_, more := <-c1
	is.True(!more)
}

// WithContext should echo everything from the input channel
// and close the resulting channel when the input one is closed.
func TestWithContext_NoCancellation(t *testing.T) {
	is := is.New(t)
	c1 := make(chan int)
	c2 := channels.WithContext(c1, context.Background())

	r := make([]int, 0)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range c2 {
			r = append(r, i)
		}
	}()

	c1 <- 11
	c1 <- 12
	c1 <- 13
	close(c1)
	wg.Wait()
	is.Equal(r, []int{11, 12, 13})
}

// WithContext should exit on cancellation.
func TestWithContext_Cancellation(t *testing.T) {
	is := is.New(t)
	c1 := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	c2 := channels.WithContext(c1, ctx)

	r := make([]int, 0)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range c2 {
			r = append(r, i)
			if i == 12 {
				cancel()
			}
		}
	}()

	c1 <- 11
	c1 <- 12
	c1 <- 13
	wg.Wait()
	is.Equal(r, []int{11, 12})
	close(c1)
}

func TestWithContext_CancelWaitingInput(t *testing.T) {
	is := is.New(t)
	c1 := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	c2 := channels.WithContext(c1, ctx)
	cancel()
	_, more := <-c2
	is.True(!more)
}
