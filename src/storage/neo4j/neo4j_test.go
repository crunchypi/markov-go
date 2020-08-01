package neo4j

import (
	"log"
	"testing"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	uri = "bolt://localhost:7687"
	usr = ""
	pwd = ""
	enc = false
)

func init() {
	if usr == "" || pwd == "" {
		panic("credentials not set")
	}

}

func tryCleanup() {

	n, err1 := New(uri, usr, pwd, enc)
	m := n.(*Manager)
	err2 := m.execute(executeParams{
		cypher:   `MATCH (x) DETACH DELETE x;`,
		bindings: make(map[string]interface{}, 0),
	})
	if err1 != nil || err2 != nil {
		log.Println("Failed while cleaning up- check credentials.")
	}
}

func (n *Manager) getRelationship(word, other string, dst int) []interface{} {
	res := make([]interface{}, 0, 10) // # 10 is arbitrary
	n.execute(executeParams{
		cypher: `
			 MATCH (x:MChain{word:$word})-[r:conn{dst:$dst}]->(y:MChain{word:$other})
			RETURN x.word AS a, y.word AS b, r.dst AS c, r.count AS d
		`,
		bindings: map[string]interface{}{
			"word":  word,
			"other": other,
			"dst":   dst,
		},
		callback: func(r neo4j.Result) {
			a, aok := r.Record().Get("a")
			b, bok := r.Record().Get("b")
			c, cok := r.Record().Get("c")
			d, dok := r.Record().Get("d")
			if aok && bok && cok && dok {
				res = append(res, a, b, c, d)
			}
			// res = append(res, r.Record().Keys())
		},
	})
	return res
}

func TestNew(t *testing.T) {
	_, err := New(uri, usr, pwd, enc)
	if err != nil {
		t.Error("failed to create new Neo4jManager instance:", err)
	}
}

// # Written to be inspected manually.
func TestModifierNew(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Manager)
	m.execute(executeParams{
		cypher:   `CREATE (_:TestNode {t:$val})`,
		bindings: map[string]interface{}{"val": 99},
	})
}

func TestModifierCallback(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Manager)

	val := 44
	// # Create data.
	m.execute(executeParams{
		cypher:   `CREATE (_:TestNode {t:$val})`,
		bindings: map[string]interface{}{"val": val},
	})
	// # Read data
	res := make([]interface{}, 0, 0)
	callback := func(r neo4j.Result) {
		res = r.Record().Values()
	}
	m.execute(executeParams{
		cypher:   `MATCH (x) WHERE x.t = $val RETURN x.t`,
		bindings: map[string]interface{}{"val": val},
		callback: callback,
	})
	// # Check (unsafely)
	if res[0].(int64) != int64(val) {
		t.Error("Unexpected fetch result.")
	}
}

// # created to be checked manually.
func TestNewNodes(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Manager)
	err := m.newNodes([]string{"testword"})
	if err != nil {
		t.Error(err)
	}
}

func TestIncrementPair(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Manager)

	// # generate data.
	nodeA, nodeB := "a", "b"
	dst := 1
	// m.newNodes([]string{nodeA, nodeB})
	m.IncrementPair(nodeA, nodeB, dst)

	// # fetch & check data.
	res := m.getRelationship(nodeA, nodeB, dst)
	if res[0].(string) != nodeA &&
		res[1].(string) != nodeB &&
		res[2].(int) != dst &&
		res[3].(int) != 2 {
		t.Error("failed to fetch:", res)
	}
}

func TestSucceedingX(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Manager)

	// # generate data.
	nodeA, nodeB, nodeC := "a", "b", "c"
	dst := 1
	m.newNodes([]string{nodeA, nodeB, nodeC})

	m.IncrementPair(nodeA, nodeB, dst)
	m.IncrementPair(nodeA, nodeC, dst)

	res := m.SucceedingX(nodeA)
	if len(res) != 2 {
		t.Error("unexpected res count", res)
	}

}

func TestFinal(t *testing.T) {
	tryCleanup()
	n, err := New(uri, usr, pwd, enc)
	if err != nil {
		t.Error("db setup failed. Credentials?")
	}

	nodeA, nodeB, dst := "a", "b", 1
	n.IncrementPair(nodeA, nodeB, dst)

	r := n.SucceedingX(nodeA)
	if r[0].Word != nodeB {
		t.Error("fetch failed", r)
	}

}
