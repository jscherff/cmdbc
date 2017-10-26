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
	`github.com/jscherff/cmdb/ci/peripheral/usb`
)

var (
	httpTransp = &http.Transport{ResponseHeaderTimeout: 10 * time.Second}
	httpClient = &http.Client{Transport: httpTransp}
)

type httpStatus int

// Accepted returns true for a successful http response status.
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

// String implements the Stringer interface for httpStatus.
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

// StatusText returns the HTTP status text associated with the status code.
func (this httpStatus) StatusText() (s string) {
	return http.StatusText(int(this))
}

// usbCiNewSnV1 obtains a serial number from the cmdbd server.
func usbCiNewSnV1(dev usb.Serializer) (string, error) {

	var s string

	url := fmt.Sprintf(`%s/%s/%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoint[`usbCiNewSnV1`],
		dev.Host(), dev.VID(), dev.PID(),
	)

	if j, err := dev.JSON(); err != nil {
		return ``, err
	} else if j, hs, err := httpPost(url, j); err != nil {
		return ``, err
	} else if !hs.Accepted() {
		return ``, fmt.Errorf(`serial number not generated - %s`, hs)
	} else if err = json.Unmarshal(j, &s); err != nil {
		return ``, err
	} else {
		sl.Printf(`serial number %q generated - %s`, s, hs)
		return s, nil
	}

}

// usbCiCheckinV1 checks a device in with the cmdbd server.
func usbCiCheckinV1(dev usb.Reporter) (error) {

	url := fmt.Sprintf(`%s/%s/%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoint[`usbCiCheckinV1`],
		dev.Host(), dev.VID(), dev.PID(),
	)

	if j, err := dev.JSON(); err != nil {
		return err
	} else if _, hs, err := httpPost(url, j); err != nil {
		return err
	} else if !hs.Accepted() {
		return fmt.Errorf(`checkin not accepted - %s`, hs)
	} else {
		sl.Printf(`checkin accepted - %s`, hs)
		return nil
	}
}

// usbCiCheckoutV1 obtains the JSON representation of a serialized device object
// from the server using the unique key combination VID+PID+SN.
func usbCiCheckoutV1(dev usb.Auditer) ([]byte, error) {

<<<<<<< HEAD:httpio.go
	if o.ID() == `` {
		slog.Printf(`device %s-%s fetch: skipping, no SN`,
			o.VID(), o.PID(),
		)
		return j, err
=======
	if dev.SN() == `` {
		sl.Printf(`device %s-%s skipping fetch, no SN`, dev.VID(), dev.PID())
		return nil, nil
>>>>>>> newlib:client.go
	}

	url := fmt.Sprintf(`%s/%s/%s/%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoint[`usbCiCheckoutV1`],
		dev.Host(), dev.VID(), dev.PID(), dev.SN(),
	)

	if j, hs, err := httpGet(url); err != nil {
		return nil, err
	} else if !hs.Accepted() {
		return nil, fmt.Errorf(`device not retreived - %s`, hs)
	} else {
		sl.Printf(`device retrieved - %s`, hs)
		return j, nil
	}
}

// usbCiAuditV1 submits changes from audit to the server in JSON format.
func usbCiAuditV1(dev usb.Auditer) (error) {

	url := fmt.Sprintf(`%s/%s/%s/%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoint[`usbCiAuditV1`],
		dev.Host(), dev.VID(), dev.PID(), dev.SN(),
	)

	if j, err := json.Marshal(dev.GetChanges()); err != nil {
		return err
	} else if _, hs, err := httpPost(url, j); err != nil {
		return err
	} else if !hs.Accepted() {
		return fmt.Errorf(`audit not accepted - %s`, hs)
	} else {
		sl.Printf(`audit accepted - %s`, hs)
		return nil
	}
}

//usbMetaVendorV1 retrieves the vendor name given the vid.
func usbMetaVendorV1(dev usb.Updater) (s string, err error) {

	url := fmt.Sprintf(`%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoint[`usbMetaVendorV1`],
		dev.VID(),
	)

	if j, hs, err := httpGet(url); err != nil {
		return ``, err
	} else if !hs.Accepted() {
		return ``, fmt.Errorf(`lookup failed - %s`, hs)
	} else if err = json.Unmarshal(j, &s); err != nil {
		return ``, err
	} else {
		sl.Printf(`lookup succeeded - %s`, hs)
		return s, nil
	}
}

//usbMetaProductV1 retrieves the product name given the vid and pid.
func usbMetaProductV1(dev usb.Updater) (s string, err error) {

	url := fmt.Sprintf(`%s/%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoint[`usbMetaProductV1`],
		dev.VID(), dev.PID(),
	)

	if j, hs, err := httpGet(url); err != nil {
		return ``, err
	} else if !hs.Accepted() {
		return ``, fmt.Errorf(`lookup failed - %s`, hs)
	} else if err = json.Unmarshal(j, &s); err != nil {
		return ``, err
	} else {
		sl.Printf(`lookup succeeded - %s`, hs)
		return s, nil
	}
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

	sl.Printf(`API call %s %s`, req.Method, req.URL)

	resp, err := httpClient.Do(req)

	if err != nil {
		return b, hs, err
	}

	defer resp.Body.Close()

	hs = httpStatus(resp.StatusCode)
	b, err = ioutil.ReadAll(resp.Body)

	return b, hs, err
}
