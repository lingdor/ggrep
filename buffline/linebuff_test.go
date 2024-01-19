package buffline

import (
	"fmt"
	"testing"
)

func TestClean(t *testing.T) {

	buff := New(10, 2, true)
	for i := 0; i < 10; i++ {
		group := fmt.Sprintf("[%d]", i)
		log1 := fmt.Sprintf("[111] read%d", i)
		if _, done := buff.Write(group, log1, 0); done {
			t.Errorf("[1]wrong buf state in %d", i)
		}
		log2 := fmt.Sprintf("[111] write%d", i)
		if bs, done := buff.Write(group, log2, 1); done {
			str := string(bs)
			expect := fmt.Sprintf("%s\n%s", log1, log2)
			if str != expect {
				t.Errorf("result:\n%s\n expect:\n%s", str, expect)
			}
		} else {
			t.Errorf("[2]wrong buf state in %d", i)
		}
	}
	if len(buff.lines) != 0 {
		t.Errorf("buff length wrong, result:%d\n expect:%d", len(buff.lines), 0)
	}
}

func TestBuff(t *testing.T) {

	buff := New(10, 3, true)
	for i := 0; i < 11; i++ {
		group := fmt.Sprintf("[%d]", i)
		log1 := fmt.Sprintf("[111] read%d", i)
		if _, done := buff.Write(group, log1, 0); done {
			t.Errorf("[1]wrong buf state in %d", i)
		}
		log2 := fmt.Sprintf("[111] write%d", i)
		if _, done := buff.Write(group, log2, 1); done {

			t.Errorf("[2]wrong buf state in %d", i)
		}
	}
	if len(buff.lines) != 10 {
		t.Errorf("buff length wrong, result:%d\n expect:%d", len(buff.lines), 0)
	}

}
