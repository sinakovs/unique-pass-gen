package generator

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"unique-pass-gen/pkg/passwordstore"
)

var (
	digits = []rune("0123456789")
	lowerC = []rune("abcdefghijklmnopqrstuvwxyz")
	upperC = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type Options struct {
	length int
	digits bool
	lowerC bool
	upperC bool
}

type Generator struct {
	store passwordstore.PasswordStore
}

type Option func(*Options)

func NewOptions(opts ...Option) Options {
	o := Options{
		length: 0,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func NewGenerator(store passwordstore.PasswordStore) *Generator {
	return &Generator{
		store: store,
	}
}

func (g *Generator) UniquePasswordGenerator(options Options) (string, error) {
	combinedSets, err := combine(options.digits, options.lowerC, options.upperC)
	if err != nil {
		return "", err
	}

	pool := buildPool(combinedSets)

	if options.length > len(pool) {
		return "", errors.New("the selected length exceeds the allowed length")
	}

	for {
		result, used, err := pickFromSets(combinedSets)
		if err != nil {
			return "", err
		}

		result, err = fillRemaining(result, used, pool, options.length)
		if err != nil {
			return "", err
		}

		password := string(result)

		if g.store.Exists(password) {
			continue
		}

		g.store.Add(password)

		return password, nil
	}
}

func randIndex(max int) (int, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, fmt.Errorf("crypto/rand.Int failed: %w", err)
	}

	return int(nBig.Int64()), nil
}

func combine(d, l, u bool) ([][]rune, error) {
	var combined [][]rune

	if d {
		combined = append(combined, digits)
	}

	if l {
		combined = append(combined, lowerC)
	}

	if u {
		combined = append(combined, upperC)
	}

	if len(combined) == 0 {
		return nil, errors.New("no options selected")
	}

	return combined, nil
}

func buildPool(combinedSets [][]rune) map[int]rune {
	pool := make(map[int]rune)

	counter := 0

	for _, set := range combinedSets {
		for _, r := range set {
			pool[counter] = r
			counter++
		}
	}

	return pool
}

func pickFromSets(combinedSets [][]rune) ([]rune, map[rune]bool, error) {
	result := make([]rune, 0, len(combinedSets))
	used := make(map[rune]bool)

	for _, set := range combinedSets {
		i, err := randIndex(len(set))
		if err != nil {
			return nil, nil, err
		}

		r := set[i]

		result = append(result, r)
		used[r] = true
	}

	return result, used, nil
}

func fillRemaining(result []rune, used map[rune]bool, pool map[int]rune, targetLength int) ([]rune, error) {
	for len(result) < targetLength {
		if len(pool) == 0 {
			return nil, errors.New("not enough unique characters")
		}

		available := availableRunes(pool, used)
		if len(available) == 0 {
			return nil, errors.New("ran out of unique characters")
		}

		i, err := randIndex(len(available))
		if err != nil {
			return nil, err
		}

		r := available[i]

		used[r] = true

		result = append(result, r)

		removeRuneFromPool(pool, r)
	}

	return result, nil
}

func availableRunes(pool map[int]rune, used map[rune]bool) []rune {
	available := make([]rune, 0, len(pool))

	for _, r := range pool {
		if !used[r] {
			available = append(available, r)
		}
	}

	return available
}

func removeRuneFromPool(pool map[int]rune, r rune) {
	for k, v := range pool {
		if v == r {
			delete(pool, k)
		}
	}
}

func WithLength(i int) Option {
	return func(o *Options) {
		o.length = i
	}
}

func WithDigits() Option {
	return func(o *Options) {
		o.digits = true
	}
}

func WithLowerC() Option {
	return func(o *Options) {
		o.lowerC = true
	}
}

func WithUpperC() Option {
	return func(o *Options) {
		o.upperC = true
	}
}
