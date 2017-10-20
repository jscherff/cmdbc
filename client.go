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
)

type httpStatus int

func (this httpStatus) Accepted() (ok bool) {

	ok = true

	switch int(this) {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusAccepted:
	case http.StatusNoContent:
	case http.StatusNotModified:
	default:
		ok = false
	}

	return ok
}

func (this httpStatus) String() (s string) {

	switch this {

	case http.StatusOK:
		s = `request processed, no errors`
	case http.StatusCreated:
		s = `request processed, object created`
	case http.StatusAccepted:
		s = `request processed, data accepted`
	case http.StatusNoContent:
		s = `request processed, no action taken`
	case http.StatusNotModified:
		s = `request processed, no changes found`
	case http.StatusBadRequest:
		s = `unsupported or malformed request`
	case http.StatusNotAcceptable:
		s = `insufficient or incorrect data`
	case http.StatusUnprocessableEntity:
		s = `unable to decode request`
	case http.StatusFailedDependency:
		s = `unsatisfied prerequisite`
	case http.StatusInternalServerError:
		s = `unable to process request`
	case http.StatusNotFound:
		s = `api endpoint not found`
	default:
		s = this.StatusText()
	}

	return s
}

func (this httpStatus) StatusText() (s string) {
	return http.StatusText(int(this))
}

// usbCiNewSnV1 obtains a serial number from the cmdbd server.
func usbCiNewSnV1(o gocmdb.Registerable) (s string, err error) {

	var (
		j []byte
		hs httpStatus
	)

	url := fmt.Sprintf(`%s/%s/%s/%s/%s`, conf.Server.URL,
		conf.Server.Endpoint[`usbCiNewSnV1`], o.Host(), o.VID(), o.PID(),
	)

	if j, err = o.JSON(); err != nil {
		return s, err
	}

	if j, hs, err = httpPost(url, j); err != nil {
		return s, err
	}

	if hs.Accepted() {
		err = json.Unmarshal(j, &s)
	} else {
		err = fmt.Errorf(`serial number not generated - %s`, hs)
	}

	if err == nil {
		slog.Printf(`serial number %q generated - %s`, s, hs)
	}

	return s, err
}

// usbCiCheckinV1 checks a device in with the cmdbd server.
func usbCiCheckinV1(o gocmdb.Registerable) (err error) {

	var (
		j []byte
		hs httpStatus
	)

	url := fmt.Sprintf(`%s/%s/%s/%s/%s`, conf.Server.URL,
		conf.Server.Endpoint[`usbCiCheckinV1`], o.Host(), o.VID(), o.PID(),
	)

	if j, err = o.JSON(); err != nil {
		return err
	}

	if _, hs, err = httpPost(url, j); err != nil {
		return err
	}

	if hs.Accepted() {
		slog.Printf(`checkin accepted - %s`, hs)
	} else {
		err = fmt.Errorf(`checkin not accepted - %s`, hs)
	}

	return err
}

// usbCiCheckoutV1 obtains the JSON representation of a serialized device object
// from the server using the unique key combination VID+PID+SN.
func usbCiCheckoutV1(o gocmdb.Auditable) (j []byte, err error) {

	var (
		hs httpStatus
	)

	if o.ID() == `` {
		slog.Print(`device %s-%s fetch: skipping, no SN`,
			o.VID(), o.PID(),
		)
		return j, err
	}

	url := fmt.Sprintf(`%s/%s/%s/%s/%s/%s`, conf.Server.URL,
		conf.Server.Endpoint[`usbCiCheckoutV1`], o.Host(), o.VID(), o.PID(), o.ID(),
	)

	if j, hs, err = httpGet(url); err != nil {
		return j, err
	}

	if hs.Accepted() {
		slog.Printf(`device retrieved - %s`, hs)
	} else {
		err = fmt.Errorf(`device not retreived - %s`, hs)
	}

	return j, err
}

// usbCiAuditV1 submits changes from audit to the server in JSON format.
func usbCiAuditV1(o gocmdb.Auditable) (err error) {

	var (
		j []byte
		hs httpStatus
	)

	url := fmt.Sprintf(`%s/%s/%s/%s/%s/%s`, conf.Server.URL,
		conf.Server.Endpoint[`usbCiAuditV1`], o.Host(), o.VID(), o.PID(), o.ID(),
	)

	if j, err = json.Marshal(o.GetChanges()); err != nil {
		return err
	}

	if _, hs, err = httpPost(url, j); err != nil {
		return err
	}

	if hs.Accepted() {
		slog.Printf(`audit accepted - %s`, hs)
	} else {
		err = fmt.Errorf(`audit not accepted - %s`, hs)
	}

	return err
}

// httpPost sends http POST requests to cmdbd server endpoints for other functions.
func httpPost(url string, j []byte ) (b []byte, hs httpStatus, err error) {

	if req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(j)); err != nil {
		return b, hs, err
	} else {
		req.Header.Add(`Content-Type`, `application/json; charset=UTF8`)
		return httpRequest(req)
	}
}

// httpGet sends http GET requests to cmdbd server endpoints for other functions.
func httpGet(url string) (b []byte, hs httpStatus, err error) {

	if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
		return b, hs, err
	} else {
		return httpRequest(req)
	}
}

// httpRequest sends http requests to cmdbd server endpoints for other functions.
func httpRequest(req *http.Request) (b []byte, hs httpStatus, err error) {

	req.Header.Add(`Accept`, `application/json; charset=UTF8`)
	req.Header.Add(`X-Custom-Header`, `cmdbc`)

	slog.Printf(`API call %s %s`, req.Method, req.URL)

	resp, err := client.Do(req)

	if err != nil {
		return b, hs, err
	}

	defer resp.Body.Close()

	hs = httpStatus(resp.StatusCode)
	b, err = ioutil.ReadAll(resp.Body)

	return b, hs, err
}
