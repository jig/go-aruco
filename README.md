# go-aruco

This Go (Golang) package `aruco` that provides a JSON stream through stdout of the coordinates of the Aruco markers on view.

The package uses Python. Tested on a Raspberry Pi 4 with Camera Module 3 WIDE infrared with Raspberry Pi OS "Bookworm".

See the sample code on [example-2d](./example-2d/) on how to use it.

# Install

_Installation process pending to be reviewed_

You need Go compiler (tested with Go compiler v1.23.2) an the Python interpreter (tested on Python v3.13) and the Python OpenCV dependencies (tested with v4.10).

To install the Python dependencies use the Raspberry OS package manager:

```bash
sudo apt update && sudo apt upgrade -y && sudo apt install -y \
 	openocd \
    opencv-data \
   	python3-opencv \
    python3-scipy
```

# Python installation verification

Check that the script runs alone (without Go) by executing:

```bash
python ./markers.py
```

It requires a few seconds to load all dependencies.

Note that every line is suffixed with few kB of spaces (used for the Python to Go communication). Look for JSON data like this:

```json
[{"id": 9, "x": -0.13673050917742155, "y": 0.06945671561055802, "z": 0.7309941550510376, "roll-x": -172.7556401203728, "pitch-y": 6.610564919377947, "yaw-z": 178.17154128756582}, {"id": 6, "x": 0.34398351230356095, "y": -0.020529717758146982, "z": 0.8951305793554741, "roll-x": -177.8320267892627, "pitch-y": -26.86036068104307, "yaw-z": 177.5710900851925}, {"id": 8, "x": -0.13820595929256171, "y": -0.03209965216632465, "z": 0.7309398765039592, "roll-x": -171.51842160897496, "pitch-y": 2.0533179575661684, "yaw-z": 176.7455522076878}, {"id": 7, "x": 0.34751714240210885, "y": 0.07921186314686002, "z": 0.8838179732333362, "roll-x": -173.4195957396569, "pitch-y": -14.31427564710096, "yaw-z": 178.3094912498296}]
```

## Camera hardware verification

Check that your camera is working with following command on the Raspberry Pi Console:

```bash
rpicam-hello
```

# Camera calibration

Script is already calibrated for a Raspberry Pi Camera Module 3 WIDE. You might need to change it for other cameras by changing the matrices on constants `camera_matrix` and `distortion_coefficients` on [markers.py](./markers.py) Python script. You can use [github.com/jig/charuco-calibration](https://github.com/jig/charuco-calibration) for that.
