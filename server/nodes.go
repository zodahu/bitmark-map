package server

import (
	"log"
	"sync"

	"github.com/boltdb/bolt"
)

const dbName = "node-map.db"
const nodesTableName = "nodes"
const centralNode = "Last"

// Nodes .
type Nodes struct {
	DB *bolt.DB
}

var instance *Nodes
var once sync.Once

// GetNodesInstance ensures only one Nodes instance are created
func GetNodesInstance() *Nodes {
	once.Do(func() {
		instance = getNodesDB()
	})
	return instance
}

// getNodesDB creates an new DB or retrieves existing DB
func getNodesDB() *Nodes {

	// Open the db file in your current directory.
	// It will be created if it doesn't exist.x
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		// try to find the table
		b := tx.Bucket([]byte(nodesTableName))

		if b == nil { // create genesis block for empty DB
			b, err = tx.CreateBucket([]byte(nodesTableName))
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Nodes{db}
}

// AddNodeToNodes inserts node to DB
func (ns *Nodes) AddNodeToNodes(n Node) {
	err := ns.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodesTableName))
		if b != nil {
			// Put new node to DB
			err := b.Put([]byte(n.PublicKey), n.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (ns *Nodes) UpdateNode(n Node) {
	err := ns.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodesTableName))
		if b != nil {
			// get node
			var nTemp *Node
			if v := b.Get([]byte(n.PublicKey)); v != nil {
				nTemp = Deserialize(v)
				nTemp.Height = n.Height
				nTemp.Timestamp = n.Timestamp
				nTemp.TimeDff = n.TimeDff
			} else {
				return nil
			}

			// update node
			err := b.Put([]byte(nTemp.PublicKey), nTemp.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (ns *Nodes) SearchNodeExisting(pk string) bool {
	var existing bool

	if err := ns.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodesTableName))
		if b.Get([]byte(pk)) != nil {
			existing = true
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}

	return existing
}
