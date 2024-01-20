package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/lingdor/ggrep/buffline"
	"github.com/lingdor/ggrep/util"
	"go.uber.org/automaxprocs/maxprocs"
	"log"
	"os"
	"regexp"
	"runtime"
	"sync"
)

func init() {

}

const (
	COLOR_NONE = "\033[0m"
	COLOR_RED  = "\033[38;5;124m"
)

// ggrep --group 'logid[\w+]' 'LOG1' 'LOG2' --buff-lines 200 --merge-lines
type command struct {
	group         string
	bufferLineMax int
	mergeLine     bool
	verbose       bool
	concurrent    int
	file          string
	inputbuf      chan string
	buffLines     *buffline.LineBuff
	errLog        *log.Logger
	regexps       []*regexp.Regexp
	groupRegex    *regexp.Regexp
	outChan       chan string
}

func (c *command) verboseLogf(content string, args ...any) {
	if c.verbose {
		log.Printf("[verbose] "+content, args...)
	}
}

func (c *command) ready() {

	flag.StringVar(&c.group, "group", "", "group match expression for regex")
	flag.IntVar(&c.bufferLineMax, "buffer-line-max", 1000, "cache lines max size of buffers")
	flag.IntVar(&c.bufferLineMax, "b", 1000, "cache lines max size of buffers(shorthand)")
	flag.BoolVar(&c.mergeLine, "merge-lines", false, "merge output of matched lines")
	flag.BoolVar(&c.mergeLine, "m", false, "merge output of matched lines(shorthand)")
	flag.BoolVar(&c.verbose, "verbose", false, "show debug detail infos")
	flag.IntVar(&c.concurrent, "concurrent", 1, "concurrent of analyse")
	flag.StringVar(&c.file, "file", "", "input file path")
	flag.StringVar(&c.file, "f", "", "input file path(shorthand)")
	c.errLog = log.New(os.Stderr, "[err]", log.Ldate|log.Ltime)
}

func (c *command) validParameters() {
	if flag.NArg() < 1 {
		fmt.Fprint(os.Stderr, "no found match expression in parameters\n")
		os.Exit(1)
	}
	c.regexps = make([]*regexp.Regexp, flag.NArg())
	var err error
	if c.groupRegex, err = regexp.Compile(c.group); err == nil {
		for i, arg := range flag.Args() {
			if c.regexps[i], err = regexp.Compile(arg); err != nil {
				c.errLog.Printf("arg (%s) expression mistake:%s", arg, err.Error())
				os.Exit(1)
			}
		}
	} else {
		c.errLog.Printf("group expression mistake:%s", err.Error())
		os.Exit(1)
	}
}

func (c *command) run() {
	if c.concurrent < 1 {
		maxprocs.Set(maxprocs.Logger(log.Printf))
		c.concurrent = runtime.NumCPU()
	} else {
		runtime.GOMAXPROCS(c.concurrent)
	}
	c.validParameters()
	c.verboseLogf("ggroup prepared, parameters:%+v", *c)
	var reader *bufio.Reader
	var file *os.File
	var err error
	if c.file != "" {
		file, err = os.Open(c.file)
		defer func() {
			if file != nil {
				file.Close()
			}
		}()
		c.verboseLogf("opened file: %s", c.file)
		reader = bufio.NewReader(file)
	} else {
		file = os.Stdin
		reader = bufio.NewReader(file)
		c.verboseLogf("stdin stream reading...")
	}
	util.CheckPanic(err)

	exps := flag.Args()
	c.buffLines = buffline.New(c.bufferLineMax, len(exps), !c.mergeLine)
	c.inputbuf = make(chan string, 100)
	c.outChan = make(chan string, 100)
	ctx := context.Background()
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(c.concurrent)
	c.verboseLogf("prepared finished")
	for i := 0; i < c.concurrent; i++ {
		go c.startScan(ctx, waitGroup, i)
	}
	var finishedChan = make(chan struct{})
	go c.startPrint(ctx, finishedChan)
	for bs, exceed, err := reader.ReadLine(); err == nil; bs, exceed, err = reader.ReadLine() {
		if exceed {
			c.errLog.Println("a line too long, Cutted to multi-lines matching....")
		}
		c.inputbuf <- string(bs)
	}
	close(c.inputbuf)

	c.verboseLogf("waiting for waitGroup")
	waitGroup.Wait()
	close(c.outChan)
	c.verboseLogf("waiting for output chan ")
	<-finishedChan
	c.verboseLogf("done")
}

func (c *command) startScan(ctx context.Context, waitGroup *sync.WaitGroup, i int) {
	c.verboseLogf("scan worker %d started", i)
	defer func() {
		waitGroup.Done()
		c.verboseLogf("scan worker %d finished", i)
	}()

	for {
		if line, ok := <-c.inputbuf; ok {

			for i, reg := range c.regexps {
				if reg.MatchString(line) {
					indexes := c.groupRegex.FindStringIndex(line)
					if len(indexes) > 1 {
						group := line[indexes[0]:indexes[1]]
						if bs, matched := c.buffLines.Write(group, line, i); matched {
							c.outChan <- string(bs)
						}
					}
				}
			}
		} else {
			break
		}
	}
}

func (c *command) startPrint(ctx context.Context, finishedChan chan struct{}) {
	for {
		if line, ok := <-c.outChan; ok {
			fmt.Println(line)
		} else {
			c.verboseLogf("output chan closted")
			close(finishedChan)
			break
		}
	}
}
