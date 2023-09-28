export function initGamepad(sendData) {
    let rafID = null;

    window.addEventListener("gamepadconnected", _ => {
        if (!rafID) {
            pollGamepad();
        }
    });

    const pollGamepad = () => {
        const gamepads = navigator.getGamepads();
        for (const gamepad of gamepads) {
            if (!gamepad) {
                continue;
            }

            const buffer = new ArrayBuffer(32)

            const reportId = new Uint8Array(buffer, 0, 1)
            reportId[0] = 3;

            let buttons = 0
            gamepad.buttons.forEach((button, index) => {
                buttons = buttons | (button.pressed << index)
            });

            const buttonsData = new Uint32Array(buffer, 28, 1)
            buttonsData[0] = buttons

            const axes = new Float32Array(buffer, 4, 4)
            gamepad.axes.forEach((axe, index) => {
                axes[index] = axe
            });

            const triggers = new Float32Array(buffer, 20, 2)
            triggers[0] = gamepad.buttons[6].value
            triggers[1] = gamepad.buttons[7].value

            sendData(buffer)
        }
        rafID = window.requestAnimationFrame(pollGamepad);
    };
}
