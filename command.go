package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/lingdor/cyclemap"
	"github.com/lingdor/ggrep/groupbuff"
	"github.com/lingdor/ggrep/util"
	"go.uber.org/automaxprocs/maxprocs"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"
)

func init() {

}

// ggrep --group 'logid[\w+]' 'LOG1' 'LOG2' --buff-lines 200 --merge-lines
type command struct {
	group          string
	cacheGroupSize int
	mergeLine      bool
	verbose        bool
	parallelCnt    int
	workerChans    []chan *inputdata
	groupWorkers   *cyclemap.CycleMap[string, int]
	errLog         *log.Logger
	regexps        []*regexp.Regexp
	groupRegex     *regexp.Regexp
	greps          []string
	outChan        chan string
	cutLineSize    int
	orderlyMatch   bool
	fullMatch      bool
	recursive      bool
	lineIndex      int
	color          string
}

const DefaultCutLineSize = 102400

func (c *command) verboseLogf(content string, args ...any) {
	if c.verbose {
		log.Printf("[verbose] "+content, args...)
	}
}

func (c *command) ready() {

	flag.StringVar(&c.group, "group", "", "group match expression for regex")
	flag.IntVar(&c.cacheGroupSize, "cache-group-size", 100, "the size of cache group number")
	flag.BoolVar(&c.mergeLine, "merge-lines", false, "merge output of matched lines")
	flag.BoolVar(&c.mergeLine, "M", false, "merge output of matched lines(shorthand)")
	flag.BoolVar(&c.verbose, "verbose", false, "show debug detail infos")
	flag.IntVar(&c.parallelCnt, "parallel-count", 1, "parallel count of tasks to running,when 0 will be full cpu number")
	flag.IntVar(&c.cutLineSize, "cut-line-size", DefaultCutLineSize, fmt.Sprintf("line max size of a line, default:%d", DefaultCutLineSize))
	c.greps = make([]string, 0, 2)
	flag.Func("grep", "match expression for scanning", func(v string) error {
		c.greps = append(c.greps, v)
		return nil
	})
	flag.BoolVar(&c.orderlyMatch, "match-orderly", false, "if true, result must for grep declared sequence")
	flag.BoolVar(&c.fullMatch, "full-match", false, "Only output checked all grep expression in group, will be output")
	//support auto todo
	flag.StringVar(&c.color, "color", "", "will highlight the matched word when you set always")
	//todo
	//flag.BoolVar(&c.recursive, "recursive", false, "Recursively search subdirectories listed")

}

func (c *command) valid() error {

	l := len(c.greps)
	if !c.orderlyMatch {
		grepMax := int(unsafe.Sizeof(l)) * 8
		if l > grepMax {
			return fmt.Errorf("grep command parameters max of %d", grepMax)
		}
	}
	if c.verbose {
		c.errLog = log.New(os.Stderr, "[err]", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	} else {
		c.errLog = log.New(os.Stderr, "[err]", log.Lmsgprefix)
	}
	if len(c.greps) < 1 {
		return fmt.Errorf("no found match expression in parameters")
	}
	if c.group == "" {
		return errors.New("no found parameter: group")
	}
	c.regexps = make([]*regexp.Regexp, len(c.greps))
	var err error
	if c.groupRegex, err = c.ggrepRegex(c.group); err == nil {
		for i, grep := range c.greps {
			if c.regexps[i], err = c.ggrepRegex(grep); err != nil {
				return fmt.Errorf("grep (%s) expression mistake:%w", grep, err)
			}
		}
	} else {
		return fmt.Errorf("group expression mistake:%w", err)
	}
	return nil
}
func (c *command) ggrepRegex(exp string) (*regexp.Regexp, error) {
	exp = strings.ReplaceAll(exp, "[:logid:]", "[0-9a-zA-z_\\-]+")
	return regexp.CompilePOSIX(exp)
}

func (c *command) run() {
	if err := c.valid(); err != nil {
		c.errLog.Printf("valid faild,%s", err.Error())
		os.Exit(1)
	}
	if c.parallelCnt < 1 {
		maxprocs.Set(maxprocs.Logger(log.Printf))
		c.parallelCnt = runtime.NumCPU()
	} else {
		runtime.GOMAXPROCS(c.parallelCnt)
	}
	start := time.Now()
	c.verboseLogf("ggroup prepared, parameters:%+v", *c)
	c.groupWorkers = cyclemap.New[string, int](c.cacheGroupSize, false)
	c.outChan = make(chan string, 100)
	c.workerChans = make([]chan *inputdata, c.parallelCnt)
	ctx := context.Background()
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(c.parallelCnt)
	c.verboseLogf("prepared finished")
	for i := 0; i < c.parallelCnt; i++ {
		c.workerChans[i] = make(chan *inputdata, 100)
		go c.startMatchWorker(ctx, waitGroup, i)

	}
	var finishedChan = make(chan struct{})
	go c.startPrint(ctx, finishedChan)

	for fi := 0; fi < 1 || fi < flag.NArg(); fi++ {
		c.readFile(ctx, fi)
	}
	c.verboseLogf("input read complete")
	for i := 0; i < len(c.workerChans); i++ {
		close(c.workerChans[i])
	}

	c.verboseLogf("waiting for waitGroup")
	waitGroup.Wait()
	close(c.outChan)
	c.verboseLogf("waiting for output chan ")
	<-finishedChan
	if c.verbose {
		now := time.Now()
		c.verboseLogf("time const: %s", now.Sub(start))
	}
}
func (c *command) readFile(ctx context.Context, fi int) {
	var reader *bufio.Reader
	var file *os.File
	var err error
	if flag.NArg() > 0 {
		file, err = os.Open(flag.Args()[fi])
		defer func() {
			if file != nil {
				file.Close()
			}
		}()
		c.verboseLogf("opened file: %s", flag.Args()[fi])
		reader = bufio.NewReaderSize(file, c.cutLineSize)
	} else {
		file = os.Stdin
		reader = bufio.NewReaderSize(file, c.cutLineSize)
		c.verboseLogf("stdin stream reading...")
	}
	var groupColorFormat ColorFormatFunc
	if c.color == "always" {
		groupColorFormat = GrepColorBLUEFormat
	}
	util.CheckPanic(err)
	for bs, exceed, err := reader.ReadLine(); err == nil; bs, exceed, err = reader.ReadLine() {
		if exceed {
			c.errLog.Println("[WARN] a line too long, Cutted to multi-lines matching....")
			c.verboseLogf("log:%s", string(bs))
		}
		line := string(bs)
		groupIndexes := c.groupRegex.FindStringIndex(line)
		if len(groupIndexes) > 1 {
			group := line[groupIndexes[0]:groupIndexes[1]]
			if groupColorFormat != nil {
				line = groupColorFormat(line, groupIndexes)
			}
			input := &inputdata{group: group, line: line}
			var goIndex int
			var ok bool
			if goIndex, ok = c.groupWorkers.Get(group); !ok {
				c.lineIndex++
				goIndex = c.lineIndex % c.parallelCnt
			}
			c.groupWorkers.Set(group, goIndex)
			c.workerChans[goIndex] <- input
		}
	}
}

func (c *command) CheckGrepMatch(ctx context.Context, matchIndex uint, line string) []int {

	if c.orderlyMatch {
		if int(matchIndex) < len(c.regexps) {
			reg := c.regexps[matchIndex]
			return reg.FindStringIndex(line)
		}
	} else {
		for i := 0; i < len(c.regexps); i++ {
			if !util.IsTrue(matchIndex, i) {
				reg := c.regexps[i]
				return reg.FindStringIndex(line)
			}
		}
	}
	c.errLog.Println("unknow error #202401222101")
	os.Exit(1)
	return nil
}

func (c *command) startMatchWorker(ctx context.Context, waitGroup *sync.WaitGroup, goIndex int) {
	c.verboseLogf("orderly scan worker %d started", goIndex)
	var colorFormat ColorFormatFunc
	if c.color == "always" {
		colorFormat = GrepColorRedFormat
	}
	var wokerCycleMaps = cyclemap.New[string, *groupbuff.GroupBuff](c.cacheGroupSize/c.parallelCnt, false)
	if !c.fullMatch {
		wokerCycleMaps.SetListenRemoveFunc(func(group string, item *groupbuff.GroupBuff) {
			c.outChan <- item.String()
		})
	}
	defer func() {
		if !c.fullMatch {
			iter := wokerCycleMaps.Iter()
			for item, ok := iter.First(); ok; item, ok = iter.Next() {
				c.outChan <- item.String()
			}
		}
		waitGroup.Done()
		c.verboseLogf("orderly scan worker %d finished", goIndex)
	}()
	for {
		if lineInfo, ok := <-c.workerChans[goIndex]; ok {
			line := lineInfo.line
			group := lineInfo.group
			groupItem, existsGroup := wokerCycleMaps.Get(group)

			var matchIndex uint = 0
			if existsGroup {
				matchIndex = groupItem.MatchIndex()
			}
			if indexes := c.CheckGrepMatch(ctx, matchIndex, line); len(indexes) > 1 {
				if colorFormat != nil {
					line = colorFormat(line, indexes)
				}
				if !existsGroup {
					groupItem = groupbuff.NewItem(goIndex, c.orderlyMatch, false)
					wokerCycleMaps.Set(group, groupItem)
				}
				if bs, finished := groupItem.Write(line, 0, len(c.regexps), c.mergeLine); finished {
					wokerCycleMaps.Remove(group)
					c.outChan <- string(bs)
				}
			}
		} else {
			c.verboseLogf("worker %d canceled", goIndex)
			break
		}
	}
}

func (c *command) startPrint(ctx context.Context, finishedChan chan struct{}) {
	for {
		if line, ok := <-c.outChan; ok {
			if _, err := fmt.Println(line); err == nil {
				continue
			} else {
				os.Exit(0)
			}
		} else {
			c.verboseLogf("output chan closted")
			close(finishedChan)
			break
		}
	}
}
