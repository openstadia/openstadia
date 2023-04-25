export DISPLAY=:99
Xvfb :99 -screen 0 640x480x24 &
sleep 5
xrandr --query
sleep 5
#nohup startxfce4 &
#./openstadia &
