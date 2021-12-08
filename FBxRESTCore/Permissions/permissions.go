package permissions

import (
	"encoding/xml"
	_apperrors "softend.de/fbrest/FBxRESTCore/apperrors"
	"io/ioutil"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

const PermissionKeyStr = "user permission key"

type PermissionType int

const (
	All       = 9
	None      = 0
	Read      = 1
	ReadWrite = 2
)

type repository struct {
	permissions map[string]Permission
	mu          sync.RWMutex
}

//UserKey and UserPassword for AppNeeds and Session Token geb´neration
//DBUser and DBPassword for datanabase login

type Permission struct {
	UserKey      string `xml:"userkey"`
	UserPassword string `xml:"userpassword"`
	/*
		DBUser string `xml:"dbuser"`
		DBPassword string `xml:"dbpassword"`
	*/
	DBScheme string         `xml:"dbscheme"`
	Type     PermissionType `xml:"type"`
}

type Permissions struct {
	Permission []Permission `xml:"Permission"`
}

func (r *repository) Add(userkey string, userpassword string, dbscheme string, ptype PermissionType) (ky Permission) {
	const funcstr = "func Permissions.Add"

	r.mu.Lock()
	defer r.mu.Unlock()
	var data Permission
	data.UserKey = userkey
	data.UserPassword = userpassword
	data.Type = ptype
	data.DBScheme = dbscheme
	r.permissions[userkey] = data

	log.WithFields(log.Fields{"userkey": userkey}).Debug(funcstr)
	return data
}

func (r *repository) Delete(token string) {
	const funcstr = "func Permissions.Delete"
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.permissions, token)
	log.WithFields(log.Fields{PermissionKeyStr: token}).Debug(funcstr)
}

func (r *repository) Get(token string) (item Permission, err error) {
	const funcstr = "func Permissions.Get"
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.permissions[token]
	if !ok {
		err = _apperrors.ErrPermissionKeyNotFound
		log.WithFields(log.Fields{PermissionKeyStr + " not found": token}).Debug(funcstr)
		return item, err
	}
	log.WithFields(log.Fields{PermissionKeyStr + " found": token}).Debug(funcstr)
	return item, nil
}

func (r *repository) Exists(token string) bool {
	const funcstr = "func Permissions.Exists"
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.permissions[token]
	if !ok {
		log.WithFields(log.Fields{PermissionKeyStr + " not found": token}).Debug(funcstr)
		return false
	}
	log.WithFields(log.Fields{PermissionKeyStr + " found": token}).Debug(funcstr)
	return true
}

var (
	pr *repository
)

func Repository() *repository {
	if pr == nil {
		pr = &repository{
			permissions: make(map[string]Permission),
		}
	}
	return pr
}

//Gets Permissiom from a Userkey and its password
//If permissions.DBUser for this key is empty than will be set to userkey, means userkey == DBLoginUser
//If permissions.DBPassword for this key is empty than will be set to userpassword
//if permissions.Type is empty it will be set to none
func GetPermissionFromRepository(userkey string, userpassword string, dbscheme string) (perm Permission, err error) {
	const funcstr = "func Permissions.GetPermissionFromRepository"
	log.WithFields(log.Fields{PermissionKeyStr: userkey}).Debug(funcstr)
	var rep = Repository()
	var result, err1 = rep.Get(userkey)
	if err1 != nil {
		return result, err1
	}

	if result.UserPassword != userpassword {
		err1 = _apperrors.ErrUserOrPasswordWrong
		return result, err1
	}

	return result, err
}

func ReadPermissions(pfile string) {
	const funcstr = "func Permissions.ReadPermissions"
	log.WithFields(log.Fields{"Permission file": pfile}).Debug(funcstr) 

	data, err := ioutil.ReadFile(pfile)
	if err != nil {
		log.WithFields(log.Fields{"File reading error": err}).Error(funcstr)
		return
	}
	xdata := &Permissions{}
	xml.Unmarshal(data, &xdata)
	var rep = Repository()
	for _, xd := range xdata.Permission {
		if len(xd.UserKey) > 0 {
			if !rep.Exists(xd.UserKey) {
				var itm = rep.Add(xd.UserKey, xd.UserPassword, xd.DBScheme, xd.Type)
				log.WithFields(log.Fields{"Added permission": "User key:" + itm.UserKey + " DBScheme:" + itm.DBScheme + " Permission:" + strconv.Itoa(int(itm.Type))}).Debug(funcstr)
			}
		}
	}
}

func WritePermissions(pfile string) {
	const funcstr = "func _permissions.WritePermissions"
	log.WithFields(log.Fields{pfile: pfile}).Debug(funcstr)
	var data Permissions
	var dt Permission
	dt.UserKey = "superuser"
	dt.UserPassword = "su"
	dt.DBScheme = "health_bln11"
	dt.Type = All
	data.Permission = append(data.Permission, dt)

	dt.UserKey = "456"
	dt.UserPassword = ""
	dt.DBScheme = "health_ffm1"
	dt.Type = None
	data.Permission = append(data.Permission, dt)

	dt.UserKey = "789"
	dt.UserPassword = ""
	dt.DBScheme = ""
	dt.Type = Read
	data.Permission = append(data.Permission, dt)

	file, _ := xml.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile(pfile, file, 0644)

}
