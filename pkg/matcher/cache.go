package matcher

import "sync"

type (
	cachedMatcher struct {
		matcher Matcher

		mtx   sync.RWMutex
		cache map[string]bool
	}
)

// WithCache WithCache
func WithCache(m Matcher) Matcher {
	switch m {
	case TRUE(), FALSE():
		return m
	default:
		return &cachedMatcher{
			matcher: m,
			cache:   map[string]bool{},
		}
	}
}

func (m *cachedMatcher) Match(b []byte) bool {
	s := string(b)
	if result, ok := m.fetch(s); ok {
		return result
	}
	result := m.matcher.Match(b)
	m.put(s, result)
	return result
}

func (m *cachedMatcher) MatchString(s string) bool {
	if result, ok := m.fetch(s); ok {
		return result
	}
	result := m.matcher.MatchString(s)
	m.put(s, result)
	return result
}

func (m *cachedMatcher) fetch(key string) (result bool, ok bool) {
	m.mtx.RLock()
	result, ok = m.cache[key]
	m.mtx.RUnlock()
	return
}

func (m *cachedMatcher) put(key string, result bool) {
	m.mtx.Lock()
	m.cache[key] = result
	m.mtx.Unlock()
}