package infer

import (
	"testing"
)

func TestCreateValidPath(t *testing.T) {
	model := "../../build/model/model.dat"
	sp := "../../build/model/shape_predictor.dat"

	inferer, err := NewInferer(model, sp)
	if err != nil {
		t.Errorf("Error while creating inferer: %s for model %s and shape predictor %s", err, model, sp)
	}
	testFreeInferer(t, inferer)
}

func TestCreateInvalidPath(t *testing.T) {
	path1 := "( ͡° ͜ʖ ͡°)"
	path2 := "this is not a real filepath lol 💩"

	_, err := NewInferer(path1, path2)
	if err == nil {
		t.Errorf("Did not give error for invalid filepaths %s and %s", path1, path2)
	}
}

func TestFree(t *testing.T) {
	inferer := testCreateInferer(t)
	testFreeInferer(t, inferer)
}

func TestDoubleFree(t *testing.T) {
	inferer := testCreateInferer(t)

	testFreeInferer(t, inferer)

	// TODO: this segfaults ¯\_(ツ)_/¯
	// if err := inferer.Free(); err == nil {
	// 	t.Errorf("Did not give error when double freeing inference")
	// }
}

func TestResultsInvalidPath(t *testing.T) {
	model := "../../build/model/model.dat"
	sp := "../../build/model/shape_predictor.dat"

	inferer, err := NewInferer(model, sp)
	if err != nil {
		t.Errorf("Error while creating inferer: %s", err)
	}

	if _, err := inferer.GetResults("thispathdoesnotexist.xyz"); err == nil {
		t.Errorf("Did not give error when trying to get results on nonexistant image path %s", err)
	}

	testFreeInferer(t, inferer)
}

func TestOneFace(t *testing.T) {
	inferer := testCreateInferer(t)
	testResultsForImage(t, inferer, "../../test-images/obama.jpg", 1)
	testFreeInferer(t, inferer)
}

func TestMultipleFaces(t *testing.T) {
	inferer := testCreateInferer(t)
	testResultsForImage(t, inferer, "../../test-images/bald_guys.jpg", 24)
	testFreeInferer(t, inferer)
}

func TestMultipleResults(t *testing.T) {
	inferer := testCreateInferer(t)
	testResultsForImage(t, inferer, "../../test-images/obama.jpg", 1)
	testResultsForImage(t, inferer, "../../test-images/bald_guys.jpg", 24)
	testFreeInferer(t, inferer)
}

// utilities
func testCreateInferer(t *testing.T) Inferer {
	model := "../../build/model/model.dat"
	sp := "../../build/model/shape_predictor.dat"

	inferer, err := NewInferer(model, sp)
	if err != nil {
		t.Errorf("Error while creating inferer: %s", err)
	}
	return inferer
}

func testResultsForImage(t *testing.T, inferer Inferer, img_path string, nfaces int) {
	res, err := inferer.GetResults(img_path)
	if err != nil {
		t.Errorf("Error while trying to find results: %s", err)
	}
	if res.Nfaces != nfaces {
		t.Errorf("Expected to find %d face, found %d instead", nfaces, res.Nfaces)
	}
	if len(res.Encodings) != 128*nfaces {
		t.Errorf("Expected length of encoding to be %d, found %d instead", 128*nfaces, len(res.Encodings))
	}
	if len(res.BBoxes) != 4*nfaces {
		t.Errorf("Expected length of bboxes to be %d, found %d instead", 4*nfaces, len(res.BBoxes))
	}
}

func testFreeInferer(t *testing.T, inferer Inferer) {
	if err := inferer.Free(); err != nil {
		t.Errorf("Failed to free inferer: %s", err)
	}
}

// inferer, _ := infer.NewInferer("build/model/model.dat", "build/model/shape_predictor.dat")
// res := inferer.GetResults("test-images/obama.jpg")
// fmt.Printf("%+v\n", res)
// inferer.Free()
