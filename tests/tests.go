package tests

import (
	"sync"

	"github.com/samtech09/api-template/internal/initializer"
	"github.com/samtech09/api-template/internal/logger"

	g "github.com/samtech09/api-template/global"
)

var initialized bool
var wg sync.WaitGroup

func InitTestConfig() {
	wg.Add(1)

	if initialized {
		return // already initialized
	}

	//initialize config
	prod := false
	initializer.Initconfig(&prod, "../../")

	//-------
	//------------
	//
	//
	// set flag to mark test environment, so JWT tokens will be bypassed
	//
	g.TestEnv = true
	//
	//---------------------

	//initialize logs
	g.Logger = logger.NewConsole(true, true, true, true)

	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	initializer.InitServices("../../")

	go func() {
		initialized = true
	}()

	// wailt till all test finished
	wg.Wait()

	initializer.AppCleanup()
	initialized = false
	g.TestEnv = false
}

func ResetTestConfig() {
	wg.Done()
}
