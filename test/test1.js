var url = "",
    ws = new WebSocket(url);

ws.onopen = function () {
    console.info("onopen")
};
ws.onmessage = function (event) {
    var text = event.data;
    console.info("onmessage: [" + text + "]");
}
ws.onerror = function () {
    console.info("onerror")
};
ws.onclose = function () {
    console.info("onclose")
};
