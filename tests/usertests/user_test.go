package usertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/samtech09/api-template/controllers"
	"github.com/samtech09/api-template/tests"
	"github.com/samtech09/api-template/viewmodels"

	apitest "github.com/samtech09/apitestengine"
)

var userTest *apitest.APITest
var userTestCases map[string]apitest.TestCase
var userSrv *httptest.Server

func inituserTest() error {
	// prepare test routes and server
	u := controllers.User{}
	r := u.SetRoutes()
	userSrv = httptest.NewServer(r)

	// post payload
	p := viewmodels.DbUser{}
	p.Name = "Krishna"
	p.ID = 2
	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("Failed marshling payload")
	}

	// prepare test cases
	userTestCases = make(map[string]apitest.TestCase)

	expected := `"Data":"2"`
	userTestCases["TestCreateuser"] = apitest.NewTestCase("TestCreateuser", "POST", "/", expected, apitest.MatchContains, "", bytes.NewBuffer(payload))

	userTestCases["TestListusers"] = apitest.NewTestCase("TestListusers", "GET", "/", `"Name\":\"Krishna\"`, apitest.MatchContains, "", nil)
	userTestCases["TestDeleteuser"] = apitest.NewTestCase("TestDeleteuser", "DELETE", "/2", `"Data":"1"`, apitest.MatchContains, "", nil)

	// initialize test engine
	userTest = apitest.NewAPITest(userSrv)
	return nil
}

func TestUserInit(t *testing.T) {
	go tests.InitTestConfig()
	// give it time to initialize
	time.Sleep(3 * time.Second)
	err := inituserTest()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateuser(t *testing.T) {
	ret := userTest.DoTest(userTestCases["TestCreateuser"])
	if ret.Err != nil {
		t.Error(ret.Err)
	}
}

func TestListusers(t *testing.T) {
	ret := userTest.DoTest(userTestCases["TestListusers"])
	if ret.Err != nil {
		t.Error(ret.Err)
	}
}

func TestDeleteuser(t *testing.T) {
	tcs := []apitest.TestCase{}
	tcs = append(tcs, userTestCases["TestDeleteuser"])

	tresult := userTest.DoTests(tcs)
	for _, ret := range tresult {
		if ret.Err != nil {
			fmt.Println(ret.ID)
			t.Error(ret.Err)
		} else {
			fmt.Println("PASS: ", ret.ID)
		}
	}
}

func TestUserCloser(t *testing.T) {
	userSrv.Close()
	tests.ResetTestConfig()
}
