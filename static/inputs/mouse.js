export function initMouse(videoEl, sendData) {
    videoEl.addEventListener('mousemove', (event) => {
        const rect = event.currentTarget.getBoundingClientRect();

        const x = event.clientX - rect.left;
        const y = event.clientY - rect.top;

        const width = rect.width
        const height = rect.height

        const buffer = new ArrayBuffer(12)

        const reportId = new Uint8Array(buffer, 0, 1)
        reportId[0] = 0

        const payload = new Float32Array(buffer, 4, 2);
        payload[0] = x / width
        payload[1] = y / height

        sendData(buffer)
    })

    videoEl.addEventListener('contextmenu', event => event.preventDefault());

    videoEl.addEventListener('mousedown', (event) => {
        const buffer = new ArrayBuffer(8)

        const reportId = new Uint8Array(buffer, 0, 1)
        reportId[0] = 6;

        const payload = new Uint8Array(buffer, 4, 1);
        payload[0] = event.button

        sendData(buffer)

        return false
    })

    videoEl.addEventListener('mouseup', (event) => {
        event.preventDefault()

        const buffer = new ArrayBuffer(8)

        const reportId = new Uint8Array(buffer, 0, 1)
        reportId[0] = 7;

        const payload = new Uint8Array(buffer, 4, 1);
        payload[0] = event.button

        sendData(buffer)

        return false
    })

    videoEl.addEventListener('wheel', (event) => {
        event.preventDefault()

        const buffer = new ArrayBuffer(12)

        const reportId = new Uint8Array(buffer, 0, 1)
        reportId[0] = 2;

        const data = new Int32Array(buffer, 4, 2);
        data[0] = event.deltaX
        data[1] = event.deltaY

        sendData(buffer)
    })
}
