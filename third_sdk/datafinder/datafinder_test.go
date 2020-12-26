package datafinder

import (
	"github.com/hero1s/golib/log"
	"testing"
)

func TestDataFinder(t *testing.T) {

	RemixAppKeyTest := "c17182141be4169dd63abde2e221b768"
	RemixAppNameTest := "remix_test"

	EventNameCrystalFinish := "mobile_crystal_game_click"

	dataFinder := NewDataFinder(RemixAppKeyTest)
	dataFinder.AddUser("272687711876308003")
	dataFinder.AddHeader(&HeaderStruct{
		AppName: RemixAppNameTest,
	})
	dataFinder.AddEvent(&EventStruct{
		EventName: EventNameCrystalFinish,
		Params:    `{"StepSec":20,"StepNum":10,"BackUpGap":300}`,
	})
	resp, err := dataFinder.Request()
	if err != nil {
		log.Error("request fail", err)
		return
	}
	log.Info("request ok", *resp)
}
