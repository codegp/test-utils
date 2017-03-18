package main

import (
	"flag"
	"log"
	"time"

	"github.com/codegp/cloud-persister"
	"github.com/codegp/env"
	job "github.com/codegp/job-client"
	"github.com/codegp/test-utils"
)

var (
	cp         *cloudpersister.CloudPersister
	testUtils  *testutils.TestUtils
	forceBuild *bool
	runGame    *bool
)

func init() {
	var err error
	for {
		cp, err = cloudpersister.NewCloudPersister()
		if err == nil || !env.IsLocal() {
			break
		}

		log.Printf("Failed to create cloud persister. Confirm DSEmulator is running.\nError: %v", err)
		time.Sleep(time.Millisecond * time.Duration(1000))
	}
	fatalOnErr(err)
	jclient, err := job.GetJobClient()
	fatalOnErr(err)
	testUtils = testutils.NewTestUtils(cp, jclient)

	forceBuild = flag.Bool("fb", false, "force a new build even if one already exists, bool")
	runGame = flag.Bool("r", false, "run test game, bool")
	flag.Parse()
}

func main() {
	if *forceBuild {
		fatalOnErr(testUtils.BuildTestGametype(true, true))
	}

	if *runGame {
		fatalOnErr(testUtils.RunTestGame(true))
	} else if !*forceBuild {
		fatalOnErr(testUtils.BuildTestGametype(true, false))
	}
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}
