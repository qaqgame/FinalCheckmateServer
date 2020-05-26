package Roles

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func init() {
	dir := `Roles/`
	filinfo, err := ioutil.ReadDir(dir)
	if err != nil {
		logrus.Warn("unable to read this dir: ",dir)
	}

	for _, fi := range filinfo {
		if !fi.IsDir()&&path.Ext(fi.Name())==".json" {
			f, err := os.Open(dir+fi.Name())
			if err != nil {
				f.Close()
				logrus.Warn("unable to open file: ", dir+fi.Name())
				continue
			}
			cnt,err := ioutil.ReadAll(f)
			if err != nil {
				f.Close()
				logrus.Warn("unable to read file: ", dir+fi.Name())
				continue
			}

			roleData := new(DataFormat.RoleData)
			err = json.Unmarshal(cnt, roleData)
			if err != nil {
				f.Close()
				logrus.Warn("unable to marshal json: ",fi.Name())
				continue
			}

			DataFormat.RolesMap[strings.Split(fi.Name(),".")[0]] = roleData
		}
	}
}