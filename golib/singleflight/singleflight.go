package singleflight

import "golang.org/x/sync/singleflight"

type Group[V any] struct {
	singleflight.Group
}

type Result[V any] struct {
	Val    V
	Err    error
	Shared bool
}

func (g *Group[V]) Do(key string, fn func() (V, error)) (V, error, bool) {
	v, err, shared := g.Group.Do(key, func() (interface{}, error) {
		return fn()
	})
	return v.(V), err, shared
}

func (g *Group[V]) DoChan(key string, fn func() (V, error)) <-chan Result[V] {
	ch := make(chan Result[V], 1)
	result := g.Group.DoChan(key, func() (interface{}, error) {
		return fn()
	})

	go func() {
		for r := range result {
			ch <- Result[V]{Val: r.Val.(V), Err: r.Err, Shared: r.Shared}
		}
	}()

	return ch

}
