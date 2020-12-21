// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// table provides infrastructure to easily print (sorted) lists in different formats
type table struct {
	headers    []string
	values     [][]string
	sortDesc   bool // used internally by sortable interface
	sortColumn uint // ↑
}

// printable can be implemented for structs to put fields dynamically into a table
type printable interface {
	FormatField(field string) string
}

// high level api to print a table of items with dynamic fields
func tableFromItems(fields []string, values []printable) table {
	t := table{headers: fields}
	for _, v := range values {
		row := make([]string, len(fields))
		for i, f := range fields {
			row[i] = v.FormatField(f)
		}
		t.addRowSlice(row)
	}
	return t
}

func tableWithHeader(header ...string) table {
	return table{headers: header}
}

// it's the callers responsibility to ensure row length is equal to header length!
func (t *table) addRow(row ...string) {
	t.addRowSlice(row)
}

// it's the callers responsibility to ensure row length is equal to header length!
func (t *table) addRowSlice(row []string) {
	t.values = append(t.values, row)
}

func (t *table) sort(column uint, desc bool) {
	t.sortColumn = column
	t.sortDesc = desc
	sort.Stable(t) // stable to allow multiple calls to sort
}

// sortable interface
func (t table) Len() int      { return len(t.values) }
func (t table) Swap(i, j int) { t.values[i], t.values[j] = t.values[j], t.values[i] }
func (t table) Less(i, j int) bool {
	const column = 0
	if t.sortDesc {
		i, j = j, i
	}
	return t.values[i][t.sortColumn] < t.values[j][t.sortColumn]
}

func (t *table) print(output string) {
	switch output {
	case "", "table":
		outputtable(t.headers, t.values)
	case "csv":
		outputdsv(t.headers, t.values, ",")
	case "simple":
		outputsimple(t.headers, t.values)
	case "tsv":
		outputdsv(t.headers, t.values, "\t")
	case "yml", "yaml":
		outputyaml(t.headers, t.values)
	default:
		fmt.Printf("unknown output type '" + output + "', available types are:\n- csv: comma-separated values\n- simple: space-separated values\n- table: auto-aligned table format (default)\n- tsv: tab-separated values\n- yaml: YAML format\n")
	}
}

// outputtable prints structured data as table
func outputtable(headers []string, values [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	if len(headers) > 0 {
		table.SetHeader(headers)
	}
	for _, value := range values {
		table.Append(value)
	}
	table.Render()
}

// outputsimple prints structured data as space delimited value
func outputsimple(headers []string, values [][]string) {
	for _, value := range values {
		fmt.Printf(strings.Join(value, " "))
		fmt.Printf("\n")
	}
}

// outputdsv prints structured data as delimiter separated value format
func outputdsv(headers []string, values [][]string, delimiterOpt ...string) {
	delimiter := ","
	if len(delimiterOpt) > 0 {
		delimiter = delimiterOpt[0]
	}
	fmt.Println("\"" + strings.Join(headers, "\""+delimiter+"\"") + "\"")
	for _, value := range values {
		fmt.Printf("\"")
		fmt.Printf(strings.Join(value, "\""+delimiter+"\""))
		fmt.Printf("\"")
		fmt.Printf("\n")
	}
}

// outputyaml prints structured data as yaml
func outputyaml(headers []string, values [][]string) {
	for _, value := range values {
		fmt.Println("-")
		for j, val := range value {
			intVal, _ := strconv.Atoi(val)
			if strconv.Itoa(intVal) == val {
				fmt.Printf("    %s: %s\n", headers[j], val)
			} else {
				fmt.Printf("    %s: '%s'\n", headers[j], val)
			}
		}
	}
}

func isMachineReadable(outputFormat string) bool {
	switch outputFormat {
	case "yml", "yaml", "csv":
		return true
	}
	return false
}
