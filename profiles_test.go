package goprofiles

import (
	"testing"
)

type MockFileReader struct {
	FileContents []byte
	Err          error
}

// ReadFile is a mock method that returns predefined file contents and error.
func (m *MockFileReader) ReadFile(filename string) ([]byte, error) {
	return m.FileContents, m.Err
}

func TestDefaultOptions(t *testing.T) {
	options := defaultOptions()
	if options.file != "profiles.yaml" {
		t.Errorf("Default file value is not 'profiles.yaml'")
	}
}

func TestWithFile(t *testing.T) {
	options := defaultOptions()
	withFile := WithFile("custom.yaml")
	withFile(&options)
	if options.file != "custom.yaml" {
		t.Errorf("WithFile did not set the 'file' option correctly")
	}
}

func TestWithProfile(t *testing.T) {
	options := defaultOptions()
	withProfile := WithProfile("profile1", "profile2")
	withProfile(&options)
	if len(options.profiles) != 2 {
		t.Errorf("WithProfile did not append profiles correctly")
	}
	if options.profiles[0] != "profile1" || options.profiles[1] != "profile2" {
		t.Errorf("WithProfile did not set profiles correctly")
	}
}

func TestGetString(t *testing.T) {
	p := &Profiles{
		values: map[interface{}]interface{}{
			"some_string": "some string",
		},
	}
	if p.GetString("some_string") != "some string" {
		t.Errorf("GetString did not return the correct value")
	}
}

func TestGetInt(t *testing.T) {
	p := &Profiles{
		values: map[interface{}]interface{}{
			"some_int": 1,
		},
	}
	if p.GetInt("some_int") != 1 {
		t.Errorf("GetInt did not return the correct value")
	}
}

func TestGetBool(t *testing.T) {
	p := &Profiles{
		values: map[interface{}]interface{}{
			"some_bool": true,
		},
	}
	if p.GetBool("some_bool") != true {
		t.Errorf("GetBool did not return the correct value")
	}
}

func TestGetFloat(t *testing.T) {
	p := &Profiles{
		values: map[interface{}]interface{}{
			"some_float": 1.1,
		},
	}
	if p.GetFloat("some_float") != 1.1 {
		t.Errorf("GetFloat did not return the correct value")
	}
}

func TestGetIntSlice(t *testing.T) {
	p := &Profiles{
		values: map[interface{}]interface{}{
			"some_int_slice": []interface{}{1, 2, 3},
		},
	}
	if p.GetIntSlice("some_int_slice")[0] != 1 {
		t.Errorf("GetIntSlice did not return the correct value")
	}
}

func TestGetStringSlice(t *testing.T) {
	p := &Profiles{
		values: map[interface{}]interface{}{
			"some_string_slice": []interface{}{"one", "two", "three"},
		},
	}
	if p.GetStringSlice("some_string_slice")[0] != "one" {
		t.Errorf("GetStringSlice did not return the correct value")
	}
}

func TestGetBoolSlice(t *testing.T) {
	p := &Profiles{
		values: map[interface{}]interface{}{
			"some_bool_slice": []interface{}{true, false, true},
		},
	}
	if p.GetBoolSlice("some_bool_slice")[0] != true {
		t.Errorf("GetBoolSlice did not return the correct value")
	}
}

func TestGetFloatSlice(t *testing.T) {
	p := &Profiles{
		values: map[interface{}]interface{}{
			"some_float_slice": []interface{}{1.1, 2.2, 3.3},
		},
	}
	if p.GetFloatSlice("some_float_slice")[0] != 1.1 {
		t.Errorf("GetFloatSlice did not return the correct value")
	}
}

func TestGetValues(t *testing.T) {

	var e error
	_, e = getValues("", "dev")
	if e.Error() != "no file specified" {
		t.Errorf("getValues did not return the correct error")
	}

	_, e = getValues("custom.json", "dev")

	if e.Error() != "unsupported file type" {
		t.Errorf("getValues did not return the correct error")
	}
}

func TestGetYamlValues(t *testing.T) {
	common, err := getYamlValues("testdata/profiles.yaml", "common")
	if err != nil {
		t.Errorf("getYamlValues returned an error" + err.Error())
	}
	if common["some_float"] != 1.2 {
		t.Errorf("getYamlValues did not return the correct value")
	}

	commonNested := common["nested"].(map[interface{}]interface{})
	if commonNested["some_int"] != 1 {
		t.Errorf("getYamlValues did not return the correct value")
	}

	dev, err := getYamlValues("testdata/profiles.yaml", "dev")
	if err != nil {
		t.Errorf("getYamlValues returned an error" + err.Error())
	}
	if dev["some_string"] != "some string" {
		t.Errorf("getYamlValues did not return the correct value")
	}

	prod, err := getYamlValues("testdata/profiles.yaml", "prod")
	if err != nil {
		t.Errorf("getYamlValues returned an error" + err.Error())
	}
	if prod["some_string"] != "some other string" {
		t.Errorf("getYamlValues did not return the correct value")
	}

	devProd, err := getYamlValues("testdata/profiles.yaml", "dev", "prod")
	if err.Error() != "💥 value conflict found for key: some_string" {
		t.Errorf("getYamlValues did not return value conflict error")
	}
	if devProd != nil {
		t.Errorf("getYamlValues did not return the correct value")
	}

	incorrectProfile, err := getYamlValues("testdata/profiles.yaml", "test")
	if err.Error() != "📝: profile not found" {
		t.Errorf("getYamlValues did not return profile not found error")
	}
	if incorrectProfile != nil {
		t.Errorf("getYamlValues did not return the correct value")
	}

	invalid, err := getYamlValues("testdata/invalid.yaml", "dev")
	if err == nil {
		t.Errorf("getYamlValues returned an error" + err.Error())
	}
	if invalid != nil {
		t.Errorf("getYamlValues did not return the correct value")
	}
}

func TestNew(t *testing.T) {
	p := New(WithProfile("common"), WithFile("testdata/profiles.yaml"))
	if p == nil {
		t.Errorf("New did not return a valid pointer")
	}

	if p.GetInt("nested.some_int") != 1 {
		t.Errorf("New did not return the correct value")
	}
}

func TestNewFileNotFoundPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("get did not panic")
		}
	}()
	New()
}

func TestNewUnsupportedFileTypePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("get did not panic")
		}
	}()
	New(WithFile("custom.json"))
}

func TestNewProfileNotFoundPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("get did not panic")
		}
	}()
	New(WithProfile("test"), WithFile("testdata/profiles.yaml"))
}

func TestNewValueConflictPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("get did not panic")
		}
	}()
	New(WithProfile("dev", "prod"), WithFile("testdata/profiles.yaml"))
}

func TestGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("get did not panic")
		}
	}()

	p := New(WithProfile("dev"), WithFile("testdata/profiles.yaml"))

	p.get("some_float")
}

func TestGetPanic2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("get did not panic")
		}
	}()

	p := New(WithProfile("dev"), WithFile("testdata/profiles.yaml"))

	p.GetString("some_float")
}
