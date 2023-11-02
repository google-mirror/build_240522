package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"android/soong/testing/test_spec_proto"
	"google.golang.org/protobuf/proto"
)

type keyToLocksMap struct {
	locks sync.Map
}

func (kl *keyToLocksMap) GetLockForKey(key string) *sync.Mutex {
	mutex, _ := kl.locks.LoadOrStore(key, &sync.Mutex{})
	return mutex.(*sync.Mutex)
}

func getSortedKeys(syncMap *sync.Map) []string {
	var allKeys []string
	syncMap.Range(func(key, _ interface{}) bool {
		allKeys = append(allKeys, key.(string))
		return true
	})

	sort.Strings(allKeys)
	return allKeys
}

func writeOutput(outputFile string, allMetadata []*test_spec_proto.TestSpec_OwnershipMetadata) {
	testSpec := &test_spec_proto.TestSpec{
		OwnershipMetadataList: allMetadata,
	}
	data, err := proto.Marshal(testSpec)
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(outputFile, data, 0644 /* rw-r--r-- */); err != nil {
		log.Fatal(err)
	}
}

func processProtobuf(filePath string, ownershipMetadataMap *sync.Map, keyLocks *keyToLocksMap, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	data, err := ioutil.ReadFile(filePath)
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
		log.Fatal(err)
	}

	filePaths := strings.Split(string(inputFileData), "\n")
	ownershipMetadataMap := &sync.Map{}
	keyLocks := &keyToLocksMap{}
	errCh := make(chan error, len(filePaths))
	var wg sync.WaitGroup

	for _, filePath := range filePaths {
		wg.Add(1)
		go processProtobuf(filePath, ownershipMetadataMap, keyLocks, errCh, &wg)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		log.Fatal(err)
	}

	allKeys := getSortedKeys(ownershipMetadataMap)
	var allMetadata []*test_spec_proto.TestSpec_OwnershipMetadata

	for _, key := range allKeys {
		value, _ := ownershipMetadataMap.Load(key)
		metadataList := value.([]*test_spec_proto.TestSpec_OwnershipMetadata)
		allMetadata = append(allMetadata, metadataList...)
	}

	writeOutput(*outputFile, allMetadata)
}
