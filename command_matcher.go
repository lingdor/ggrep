package main

import (
	"github.com/lingdor/cyclemap"
	"github.com/lingdor/ggrep/groupbuff"
	"github.com/lingdor/ggrep/util"
	"strings"
)

type matchFunc func(c *command, lineInfo *inputdata, wokerCycleMaps *cyclemap.CycleMap[string, *groupbuff.GroupBuff])

func matchOrderlyOut(c *command, lineInfo *inputdata, wokerCycleMaps *cyclemap.CycleMap[string, *groupbuff.GroupBuff]) {
	line := lineInfo.line
	group := lineInfo.group
	groupItem, existsGroup := wokerCycleMaps.Get(lineInfo.group)
	var matchIndex uint = 0
	var isOutputted = false
	if existsGroup {
		matchIndex = groupItem.MatchIndex()
	}
	seekline := line
	if c.ignoreCase {
		seekline = strings.ToLower(seekline)
	}
	if int(matchIndex) < len(c.regexps) {
		reg := c.regexps[matchIndex]
		if indexes := reg.FindStringIndex(seekline); len(indexes) > 1 {
			if c.colorFormat != nil {
				line = c.colorFormat(line, indexes)
			}
			if !existsGroup {
				groupItem = groupbuff.NewItem()
				wokerCycleMaps.Set(group, groupItem)
			}
			groupItem.Increment()
			isOutputted = true
			//check buff size, and break todo
			groupItem.Write([]byte(line), c.mergeLine)
			if groupItem.MatchIndex() >= uint(len(c.greps)) {
				groupItem.SetFullMatched(true)
				if !c.printGroup {
					wokerCycleMaps.Remove(group)
					c.outChan <- groupItem.String()
				}
			}
			return
		}
	}

	if c.printGroup && !isOutputted {
		if !existsGroup {
			groupItem = groupbuff.NewItem()
			wokerCycleMaps.Set(group, groupItem)
		}
		groupItem.Write([]byte(line), c.mergeLine)
	}
}

func matchDirectOut(c *command, lineInfo *inputdata, wokerCycleMaps *cyclemap.CycleMap[string, *groupbuff.GroupBuff]) {

	c.outChan <- lineInfo.line

}
func matchUnOrderlyOut(c *command, lineInfo *inputdata, wokerCycleMaps *cyclemap.CycleMap[string, *groupbuff.GroupBuff]) {
	line := lineInfo.line
	group := lineInfo.group
	groupItem, existsGroup := wokerCycleMaps.Get(lineInfo.group)
	var matchIndex uint = 0
	var isOutputted = false
	if existsGroup {
		matchIndex = groupItem.MatchIndex()
	}
	seekline := line
	if c.ignoreCase {
		seekline = strings.ToLower(seekline)
	}
	for i := 0; i < len(c.regexps); i++ {
		if !util.IsTrue(matchIndex, i) {
			reg := c.regexps[i]
			if indexes := reg.FindStringIndex(seekline); len(indexes) > 1 {
				if c.colorFormat != nil {
					line = c.colorFormat(line, indexes)
				}
				if !existsGroup {
					groupItem = groupbuff.NewItem()
					wokerCycleMaps.Set(group, groupItem)
				}
				isOutputted = true
				matchIndex = util.SetTrue(matchIndex, int(matchIndex))
				groupItem.SetIndex(matchIndex)
				groupItem.Write([]byte(line), c.mergeLine)
				full := util.FullIntBinary(len(c.regexps))
				if matchIndex == full {
					groupItem.SetFullMatched(true)
					if !c.printGroup {
						wokerCycleMaps.Remove(group)
						c.outChan <- groupItem.String()
					}
				}
				return
			}
		}
	}
	if c.printGroup && !isOutputted {
		if !existsGroup {
			groupItem = groupbuff.NewItem()
			wokerCycleMaps.Set(group, groupItem)
		}
		groupItem.Write([]byte(line), c.mergeLine)
	}
}

//func (c *command) CheckGrepMatch(ctx context.Context, matchIndex uint, line string) []int {
//	if c.orderlyMatch {
//		if int(matchIndex) < len(c.regexps) {
//			reg := c.regexps[matchIndex]
//			return reg.FindStringIndex(line)
//		}
//	} else {
//		for i := 0; i < len(c.regexps); i++ {
//			if !util.IsTrue(matchIndex, i) {
//				reg := c.regexps[i]
//				return reg.FindStringIndex(line)
//			}
//		}
//	}
//	if len(c.regexps) < 1 {
//		return []int{0, 0}
//	}
//	c.errLog.Println("unknow error #202401222101")
//	os.Exit(1)
//	return nil
//
//}
