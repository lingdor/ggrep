package buffline

import (
	"bytes"
	"sync"
)

type LineBuffItem struct {
	buf     *bytes.Buffer
	matches []bool
}

func (l *LineBuffItem) isFull() bool {
	for _, v := range l.matches {
		if !v {
			return false
		}
	}
	return true
}

type LineBuff struct {
	lines     map[string]*LineBuffItem
	keys      []string
	keysIndex int
	*sync.Mutex
	isWrap     bool
	bufSize    int
	matchCount int
	isSafe     bool
}

func (l *LineBuff) IsMatch(group string, index int) bool {

	if l.isSafe {
		l.Mutex.Lock()
		defer l.Unlock()
	}
	if item, ok := l.lines[group]; ok {
		return item.matches[index]
	}
	return false
}

func (l *LineBuff) Write(group string, line string, index int) ([]byte, bool) {
	if l.matchCount == 1 {
		return []byte(line), true
	}
	if l.isSafe {
		l.Mutex.Lock()
		defer l.Unlock()
	}
	var item *LineBuffItem
	var ok bool
	if item, ok = l.lines[group]; ok {
		if item.matches[index] {
			return nil, false
		}
		if l.isWrap {
			item.buf.WriteString("\n")
		}
		item.buf.WriteString(line)
		item.matches[index] = true
		if item.isFull() {
			delete(l.lines, group)
			return item.buf.Bytes(), true
		}
	} else {
		item := &LineBuffItem{
			buf:     &bytes.Buffer{},
			matches: make([]bool, l.matchCount),
		}
		item.buf.WriteString(line)
		item.matches[index] = true
		l.lines[group] = item
		if len(l.keys) == l.bufSize {
			l.keysIndex++
			if l.keysIndex > l.bufSize-1 {
				l.keysIndex = 0
			}
			delete(l.lines, l.keys[l.keysIndex])
			l.keys[l.keysIndex] = group
		} else {
			l.keys = append(l.keys, group)
			l.keysIndex++
		}
	}
	return nil, false
}

func New(buffSize, matchCount int, isWrap, isSafe bool) *LineBuff {
	return &LineBuff{
		lines:      make(map[string]*LineBuffItem, buffSize),
		keys:       make([]string, 0, buffSize),
		keysIndex:  -1,
		Mutex:      &sync.Mutex{},
		isWrap:     isWrap,
		bufSize:    buffSize,
		matchCount: matchCount,
		isSafe:     isSafe,
	}
}
