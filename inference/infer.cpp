#include "infer.h"

#include <dlib/dnn.h>
#include <dlib/image_io.h>
#include <dlib/image_processing/frontal_face_detector.h>

Inferer::Inferer(char* model_file, char* sp_file) {
#ifndef NDEBUG
  std::cout << "[c++] Inferer(\"" << model_file << "\", \"" << sp_file << "\")"
            << std::endl;
#endif

  this->detector = dlib::get_frontal_face_detector();

  dlib::deserialize(sp_file) >> this->sp;
  dlib::deserialize(model_file) >> this->net;
}

Inferer::~Inferer() {
#ifndef NDEBUG
  std::cout << "[c++] ~Inferer()" << std::endl;
#endif
  if (this->result) {
    this->destroy_result();
  }
}

Result* Inferer::get_result(char* img_path) {
#ifndef NDEBUG
  std::cout << "[c++] Inferer->get_encoding(\"" << img_path << "\", \""
            << img_path << "\")" << std::endl;
#endif

  dlib::matrix<dlib::rgb_pixel> img;
  dlib::load_image(img, img_path);

  std::vector<dlib::matrix<dlib::rgb_pixel>> faces;
  std::vector<dlib::rectangle> bboxes_vec;

  for (auto face : this->detector(img)) {
    bboxes_vec.push_back(face);

    auto shape = this->sp(img, face);
    dlib::matrix<dlib::rgb_pixel> face_chip;
    dlib::extract_image_chip(img, dlib::get_face_chip_details(shape, 150, 0.25),
                             face_chip);
    faces.push_back(std::move(face_chip));
  }

  std::vector<dlib::matrix<float, 0, 1>> face_descriptors = this->net(faces);

  int nfaces = faces.size();
  std::cout << "[c++] " << nfaces << " faces found in image " << img_path
            << std::endl;

  float* encodings = new float[128 * nfaces];
  int encoding_idx = 0;

  float* bboxes = new float[4 * nfaces];

  for (int i = 0; i < nfaces; i++) {
    dlib::matrix<dlib::rgb_pixel> face = faces[i];
    dlib::matrix<float, 0, 1> fd = face_descriptors[i];

    for (auto j = fd.begin(); j != fd.end(); j++) {
      encodings[encoding_idx] = *j;
      encoding_idx++;
    }

    bboxes[i * 4] = bboxes_vec[i].top();
    bboxes[i * 4 + 1] = bboxes_vec[i].left();
    bboxes[i * 4 + 2] = bboxes_vec[i].width();
    bboxes[i * 4 + 3] = bboxes_vec[i].height();
  }

  if (this->result) {
    this->destroy_result();
  }
  this->result = (Result*)malloc(sizeof(Result));
  this->result->nfaces = nfaces;
  this->result->encodings = encodings;
  this->result->bboxes = bboxes;

  return this->result;
}

void Inferer::destroy_result() {
  if (!this->result) {
    return;
  }
  if (this->result->encodings) {
    delete[] this->result->encodings;
  }
  if (this->result->bboxes) {
    delete[] this->result->bboxes;
  }

  this->result->nfaces = 0;
  this->result->encodings = nullptr;
  this->result->bboxes = nullptr;
}