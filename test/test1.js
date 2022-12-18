var url = "",
    ws = new WebSocket(url);

ws.onopen = function () {
    console.info("onopen.")
};
ws.onmessage = function (event) {
    var text = event.data;
    console.info(`onmessage: [${text}].`);
}
ws.onerror = function (ev) {
    console.error("onerror.")
    console.error(ev)
};
ws.onclose = function (ev) {
    console.info(`onclose with code(${ev.code}), reason(${ev.reason}) and wasClean(${ev.wasClean}).`)
    console.error(ev)
};
