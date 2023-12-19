package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"android/soong/android/team_proto"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// First read the inputFile that list all the textproto file to read.
// Then read each textproto and aggregate into a final proto and serialize it to outputFile
func main() {
	inputFile := flag.String("inputFile", "", "Input file path, file of filenames")
	outputFile := flag.String("outputFile", "", "Output file path")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Usage: team_collector -inputFile <input file path> -outputFile <output file path>")
		os.Exit(1)
	}

	listOfFiles, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error processing file: %s.  %q\n", *inputFile, err)
	}

	filePaths := strings.Split(string(listOfFiles), " ")
	teams := make([]*team_proto.Team, len(filePaths))
	for i, filePath := range filePaths {

		teamProtoBytes, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Error processing file: %s. %q\n", filePath, err)
		}

		c := new(team_proto.Team)
		if err := prototext.Unmarshal(teamProtoBytes, c); err != nil {
			log.Fatalf("Failed to parse team file: %s. %q", filePath, err)
		}

		teams[i] = c
	}

	allTeams := team_proto.AllTeams{Teams: teams}
	data, err := proto.Marshal(&allTeams)
	if err != nil {
		log.Fatalf("Unable to marshal proto out for input: %s.  %q\n", *inputFile, err)
	}
	os.WriteFile(*outputFile, data, 0755)
}
