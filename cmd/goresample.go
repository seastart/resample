/*
	Copyright (C) 2016 - 2017, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.
*/

// The program takes as input a WAV or RAW PCM sound file
// and resamples it to the desired sampling rate.
// The output is RAW PCM data.
// Usage: goresample [flags] input_file output_file

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/seastart/resample"
)

const wavHeader = 44

var (
	format = flag.String("format", "i16", "PCM format")
	ch     = flag.Int("ch", 2, "Number of channels")
	ir     = flag.Int("ir", 44100, "Input sample rate")
	or     = flag.Int("or", 0, "Output sample rate")
)

func main() {
	flag.Parse()
	var frmt int
	switch *format {
	case "i16":
		frmt = resample.I16
	case "i32":
		frmt = resample.I32
	case "f32":
		frmt = resample.F32
	case "f64":
		frmt = resample.F64
	default:
		log.Fatalln("Invalid Format")
	}
	if *ch < 1 {
		log.Fatalln("Invalid channel number")
	}
	if *ir <= 0 || *or <= 0 {
		log.Fatalln("Invalid input or output sample rate")
	}
	if flag.NArg() < 2 {
		log.Fatalln("No input or output files given")
	}
	inputFile := flag.Arg(0)
	outputFile := flag.Arg(1)
	var err error
	var input []byte

	input, err = ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	output, err := os.Create(outputFile)
	if err != nil {
		log.Fatalln(err)
	}
	// Create a Reampler
	res, err := resample.New(output, float64(*ir), float64(*or), *ch, frmt, resample.HighQ)
	if err != nil {
		output.Close()
		os.Remove(outputFile)
		log.Fatalln(err)
	}
	// Skip WAV file header
	if strings.ToLower(filepath.Ext(inputFile)) == ".wav" {
		input = input[wavHeader:]
	}
	// Resample data and wrte to output file
	i, err := res.Write(input)
	res.Close()
	output.Close()
	if err != nil {
		os.Remove(outputFile)
		log.Fatalln(err)
	}
	if i < len(input) {
		log.Fatalln("Short write")
	}
}
