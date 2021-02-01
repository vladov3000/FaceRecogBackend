package infer

// #cgo LDFLAGS: -L ../../build/lib -l infer
// #include "../../inference/infer-bridge.h"
import "C"
import (
	"errors"
	"os"
	"reflect"
	"unsafe"
)

type Result struct {
	Nfaces    int       `json:"nfaces"`
	Encodings []float32 `json:"encodings"`
	BBoxes    []float32 `json:"bboxes"`
}

type Inferer struct {
	ptr   unsafe.Pointer
	freed bool
}

func NewInferer(model_file string, sp_file string) (Inferer, error) {
	var inferer Inferer

	if _, err := os.Stat(model_file); err != nil {
		return inferer, err
	}
	if _, err := os.Stat(sp_file); err != nil {
		return inferer, err
	}

	inferer = Inferer{
		ptr:   C.LIB_NewInferer(C.CString(model_file), C.CString(sp_file)),
		freed: false,
	}
	return inferer, nil
}

func (inferer *Inferer) Free() error {
	if inferer.freed {
		return errors.New("Inferer has already been destroyed.")
	}
	C.LIB_DestroyInferer(inferer.ptr)

	inferer.freed = true
	return nil
}

func (inferer Inferer) GetResults(img_file string) (Result, error) {
	var ret Result

	if _, err := os.Stat(img_file); err != nil {
		return ret, err
	}

	res := *C.LIB_InfererGetResult(inferer.ptr, C.CString(img_file))
	nfaces := int(res.nfaces)
	encodings := cFloatArrToSlice(uintptr(unsafe.Pointer(res.encodings)), 128*nfaces)
	bboxes := cFloatArrToSlice(uintptr(unsafe.Pointer(res.bboxes)), 4*nfaces)

	ret = Result{
		Nfaces:    nfaces,
		Encodings: encodings,
		BBoxes:    bboxes,
	}

	return ret, nil
}

func cFloatArrToSlice(arr uintptr, len int) []float32 {
	var res []float32

	// not reccomended, we want our C and Go arrays seperate
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&res)))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = arr

	// deep copy
	// for i := 0; i < len; i++ {
	// 	res = append(res, *(*float32)(unsafe.Pointer(arr + uintptr(i)*4)))
	// }

	return res
}
