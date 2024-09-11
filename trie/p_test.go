package trie

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/internal/testrand"
	"github.com/ethereum/go-ethereum/trie/trienode"
)

func TestCommit(t *testing.T) {
	// testCommit(100)
	// testCommit(200)
	// testCommit(300)
	// testCommit(400)
	// testCommit(500)
	testCommit(1000)
	testCommit(2000)
	testCommit(3000)
}

func testCommit(n int) {
	trA := NewEmpty(nil)
	trB := NewEmpty(nil)
	for i := 0; i < n; i++ {
		key := testrand.Bytes(32)
		val := testrand.Bytes(32)
		trA.Update(key, val)
		trB.Update(key, val)
	}
	trA.Hash()
	trB.Hash()

	// non-sync mode
	trA.unhashed = 0
	ss := time.Now()
	trA.Commit(true)
	fmt.Printf("item: %d commit time (single mode): %v\n", n, time.Since(ss))

	ss = time.Now()
	trB.Commit(true)
	fmt.Printf("item: %d commit time (parallel mode): %v\n", n, time.Since(ss))
}

func TestAddNode(t *testing.T) {
	testAddNode(100)
	testAddNode(200)
	testAddNode(300)
	testAddNode(400)
	testAddNode(500)

}

func testAddNode(n int) {
	toBeAdded := make(map[string]*trienode.Node)
	for i := 0; i < n; i++ {
		toBeAdded[string(testrand.Bytes(32))] = trienode.NewDeleted()
	}

	testAddNode_single(n, toBeAdded)
	testAddNode_parallel(n, toBeAdded)
}

func testAddNode_single(n int, toBeAdded map[string]*trienode.Node) {
	nodes := trienode.NewNodeSet(common.Hash{})

	start := time.Now().UnixMicro()
	for path, node := range toBeAdded {
		nodes.AddNode([]byte(path), node)
	}
	end := time.Now().UnixMicro()
	fmt.Println("number: ", n, "  single cost: ", end-start)
}

func testAddNode_parallel(n int, toBeAdded map[string]*trienode.Node) {
	nodes := trienode.NewNodeSet(common.Hash{})

	start := time.Now().UnixMicro()
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for path, node := range toBeAdded {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			nodes.AddNode([]byte(path), node)
		}()
	}
	wg.Wait()
	end := time.Now().UnixMicro()
	fmt.Println("number: ", n, "  parallel cost: ", end-start)
}
