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
	m := n.(*Neo4jManager)
	err2 := m.execute(executeParams{
		cypher:   `MATCH (x) DETACH DELETE x;`,
		bindings: make(map[string]interface{}, 0),
	})
	if err1 != nil || err2 != nil {
		log.Println("Failed while cleaning up")
	}
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
	m := n.(*Neo4jManager)
	m.execute(executeParams{
		cypher:   `CREATE (_:TestNode {t:$val})`,
		bindings: map[string]interface{}{"val": 99},
	})
}

func TestModifierCallback(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Neo4jManager)

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
		cypher:       `MATCH (x) WHERE x.t = $val RETURN x.t`,
		bindings:     map[string]interface{}{"val": val},
		callbackMode: true,
		callback:     callback,
	})
	// # Check (unsafely)
	if res[0].(int64) != int64(val) {
		t.Error("Unexpected fetch result.")
	}
}

// # created to be checked manually.
func TestNewNode(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Neo4jManager)
	err := m.newNode("testword")
	if err != nil {
		t.Error(err)
	}
}

// # created to be checked manually.
func TestNodeExists(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Neo4jManager)

	word := "test"
	if m.nodeExists(word) {
		t.Error("unexpected: node exists")
	}
	m.newNode(word)
	if !m.nodeExists(word) {
		t.Error("unexpected: node does not exist")
	}
}

func TestGetNode(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Neo4jManager)

	word := "test"
	m.newNode(word)
	r := m.getNode(word)
	if r[0].(string) != word {
		t.Error("not found")
	}
}

// # Written to be tested manually.
func TestNewRelationship(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Neo4jManager)

	x, y := "x", "y"
	m.newNode(x)
	m.newNode(y)
	m.newRelationship(x, y, 1)
}

func TestGetRelationship(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Neo4jManager)

	nodeA, nodeB := "a", "b"
	dst := 1
	m.newNode(nodeA)
	m.newNode(nodeB)
	m.newRelationship(nodeA, nodeB, dst)

	res := m.getRelationship(nodeA, nodeB, dst)
	// t.Log(err, res)
	if res[0].(string) != nodeA &&
		res[1].(string) != nodeB &&
		res[2].(int) != dst &&
		res[3].(int) != 1 {
		t.Error("didnt get rel:", res)
	}
}

func TestRelationshipExists(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Neo4jManager)

	nodeA, nodeB := "a", "b"
	dst := 1
	m.newNode(nodeA)
	m.newNode(nodeB)
	m.newRelationship(nodeA, nodeB, dst)

	if !m.relationshipExists(nodeA, nodeB, dst) {
		t.Error("relship should exist")
	}
}

func TestIncrementPair(t *testing.T) {
	tryCleanup()
	n, _ := New(uri, usr, pwd, enc)
	m := n.(*Neo4jManager)

	// # generate data.
	nodeA, nodeB := "a", "b"
	dst := 1
	m.newNode(nodeA)
	m.newNode(nodeB)
	m.newRelationship(nodeA, nodeB, dst)
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
	m := n.(*Neo4jManager)

	// # generate data.
	nodeA, nodeB, nodeC := "a", "b", "c"
	dst := 1
	m.newNode(nodeA)
	m.newNode(nodeB)
	m.newNode(nodeC)

	m.newRelationship(nodeA, nodeB, dst)
	m.newRelationship(nodeA, nodeC, dst)

	res := m.SucceedingX(nodeA)
	if len(res) != 2 {
		t.Error("unexpected res count", res)
	}

}
