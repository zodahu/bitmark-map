package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/zodahu/bitmark-map/server"
)

const nodesTableName = "nodes"

func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run the other tests
	os.Exit(m.Run())
}

func TestRegister(t *testing.T) {
	w := httptest.NewRecorder()
	r := getRouter(false)
	r.Use(nodesMiddleware())
	r.POST("/register", server.Register)

	registerPayload := getRegisterPOSTPayload()
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(registerPayload))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	printNodes()
}

func TestFillLatLng(t *testing.T) {
	n := server.Node{
		PublicKey: "be44ca9a1",
		Ip:        "54.199.155.101",
		Height:    11,
		Lat:       0,
		Lng:       0,
		Timestamp: time.Now(),
	}

	n.FillLatLng()
	fmt.Println(n)
}

func getRegisterPOSTPayload() []byte {
	params := make(gin.H)
	params["publicKey"] = "be44ca9a1249b6ce095f0867e9ec9d5c6c2e4b6043fd55349a5a9b0c7b9a4813"
	params["ip"] = "54.199.155.101"
	params["height"] = 666
	params["lat"] = 1.2931
	params["lng"] = 103.8558
	params["timestamp"] = time.Now()
	params["timediff"] = "4440h0m0s"
	jsonByte, _ := json.Marshal(params)
	return jsonByte
}

// Helper function to create a router during testing
func getRouter(withTemplates bool) *gin.Engine {
	r := gin.Default()
	if withTemplates {
		r.LoadHTMLGlob("views/*")
	}
	return r
}

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	// Create a response recorder
	w := httptest.NewRecorder()
	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

// PrintNodes prints blocks
func printNodes() {
	ns := server.GetNodesInstance()
	if err := ns.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodesTableName))
		if err := b.ForEach(func(k, v []byte) error {
			n := server.Deserialize(v)
			fmt.Printf("%v\n", n.PublicKey)
			fmt.Printf("%v\n", n.Ip)
			fmt.Printf("%v\n", n.Lat)
			fmt.Printf("%v\n", n.Lng)
			fmt.Printf("%v\n", n.Timestamp)
			fmt.Printf("%v\n", n.TimeDff)
			fmt.Println()
			return nil
		}); err != nil {
			log.Panic(err)
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}
}
