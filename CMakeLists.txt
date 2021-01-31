
cmake_minimum_required(VERSION 3.14)

project(face-recog-backend)

# compiler
SET(CMAKE_CXX_FLAGS_DEBUG "-O3")

# setup output folders
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)

# functions to download and extract model files
function(download_file url filename hash_type hash)

if(NOT EXISTS ${filename})
  file(DOWNLOAD ${url} ${filename}
       TIMEOUT 60  # seconds
       EXPECTED_HASH ${hash_type}=${hash}
       TLS_VERIFY ON)
endif()

endfunction(download_file)

function(get_model_file url filename hash_type hash)

download_file(${url} ${filename}.bz2 ${hash_type} ${hash})
execute_process(
	COMMAND bzip2 -d 
	INPUT_FILE ${filename}.bz2
	OUTPUT_FILE ${filename}
)

endfunction(get_model_file)

file(MAKE_DIRECTORY model)

# get serialized model files
get_model_file(
	http://dlib.net/files/shape_predictor_5_face_landmarks.dat.bz2 
	model/shape_predictor.dat
	SHA1
	cd47b9dd2c67052e8695f693b50d3e7c828290f6
)

get_model_file(
	http://dlib.net/files/dlib_face_recognition_resnet_model_v1.dat.bz2
	model/model.dat
	SHA1
	c93227f4b3fbc60cf3b32a565ec22ed37217ad03
)

# fetch dlib
include(FetchContent)
FetchContent_Declare(dlib
	GIT_REPOSITORY https://github.com/davisking/dlib.git
	GIT_TAG        v19.18
)
FetchContent_MakeAvailable(dlib)

# compile cpp to static library and link to dlib
add_library(infer SHARED inference/infer-bridge.cpp inference/infer.cpp)
target_link_libraries(infer dlib::dlib)