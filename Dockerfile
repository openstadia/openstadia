FROM ubuntu:jammy

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update
RUN apt-get install -y keyboard-configuration xfce4 xfce4-goodies xvfb
RUN apt-get install -y libx11-dev libxext-dev libvpx-dev
RUN apt-get install -y xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
RUN apt-get install -y libx11-dev xorg-dev libxtst-dev
RUN apt-get install -y libgl1-mesa-dev libsdl2-dev libvulkan-dev

ENV DISPLAY :99

WORKDIR /app
ADD . .
RUN chmod a+x display.sh

CMD { bash display.sh; }
