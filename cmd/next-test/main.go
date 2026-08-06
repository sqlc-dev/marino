// Command next-test prints the inputs the recursive-descent parser most
// recently fell back to goyacc for, ranked, so the next construct to port
// is easy to pick (the family dev loop; see PLAN.md and CLAUDE.md).
//
// Usage:
//
//	rm -f /tmp/rd.log
//	MARINO_RD_LOG=/tmp/rd.log go test ./parser -count=1
//	go run ./cmd/next-test /tmp/rd.log
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	path := "/tmp/rd.log"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "next-test: %v\nrun: MARINO_RD_LOG=%s go test ./parser -count=1\n", err, path)
		os.Exit(1)
	}
	// Records are reason\x00sql\x00\n.
	type entry struct {
		reason string
		sqls   map[string]bool
	}
	byReason := map[string]*entry{}
	total := 0
	seen := map[string]bool{}
	for _, rec := range strings.Split(string(data), "\x00\n") {
		parts := strings.SplitN(rec, "\x00", 2)
		if len(parts) != 2 {
			continue
		}
		reason, sql := parts[0], parts[1]
		key := reason + "\x00" + sql
		if seen[key] {
			continue
		}
		seen[key] = true
		total++
		e := byReason[reason]
		if e == nil {
			e = &entry{reason: reason, sqls: map[string]bool{}}
			byReason[reason] = e
		}
		e.sqls[sql] = true
	}
	entries := make([]*entry, 0, len(byReason))
	for _, e := range byReason {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		if len(entries[i].sqls) != len(entries[j].sqls) {
			return len(entries[i].sqls) > len(entries[j].sqls)
		}
		return entries[i].reason < entries[j].reason
	})
	fmt.Printf("%d distinct fallbacks in %d reason groups\n\n", total, len(entries))
	for i, e := range entries {
		if i >= 30 {
			fmt.Printf("... and %d more groups\n", len(entries)-i)
			break
		}
		shortest := ""
		for s := range e.sqls {
			if shortest == "" || len(s) < len(shortest) {
				shortest = s
			}
		}
		if len(shortest) > 120 {
			shortest = shortest[:120] + "..."
		}
		fmt.Printf("%6d  %s\n        e.g. %s\n", len(e.sqls), e.reason, strings.ReplaceAll(shortest, "\n", " "))
	}
}
