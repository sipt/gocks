package selector

import "github.com/sipt/gocks/core"

func init() {
	core.RegisterSelector(HandSelector, func(group *core.ServerGroup) (core.ISelector, error) {
		selector := &handSelector{
			group: group,
		}
		return selector, nil
	})
}

const HandSelector = "select"

type handSelector struct {
	selected interface{}
	group    *core.ServerGroup
}

func (h *handSelector) Get(serverName string) (*core.Server, error) {
	if len(serverName) > 0 {
		for _, v := range h.group.Servers {
			switch v.(type) {
			case *core.Server:
				s := v.(*core.Server)
				if s.Name == serverName {
					h.selected = s
					return s, nil
				}
			case *core.ServerGroup:
				sg := v.(*core.ServerGroup)
				if sg.Name == serverName {
					h.selected = sg
					return sg.Selector.Get("")
				}
			}
		}
	} else if h.selected == nil {
		h.selected = h.group.Servers[0]
	}

	if sg, ok := h.selected.(*core.ServerGroup); ok {
		return sg.Selector.Get("")
	}
	s := h.selected.(*core.Server)
	return s, nil
}

func (h *handSelector) Refresh() {
	h.selected = 0
}

func (h *handSelector) Reset(group *core.ServerGroup) {
	h.selected = 0
	h.group = group
}

func (h *handSelector) Start() {
	h.selected = 0
}

func (h *handSelector) Stop() {
	h.selected = 0
}
