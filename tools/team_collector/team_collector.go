package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"android/soong/android/team_proto"

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
		fmt.Println("Usage: team_collector -inputFile <input file path> -outputFile <output file path>")
		os.Exit(1)
	}

	// TODO(ron): is single input special cased?
	inputFileData := strings.TrimRight(readFileToString(*inputFile), "\n")
	filePaths := strings.Split(inputFileData, " ")
	teams := make([]*team_proto.Team, len(filePaths))
	for i, filePath := range filePaths {
		// or ioutil.ReadFile?
		teamProtoBytes, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("Error processing file") // TOOD(ron) add filePath to msg.
		}

		c := new(team_proto.Team)
		fmt.Printf("team file: %s\n", filePath)
		if err := prototext.Unmarshal(teamProtoBytes, c); err != nil {
			log.Fatalln("Failed to parse team file to B: ", err)
		}

		teams[i] = c
	}

	allTeams := team_proto.AllTeams{Teams: teams}

	for _, p := range allTeams.Teams {
		fmt.Println(p.String())
	}

	// Marshal allTeams  to outputFile.
	// TODO(rbraunstein): This is writing a textproto, not serialized proto.
	// Decide which we want.
	os.WriteFile(*outputFile, []byte(allTeams.String()), 0755)
}
