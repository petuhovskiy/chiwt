function sendChat() {
    const message = chatInput.value
    sendChatRequest(message)
}

function sendChatRequest(message) {
    const data = {
        message: message,
    }

    fetch('/chat/' + chatName + '/send', {
        method: 'POST',
        mode: 'cors', // no-cors, *cors, same-origin
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        credentials: 'same-origin', // include, *same-origin, omit
        headers: {
            'Content-Type': 'application/json'
            // 'Content-Type': 'application/x-www-form-urlencoded',
        },
        redirect: 'follow', // manual, *follow, error
        referrerPolicy: 'no-referrer', // no-referrer, *client
        body: JSON.stringify(data) // body data type must match "Content-Type" header
    })
}

function subscribeToChat(dest, url) {
    const ws = new WebSocket(url)
    ws.onmessage = (e) => {
        const msg = JSON.parse(e.data)
        appendChatMessage(dest, msg)
    }

    ws.onopen = () => {
        // TODO: show connected
    }

    ws.onclose = () => {

    }
}

function appendChatMessage(dest, msg) {
    console.log(dest, msg)

    const p = document.createElement("p")
    p.innerHTML = msg.From + ": " + msg.Data

    dest.appendChild(p)
}

if (chatName) {
    subscribeToChat(chatContent, "ws://" + location.host + "/chat/" + chatName + "/subscribe")
}
