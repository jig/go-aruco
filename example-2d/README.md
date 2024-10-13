This sample code shows the camera coordinates of an 6x6 Aruco marker #7 of size 38 mm that is vertically oriented but with some pitch (Y axis) tilt.

That is a typical situation of a 2D moving camera.

That is the result of static marker #7 on a Raspberry Pi 4 with 2GB of RAM:

```bash
go run ./example-2d/main.go
2024/10/13 18:17:00 Launching Python script...
2024/10/13 18:17:00 Starting to run the task...
2024/10/13 18:17:00 Waiting for samples...
2024/10/13 18:17:10 Marker 7:   Z=87.8cm  X=34.4cm  pose=-24°
2024/10/13 18:17:10 Marker 7:   Z=87.8cm  X=34.4cm  pose=-24°
2024/10/13 18:17:10 Marker 7:   Z=88.3cm  X=34.6cm  pose=-22°
2024/10/13 18:17:11 Marker 7:   Z=88.3cm  X=34.6cm  pose=-22°
2024/10/13 18:17:11 Marker 7:   Z=88.3cm  X=34.6cm  pose=-22°
2024/10/13 18:17:11 Marker 7:   Z=88.3cm  X=34.6cm  pose=-22°
2024/10/13 18:17:11 Marker 7:   Z=88.3cm  X=34.6cm  pose=-22°
2024/10/13 18:17:12 Marker 7:   Z=88.3cm  X=34.6cm  pose=-22°
2024/10/13 18:17:12 Marker 7:   Z=87.6cm  X=34.3cm  pose=-13°
^Csignal: interrupt
```

Notice it requires around 10 seconds to start on said computer. Refresh rate is set at 250 ms which is actually higher than the rate of frames processed for a Raspberry Pi 4 (around 1.5 frames per second). On a Raspberry Pi 5 the rate is ten times faster (around 15 frames per second).

To stop the Go program press <kbd>Ctrl</kbd> + <kbd>C</kbd>.
