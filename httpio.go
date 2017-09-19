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
	`github.com/jscherff/gocmdb`
)

var HttpStatusMap = map[int]string {
	http.StatusOK:			`request processed, no errors`,		// 200
	http.StatusCreated:		`request processed, object created`,	// 201
	http.StatusAccepted:		`request processed, data accepted`,	// 202
	http.StatusNoContent:		`request processed, no action taken`,	// 204
	http.StatusNotModified:		`request processed, no changes found`,	// 302
	http.StatusBadRequest:		`unsupported or malformed request`,	// 400
	http.StatusNotAcceptable:	`insufficient or incorrect data`,	// 406
	http.StatusUnprocessableEntity:	`unable to decode request`,		// 422
	http.StatusFailedDependency:	`request condition not satisified`,	// 424
	http.StatusInternalServerError:	`unable to process request`,		// 500
}

// newSNRequest obtains a serial number from the gocmdbd server.
func newSNRequest(o gocmdb.Registerable) (s string, err error) {

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

	if j, sc, err = httpRequest(url, j); err != nil {
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

// checkinRequest checks a device in with the gocmdbd server.
func checkinRequest(o gocmdb.Registerable) (err error) {

	var j []byte

	url := fmt.Sprintf(`%s/%s/%s/%s/%s`, conf.Server.URL,
		conf.Server.CheckinPath, o.Host(), o.VID(), o.PID())

	if j, err = o.JSON(); err != nil {
		elog.Print(err)
		return err
	}

	_, _, err = httpRequest(url, j)

	return err
}

// auditRequest requests a server-side audit against the previous checkin.
func auditRequest(o gocmdb.Auditable) (ss [][]string, err error) {

	if o.ID() == `` {
		slog.Print(`skipping audit for VID %q PID %q: no serial number`, o.VID(), o.PID())
		return ss, err
	}

	var (
		j []byte
		sc int
	)

	url := fmt.Sprintf(`%s/%s/%s/%s/%s/%s`, conf.Server.URL,
		conf.Server.AuditPath, o.Host(), o.VID(), o.PID(), o.ID())

	if j, err = o.JSON(); err != nil {
		elog.Print(err)
		return ss, err
	}

	if j, sc, err = httpRequest(url, j); err != nil {
		return ss, err
	}

	if sc == http.StatusNoContent {
		slog.Print(`empty content from server`)
		return ss, err
	}

	if err = json.Unmarshal(j, &ss); err != nil {
		elog.Print(err)
	}

	return ss, err
}

// httpRequest sends JSON requests to the gocmdbd server for other functions.
// Error decoration will be handled by caller functions.
func httpRequest(url string, jreq []byte ) (jrsp []byte, sc int, err error) {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jreq))

	if err != nil {
		elog.Print(err)
		return jrsp, sc, err
	}

	req.Header.Add(`Content-Type`, `application/json; charset=UTF8`)
	req.Header.Add(`Accept`, `application/json; charset=UTF8`)
	req.Header.Add(`X-Custom-Header`, `gocmdb`)

	rsp, err := client.Do(req)

	if err != nil {
		elog.Print(err)
		return jrsp, sc, err
	}

	defer rsp.Body.Close()

	msg := fmt.Sprintf(`http status %q: %s`,
		rsp.Status, HttpStatusMap[rsp.StatusCode])

	if rsp.StatusCode < 400 {
		slog.Print(msg)
	} else {
		elog.Print(msg)
	}

	jrsp, err = ioutil.ReadAll(rsp.Body)

	if err != nil {
		elog.Print(err)
	}

	return jrsp, rsp.StatusCode, err
}
