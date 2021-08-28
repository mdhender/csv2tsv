// Copyright (c) 2021 Michael D Henderson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package main implements the converter.
package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	input := "input.csv"
	output := "output.tsv"

	flag.StringVar(&input, "i", input, fmt.Sprintf("name of input file [default is %q]", input))
	flag.StringVar(&output, "o", output, fmt.Sprintf("name of output file [default is %q]", output))
	flag.Parse()
	if input == "" || output == "" {
		log.Fatal("please specify both input and output file names")
	}

	if err := run(input, output); err != nil {
		log.Fatal(err)
	}
}

// run the things that need to be run
func run(input, output string) error {
	rows, err := loadFile(input)
	if err != nil {
		return err
	}
	log.Printf("loaded %q\n", input)

	// scrub the input of tabs so that we don't have to worry about quoting fields
	for _, row := range rows {
		for col := range row {
			row[col] = strings.ReplaceAll(row[col], "\t", " ")
		}
	}

	fp, err := os.Create(output)
	if err != nil {
		return err
	}

	for _, row := range rows {
		_, _ = fmt.Fprintln(fp, strings.Join(row, "\t"))
	}

	log.Printf("created %q\n", output)

	return nil
}

// dos2unix accepts CR+LF, CR, and LF+CR endings and replaces them with LF.
func dos2unix(b []byte) []byte {
	const cr, lf = '\r', '\n'

	return bytes.ReplaceAll(bytes.ReplaceAll(bytes.ReplaceAll(b, []byte{cr, lf}, []byte{lf}), []byte{lf, cr}, []byte{lf}), []byte{cr}, []byte{lf})
}

// loadFile loads a CSV file and returns the slice of rows and columns.
func loadFile(name string) ([][]string, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	rdr := csv.NewReader(bytes.NewReader(dos2unix(b)))
	rdr.FieldsPerRecord = -1 // we expect a variable number of fields per line
	return rdr.ReadAll()
}
