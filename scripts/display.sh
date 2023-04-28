export DISPLAY=:99
Xvfb :99 -screen 0 640x480x24 &
#export LIBVA_DRIVER_NAME=iHD
sleep 5
xrandr --query
sleep 5
#./openstadia &
