export function initChannel() {
    let sendChannel;

    return {
        sendData(data) {
            if (sendChannel && sendChannel.readyState === "open") {
                sendChannel.send(data)
            }
        },
        setChannel(channel) {
            sendChannel = channel
        }
    }
}