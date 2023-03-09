package load_balance

import (
	"errors"
)

type RoundRobinBalance struct {
	curIndex int
	rss      []string
}

func (r *RoundRobinBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

func (r *RoundRobinBalance) next() string {
	if len(r.rss) == 0 {
		return ""
	}
	lens := len(r.rss) //5
	if r.curIndex >= lens {
		r.curIndex = 0
	}
	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}

func (r *RoundRobinBalance) Get(key string) (string, error) {
	return r.next(), nil
}

func (r *RoundRobinBalance) Del(key string) error {
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
