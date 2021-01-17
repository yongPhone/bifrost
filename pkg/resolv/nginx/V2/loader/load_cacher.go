package loader

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/context"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/loop_preventer"
	"strings"
	"sync"
)

type LoadCacher interface {
	MainConfig() *context.Config
	GetConfig(configName string) *context.Config
	CheckIncludeConfig(src, dst string) error
	SetConfig(config *context.Config) error
}

type loadCache struct {
	mainConfig    string
	cache         map[string]*context.Config
	loopPreventer loop_preventer.LoopPreventer
	rwLocker      *sync.RWMutex
}

func (l loadCache) MainConfig() *context.Config {
	l.rwLocker.RLock()
	defer l.rwLocker.RUnlock()
	return l.GetConfig(l.mainConfig)
}

func (l loadCache) GetConfig(configName string) *context.Config {
	l.rwLocker.RLock()
	defer l.rwLocker.RUnlock()
	config, ok := l.cache[configName]
	if ok {
		return config
	}
	return nil
}

func (l loadCache) CheckIncludeConfig(src, dst string) error {
	return l.loopPreventer.AddStringElement(src, dst)
}

func (l *loadCache) SetConfig(config *context.Config) error {
	configName := config.BasicContext.Position.ConfigAbsPath()
	if strings.EqualFold(configName, "") {
		return ErrInvalidConfig
	}
	l.rwLocker.Lock()
	defer l.rwLocker.Unlock()

	l.cache[configName] = config
	return nil
}

func NewLoadCacher(configAbsPath string) LoadCacher {
	return &loadCache{
		mainConfig:    configAbsPath,
		cache:         make(map[string]*context.Config),
		loopPreventer: loop_preventer.NewLoopPreverter(configAbsPath),
		rwLocker:      new(sync.RWMutex),
	}
}
