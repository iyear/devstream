package harbor

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	. "github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	// get harbor running status
	retState, err := GetDynamicState(&opts)
	if err != nil {
		return nil, err
	}
	retMap := retState.ToStringInterfaceMap()
	log.Debugf("Return map: %v.", retMap)

	return retMap, nil
}
