package singleflight

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Foo struct {
	Name string
}

func TestGroup_Do(t *testing.T) {
	var g Group[Foo]

	var (
		v      Foo
		err    error
		shared bool
	)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			v, err, shared = g.Do("foo", func() (Foo, error) {
				time.Sleep(10 * time.Millisecond)
				return Foo{Name: "test"}, nil
			})
			wg.Done()
		}()
	}

	wg.Wait()
	assert.NoError(t, err)
	assert.Equal(t, "test", v.Name)
	assert.Equal(t, true, shared)
}

func TestGroup_DoChan(t *testing.T) {
	var g Group[Foo]
	blocked := make(chan struct{})
	unblock := make(chan struct{})

	go func() {
		_, _, _ = g.Do("foo", func() (Foo, error) {
			close(blocked)
			<-unblock
			return Foo{Name: "test"}, nil
		})
	}()

	<-blocked
	ch := g.DoChan("foo", func() (Foo, error) {
		return Foo{Name: "test"}, nil
	})
	close(unblock)
	r := <-ch

	assert.NoError(t, r.Err)
	assert.Equal(t, "test", r.Val.Name)
	assert.Equal(t, true, r.Shared)
}
