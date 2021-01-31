#pragma once

#include "infer-bridge.h"
#include "model.h"

class Inferer {
public:
  Inferer(char* model_file, char* sp_file);
  ~Inferer();
  Result* get_result(char* img_path);

private:
  void destroy_result();
  dlib::frontal_face_detector detector;
  dlib::shape_predictor sp;
  anet_type net;
  Result* result;
};