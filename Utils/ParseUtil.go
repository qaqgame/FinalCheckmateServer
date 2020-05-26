package Utils

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func ParsePositionToV3i(position string) *DataFormat.V3I {
	target := new(DataFormat.V3I)

	pos := position[1:len(position)-1]
	poss := strings.Split(pos,",")
	x,err := strconv.Atoi(poss[0])
	if !checkAtoiErr(err) {
		return nil
	}
	y,err := strconv.Atoi(poss[1])
	if !checkAtoiErr(err) {
		return nil
	}
	z,err := strconv.Atoi(poss[2])
	if !checkAtoiErr(err) {
		return nil
	}
	target.X = int32(x)
	target.Y = int32(y)
	target.Z = int32(z)
	return target
}

func checkAtoiErr(err error) bool {
	if err != nil {
		logrus.Warn("Position format not correct")
		return false
	}
	return true
}