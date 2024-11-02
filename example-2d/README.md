This sample code shows the camera coordinates of the 6x6 Aruco markers in view of size 38 mm that are vertically oriented but with some pitch (Y axis) tilt.

That is a typical situation of a 2D moving camera.

That is the result of static marker #7 on a Raspberry Pi 4 with 2GB of RAM:

```bash
go run example-2d/main.go
[0:40:19.143080485] [8436]  INFO Camera camera_manager.cpp:325 libcamera v0.3.2+27-7330f29b
[0:40:19.332158538] [8464]  WARN RPiSdn sdn.cpp:40 Using legacy SDN tuning - please consider moving SDN inside rpi.denoise
[0:40:19.337509759] [8464]  INFO RPI vc4.cpp:447 Registered camera /base/soc/i2c0mux/i2c@1/imx708@1a to Unicam device /dev/media3 and ISP device /dev/media0
[0:40:19.337635443] [8464]  INFO RPI pipeline_base.cpp:1126 Using configuration file '/usr/share/libcamera/pipeline/rpi/vc4/rpi_apps.yaml'
[0:40:19.367103990] [8436]  INFO Camera camera.cpp:1197 configuring streams: (0) 4608x2592-XBGR8888 (1) 4608x2592-YUV420 (2) 2304x1296-SBGGR10_CSI2P
[0:40:19.371836349] [8464]  INFO RPI vc4.cpp:622 Sensor: /base/soc/i2c0mux/i2c@1/imx708@1a - Selected sensor format: 2304x1296-SBGGR10_1X10 - Selected unicam format: 2304x1296-pBAA
2024/11/02 20:15:03 Marker 6:   Z=99.3cm  X=-4.6cm  pose=5°
2024/11/02 20:15:03 Marker 6:   Z=101.8cm  X=-5.1cm  pose=-9°
2024/11/02 20:15:04 Marker 6:   Z=99.5cm  X=-4.6cm  pose=8°
2024/11/02 20:15:04 Marker 6:   Z=102.6cm  X=-5.2cm  pose=-4°
2024/11/02 20:15:05 Marker 6:   Z=99.6cm  X=-4.6cm  pose=11°
2024/11/02 20:15:05 Marker 6:   Z=102.4cm  X=-5.2cm  pose=-0°
^Csignal: interrupt
```

Notice it requires around 10 seconds to start on said computer. Refresh rate is around 2 frames per second on a Raspberry Pi 4.

To stop the Go program press <kbd>Ctrl</kbd> + <kbd>C</kbd>.
