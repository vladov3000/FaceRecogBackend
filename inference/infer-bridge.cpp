#include <iostream>

#include "infer-bridge.h"
#include "infer.h"

void* LIB_NewInferer(char* model_file, char* sp_file) {
  return new Inferer(model_file, sp_file);
}

Inferer* AsInferer(void* inferer) {
  return reinterpret_cast<Inferer*>(inferer);
}

void LIB_DestroyInferer(void* inferer) { AsInferer(inferer)->~Inferer(); }

Result* LIB_InfererGetResult(void* inferer, char* img_path) {
  return AsInferer(inferer)->get_result(img_path);
}