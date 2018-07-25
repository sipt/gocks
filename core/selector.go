package core

import (
	"errors"
)

type NewSelector func(*ServerGroup) (ISelector, error)

type ISelector interface {
	Start()
	Stop()
	Get(string) (*Server, error)
	Refresh()
	Reset(group *ServerGroup)
}

var (
	selectorMap = make(map[string]NewSelector)

	ErrSelectTypeUndefined = errors.New("SelectType undefined")
	ErrSelectTypeExist     = errors.New("SelectType exist")
)

//选择器配置检查
func CheckSelectorType(selectorType string) bool {
	_, ok := selectorMap[selectorType]
	return ok
}

// 注册选择器
func RegisterSelector(name string, selector NewSelector) error {
	_, ok := selectorMap[name]
	if ok {
		return ErrSelectTypeExist
	}
	selectorMap[name] = selector
	return nil
}
