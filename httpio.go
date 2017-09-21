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
	`bytes`
	`encoding/json`
	`fmt`
	`io/ioutil`
	`net/http`
	`time`
	`github.com/jscherff/gocmdb`
)

var (
	transport = &http.Transport{ResponseHeaderTimeout: 10 * time.Second}
	client = &http.Client{Transport: transport}

	HttpStatusMap = map[int]string {
		http.StatusOK:			`request processed, no errors`,		// 200
		http.StatusCreated:		`request processed, object created`,	// 201
		http.StatusAccepted:		`request processed, data accepted`,	// 202
		http.StatusNoContent:		`request processed, no action taken`,	// 204
		http.StatusNotModified:		`request processed, no changes found`,	// 302
		http.StatusBadRequest:		`unsupported or malformed request`,	// 400
		http.StatusNotAcceptable:	`insufficient or incorrect data`,	// 406
		http.StatusUnprocessableEntity:	`unable to decode request`,		// 422
		http.StatusFailedDependency:	`unsatisfied prerequisite`,		// 424
		http.StatusInternalServerError:	`unable to process request`,		// 500
	}
)

// GetNewSN obtains a serial number from the gocmdbd server.
func GetNewSN(o gocmdb.Registerable) (s string, err error) {

	var (
		j []byte
		sc int
	)

	url := fmt.Sprintf(`%s/%s/%s/%s/%s`, conf.Server.URL,
		conf.Server.NewSNPath, o.Host(), o.VID(), o.PID())

	if j, err = o.JSON(); err != nil {
		elog.Print(err)
		return s, err
	}

	if j, sc, err = httpPost(url, j); err != nil {
		return s, err
	}

	if sc == http.StatusNoContent {
		slog.Print(`empty content from server`)
		return s, err
	}

	if err = json.Unmarshal(j, &s); err != nil {
		elog.Print(err)
	}

	return s, err
}

// SubmitCheckin checks a device in with the gocmdbd server.
func SubmitCheckin(o gocmdb.Registerable) (err error) {

	var j []byte

	url := fmt.Sprintf(`%s/%s/%s/%s/%s`, conf.Server.URL,
		conf.Server.CheckinPath, o.Host(), o.VID(), o.PID())

	if j, err = o.JSON(); err != nil {
		elog.Print(err)
		return err
	}

	_, _, err = httpPost(url, j)

	return err
}

// GetDevice obtains the JSON representation of a serialized device object
// from the server using the unique key combination VID+PID+SN.
func GetDevice(o gocmdb.Registerable) (j []byte, err error) {

	if o.ID() == `` {
		slog.Print(`device %s-%s fetch: skipping, no serial number`,
			o.VID(), o.PID(),
		)
		return j, err
	}

	url := fmt.Sprintf(`%s/%s/%s/%s/%s/%s`, conf.Server.URL,
		conf.Server.FetchPath, o.Host(), o.VID(), o.PID(), o.ID())

	j, _, err = httpGet(url)

	return j, err
}

// SubmitAudit submits changes from audit to the server in JSON format.
func SubmitAudit(o gocmdb.Auditable) (err error) {

	url := fmt.Sprintf(`%s/%s/%s/%s/%s/%s`, conf.Server.URL,
		conf.Server.AuditPath, o.Host(), o.VID(), o.PID(), o.ID())

	j, err := json.Marshal(o.GetChanges())

	if err == nil {
		_, _, err = httpPost(url, j)
	}

	if err != nil {
		elog.Print(err)
	}

	return err
}

// httpPost sends http POST requests to gocmdbd server endpoints for other functions.
func httpPost(url string, j []byte ) (b []byte, sc int, err error) {

	if req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(j)); err != nil {
		return b, sc, err
	} else {
		req.Header.Add(`Content-Type`, `application/json; charset=UTF8`)
		return httpRequest(req)
	}
}

// httpGet sends http GET requests to gocmdbd server endpoints for other functions.
func httpGet(url string) (b []byte, sc int, err error) {

	if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
		return b, sc, err
	} else {
		return httpRequest(req)
	}
}

// httpRequest sends http requests to gocmdbd server endpoints for other functions.
func httpRequest(req *http.Request) (b []byte, sc int, err error) {

	req.Header.Add(`Accept`, `application/json; charset=UTF8`)
	req.Header.Add(`X-Custom-Header`, `gocmdb`)

	resp, err := client.Do(req)

	if err == nil {

		defer resp.Body.Close()

		sc = resp.StatusCode
		msg := fmt.Sprintf(`%s - %s`, resp.Status, HttpStatusMap[sc])

		if sc < 400 {
			slog.Print(msg)
			b, err = ioutil.ReadAll(resp.Body)
		} else {
			elog.Print(msg)
		}
	}

	if err != nil {
		elog.Print(err)
	}

	return b, sc, err
}
