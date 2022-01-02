package keymange

import (
	"context"
	"errors"
	"sync"
)

type Selector interface {
	SecretKey(ctx context.Context, keyName string) (secretKey interface{}, err error)
	NextSecretKey(ctx context.Context) (keyName string, secretKey interface{}, err error)
}

type Key struct {
	Name  string
	Value interface{}
}

type RoundRobin struct {
	keys []Key
	head int
	mu   *sync.Mutex
}

func NewRoundRobin(keys []Key) *RoundRobin {
	return &RoundRobin{
		keys: keys,
		head: 0,
		mu:   &sync.Mutex{},
	}
}

func (r *RoundRobin) SecretKey(ctx context.Context, keyName string) (secretKey interface{}, err error) {
	keyChan := make(chan Key, 1)
	errChan := make(chan error, 1)
	go func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		for _, k := range r.keys {
			if k.Name == keyName {
				keyChan <- k
				return
			}
		}

		errChan <- errors.New("key is not found")
	}()

	select {
	case key := <-keyChan:
		return key.Value, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *RoundRobin) NextSecretKey(ctx context.Context) (keyName string, secretKey interface{}, err error) {
	keyChan := make(chan Key, 1)
	go func() {
		r.mu.Lock()
		defer r.mu.Unlock()

		key := r.keys[r.head]

		r.head = (r.head + 1) % len(r.keys)

		keyChan <- key
	}()

	select {
	case key := <-keyChan:
		return key.Name, key.Value, nil
	case <-ctx.Done():
		return "", nil, ctx.Err()
	}
}
