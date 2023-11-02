package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"android/soong/testing/test_spec_proto"
	"google.golang.org/protobuf/proto"
)

type KeyLocks struct {
	locks map[string]*sync.Mutex
	mu    sync.Mutex
}

func (kl *KeyLocks) GetLockForKey(key string) *sync.Mutex {
	kl.mu.Lock()
	defer kl.mu.Unlock()

	if _, exists := kl.locks[key]; !exists {
		kl.locks[key] = &sync.Mutex{}
	}

	return kl.locks[key]
}

func main() {
	inputFile := flag.String("inputFile", "", "Input file path")
	outputFile := flag.String("outputFile", "", "Output file path")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Usage: metadata -inputFile <input file path> -outputFile <output file path>")
		os.Exit(1)
	}

	inputFileData, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filePaths := strings.Split(string(inputFileData), "\n")
	ownershipMetadataMap := &sync.Map{}
	keyLocks := &KeyLocks{locks: make(map[string]*sync.Mutex)}
	errCh := make(chan error, len(filePaths))
	var wg sync.WaitGroup

	for _, filePath := range filePaths {
		wg.Add(1)
		go func(fp string) {
			defer wg.Done()

			data, err := ioutil.ReadFile(fp)
			if err != nil {
				errCh <- err
				return
			}

			fileContent := strings.TrimRight(string(data), "\n")
			testData := test_spec_proto.TestSpec{}
			err = proto.Unmarshal([]byte(fileContent), &testData)
			if err != nil {
				errCh <- err
				return
			}

			ownershipMetadata := testData.GetOwnershipMetadataList()
			for _, metadata := range ownershipMetadata {
				key := metadata.GetTargetName()
				lock := keyLocks.GetLockForKey(key)
				lock.Lock()

				value, loaded := ownershipMetadataMap.LoadOrStore(key, []*test_spec_proto.TestSpec_OwnershipMetadata{metadata})
				if loaded {
					existingMetadata := value.([]*test_spec_proto.TestSpec_OwnershipMetadata)
					isDuplicate := false
					for _, existing := range existingMetadata {
						if metadata.GetTrendyTeamId() != existing.GetTrendyTeamId() {
							errCh <- fmt.Errorf("error: Conflicting TrendyTeamId for %s", key)
							lock.Unlock()
							return
						}
						if metadata.GetTrendyTeamId() == existing.GetTrendyTeamId() && metadata.GetPath() == existing.GetPath() {
							isDuplicate = true
							break
						}
					}
					if !isDuplicate {
						existingMetadata = append(existingMetadata, metadata)
						ownershipMetadataMap.Store(key, existingMetadata)
					}
				}

				lock.Unlock()
			}
		}(filePath)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		fmt.Println(err)
		os.Exit(1)
	}

	var allMetadata []*test_spec_proto.TestSpec_OwnershipMetadata
	ownershipMetadataMap.Range(func(_, value interface{}) bool {
		metadataList := value.([]*test_spec_proto.TestSpec_OwnershipMetadata)
		allMetadata = append(allMetadata, metadataList...)
		return true
	})

	testSpec := &test_spec_proto.TestSpec{
		OwnershipMetadataList: allMetadata,
	}
	data, err := proto.Marshal(testSpec)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(*outputFile, data, 0644 /* rw-r--r-- */); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
