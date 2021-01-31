package inferer

// #cgo LDFLAGS: -L ../build/lib -l infer
// #include "../inference/infer-bridge.h"
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

type Result struct {
	Nfaces    int
	Encodings []float32
	BBoxes    []float32
}

type Inferer struct {
	ptr unsafe.Pointer
}

func NewInferer(model_file string, sp_file string) Inferer {
	var inferer Inferer
	inferer.ptr = C.LIB_NewInferer(C.CString(model_file), C.CString(sp_file))
	return inferer
}

func (inferer Inferer) Free() {
	C.LIB_DestroyInferer(inferer.ptr)
}

func (inferer Inferer) GetResults(img_file string) Result {
	res := *C.LIB_InfererGetResult(inferer.ptr, C.CString(img_file))
	nfaces := int(res.nfaces)
	encodings := cFloatArrToSlice(uintptr(unsafe.Pointer(res.encodings)), 128*nfaces)
	bboxes := cFloatArrToSlice(uintptr(unsafe.Pointer(res.bboxes)), 4*nfaces)

	return Result{
		Nfaces:    nfaces,
		Encodings: encodings,
		BBoxes:    bboxes,
	}
}

func cFloatArrToSlice(arr uintptr, len int) []float32 {
	var res []float32
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&res)))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = arr

	return res
}
