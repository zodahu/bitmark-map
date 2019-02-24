# Nodemap service
Nodemap collects Bitmark peers information and draws them on the Google map, this project is based on https://github.com/zodahu/bitmark-node/tree/release_v0.97
    
* Backend
  * provides register and get node api
  * stores nodes data in boltdb

* Fronted
  * displays details by clicking marker
  * binds server ip to 34.80.48.90:8080

## Installation

```
docker pull zodahu/bitmark-map
```

## Run
```
docker run -d --name bitmarkMap -p 8080:8080 zodahu/bitmark-map
```
Note that nodemap server ip is bound to http://34.80.48.90:8080 currently, so the map will retrieve data from http://34.80.48.90:8080 instead of your local nodemap service.
