export function initRtc(videoEl, audioEl, log) {
    const pc = new RTCPeerConnection({
        iceServers: [{
            urls: 'stun:stun.l.google.com:19302'
        }],
    })

    pc.addEventListener('track', (event) => {
        console.log(event)
        if (event.track.kind === "video") {
            videoEl.srcObject = event.streams[0]
        }

        if (event.track.kind === "audio") {
            audioEl.srcObject = event.streams[0]
            audioEl.autoplay = true
            audioEl.controls = true
        }

    })

    pc.addEventListener('iceconnectionstatechange', _ => log(pc.iceConnectionState))

    pc.addTransceiver('video', {direction: 'recvonly'});
    pc.addTransceiver('audio', {direction: 'recvonly'});

    const sendChannel = pc.createDataChannel('sendDataChannel');
    sendChannel.onopen = () => console.log('sendChannel.onopen');
    sendChannel.onmessage = () => console.log('sendChannel.onmessage');
    sendChannel.onclose = () => console.log('sendChannel.onclose');

    const getLocalDescription = async () => {
        const offer = await pc.createOffer()
        await pc.setLocalDescription(offer)
        return pc.localDescription
    }

    const setRemoteDescription = async (answer) => {
        await pc.setRemoteDescription(new RTCSessionDescription(answer))
    }

    return {
        sendChannel,
        getLocalDescription,
        setRemoteDescription
    }
}