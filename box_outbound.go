package box

import (
	"fmt"
	"strings"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/common/taskmonitor"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	F "github.com/sagernet/sing/common/format"
)

func (s *Box) startProviderOutbounds() error {
	monitor := taskmonitor.New(s.logger, C.StartTimeout)
	outboundTag := make(map[string]int)
	for _, out := range s.outbounds {
		tag := out.Tag()
		outboundTag[tag] = 0
	}
	for i, p := range s.providers {
		var pTag string
		if p.Tag() == "" {
			pTag = F.ToString(i)
		} else {
			pTag = p.Tag()
		}
		for j, out := range p.Outbounds() {
			var tag string
			if out.Tag() == "" {
				out.SetTag(fmt.Sprint("[", pTag, "]", F.ToString(j)))
			}
			tag = out.Tag()
			if _, exists := outboundTag[tag]; exists {
				count := outboundTag[tag] + 1
				tag = fmt.Sprint(tag, "[", count, "]")
				out.SetTag(tag)
				outboundTag[tag] = count
			}
			outboundTag[tag] = 0
			if starter, isStarter := out.(common.Starter); isStarter {
				monitor.Start("initialize outbound provider[", pTag, "]", " outbound/", out.Type(), "[", tag, "]")
				err := starter.Start()
				monitor.Finish()
				if err != nil {
					return E.Cause(err, "initialize outbound provider[", pTag, "]", " outbound/", out.Type(), "[", tag, "]")
				}
			}
		}
		p.UpdateOutboundByTag()
	}
	return nil
}

func (s *Box) startOutbounds() error {
	monitor := taskmonitor.New(s.logger, C.StartTimeout)
	outboundTags := make(map[adapter.Outbound]string)
	outbounds := make(map[string]adapter.Outbound)
	for i, outboundToStart := range s.outbounds {
		var outboundTag string
		if outboundToStart.Tag() == "" {
			outboundToStart.SetTag(F.ToString(i))
		}
		outboundTag = outboundToStart.Tag()
		if _, exists := outbounds[outboundTag]; exists {
			return E.New("outbound tag ", outboundTag, " duplicated")
		}
		outboundTags[outboundToStart] = outboundTag
		outbounds[outboundTag] = outboundToStart
	}
	err := s.startProviderOutbounds()
	if err != nil {
		return nil
	}
	started := make(map[string]bool)
	for {
		canContinue := false
	startOne:
		for _, outboundToStart := range s.outbounds {
			outboundTag := outboundTags[outboundToStart]
			if started[outboundTag] {
				continue
			}
			dependencies := outboundToStart.Dependencies()
			for _, dependency := range dependencies {
				if !started[dependency] {
					continue startOne
				}
			}
			started[outboundTag] = true
			canContinue = true
			if starter, isStarter := outboundToStart.(common.Starter); isStarter {
				monitor.Start("initialize outbound/", outboundToStart.Type(), "[", outboundTag, "]")
				err := starter.Start()
				monitor.Finish()
				if err != nil {
					return E.Cause(err, "initialize outbound/", outboundToStart.Type(), "[", outboundTag, "]")
				}
			}
		}
		if len(started) == len(s.outbounds) {
			break
		}
		if canContinue {
			continue
		}
		currentOutbound := common.Find(s.outbounds, func(it adapter.Outbound) bool {
			return !started[outboundTags[it]]
		})
		var lintOutbound func(oTree []string, oCurrent adapter.Outbound) error
		lintOutbound = func(oTree []string, oCurrent adapter.Outbound) error {
			problemOutboundTag := common.Find(oCurrent.Dependencies(), func(it string) bool {
				return !started[it]
			})
			if common.Contains(oTree, problemOutboundTag) {
				return E.New("circular outbound dependency: ", strings.Join(oTree, " -> "), " -> ", problemOutboundTag)
			}
			problemOutbound := outbounds[problemOutboundTag]
			if problemOutbound == nil {
				return E.New("dependency[", problemOutboundTag, "] not found for outbound[", outboundTags[oCurrent], "]")
			}
			return lintOutbound(append(oTree, problemOutboundTag), problemOutbound)
		}
		return lintOutbound([]string{outboundTags[currentOutbound]}, currentOutbound)
	}
	return nil
}
