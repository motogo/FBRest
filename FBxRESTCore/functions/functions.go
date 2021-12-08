// Package functions implements functions for returning URL and BODY content as struced data
package _functions

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"encoding/json"
	_dbscheme "softend.de/fbrest/FBxRESTCore/dbscheme"
	_httpstuff "softend.de/fbrest/FBxRESTCore/httpstuff"
	_sessions "softend.de/fbrest/FBxRESTCore/sessions"
	_struct "softend.de/fbrest/FBxRESTCore/struct"

	"net/http"

	"softend.de/fbrest/FBxRESTCore/config"
	"html/template"
	_ "image/png"
	"path"
)

//Returns respone for HTML site of FBRest usage
func ResponseHelpHTML(w http.ResponseWriter, code int) {

	const funcstr = "func Functions.ResponseHelpHTML"
	log.Debug(funcstr)
	profile := _struct.Profile{Appname: config.AppName, Version: config.Version, Copyright: config.Copyright, Key: "-MNhE7Yf50sz6U9Hgqae", Duration: _sessions.MaxDuration}
	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func ResponseDataHTML(w http.ResponseWriter, code int) {

}

func ResponseHelpDesignHTML(w http.ResponseWriter, code int) {

	/*
		reader, err := os.Open("templates/selfhtml.png")
		if err != nil {
		     log.Fatal(err)
		}
		m, _, err := image.Decode(reader)
		defer reader.Close()
	*/
	const funcstr = "func Functions.ResponseHelpDesignHTML"
	log.Debug(funcstr)
	profile := _struct.Profile{Appname: config.AppName, Version: config.Version, Copyright: config.Copyright, Key: "-MNhE7Yf50sz6U9Hgqae", Duration: _sessions.MaxDuration}
	fp := path.Join("templates", "design.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func ResponseHelpCommandsHTML(w http.ResponseWriter, code int) {
	const funcstr = "func Functions.ResponseHelpCommandsHTML"
	log.Debug(funcstr)
	profile := _struct.Profile{Appname: config.AppName, Version: config.Version, Copyright: config.Copyright, Key: "-MNhE7Yf50sz6U9Hgqae", Duration: _sessions.MaxDuration}
	fp := path.Join("templates", "commands.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ResponseInfoBusyText(w http.ResponseWriter, code int) {
	const funcstr = "func Functions.ResponseInfoBusyText"
	log.Debug(funcstr)
	profile := _struct.Profile{Appname: config.AppName, Version: config.Version, Copyright: config.Copyright, Key: "-MNhE7Yf50sz6U9Hgqae", Duration: _sessions.MaxDuration}
	fp := path.Join("templates", "busy.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func OutParameters(entitiesData _struct.SQLAttributes) {
	const funcstr = "func Functions.OutParameters"
	log.Debug(funcstr)

	log.WithFields(log.Fields{"Given command     ": entitiesData.Cmd}).Info("")
	log.WithFields(log.Fields{"Given sepfill     ": entitiesData.Sepfill}).Info("")
	log.WithFields(log.Fields{"Given info        ": entitiesData.Info}).Info("")
}

func OutTableParameters(entitiesData _struct.GetTABLEAttributes) {
	const funcstr = "func Functions.OutTableParameters"
	log.Debug(funcstr)

	log.WithFields(log.Fields{"Given fields      ": entitiesData.Fields}).Info("")
	log.WithFields(log.Fields{"Given order by    ": entitiesData.OrderBy}).Info("")
	log.WithFields(log.Fields{"Given group by    ": entitiesData.GroupBy}).Info("")
	log.WithFields(log.Fields{"Given sepfill     ": entitiesData.Sepfill}).Info("")
	log.WithFields(log.Fields{"Given info        ": entitiesData.Info}).Info("")
}

func GetSQLParamsFromBODY(r *http.Request, entitiesData *_struct.SQLAttributes) bool {
	const funcstr = "func Functions.GetSQLParamsFromBODY"
	log.WithFields(log.Fields{"body": r.Body}).Debug(funcstr)
	var xdata _struct.GetUrlSQLAttributes
	err2 := json.NewDecoder(r.Body).Decode(&xdata)
	if err2 != nil {
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error()}).Error(funcstr)
		return false
	}
	entitiesData.Cmd = ReplaceAritmetics(xdata.Cmd)
	entitiesData.Info = xdata.Info
	return true
}

func GetTableParamsFromBODY(r *http.Request, entitiesData *_struct.GetTABLEAttributes) (ok bool) {
	const funcstr = "func Functions.GetTableParamsFromBODY"
	log.WithFields(log.Fields{"body": r.Body}).Debug(funcstr)

	var xdata _struct.GetUrlTABLEAttributes
	err2 := json.NewDecoder(r.Body).Decode(&xdata)
	if err2 != nil {
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error()}).Error(funcstr)
		return false
	}
	entitiesData.Fields = strings.Join(xdata.Fields, ",")
	entitiesData.Filter = xdata.Filter
	entitiesData.GroupBy = strings.Join(xdata.GroupBy, ",")
	entitiesData.OrderBy = strings.Join(xdata.OrderBy, ",")
	entitiesData.Skip = xdata.Skip
	entitiesData.First = xdata.First
	return true
}

func GetFieldParamsFromBODY(r *http.Request, entitiesData *_struct.GetTABLEAttributes) (ok bool) {
	const funcstr = "func Functions.GetFieldParamsFromBODY"
	log.WithFields(log.Fields{"body": r.Body}).Debug(funcstr)

	var xdata _struct.GetUrlTABLEAttributes
	err2 := json.NewDecoder(r.Body).Decode(&xdata)
	if err2 != nil {
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error()}).Error(funcstr)
		return false
	}
	entitiesData.Fields = strings.Join(xdata.Fields, ",")
	entitiesData.Filter = xdata.Filter
	entitiesData.GroupBy = strings.Join(xdata.GroupBy, ",")
	entitiesData.OrderBy = strings.Join(xdata.OrderBy, ",")
	entitiesData.Skip = xdata.Skip
	entitiesData.First = xdata.First
	return true
}

func GetSQLParamsFromURL(r *http.Request, entitiesData *_struct.SQLAttributes) (ok bool) {
	const funcstr = "func Functions.GetSQLParamsFromURL"

	log.WithFields(log.Fields{"url": r.URL}).Debug(funcstr)

	var u = r.URL

	if strings.HasPrefix(u.RawQuery, _struct.FormatJson) {
		var par = u.RawQuery[len(_struct.FormatJson)+1:]
		par = _httpstuff.UnEscape(par)

		if len(par) > 0 {
			xdata := &_struct.GetUrlSQLAttributes{}

			err := json.Unmarshal([]byte(par), &xdata)

			if err != nil {
				return false
			}

			log.Debug(xdata)
			entitiesData.Cmd = xdata.Cmd
			entitiesData.Info = xdata.Info
			return true
		}
	}
	var okret bool
	urlparams, ok := r.URL.Query()["cmd"]
	infoparams, okinfo := r.URL.Query()["info"]
	log.WithFields(log.Fields{"URL params length": len(urlparams)}).Debug(funcstr)

	if ok && len(urlparams[0]) > 0 {

		var cmd string = string(urlparams[0])
		if strings.HasPrefix(cmd, "'") || strings.HasPrefix(cmd, "-") {

			entitiesData.Sepfill = cmd[:1]
			cmd = cmd[1:]
			cmd = strings.ReplaceAll(cmd, entitiesData.Sepfill, " ")
			log.WithFields(log.Fields{"Sepfill": entitiesData.Sepfill}).Debug(funcstr + "->set key")
		}

		log.WithFields(log.Fields{"Cmd": cmd}).Debug(funcstr + "->set key")
		entitiesData.Cmd = cmd
		okret = true
	}
	if okinfo && len(infoparams[0]) > 0 {

		var info string = string(infoparams[0])
		log.WithFields(log.Fields{"Info": info}).Debug(funcstr + "->set key")
		entitiesData.Info = info
		okret = true
	}

	return okret
}

func GetSessionParamsFromBODY(r *http.Request, entitiesData *_dbscheme.DatabaseAttributes) bool {
	const funcstr = "func Functions.GetSessionParamsFromBODY"
	log.Debug(funcstr)
	var xdata _dbscheme.DatabaseAttributes

	err2 := json.NewDecoder(r.Body).Decode(&xdata)
	if err2 != nil {
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error()}).Error(funcstr)
		return false
	}

	entitiesData.Database = xdata.Database
	entitiesData.Location = xdata.Location
	entitiesData.Password = xdata.Password
	entitiesData.User = xdata.User
	entitiesData.Port = xdata.Port

	return true
}

func GetSessionSchemeParamsFromBODY(r *http.Request, entitiesData *_dbscheme.GetUrlSessionSchemeAttributes) bool {
	const funcstr = "func Functions.GetSessionSchemeParamsFromBODY"
	log.Debug(funcstr)
	var xdata _dbscheme.GetUrlSessionSchemeAttributes

	err2 := json.NewDecoder(r.Body).Decode(&xdata)
	if err2 != nil {
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error()}).Error(funcstr)
		return false
	}

	entitiesData.Password = xdata.Password
	entitiesData.User = xdata.User
	entitiesData.DBScheme = xdata.DBScheme

	return true
}

func GetSessionParamsFromURL(r *http.Request, entitiesData *_dbscheme.DatabaseAttributes) bool {

	const funcstr = "func Functions.GetSessionParamsFromURL"

	var u = r.URL
	if strings.HasPrefix(u.RawQuery, _struct.FormatJson) {
		var par = u.RawQuery[len(_struct.FormatJson)+1:]
		par = _httpstuff.UnEscape(par)
		//	 par = "{\"location\":\"localhost\",\"database\":\"D:/Data/Dokuments/DOKUMENTS30.FDB\",\"port\":3050,\"password\":\"su\",\"user\":\"superuser\"}"
		if len(par) > 0 {
			xdata := &_dbscheme.DatabaseAttributes{}

			err := json.Unmarshal([]byte(par), &xdata)

			if err != nil {
				return false
			}

			log.Info(xdata)
			entitiesData.Database = xdata.Database
			entitiesData.Port = xdata.Port
			entitiesData.Password = xdata.Password
			entitiesData.User = xdata.User
			entitiesData.Location = xdata.Location
			return true
		}
	}

	var okret bool
	databaseparams, databaseok := u.Query()["Database"]
	locationparams, locationok := u.Query()["Location"]
	portparams, portok := u.Query()["Port"]
	userparams, userok := u.Query()["User"]
	passwordparams, passwordok := u.Query()["Password"]

	if databaseok && len(databaseparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(databaseparams)}).Debug(funcstr)
		urlparam := databaseparams[0]
		if urlparam[:1] == "(" {
			entitiesData.Database = urlparam[1 : len(urlparam)-1]
		} else {
			entitiesData.Database = urlparam
		}
		okret = true
	}

	if locationok && len(locationparams[0]) > 0 {
		log.WithFields(log.Fields{"URL location length": len(locationparams)}).Debug(funcstr)
		urlparam := locationparams[0]
		if urlparam[:1] == "(" {
			entitiesData.Location = urlparam[1 : len(urlparam)-1]
		} else {
			entitiesData.Location = urlparam
		}
		okret = true
	}

	if portok && len(portparams[0]) > 0 {
		log.WithFields(log.Fields{"URL port length": len(portparams)}).Debug(funcstr)
		urlparam := portparams[0]
		if urlparam[:1] == "(" {
			entitiesData.Port, _ = strconv.Atoi(urlparam[1 : len(urlparam)-1])
		} else {
			entitiesData.Port, _ = strconv.Atoi(urlparam)
		}
		okret = true
	}

	if userok && len(userparams[0]) > 0 {
		log.WithFields(log.Fields{"URL user length": len(userparams)}).Debug(funcstr)
		urlparam := userparams[0]
		if urlparam[:1] == "(" {
			entitiesData.User = urlparam[1 : len(urlparam)-1]
		} else {
			entitiesData.User = urlparam
		}
		okret = true
	}

	if passwordok && len(passwordparams[0]) > 0 {
		log.WithFields(log.Fields{"URL password length": len(passwordparams)}).Debug(funcstr)
		urlparam := passwordparams[0]
		if urlparam[:1] == "(" {
			entitiesData.Password = urlparam[1 : len(urlparam)-1]
		} else {
			entitiesData.Password = urlparam
		}
		okret = true
	}

	return okret
}

func GetSessionSchemeParamsFromURL(r *http.Request, entitiesData *_dbscheme.GetUrlSessionSchemeAttributes) bool {

	const funcstr = "func Functions.GetSessionSchemeParamsFromURL"

	var u = r.URL
	if strings.HasPrefix(u.RawQuery, _struct.FormatJson) {
		var par = u.RawQuery[len(_struct.FormatJson)+1:]
		par = _httpstuff.UnEscape(par)
		//old	 par = "{\"location\":\"localhost\",\"database\":\"D:/Data/Dokuments/DOKUMENTS30.FDB\",\"port\":3050,\"password\":\"su\",\"user\":\"superuser\"}"
		//old	 par = "{\"dbscheme\":\"health_ffm1\",\"password\":\"su\",\"user\":\"superuser\"}"
		if len(par) > 0 {
			xdata := &_dbscheme.GetUrlSessionSchemeAttributes{}

			err := json.Unmarshal([]byte(par), &xdata)

			if err != nil {
				return false
			}

			log.Info(xdata)
			entitiesData.Password = xdata.Password
			entitiesData.User = xdata.User
			entitiesData.DBScheme = xdata.DBScheme
			return true
		}
	}

	var okret bool
	dbschemeparams, dbschemeok := u.Query()["DBScheme"]

	userparams, userok := u.Query()["User"]
	passwordparams, passwordok := u.Query()["Password"]

	if dbschemeok && len(dbschemeparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(dbschemeparams)}).Debug(funcstr)
		urlparam := dbschemeparams[0]
		if urlparam[:1] == "(" {
			entitiesData.DBScheme = urlparam[1 : len(urlparam)-1]
		} else {
			entitiesData.DBScheme = urlparam
		}
		okret = true
	}

	if userok && len(userparams[0]) > 0 {
		log.WithFields(log.Fields{"URL user length": len(userparams)}).Debug(funcstr)
		urlparam := userparams[0]
		if urlparam[:1] == "(" {
			entitiesData.User = urlparam[1 : len(urlparam)-1]
		} else {
			entitiesData.User = urlparam
		}
		okret = true
	}

	if passwordok && len(passwordparams[0]) > 0 {
		log.WithFields(log.Fields{"URL password length": len(passwordparams)}).Debug(funcstr)
		urlparam := passwordparams[0]
		if urlparam[:1] == "(" {
			entitiesData.Password = urlparam[1 : len(urlparam)-1]
		} else {
			entitiesData.Password = urlparam
		}
		okret = true
	}

	return okret
}

func GetFIELDPayloadFromString2(params string, entitiesData *_struct.FIELDVALUEAttributes) {

	// payload=(id:1, username: 'admin', email: 'email@example.org')

	const funcstr = "func Functions.GetFIELDPayloadFromString"
	var psplit = "&"
	var csplit = "="
	var par = strings.Split(params, psplit)

	log.WithFields(log.Fields{"URL params length": len(par)}).Debug(funcstr)

	if len(par) > 0 {
		for _, pars := range par {
			params = _httpstuff.UnEscape(pars)
			keyval := strings.SplitN(params, csplit, 2)

			log.WithFields(log.Fields{"Key": string(keyval[0])}).Debug(funcstr + "->found key")
			log.WithFields(log.Fields{"Val": string(keyval[1])}).Debug(funcstr + "->found value")

			if strings.EqualFold(string(keyval[0]), string("FIELDS")) {
				log.WithFields(log.Fields{"Fields": string(keyval[1])}).Debug(funcstr + "->set Fields")

			}
		}
	}

}

func GetSQLParamsFromString2(params string, entitiesData *_struct.SQLAttributes) {

	const funcstr = "func Functions.GetSQLParamsFromString2"
	var csplit = "="
	var par = params

	log.WithFields(log.Fields{"URL params length": len(par)}).Debug(funcstr)
	params = _httpstuff.UnEscape(par)
	log.WithFields(log.Fields{"SQL": params}).Debug(funcstr + "->set key")

	keyval := strings.SplitN(params, csplit, 2)

	log.WithFields(log.Fields{"Key": string(keyval[0])}).Debug(funcstr + "->found key")
	log.WithFields(log.Fields{"Val": string(keyval[1])}).Debug(funcstr + "->found value")

	if strings.EqualFold(string(keyval[0]), string("CMD")) {
		var cmd string = string(keyval[1])
		if strings.HasPrefix(cmd, "'") || strings.HasPrefix(cmd, "-") {

			entitiesData.Sepfill = cmd[:1]
			cmd = cmd[1:]
			cmd = strings.ReplaceAll(cmd, entitiesData.Sepfill, " ")
			log.WithFields(log.Fields{"Sepfill": entitiesData.Sepfill}).Debug(funcstr)
		}

		log.WithFields(log.Fields{"Set CMD": cmd}).Debug(funcstr)
		entitiesData.Cmd = cmd
	}

	if strings.EqualFold(string(keyval[0]), string("INFO")) {
		log.WithFields(log.Fields{"Set INFO": string(keyval[1])}).Debug(funcstr)
		entitiesData.Info = string(keyval[1])
	}

}

//Returns the last-nLeft slice from URL
//e.g. when nLeft == 0 returns the last slice
func GetRightPathSliceFromURL(r *http.Request, nLeft int) (key string) {
	const funcstr = "func Functions.GetRightPathSliceFromURL"
	urlstr := string(r.URL.String())
	var keyval = strings.SplitN(urlstr, "?", 2)

	urlstr = keyval[0]
	t2 := strings.Split(urlstr, "/")
	key = t2[len(t2)-1-nLeft]
	log.WithFields(log.Fields{_struct.URLKeyStr: key}).Debug(funcstr)
	return key
}

func GetLeftPathSliceFromURL(r *http.Request, nLeft int) (key string) {
	const funcstr = "func Functions.GetLeftPathSliceFromURL"
	urlstr := string(r.URL.String())
	var keyval = strings.SplitN(urlstr, "?", 2)

	urlstr = keyval[0]

	t2 := strings.Split(urlstr, "/")
	key = t2[nLeft+1]
	log.WithFields(log.Fields{_struct.URLKeyStr: key}).Debug(funcstr)
	return key
}

func GetTableParamsFromURL(r *http.Request, entitiesData *_struct.GetTABLEAttributes) bool {

	//  http://localhost:4488/{{.Key}}/rest/get/TSTANDORT?fjson={"table": "TSTANDORT","fields": ["ID","BEZ","GUELTIG"],"filter":"ID=1 AND BEZ like 'x%'","order": ["BEZ ASC","ID DESC"],"groupby": ["ID","BEZ"],"first": 0}

	//const funcstr = "func Functions.GetTableParamsFromURL"

	var u = r.URL

	if strings.HasPrefix(u.RawQuery, _struct.FormatJson) {
		var par = u.RawQuery[len(_struct.FormatJson)+1:]
		par = _httpstuff.UnEscape(par)

		if len(par) > 0 {
			xdata := &_struct.GetUrlTABLEAttributes{}

			err := json.Unmarshal([]byte(par), &xdata)

			if err != nil {
				return false
			}

			log.Info(xdata)
			entitiesData.Fields = strings.Join(xdata.Fields, ",")
			entitiesData.Filter = xdata.Filter
			entitiesData.GroupBy = strings.Join(xdata.GroupBy, ",")
			entitiesData.OrderBy = strings.Join(xdata.OrderBy, ",")
			entitiesData.Skip = xdata.Skip
			entitiesData.First = xdata.First
			return true
		}
	}

	var s = _httpstuff.UnEscape(u.RawQuery)

	var pars = strings.Split(s, "&")
	var okret bool
	for _, par := range pars {
		if strings.HasPrefix(par, _struct.Fields+"=") {
			var par = par[len(_struct.Fields)+1:]
			if par[:1] == "(" {
				entitiesData.Fields = par[1 : len(par)-1]
			} else {
				entitiesData.Fields = par
			}
			if len(entitiesData.Fields) < 1 {
				entitiesData.Fields = "*"
			}
			okret = true
		} else if strings.HasPrefix(par, _struct.Filter+"=") {
			var par = par[len(_struct.Filter)+1:]
			if par[:1] == "(" {
				entitiesData.Filter = par[1 : len(par)-1]
			} else {
				entitiesData.Filter = par
			}
			entitiesData.Filter = ReplaceAritmetics(entitiesData.Filter)
			okret = true
		} else if strings.HasPrefix(par, _struct.Order+"=") {
			var par = par[len(_struct.Order)+1:]
			if par[:1] == "(" {
				entitiesData.OrderBy = par[1 : len(par)-1]
			} else {
				entitiesData.OrderBy = par
			}
			okret = true
		} else if strings.HasPrefix(par, _struct.Group+"=") {
			var par = par[len(_struct.Group)+1:]
			if par[:1] == "(" {
				entitiesData.GroupBy = par[1 : len(par)-1]
			} else {
				entitiesData.GroupBy = par
			}
			okret = true

		} else if strings.HasPrefix(par, _struct.Info+"=") {
			var par = par[len(_struct.Info)+1:]
			if par[:1] == "(" {
				entitiesData.Info = par[1 : len(par)-1]
			} else {
				entitiesData.Info = par
			}
			okret = true

		} else if strings.HasPrefix(par, _struct.Limit+"=") {
			var par = par[len(_struct.Limit)+1:]
			var limit string
			if par[:1] == "(" {
				limit = par[1 : len(par)-1]
			} else {
				limit = par
			}
			var lm = strings.Split(limit, ",")
			if len(lm) == 1 {
				entitiesData.First, _ = strconv.Atoi(lm[0])
				entitiesData.Skip = 0
			} else if len(lm) == 2 {
				entitiesData.First, _ = strconv.Atoi(lm[0])
				entitiesData.Skip, _ = strconv.Atoi(lm[1])
			}
			okret = true
		}
	}

	return okret
}

func GetFieldParamsFromURL(r *http.Request, entitiesData *_struct.GetTABLEAttributes) bool {

	//  http://localhost:4488/{{.Key}}/rest/get/TSTANDORT?fjson={"table": "TSTANDORT","fields": ["ID","BEZ","GUELTIG"],"filter":"ID=1 AND BEZ like 'x%'","order": ["BEZ ASC","ID DESC"],"groupby": ["ID","BEZ"],"first": 0}

	//const funcstr = "func functions.GetFieldParamsFromURL"

	var u = r.URL
	if strings.HasPrefix(u.RawQuery, _struct.FormatJson) {
		var par = u.RawQuery[len(_struct.FormatJson)+1:]
		par = _httpstuff.UnEscape(par)

		if len(par) > 0 {
			xdata := &_struct.GetUrlTABLEAttributes{}

			err := json.Unmarshal([]byte(par), &xdata)

			if err != nil {
				return false
			}

			log.Info(xdata)
			entitiesData.Fields = strings.Join(xdata.Fields, ",")
			entitiesData.Filter = xdata.Filter
			entitiesData.GroupBy = strings.Join(xdata.GroupBy, ",")
			entitiesData.OrderBy = strings.Join(xdata.OrderBy, ",")
			entitiesData.Skip = xdata.Skip
			entitiesData.First = xdata.First
			return true
		}
	}

	var s = _httpstuff.UnEscape(u.RawQuery)

	var pars = strings.Split(s, "&")
	var okret bool
	for _, par := range pars {
		if strings.HasPrefix(par, _struct.Filter+"=") {
			var par = par[len(_struct.Filter)+1:]
			if par[:1] == "(" {
				entitiesData.Filter = par[1 : len(par)-1]
			} else {
				entitiesData.Filter = par
			}
			entitiesData.Filter = ReplaceAritmetics(entitiesData.Filter)
			okret = true
		} else if strings.HasPrefix(par, _struct.Order+"=") {
			var par = par[len(_struct.Order)+1:]
			if par[:1] == "(" {
				entitiesData.OrderBy = par[1 : len(par)-1]
			} else {
				entitiesData.OrderBy = par
			}
			okret = true
		} else if strings.HasPrefix(par, _struct.Group+"=") {
			var par = par[len(_struct.Group)+1:]
			if par[:1] == "(" {
				entitiesData.GroupBy = par[1 : len(par)-1]
			} else {
				entitiesData.GroupBy = par
			}
			okret = true

		} else if strings.HasPrefix(par, _struct.Info+"=") {
			var par = par[len(_struct.Info)+1:]
			if par[:1] == "(" {
				entitiesData.Info = par[1 : len(par)-1]
			} else {
				entitiesData.Info = par
			}
			okret = true

		} else if strings.HasPrefix(par, _struct.Limit+"=") {
			var par = par[len(_struct.Limit)+1:]
			var limit string
			if par[:1] == "(" {
				limit = par[1 : len(par)-1]
			} else {
				limit = par
			}
			var lm = strings.Split(limit, ",")
			if len(lm) == 1 {
				entitiesData.First, _ = strconv.Atoi(lm[0])
				entitiesData.Skip = 0
			} else if len(lm) == 2 {
				entitiesData.First, _ = strconv.Atoi(lm[0])
				entitiesData.Skip, _ = strconv.Atoi(lm[1])
			}
			okret = true
		}
	}

	return okret
}

func GetFIELDPayloadFromBODY(r *http.Request, entitiesData *_struct.FIELDVALUEAttributes) bool {
	const funcstr = "func Functions.GetFIELDPayloadFromBODY"

	var xdata _struct.GetUrlPayloadAttributes
	//body, err2 := ioutil.ReadAll(r.Body)
	//var str string = string(body)
	//log.Info(str)
	err2 := json.NewDecoder(r.Body).Decode(&xdata)
	if err2 != nil {
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error()}).Error(funcstr)
		return false
	} else {
		log.WithFields(log.Fields{"Body Payload": xdata.Payload}).Debug(funcstr)
	}

	entitiesData.FieldValue = xdata.Payload
	xdata.Filter = ReplaceAritmetics(xdata.Filter)

	entitiesData.Filter = xdata.Filter
	return true
}

func GetFIELDPayloadFromURL(r *http.Request, entitiesData *_struct.FIELDVALUEAttributes) bool {
	//  http://localhost:4488/{{.Key}}/rest/put/TSTANDORT?payload=(bez='Nürnberg2')&filter=(bez='Nürnberg1')
	//  http://localhost:4488/{{.Key}}/rest/put/TSTANDORT?ftext="payload=(bez='Nürnberg2')&filter=(bez='Nürnberg1')"
	// http://localhost:4488/{{.Key}}/rest/put/TSTANDORT?fjson={"payload":["ID='123'","BEZ='test'","GUELTIG=1"], "filter": "ID=1 AND BEZ like 'x%'"}

	const funcstr = "func Functions.GetFIELDPayloadFromURL"
	var u = r.URL
	log.Debug(funcstr)
	if strings.HasPrefix(u.RawQuery, _struct.FormatJson) {
		var par = u.RawQuery[len(_struct.FormatJson)+1:]
		par = _httpstuff.UnEscape(par)
		if len(par) > 0 {
			xdata := &_struct.GetUrlPayloadAttributes{}
			err := json.Unmarshal([]byte(par), &xdata)
			if err != nil {
				return false
			}

			log.WithFields(log.Fields{"xdata": xdata}).Debug(funcstr)

			for _, vals := range xdata.Payload {
				entitiesData.FieldValue = append(entitiesData.FieldValue, vals)
			}
			entitiesData.Filter = xdata.Filter
			return true
		}
	}

	var s = _httpstuff.UnEscape(u.RawQuery)

	var pars = strings.Split(s, "&")
	var okret bool

	log.WithFields(log.Fields{"Parameters": pars}).Debug(funcstr)

	for _, par := range pars {
		if strings.HasPrefix(par, _struct.Filter+"=") {
			var par = par[len(_struct.Filter)+1:]
			if par[:1] == "(" {
				entitiesData.Filter = par[1 : len(par)-1]
			} else {
				entitiesData.Filter = par
			}
			entitiesData.Filter = ReplaceAritmetics(entitiesData.Filter)
			okret = true
		} else if strings.HasPrefix(par, _struct.Payload+"=") {
			var pars = par[len(_struct.Payload)+1:]
			var st string = pars[:1]

			//log.WithFields(log.Fields{"st": st}).Debug(funcstr)
			if st == "(" {
				pars = pars[1:]
			}
			//log.Debug(funcstr + "->pars:" + pars)
			st = pars[len(pars)-1:]
			//log.Debug(funcstr + "->st:" + st)
			if st == ")" {
				pars = pars[:len(pars)-1]
			}
			//log.Debug(funcstr + "->pars:" + pars)
			//keyval :=  strings.SplitN(pars,",",-1)
			keyval := SplitPars(pars, ",")

			for _, pars := range keyval {
				entitiesData.FieldValue = append(entitiesData.FieldValue, pars)
			}

			okret = true

		}
	}

	return okret
}
