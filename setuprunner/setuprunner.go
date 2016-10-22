package main

import (
  "flag"
  "log"

  "github.com/codegp/test-utils"
  "github.com/codegp/cloud-persister"
	"github.com/codegp/kube-client"
)

var (
  testUtils *testutils.TestUtils
	forceBuild   *bool
	runGame *bool
)

func init() {
	cp, err := cloudpersister.NewCloudPersister()
	fatalOnErr(err)
	kclient, err := kubeclient.NewClient()
	fatalOnErr(err)
  testUtils = testutils.NewTestUtils(cp, kclient)

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
