package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

func pprint(ss pairList) {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"

	tw := tabwriter.NewWriter(os.Stdout, 7, 2, 0, ' ', tabwriter.AlignRight)
	for lines, kv := range ss {
		if lines > 50 {
			break
		}
		_, _ = fmt.Fprintf(tw, "%s%d\t %s%s\n", colorGreen, kv.Value, colorReset, kv.Key)
	}
	_ = tw.Flush()
}

func readInput(reader io.Reader, s *store, minWordLength int) error {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		t := scanner.Text()
		if len(t) >= minWordLength {
			v := s.load(t)
			s.store(t, v+1)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	plist := s.loadAll(true)
	pprint(plist)
	return nil
}
