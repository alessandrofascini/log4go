package context

import "sync"

type LoggerContext struct {
	persistentMap sync.Map
	oneTimeUseMap *sync.Map
}

func NewLoggerContext() *LoggerContext {
	return &LoggerContext{
		persistentMap: sync.Map{},
		oneTimeUseMap: nil,
	}
}

func (l *LoggerContext) SetContext(key string, value any) {
	l.persistentMap.Store(key, value)
}

func (l *LoggerContext) RemoveContext(key string) {
	l.persistentMap.Delete(key)
}

func (l *LoggerContext) ChangeOneContext(key string, value any) {
	if l.oneTimeUseMap == nil {
		// This means that we not already used this
		l.oneTimeUseMap = &sync.Map{}
		l.persistentMap.Range(func(key, value any) bool {
			l.oneTimeUseMap.Store(key, value)
			return true
		})
	}
	l.oneTimeUseMap.Store(key, value)
}

func (l *LoggerContext) Consume() map[string]any {
	syncMap := &l.persistentMap
	if l.oneTimeUseMap != nil {
		syncMap = l.oneTimeUseMap
		// Consumed after copy my otu map
		defer func() {
			l.oneTimeUseMap = nil
		}()
	}
	m := make(map[string]any)
	syncMap.Range(func(key, value any) bool {
		switch k := key.(type) {
		case string:
			m[k] = value
		}
		return true
	})
	return m
}

func (l *LoggerContext) Clear() {
	l.persistentMap = sync.Map{}
	l.oneTimeUseMap = nil
}
