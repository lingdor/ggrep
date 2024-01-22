package groupbuff

import (
	"bytes"
	"github.com/lingdor/ggrep/util"
	"sync"
)

type GroupBuff struct {
	buf        *bytes.Buffer
	matchsize  int
	matchindex uint
	goIndex    int
	orderly    bool
	*sync.Mutex
}

func (l *GroupBuff) GoIndex() int {
	return l.goIndex
}

func (l *GroupBuff) MatchIndex() uint {
	return l.matchindex
}

func (l *GroupBuff) String() string {
	return l.buf.String()
}

func (l *GroupBuff) Write(line string, matchindex, matchsize int, mergeLines bool) ([]byte, bool) {

	if !mergeLines && l.matchindex > 0 {
		l.buf.WriteString("\n")
	}
	l.buf.WriteString(line)
	if l.orderly {
		l.matchindex++
		if l.matchindex >= uint(matchsize) {
			return l.buf.Bytes(), true
		}
	} else {
		full := util.FullIntBinary(matchsize)
		l.matchindex = util.SetTrue(l.matchindex, matchindex)
		if l.matchindex == full {
			return l.buf.Bytes(), true
		}
	}

	return nil, false
}

func NewItem(goIndex int, orderly bool, isSafe bool) *GroupBuff {

	ret := &GroupBuff{
		buf:        &bytes.Buffer{},
		matchindex: 0,
		goIndex:    goIndex,
		orderly:    orderly,
	}
	if isSafe {
		ret.Mutex = &sync.Mutex{}
	}
	return ret
}
