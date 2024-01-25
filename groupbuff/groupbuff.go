package groupbuff

import (
	"bytes"
)

type GroupBuff struct {
	buf         *bytes.Buffer
	matchindex  uint
	fullMatched bool
}

func (l *GroupBuff) MatchIndex() uint {
	return l.matchindex
}
func (l *GroupBuff) SetIndex(v uint) {
	l.matchindex = v
}
func (l *GroupBuff) Increment() {
	l.matchindex++
}
func (l *GroupBuff) FullMatched() bool {
	return l.fullMatched
}

func (l *GroupBuff) String() string {
	return l.buf.String()
}

func (l *GroupBuff) Write(bs []byte, mergeLine bool) {
	if !mergeLine && l.buf.Len() > 0 {
		l.buf.Write([]byte{'\n'})
	}
	l.buf.Write(bs)
}

func (l *GroupBuff) Len() int {
	return l.buf.Len()
}

func NewItem() *GroupBuff {

	ret := &GroupBuff{
		buf:         &bytes.Buffer{},
		matchindex:  0,
		fullMatched: false,
	}

	return ret
}
