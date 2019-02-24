package server

import (
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var n Node
	err := c.BindJSON(&n)
	if err != nil {
		c.String(http.StatusBadRequest, "can not parse json")
		return
	}

	ns, ok := c.MustGet("nodes").(*Nodes)
	if !ok {
		c.String(http.StatusInternalServerError, "db error")
		return
	}

	// check if publickey exists
	if ns.SearchNodeExisting(n.PublicKey) {
		// update height and timestamp
		ns.UpdateNode(n)
	} else {
		// fill in lat and lng and add node
		if err = n.FillLatLng(); err != nil {
			c.String(http.StatusInternalServerError, "fillLatLng error")
			return
		}
		ns.AddNodeToNodes(n)
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": 1,
	})
}

func GetNodes(c *gin.Context) {
	ns, ok := c.MustGet("nodes").(*Nodes)
	if !ok {
		c.String(http.StatusInternalServerError, "db error")
		return
	}

	var nodesArray []Node
	if err := ns.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodesTableName))
		if err := b.ForEach(func(k, v []byte) error {
			n := Deserialize(v)

			// rewrite timediff corresponding to server time
			timeDiff := time.Since(n.Timestamp).String()
			n.TimeDff = timeDiff

			nodesArray = append(nodesArray, *n)
			return nil
		}); err != nil {
			log.Panic(err)
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}

	c.JSON(http.StatusOK, nodesArray)
}
