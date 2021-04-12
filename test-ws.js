let webSocket = new WebSocket('wss://localhost/api/ws');
webSocket.onmessage = function(e) { console.log(e)}
webSocket.send("test")
