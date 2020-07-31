package neo4j

import (
	"github.com/crunchypi/markov-go-sql.git/src/protocols"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// # Neo4jManager implements protocols.DBAbstracter
var _ protocols.DBAbstracter = (*Neo4jManager)(nil)

type Neo4jManager struct {
	db neo4j.Driver
}

// Used internal as args in main execution function.
type executeParams struct {
	cypher       string
	bindings     map[string]interface{}
	callbackMode bool
	callback     func(neo4j.Result) // Optional
}

func New(uri, user, pwd string, encr bool) (protocols.DBAbstracter, error) {
	new := Neo4jManager{}

	driver, err := neo4j.NewDriver(
		uri,
		neo4j.BasicAuth(user, pwd, ""),
		func(c *neo4j.Config) {
			c.Encrypted = encr
		})
	if err != nil {
		return &new, err
	}
	new.db = driver
	return &new, nil
}

// modifier is the point of contact of the neo4j db/driver.
// @ TODO: transactions + batched commands?
func (n *Neo4jManager) execute(x executeParams) error {
	// # Open.
	session, err := n.db.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	// # Execute.
	res, err := session.Run(x.cypher, x.bindings)
	if err != nil {
		return err
	}

	// # Optional callback.
	if x.callbackMode {
		for res.Next() {
			x.callback(res)
		}
	}
	return nil
}

// bindings does a common task in this file: converts
// arguments into a map suitable for the sql pkg
func (n *Neo4jManager) bindings(word, other string, dst int) map[string]interface{} {
	return map[string]interface{}{
		"word":  word,
		"other": other,
		"dst":   dst,
	}
}

func (n *Neo4jManager) newNode(word string) error {
	return n.execute(executeParams{
		cypher: `CREATE (:MChain {word:$word})`,
		bindings: map[string]interface{}{
			"word": word,
		},
	})
}

func (n *Neo4jManager) getNode(word string) []interface{} {
	res := make([]interface{}, 0, 10) // 10 is arbitrary.
	// res := []string{}
	n.execute(executeParams{
		cypher: `MATCH (x) WHERE x.word = $word RETURN x.word AS w`,
		bindings: map[string]interface{}{
			"word": word,
		},
		callbackMode: true,
		callback: func(r neo4j.Result) {
			item, ok := r.Record().Get("w")
			if ok {
				res = append(res, item)
			}
		},
	})
	return res
}

func (n *Neo4jManager) nodeExists(word string) bool {
	return len(n.getNode(word)) > 0
}

func (n *Neo4jManager) newRelationship(word, other string, dst int) error {
	return n.execute(executeParams{
		cypher: `
			MATCH (x), (y)
			WHERE x.word = $word
			  AND y.word = $other
		   CREATE (x)-[:conn{dst:$dst, count:1}]->(y)
		`,
		bindings: map[string]interface{}{
			"word":  word,
			"other": other,
			"dst":   dst,
		},
	})
}

func (n *Neo4jManager) getRelationship(word, other string, dst int) []interface{} {
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
		callbackMode: true,
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

func (n *Neo4jManager) relationshipExists(word, other string, dst int) bool {
	return len(n.getRelationship(word, other, dst)) > 0
}

func (n *Neo4jManager) IncrementPair(word, other string, dst int) {
	for _, v := range []string{word, other} {
		if !n.nodeExists(v) {
			n.newNode(v)
		}
	}

	if !n.relationshipExists(word, other, dst) {
		n.newRelationship(word, other, dst)
		return
	}
	n.execute(executeParams{
		cypher: `
			MATCH (x:MChain{word:$word})-[r:conn{dst:$dst}]->(y:MChain{word:$other})
			SET r.count = r.count + 1
		`,
		bindings: map[string]interface{}{
			"word":  word,
			"other": other,
			"dst":   dst,
		},
	})
}

func (n *Neo4jManager) SucceedingX(word string) []protocols.WordRelationship {

	res := make([]protocols.WordRelationship, 0, 10) // # 10 is arbitrary
	n.execute(executeParams{
		cypher: `
			 MATCH (x:MChain{word:$word})-[r:conn]->(y)
			RETURN y.word AS a, x.word AS b, r.dst AS c, r.count AS d 
		`,
		bindings: map[string]interface{}{
			"word": word,
		},
		callbackMode: true,
		callback: func(r neo4j.Result) {
			a, aok := r.Record().Get("a")
			b, bok := r.Record().Get("b")
			c, cok := r.Record().Get("c")
			d, dok := r.Record().Get("d")
			// log.Printf("####, %s, %s, %d, %d", a, b, c, d)
			if aok && bok && cok && dok {
				newNode := protocols.WordRelationship{
					Word:     a.(string),
					Other:    b.(string),
					Distance: int(c.(int64)),
					Count:    int(d.(int64)),
				}
				res = append(res, newNode)

			}
		},
	})
	return res
}
