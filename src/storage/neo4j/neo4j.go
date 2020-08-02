package neo4j

import (
	"fmt"
	"strconv"

	"github.com/crunchypi/markov-go-sql.git/src/storage"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// Manager implements protocols.DBAbstracter
var _ storage.DBAbstracter = (*Manager)(nil)

// namings for cypher
var (
	nodeLabel  = "MChain"
	nodeProp   = "word"
	relLabel   = "conn"
	relPropDst = "dst"
	relPropCnt = "count"
)

// creates a node string with format: (alias:nodeLabel {nodeProp:$bindName}),
// where nodeLabel and nodeProp are defined in the var block at the top of
// this file.
func nodeStr(alias, bindName string) string {
	return fmt.Sprintf(
		"(%s:%s {%s:$%s})",
		alias, nodeLabel,
		nodeProp, bindName,
	)
}

// creates an edge string with format: [alias,relLabel{relPropDst:$bindName}]
// where relLabel and relPropDst are defined in the var block at the top of
// this file.
func edgeStr(alias, bindName string) string {
	return fmt.Sprintf(
		"-[%s:%s {%s:$%s}]->",
		alias, relLabel,
		relPropDst, bindName,
	)
}

// Manager manages a neo4j connection and holds the set of
// methods required to implement protocols.DBAbstracter.
type Manager struct {
	db neo4j.Driver
}

// Used internal as args in main execution function.
type executeParams struct {
	cypher   string
	bindings map[string]interface{}
	callback func(neo4j.Result) // Optional
}

// New attempts to contact a Neo4j DB, given the params. Returns
// A Neo4jManager type (found in this file), which implements
// protocols.DBAbstracter.
func New(uri, user, pwd string, encr bool) (storage.DBAbstracter, error) {
	new := Manager{}

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

// execute is the point of contact of the neo4j db/driver.
// Takes in 'executeParams' struct (found in this file),
// where 'callback' property (function) is optional.
func (n *Manager) execute(x executeParams) error {
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
	if x.callback != nil {
		for res.Next() {
			x.callback(res)
		}
	}
	return nil
}

// newNodes attempts to create a node for each element in 'words'.
// Done with MERGE instead of CREATE (neo4j syntax) for an
// all-or-nothing execution.
func (n *Manager) newNodes(words []string) error {
	cypher := ``
	bindings := make(map[string]interface{})
	for i, v := range words {
		cypher += `MERGE ` + nodeStr("", fmt.Sprintf("%d", i))
		bindings[strconv.Itoa(i)] = v
	}

	return n.execute(executeParams{
		cypher:   cypher,
		bindings: bindings,
	})
}

// IncrementPair attempts to increment the 'count' property on
// the relationship between 'word' and 'other' with a certain 'dst'.
// If this is not possible, a new relationship will be created.
// Note: 'word' and 'other' do not have to be in the db.
func (n *Manager) IncrementPair(word, other string, dst int) error {
	// # Add nodes if they do not already exist - makes cypher simpler.
	n.newNodes([]string{word, other})

	cypher := `
		// If relationship exists, increment it
		OPTIONAL MATCH 
			` + nodeStr("a", "word") + `
			` + edgeStr("b", "dst") + `
			` + nodeStr("c", "other") + ` 
		SET b.` + relPropCnt + ` = b.` + relPropCnt + ` + 1

		// Else - create relship only if nodes exist.
		WITH a AS _
		MATCH 
			` + nodeStr("x", "word") + `,
			` + nodeStr("y", "other") + `
		WHERE NOT 
			(x)-[:conn]->(y)
		CREATE 
			(x)-[z:conn{dst:$dst, count:1}]->(y)
	`
	bindings := map[string]interface{}{
		"word":  word,
		"other": other,
		"dst":   dst,
	}
	return n.execute(executeParams{
		cypher:   cypher,
		bindings: bindings,
	})
}

// SucceedingX retrieves all nodes connected from `word`.
func (n *Manager) SucceedingX(word string) ([]storage.WordRelationship, error) {

	res := make([]storage.WordRelationship, 0, 100) // # 100 is arbitrary
	wrd, othr, dst, cnt := "w", "o", "d", "c"       // # Aliases
	return res, n.execute(executeParams{
		cypher: `
			 MATCH  
			 	  ` + nodeStr("x", "word") + `-[r:conn]->(y)

			RETURN x.word  AS ` + wrd + `, 
				   y.word  AS ` + othr + `, 
				   r.dst   AS ` + dst + `, 
				   r.count AS ` + cnt + `;
		`,
		bindings: map[string]interface{}{"word": word},
		callback: func(r neo4j.Result) {
			// # Unrap relationship @unsafely
			newNode := storage.WordRelationship{}
			if v, ok := r.Record().Get(wrd); ok {
				newNode.Word = v.(string)
			}
			if v, ok := r.Record().Get(othr); ok {
				newNode.Other = v.(string)
			}
			if v, ok := r.Record().Get(dst); ok {
				newNode.Distance = int(v.(int64))
			}
			if v, ok := r.Record().Get(cnt); ok {
				newNode.Count = int(v.(int64))
			}
			res = append(res, newNode)
		},
	})
}
