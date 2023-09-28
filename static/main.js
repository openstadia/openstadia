import {getApps} from "./apis/apps.js";
import {sendOffer} from "./apis/offer.js";

import {initGamepad} from "./inputs/gamepad.js";
import {initKeyboard} from "./inputs/keyboard.js";
import {initMouse} from "./inputs/mouse.js";

import {initRtc} from "./rtc/rtc.js";
import {initChannel} from "./rtc/channel.js";

const videoEl = document.querySelector('#video');
const audioEl = document.querySelector('#audio');
const fullscreen = document.querySelector('#fullscreen');
const startEl = document.querySelector('#start');
const videoWrapper = document.querySelector('#videoWrapper');

const {sendData, setChannel} = initChannel()

initGamepad(sendData)
initKeyboard(sendData)
initMouse(videoEl, sendData)

function enterFullScreen() {
    if (!document.fullscreenElement) {
        videoWrapper.requestFullscreen();
    }
}

fullscreen.addEventListener('click', enterFullScreen)

const log = msg => {
    document.getElementById('logs').innerHTML += msg + '<br>'
}

const startHttpSession = async () => {
    const {sendChannel, setRemoteDescription, getLocalDescription} = initRtc(videoEl, audioEl, log)
    setChannel(sendChannel)

    const formElements = document.forms['connect'].elements
    const codecType = formElements.type.value
    const bitrate = formElements.bitrate.valueAsNumber
    const app = formElements.app.value

    const localDescription = await getLocalDescription()

    const body = {
        sdp: localDescription.sdp,
        type: localDescription.type,
        app: {name: app},
        codec: {type: codecType, bitrate: bitrate}
    }

    console.log(body)

    const answer = await sendOffer(body)

    try {
        await setRemoteDescription(answer)
    } catch (e) {
        alert(e)
    }
}

startEl.addEventListener('click', startHttpSession)

const loadApps = async () => {
    const apps = await getApps()
    const appsEl = document.forms['connect'].elements.app

    apps.map((app) => {
        const option = new Option(app, app)
        appsEl.options.add(option)
    })
}
loadApps()
