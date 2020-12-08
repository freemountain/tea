// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var (
	showLog bool
)

// errorf printf content as an error information
func errorf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
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

// outputList provides general function to convert given list of items
// into several outputs (table, csv, simple, tsv, yaml)
func outputList(output string, headers []string, values [][]string) {
	switch {
	case output == "" || output == "table":
		outputtable(headers, values)
	case output == "csv":
		outputdsv(headers, values, ",")
	case output == "simple":
		outputsimple(headers, values)
	case output == "tsv":
		outputdsv(headers, values, "\t")
	case output == "yaml":
		outputyaml(headers, values)
	default:
		errorf("unknown output type '" + output + "', available types are:\n- csv: comma-separated values\n- simple: space-separated values\n- table: auto-aligned table format (default)\n- tsv: tab-separated values\n- yaml: YAML format\n")
	}
}
