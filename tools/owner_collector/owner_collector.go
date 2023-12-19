package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"android/soong/android/owner_team_proto"

	"google.golang.org/protobuf/encoding/prototext"
)

// TODO(ron): remove?
func readFileToString(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func main() {
	inputFile := flag.String("inputFile", "", "Input file path, file of filenames")
	outputFile := flag.String("outputFile", "", "Output file path")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Usage: owner_collector -inputFile <input file path> -outputFile <output file path>")
		os.Exit(1)
	}

	// TODO(ron): is single input special cased?
	inputFileData := strings.TrimRight(readFileToString(*inputFile), "\n")
	filePaths := strings.Split(inputFileData, " ")
	allTeams := make([]*owner_team_proto.OwnerTeam, len(filePaths))
	for i, filePath := range filePaths {
		// or ioutil.ReadFile?
		moduleOwnerProtoBytes, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("Error processing file") // TOOD(ron) add filePath to msg.
		}

		c := new(owner_team_proto.OwnerTeam)
		fmt.Printf("test owner file: %s\n", filePath)
		if err := prototext.Unmarshal(moduleOwnerProtoBytes, c); err != nil {
			log.Fatalln("Failed to parse owner file to B: ", err)
		}

		allTeams[i] = c
	}

	allOwnerTeams := owner_team_proto.AllOwnerTeams{OwnerTeams: allTeams}

	for _, p := range allOwnerTeams.OwnerTeams {
		fmt.Println(p.String())
	}

	// Marshal allOwnerTeams  to outputFile.
	os.WriteFile(*outputFile, []byte(allOwnerTeams.String()), 0755)
}
