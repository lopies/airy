package load_balance

import (
	"errors"
	"math/rand"
)

type RandomBalance struct {
	curIndex int
	rss      []string
}

func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

func (r *RandomBalance) next() string {
	if len(r.rss) == 0 {
		return ""
	}
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}

func (r *RandomBalance) Get(key string) (string, error) {
	return r.next(), nil
}

func (r *RandomBalance) Del(key string) error {
	index := -1
	for i := range r.rss {
		if r.rss[i] == key {
			index = i
		}
	}
	if index != -1 {
		r.rss[index], r.rss[len(r.rss)-1] = r.rss[len(r.rss)-1], r.rss[index]
		r.rss = r.rss[:len(r.rss)-1]
	}
	return nil
}
