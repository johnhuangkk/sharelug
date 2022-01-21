package model

import "api/services/util/log"

type Dome struct {}

func (*Dome) Test() (bool, error){

	log.Debug("GGGGGGGG")

	return true, nil
}
