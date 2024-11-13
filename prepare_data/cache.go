package prepare_data

import (
	"ai-dijkstra/tokenize_code"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

func hashBytes(bytes []byte) (key string) {
	hasher := sha256.New()
	hasher.Write(bytes)
	return hex.EncodeToString(hasher.Sum(nil))
}

func calcSha(path string) (key string, err error) {
	s, err := os.ReadFile(path)
	if err != nil {
		return
	}
	return hashBytes(s), nil
}

var cache map[string][]string
var cacheChanged bool
var cacheLock = sync.RWMutex{}

// const cacheFile = "/home/kajtek/ai-dijkstra/cache.json"
const cacheFile = "/home/kajtek/ai-dijkstra/cache.gob"

func readCache() error {
	start := time.Now()
	b, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			color.Green("No cache file creating new cache")
			cache = make(map[string][]string)
			return nil
		}
		return err
	}
	d := gob.NewDecoder(bytes.NewReader(b))
	err = d.Decode(&cache)
	if err != nil {
		return err
	}
	// cacheByte, err := os.ReadFile(cacheFile)
	// if err != nil {
	// 	return err
	// }
	// err = json.Unmarshal(cacheByte, &cache)
	// if err != nil {
	// 	return err
	// }
	color.Green("Cache read from %v (took: %v)", cacheFile, time.Since(start).Round(time.Millisecond).String())
	cacheChanged = false
	return nil
}

func writeCache() error {
	if !cacheChanged {
		return nil
	}
	start := time.Now()
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err := e.Encode(cache)
	if err != nil {
		return err
	}
	err = os.WriteFile(cacheFile, b.Bytes(), 0644)
	if err != nil {
		return err
	}
	// err := writeJsonToFile(cache, cacheFile)
	// if err != nil {
	// 	return err
	// }
	color.Green("Cache wrote to %v (took: %v)", cacheFile, time.Since(start).Round(time.Millisecond).String())
	cacheChanged = false
	return err
}

func readFromCache(key string) (val []string, ok bool) {
	cacheLock.Lock()
	val, ok = cache[key]
	cacheLock.Unlock()
	return
}

func writeToCache(key string, val []string) {
	cacheLock.Lock()
	cache[key] = val
	cacheChanged = true
	cacheLock.Unlock()
}

func tokenizeCodeCached(filename string) (tokens []string, err error) {
	key, err := calcSha(filename)
	if err != nil {
		return
	}
	if val, ok := readFromCache(key); ok {
		if len(val) == 0 {
			return val, errors.New("empty token list")
		}
		return val, nil
	}
	tokens, err = tokenize_code.TokenizeCode(filename)
	writeToCache(key, tokens)
	return
}
