package sessions

import (
	"encoding/xml"
	_apperrors "softend.de/fbrest/FBxRESTCore/apperrors"
	_httpstuff "softend.de/fbrest/FBxRESTCore/httpstuff"
	_permissions "softend.de/fbrest/FBxRESTCore/permissions"
	_struct "softend.de/fbrest/FBxRESTCore/struct"

	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const MaxDuration = 30 * 60 * 1e9 //ns

type Item struct {
	Token      string                      `json:"Token"`
	Value      string                      `json:"Value"`
	Permission _permissions.PermissionType `json:"Permission"`
	Start      time.Time                   `json:"Time"`
	Duration   time.Duration               `json:"Duration"`
	Valid      bool                        `json:"Valid"`
}

type Sessions struct {
	Session []Item `xml:"Session"`
}

type repository struct {
	items map[string]Item
	mu    sync.RWMutex
}

func (r *repository) Add(token string, permission _permissions.PermissionType, conn string) (ky Item) {
	const funcStr = "func Sessions.Add"

	r.mu.Lock()
	defer r.mu.Unlock()
	var data Item
	data.Token = token
	data.Value = conn
	data.Start = time.Now()
	data.Duration = MaxDuration
	data.Permission = permission
	r.items[token] = data
	log.WithFields(log.Fields{_struct.SessionTokenStr: token}).Debug(funcStr)
	return data
}

func (r *repository) AddStruct(data Item) (ky Item) {
	const funcStr = "func Sessions.AddStruct"
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[data.Token] = data
	log.WithFields(log.Fields{_struct.SessionTokenStr + "  ": data.Token}).Debug(funcStr)
	return data
}

func (r *repository) Delete(token string) {
	const funcStr = "func Sessions.Delete"
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, token)
	log.WithFields(log.Fields{_struct.SessionTokenStr + "  ": token}).Debug(funcStr)
}

func (r *repository) Get(token string) (Item, error) {
	const funcStr = "func Sessions.Get"
	r.mu.RLock()
	defer r.mu.RUnlock()

	//log.Debug(r.items)
	item, ok := r.items[token]
	if !ok {
		var err error = _apperrors.ErrPermissionItemNotFound
		log.WithFields(log.Fields{_struct.SessionTokenStr + " not found": token}).Debug(funcStr)
		return item, err
	}
	log.WithFields(log.Fields{_struct.SessionTokenStr + " found": token}).Debug(funcStr)
	return item, nil
}

func GetTokenDataFromRepository(token string) (kval Item) {
	const funcStr = "func Sessions.GetTokenDataFromRepository"
	log.WithFields(log.Fields{_struct.SessionTokenStr: token}).Debug(funcStr)
	var rep = Repository()
	var result, _ = rep.Get(token)
	result.Token = token
	return result
}

func TokenValid(response http.ResponseWriter, key string) (kv Item) {
	const funcStr = "func Sessions.TokenValid"

	log.WithFields(log.Fields{_struct.SessionTokenStr: key}).Debug(funcStr)

	var Response _struct.ResponseData
	kv = GetTokenDataFromRepository(key)
	if len(kv.Value) < 1 {
		Response.Status = http.StatusForbidden
		Response.Message = "No valid database connection found by " + _struct.SessionTokenStr + " " + kv.Token
		Response.Data = _apperrors.DataNil
		log.WithFields(log.Fields{_struct.SessionTokenStr + " response:": Response.Message}).Error(funcStr)
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		kv.Valid = false
		return kv
	}
	var duration time.Duration = kv.Duration
	var end = time.Now()
	difference := end.Sub(kv.Start)
	if difference > duration {
		//Zeit für Key abgelaufen
		Response.Status = http.StatusForbidden
		Response.Message = _struct.SessionTokenStr + " " + kv.Token + " has expired after " + strconv.Itoa(MaxDuration/1e9) + " seconds"
		Response.Data = _apperrors.DataNil
		log.WithFields(log.Fields{"Session expired:": Response.Message}).Error(funcStr)
		var rep = Repository()
		rep.Delete(key)
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		kv.Valid = false

		return kv
	}
	var t = kv.Start.Add(kv.Duration)
	log.WithFields(log.Fields{"Session valid unitl " + t.Format(time.RFC822Z) + "->remaining " + _struct.SessionTokenStr + " duration (s)": (duration - difference)}).Debug(funcStr)
	kv.Valid = true
	return kv
}

func WritePermanantSessions(pfile string) {
	const funcStr = "func Sessions.WritePermanantSessions"
	log.WithFields(log.Fields{pfile: pfile}).Debug(funcStr)
	var dataarr Sessions
	var data Item

	data.Token = "ffm1health"
	data.Value = "SYSDBA:masterkey@localhost:3050/D:/Data/LAR/HEALTHDATAS304.FDB"
	data.Permission = _permissions.All
	data.Start = time.Now()
	data.Duration = 100000
	data.Valid = true
	dataarr.Session = append(dataarr.Session, data)

	data.Token = "bln11health"
	data.Value = "SYSDBA:masterkey@localhost:3050/D:/Data/LAR/HEALTHDATAS304.FDB"
	data.Permission = _permissions.All
	data.Start = time.Now()
	data.Duration = 100000
	data.Valid = true
	dataarr.Session = append(dataarr.Session, data)
	file, _ := xml.MarshalIndent(dataarr, "", " ")

	_ = ioutil.WriteFile(pfile, file, 0644)

}

func ReadPermanantSessions(pfile string) {

	const funcStr = "func Sessions.ReadPermanentSessions"
	log.WithFields(log.Fields{"SessionFile": pfile}).Debug(funcStr)

	dataarr, err := ioutil.ReadFile(pfile)
	if err != nil {
		log.WithFields(log.Fields{"File reading error": err}).Error(funcStr)
		return
	}
	//data := &[]Item{}
	data := &Sessions{}
	xml.Unmarshal(dataarr, &data)
	var rep = Repository()
	for _, pars := range data.Session {
		pars.Duration = pars.Duration * 1e9
		var t = pars.Start.Add(pars.Duration)
		var now = time.Now()
		log.Debug(" ")
		log.WithFields(log.Fields{"Read Token     ": pars.Token}).Debug(funcStr)
		log.WithFields(log.Fields{"Read Connection": pars.Value}).Debug(funcStr)
		log.WithFields(log.Fields{"Read Start     ": pars.Start}).Debug(funcStr)
		log.WithFields(log.Fields{"Read Duration  ": pars.Duration}).Debug(funcStr)
		log.WithFields(log.Fields{"Ends           ": t}).Debug(funcStr)
		log.Debug(" ")
		g1 := t.Before(now)
		if g1 {
			log.WithFields(log.Fields{_struct.SessionTokenStr + " " + pars.Token + " has exceeded date/time ": t}).Error(funcStr)
			log.Debug(" ")
		}

		rep.AddStruct(pars)
	}
}

var instance *repository

func Repository() *repository {
	if instance == nil {
		instance = &repository{
			items: make(map[string]Item),
		}
	}
	return instance
}
