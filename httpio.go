// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"path/filepath"
	"io/ioutil"
	"net/http"
	"bytes"
	"fmt"

	"github.com/jscherff/gocmdb/webapi"
	"github.com/jscherff/gocmdb"
)

// Process checkin action.
func serialRequest(o gocmdb.Registerable) (s string, e error) {

	var j []byte

	if j, e = o.JSON(); e != nil {
		return s, e
	}

	wd, e := webapi.NewDevice(j)

	if e != nil {
		return s, e
	}

	if j, e = wd.JSON(); e == nil {
		url := fmt.Sprintf("%s/%s/%s", config.ServerURL, config.SerialPath, o.Type())
		j, e = httpRequest(url, j)
	}

	if e == nil {
		e = json.Unmarshal(j, &wd)
	}

	return wd.ID(), e
}

// Process checkin action.
func checkinRequest(o gocmdb.Registerable) (e error) {

	var j []byte

	if j, e = o.JSON(); e != nil {
		return e
	}

	wd, e := webapi.NewDevice(j)

	if e != nil {
		return e
	}

	if j, e = wd.JSON(); e == nil {
		url := fmt.Sprintf("%s/%s/%s", config.ServerURL, config.CheckinPath, o.Type())
		_, e = httpRequest(url, j)
	}

	return e
}

// Process audit action.
func auditRequest(o gocmdb.Auditable) (e error) {

	var j []byte

	id, e := o.ID()

	if e != nil {
		return e
	}

	c, e := o.Compare(filepath.Join(config.AuditDir, id + ".json"))

	if e == nil && len(c) != 0 {

		wc := webapi.NewChanges(c)
		j, e = wc.JSON()

		if e == nil {
			url := fmt.Sprintf("%s/%s/%s", config.ServerURL, config.AuditPath, id)
			_, e = httpRequest(url, j)
		}
	}

	_ = o.Save(filepath.Join(config.AuditDir, id + ".json"))

	return e
}

func httpRequest(url string, jreq []byte ) (jresp []byte, e error) {

	client := &http.Client{}

	req, e := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jreq))

	req.Header.Add("Content-Type", "application/json; charset=UTF8")
	req.Header.Add("Accept", "application/json; charset=UTF8")
	req.Header.Add("X-Custom-Header", "gocmdb")

	resp, e := client.Do(req)

	if e != nil {
		return jresp, e
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusAccepted:
	default:
		e = fmt.Errorf("http response status %s", resp.Status)
	}

	if e == nil {
		jresp, e = ioutil.ReadAll(resp.Body)
	}

	return jresp, e
}
