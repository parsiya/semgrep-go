package output

import (
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// HitMap is a map of strings to ints. It's used to track hits for a given rule
// and/or file. The index is the ruleID or filename. The value is the number of
// times the rule was hit or the number of findings in the file.
type HitMap map[string]int

// ruleHitMap creates a HitMap where the key is ruleID and the value is the
// number of hits for that rule.
func ruleHitMap(o Output) HitMap {
	// Create the hitmap.
	hm := make(HitMap)
	// Go through the results and create the hitmap.
	for _, result := range o.Results {
		hm[result.RuleID()]++
	}
	return hm
}

// fileHitMap creates a HitMap where the key is file path and the value is
// the number of hits in that file.
func fileHitMap(o Output) HitMap {
	// Create the hitmap.
	hm := make(HitMap)
	// Go through the results and create the hitmap.
	for _, result := range o.Results {
		hm[result.FilePath()]++
	}
	return hm
}

// -----

// HitMapRow is a single row in the HitMap.
type HitMapRow struct {
	Key   string
	Value string
}

// SortedData returns the HitMap data as a series of rows. If sortByCount is true,
// the data will be sorted by count, otherwise they will be sorted by key.
//
// Note: Sorting by count is in descending order because most of the time we
// assume we want rules/files with higher hits first. Sorting by key is
// ascending order alphabetically.
func (m HitMap) SortedData(sortByCount bool) []HitMapRow {
	// First get a slice of keys.
	mapKeys := make([]string, len(m))
	// We're using an index instead of append because we know the size of the
	// array and it's probably a bit faster.
	index := 0
	for k := range m {
		mapKeys[index] = k
		index++
	}
	// Sort by value.
	if sortByCount {
		// We can sort using sort.Slice and providing our own less function that
		// does the comparison by value of keys. Note: We're sorting in
		// descending order here because most of the time we assume we want
		// findings/files with higher hits first.
		sort.Slice(mapKeys, func(i, j int) bool {
			return m[mapKeys[i]] > m[mapKeys[j]]
		})
	} else {
		// Sort by key. This is in ascending order.
		sort.Strings(mapKeys)
	}

	// mapKeys is sorted. We can create the data.
	data := make([]HitMapRow, len(m))
	for index, key := range mapKeys {
		data[index] = HitMapRow{Key: key, Value: strconv.Itoa(m[key])}
	}
	return data
}

// Returns the HitMap data as a text table. This is useful for text printing.
// It's the caller's responsibility to pass the correct number of headers.
func (m HitMap) ToStringTable(headers []string, sortByCount bool) string {
	// Convert the HitMap data into a [][]string to pass to tablewriter.
	data := make([][]string, len(m))
	for index, row := range m.SortedData(sortByCount) {
		data[index] = []string{row.Key, row.Value}
	}

	// Create the string builder to hold the result.
	var final strings.Builder
	// Create the table writer and set the destination to the string builder.
	table := tablewriter.NewWriter(&final)
	// Set the headers.
	table.SetHeader(headers)
	// Append the data.
	table.AppendBulk(data)
	// Render the table.
	table.Render()
	// Return the data.
	return final.String()
}
