export function initKeyboard(sendData) {
    window.addEventListener('keydown', event => {
        event.preventDefault()

        const key = event.code
        const encoder = new TextEncoder()
        const payload = encoder.encode(key)

        const report = new Uint8Array(4 + payload.length)
        report[0] = 4;
        report.set(payload, 4)

        sendData(report)
    })

    window.addEventListener('keyup', event => {
        event.preventDefault()

        const key = event.code
        const encoder = new TextEncoder()
        const payload = encoder.encode(key)

        const report = new Uint8Array(4 + payload.length)
        report[0] = 5;
        report.set(payload, 4)

        sendData(report)
    })
}
