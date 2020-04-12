package chunk_cache

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestOnDisk(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "c")
	defer os.RemoveAll(tmpDir)

	totalDiskSizeMb := int64(6)
	segmentCount := 2

	cache := NewChunkCache(0, tmpDir, totalDiskSizeMb, segmentCount)

	writeCount := 5
	type test_data struct {
		data []byte
		fileId string
	}
	testData := make([]*test_data, writeCount)
	for i:=0;i<writeCount;i++{
		buff := make([]byte, 1024*1024)
		rand.Read(buff)
		testData[i] = &test_data{
			data:   buff,
			fileId: fmt.Sprintf("1,%daabbccdd", i+1),
		}
		cache.SetChunk(testData[i].fileId, testData[i].data)
	}

	for i:=0;i<writeCount;i++{
		data := cache.GetChunk(testData[i].fileId)
		if bytes.Compare(data, testData[i].data) != 0 {
			t.Errorf("failed to write to and read from cache: %d", i)
		}
	}

	cache.Shutdown()

	cache = NewChunkCache(0, tmpDir, totalDiskSizeMb, segmentCount)

	for i:=0;i<writeCount;i++{
		data := cache.GetChunk(testData[i].fileId)
		if bytes.Compare(data, testData[i].data) != 0 {
			t.Errorf("failed to write to and read from cache: %d", i)
		}
	}

	cache.Shutdown()

}
