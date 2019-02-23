// Initialize and add the map
function initMap() {
    // The location of Taipei
    var taipei = { lat: 25.105497, lng: 125.105497 };
    // The map, centered at Taipei
    var map = new google.maps.Map(
        document.getElementById('map'), { zoom: 4, center: taipei });

    fetchNods(map)
}

function fetchNods(map) {
    fetch('http://34.80.48.90:8080/nodes', {})
        .then((response) => {
            // response is a readableStream object,
            // converts response by blob(), json(), text()
            return response.json();
        }).then((nodesJson) => {
            for (var i = 0; i < nodesJson.length; i++) {
                var obj = nodesJson[i];

                // show shadow maker if no updates more than one hour
                if (obj.timediff.includes("h")) {
                    iconUrl = "http://maps.google.com/mapfiles/ms/icons/msmarker.shadow.png";
                    status = "off"
                } else {
                    iconUrl = "http://maps.google.com/mapfiles/ms/icons/red-dot.png";
                    status = "on"
                }

                // render makers
                var marker = new google.maps.Marker({
                    position: { lat: obj.lat, lng: obj.lng },
                    map: map,
                    icon: { url: iconUrl }
                });

                var infowindow = new google.maps.InfoWindow();
                bindInfoWindow(marker, map, infowindow, createMakerContent(obj, status));
            }
            console.log(nodesJson);
        }).catch((err) => {
            console.log(err);
        });
}

function bindInfoWindow(marker, map, infowindow, markerContent) {
    marker.addListener('click', function () {
        infowindow.setContent(markerContent);
        infowindow.open(map, marker);
    });
}

function createMakerContent(obj, status) {
    if (obj.height == 0) {
        obj.height = "Unknown"
    }
    s = '<div id="content">' +
        'Latitude: ' + obj.lat + '<br>' +
        'Longitude: ' + obj.lng + '<br>' +
        'Block Height: ' + obj.height + '<br>' +
        'Last update: ' + obj.timediff + '<br>' +
        '</div>';
    return s
}