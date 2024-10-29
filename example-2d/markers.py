#!/usr/bin/python3

import cv2
from picamera2 import MappedArray, Picamera2, Preview
import json
import numpy as np
from scipy.spatial.transform import Rotation as R
import math

this_aruco_dictionary = cv2.aruco.Dictionary_get(cv2.aruco.DICT_6X6_50)
this_aruco_parameters = cv2.aruco.DetectorParameters_create()

# Side length of the ArUco marker in meters
aruco_marker_side_length = 0.038
# aruco_marker_side_length = 0.0146

# Calibration parameters yaml file
camera_calibration_parameters_filename = 'calibration_chessboard.yaml'

def euler_from_quaternion(x, y, z, w):
  """
  Convert a quaternion into euler angles (roll, pitch, yaw)
  roll is rotation around x in radians (counterclockwise)
  pitch is rotation around y in radians (counterclockwise)
  yaw is rotation around z in radians (counterclockwise)
  """
  t0 = +2.0 * (w * x + y * z)
  t1 = +1.0 - 2.0 * (x * x + y * y)
  roll_x = math.atan2(t0, t1)

  t2 = +2.0 * (w * y - z * x)
  t2 = +1.0 if t2 > +1.0 else t2
  t2 = -1.0 if t2 < -1.0 else t2
  pitch_y = math.asin(t2)

  t3 = +2.0 * (w * z + x * y)
  t4 = +1.0 - 2.0 * (y * y + z * z)
  yaw_z = math.atan2(t3, t4)

  return roll_x, pitch_y, yaw_z # in radians

picam2 = Picamera2()
picam2.start_preview(Preview.DRM, x=0, y=0, width=1024, height=600)


# 4608 × 2592
config = picam2.create_preview_configuration(
  main={"size": (4608, 2592)},
  lores={"size": (4608, 2592), "format": "YUV420"}
)

picam2.configure(config)

(w0, h0) = picam2.stream_configuration("main")["size"]
(w1, h1) = picam2.stream_configuration("lores")["size"]
s1 = picam2.stream_configuration("lores")["stride"]

picam2.start()

# Load the camera parameters from the saved file
# cv_file = cv2.FileStorage(
#   camera_calibration_parameters_filename,
#   cv2.FILE_STORAGE_READ
# )
# camera_matrix = cv_file.getNode('K').mat()
# distortion_coefficients = cv_file.getNode('D').mat()
# cv_file.release()

# This are the calibration values for a Raspberry Pi Camera Module 3
camera_matrix = np.array([
    [2064.2660695692675, 0., 2305.194723859388],
    [0., 2062.8376670432476, 1317.6585179294466],
    [0., 0., 1.]
  ])

distortion_coefficients = np.array([ -0.046250204900546286, 0.15249055104559342, 0.001339977338210037, -0.0011487231459689135, -0.1079298904091414 ])

# padding to avoid stdout Linux buffering to delay the delivery
# see http://www.pixelbeat.org/programming/stdio_buffering/ for an explanation an potential improvements on this code
padding = ' ' * 13642
# padding = ''
c = 0
while True:
    buffer = picam2.capture_buffer("lores")
    frame = buffer[:s1 * h1].reshape((h1, s1))
    c=c+1
    (corners, marker_ids, rejected) = cv2.aruco.detectMarkers(
       frame,
       this_aruco_dictionary,
       parameters=this_aruco_parameters
    )
    # Check that at least one ArUco marker was detected
    if marker_ids is not None:

      # Draw a square around detected markers in the video frame
      cv2.aruco.drawDetectedMarkers(frame, corners, marker_ids)

      # Get the rotation and translation vectors
      rvecs, tvecs, obj_points = cv2.aruco.estimatePoseSingleMarkers(
        corners,
        aruco_marker_side_length,
        camera_matrix,
        distortion_coefficients)

      # Print the pose for the ArUco marker
      # The pose of the marker is with respect to the camera lens frame.
      # Imagine you are looking through the camera viewfinder,
      # the camera lens frame's:
      # x-axis points to the right
      # y-axis points straight down towards your toes
      # z-axis points straight ahead away from your eye, out of the camera

      if len(marker_ids) > 0:
        data = []
        for i, marker_id in enumerate(marker_ids):
          # Store the translation (i.e. position) information
          transform_translation_x = tvecs[i][0][0]
          transform_translation_y = tvecs[i][0][1]
          transform_translation_z = tvecs[i][0][2]

          # Store the rotation information
          rotation_matrix = np.eye(4)
          rotation_matrix[0:3, 0:3] = cv2.Rodrigues(np.array(rvecs[i][0]))[0]
          r = R.from_matrix(rotation_matrix[0:3, 0:3])
          quat = r.as_quat()

          # Quaternion format
          transform_rotation_x = quat[0]
          transform_rotation_y = quat[1]
          transform_rotation_z = quat[2]
          transform_rotation_w = quat[3]

          # Euler angle format in radians
          roll_x, pitch_y, yaw_z = euler_from_quaternion(transform_rotation_x,
                                                        transform_rotation_y,
                                                        transform_rotation_z,
                                                        transform_rotation_w)
          roll_x = math.degrees(roll_x)
          pitch_y = math.degrees(pitch_y)
          yaw_z = math.degrees(yaw_z)

          data.append({
            "id": int(marker_id[0]),
            "x": float(transform_translation_x),
            "y": float(transform_translation_y),
            "z": float(transform_translation_z),
            "roll-x": roll_x,
            "pitch-y": pitch_y,
            "yaw-z": yaw_z
          })
        print(json.dumps(data), padding)
    else:
      print("[]", padding)
