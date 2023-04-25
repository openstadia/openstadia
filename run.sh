docker run -it --rm -p 8080:8080 -v ~/games:/games --device=/dev/uinput:/dev/uinput -v /sys/devices/virtual/input:/sys/devices/virtual/input openstadia
