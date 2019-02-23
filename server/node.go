package server

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Node struct {
	PublicKey string    `json:"publicKey"`
	Ip        string    `json:"ip"`
	Height    uint64    `json:"height"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Timestamp time.Time `json:"timestamp"`
	TimeDff   string    `json:"timediff"`
}

type Geo struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (n *Node) FillLatLng() error {

	var geo Geo
	ipGeoUrl := "http://api.ipstack.com/" + n.Ip + "?access_key=d3453a877bd6805e526f77c249888b7d&fields=latitude,longitude"

	response, err := http.Get(ipGeoUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &geo); err != nil {
		fmt.Println(err)
		return err
	}

	n.Lat = geo.Latitude
	n.Lng = geo.Longitude

	return nil
}

// Serialize converts block to []byte
func (n *Node) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(n)
	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

// Deserialize converts []byte to block
func Deserialize(blockBytes []byte) *Node {
	var n Node

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&n)
	if err != nil {
		log.Panic(err)
	}

	return &n
}
