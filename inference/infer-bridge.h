#pragma once

#ifdef __cplusplus
extern "C" {
#endif

typedef struct Result {
  int nfaces;
  float* encodings;
  float* bboxes;
} Result;

void* LIB_NewInferer(char* model_file, char* sp_file);
void LIB_DestroyInferer(void* inferer);
Result* LIB_InfererGetResult(void* inferer, char* img_path);

#ifdef __cplusplus
} // extern "C"
#endif