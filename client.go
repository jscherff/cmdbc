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
	`github.com/jscherff/cmdb/ci/peripheral/usb`
)

// authenticated tracks whether or not client has authenbticated
// with server so that functions calling protected API endpoints
// can determine whether or not they need to call auth().
var authenticated = false

// httpResult contains the results of an http request/response.
type httpResult struct {
	status httpStatus
	content httpContent
}

// Status returns the status of the http response.
func (this *httpResult) Status() (httpStatus) {
	return this.status
}

// Content returns the body of the http response.
func (this *httpResult) Content() (httpContent) {
	return this.content
}

// String implements the Stringer interface for httpResult.
func (this *httpResult) String() (string) {
	return fmt.Sprintf(`%s: %s`, this.Status(), this.Content())
}

// httpStatus represents an http response status code.
type httpStatus int

// Accepted returns true for a successful http response status.
func (this httpStatus) Accepted() (bool) {

	switch int(this) {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusAccepted:
	case http.StatusNoContent:
	case http.StatusNotModified:
	default: return false
	}

	return true
}

// Rejected returns true for an unsuccessful http response status.
func (this httpStatus) Rejected() (bool) {
	return !this.Accepted()
}

// String implements the Stringer interface for httpStatus.
func (this httpStatus) String() (string) {
	return this.StatusText()
}

// StatusText returns the HTTP status text associated with the status code.
func (this httpStatus) StatusText() (string) {
	return http.StatusText(int(this))
}

// httpContent represents the body of an http response.
type httpContent []byte

// Decode unmarshals content into the provided object.
func (this httpContent) Decode(t interface{}) (error) {
	return json.Unmarshal(this, t)
}

// String implements the Stringer interface for httpStatus.
func (this httpContent) String() (s string) {
	if err := this.Decode(&s); err != nil {
		return string(this)
	} else {
		return s
	}
}

// auth authenticates with the server using basic authentication and, if
// successful, obtains JWT for API authentication in a cookie.
func auth() error {

	if authenticated {
		return nil
	}

	url := fmt.Sprintf(conf.Server.Endpoints[`cmdb_auth`],
		conf.Client.HostName,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(
		conf.Server.Auth.Username,
		conf.Server.Auth.Password,
	)

	if hr, err := httpRequest(req); err != nil {
		return err
	} else if hr.Status().Rejected() {
		return fmt.Errorf(`authentication failure - %s`, hr)
	} else {
		sl.Printf(`authentication success - %s`, hr.Status())
	}

	authenticated = true
	return nil
}

// newSn obtains a serial number from the cmdbd server.
func newSn(dev usb.Serializer) (string, error) {

	if err := auth(); err != nil {
		return ``, err
	}

	url := fmt.Sprintf(conf.Server.Endpoints[`usb_ci_newsn`],
		conf.Client.HostName,
		dev.VID(),
		dev.PID(),
	)

	var s string

	if j, err := dev.JSON(); err != nil {
		return ``, err
	} else if hr, err := httpPost(url, j); err != nil {
		return ``, err
	} else if hr.Status().Rejected() {
		return ``, fmt.Errorf(`serial number not generated - %s`, hr)
	} else if err := hr.Content().Decode(&s); err != nil {
		return ``, err
	} else {
		sl.Printf(`serial number '%s' generated - %s`, s, hr.Status())
		return s, nil
	}
}

// checkin checks a device in with the cmdbd server.
func checkin(dev usb.Reporter) (error) {

	if err := auth(); err != nil {
		return err
	}

	url := fmt.Sprintf(conf.Server.Endpoints[`usb_ci_checkin`],
		conf.Client.HostName,
		dev.VID(),
		dev.PID(),
	)

	if j, err := dev.JSON(); err != nil {
		return err
	} else if hr, err := httpPost(url, j); err != nil {
		return err
	} else if hr.Status().Rejected() {
		return fmt.Errorf(`checkin not accepted - %s`, hr)
	} else {
		sl.Printf(`checkin accepted - %s`, hr.Status())
		return nil
	}
}

// checkout obtains the JSON representation of a serialized device object
// from the server using the unique key combination VID+PID+SN.
func checkout(dev usb.Auditer) ([]byte, error) {

	if err := auth(); err != nil {
		return nil, err
	}

	if dev.SN() == `` {
		sl.Printf(`device %s-%s skipping fetch, no SN`, dev.VID(), dev.PID())
		return nil, nil
	}

	url := fmt.Sprintf(conf.Server.Endpoints[`usb_ci_checkout`],
		conf.Client.HostName,
		dev.VID(),
		dev.PID(),
		dev.SN(),
	)

	if hr, err := httpGet(url); err != nil {
		return nil, err
	} else if hr.Status().Rejected() {
		return nil, fmt.Errorf(`device not retreived - %s`, hr)
	} else {
		sl.Printf(`device retrieved - %s`, hr.Status())
		return hr.Content(), nil
	}
}

// sendAudit submits changes from audit to the server in JSON format.
func sendAudit(dev usb.Auditer) (error) {

	if err := auth(); err != nil {
		return err
	}

	url := fmt.Sprintf(conf.Server.Endpoints[`usb_ci_audit`],
		conf.Client.HostName,
		dev.VID(),
		dev.PID(),
		dev.SN(),
	)

	if j, err := json.Marshal(dev.GetChanges()); err != nil {
		return err
	} else if hr, err := httpPost(url, j); err != nil {
		return err
	} else if hr.Status().Rejected() {
		return fmt.Errorf(`audit not accepted - %s`, hr)
	} else {
		sl.Printf(`audit accepted - %s`, hr.Status())
		return nil
	}
}

// vendor retrieves the vendor name given the vid.
func vendor(dev usb.Updater) (string, error) {

	url := fmt.Sprintf(conf.Server.Endpoints[`usb_meta_vendor`],
		dev.VID(),
	)

	var s string

	if hr, err := httpGet(url); err != nil {
		return ``, err
	} else if hr.Status().Rejected() {
		return ``, fmt.Errorf(`vendor lookup failed - %s`, hr)
	} else if err := hr.Content().Decode(&s); err != nil {
		return ``, err
	} else {
		sl.Printf(`vendor lookup succeeded - %s`, hr.Status())
		return s, nil
	}
}

// product retrieves the product name given the vid and pid.
func product(dev usb.Updater) (string, error) {

	url := fmt.Sprintf(conf.Server.Endpoints[`usb_meta_product`],
		dev.VID(),
		dev.PID(),
	)

	var s string

	if hr, err := httpGet(url); err != nil {
		return ``, err
	} else if hr.Status().Rejected() {
		return ``, fmt.Errorf(`product lookup failed - %s`, hr)
	} else if err := hr.Content().Decode(&s); err != nil {
		return ``, err
	} else {
		sl.Printf(`product lookup succeeded - %s`, hr.Status())
		return s, nil
	}
}

// httpPost sends http POST requests to cmdbd server endpoints for other functions.
func httpPost(url string, data []byte ) (*httpResult, error) {

	if req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data)); err != nil {
		return nil, err
	} else {
		req.Header.Add(`Content-Type`, `application/json; charset=UTF8`)
		return httpRequest(req)
	}
}

// httpGet sends http GET requests to cmdbd server endpoints for other functions.
func httpGet(url string) (*httpResult, error) {

	if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	} else {
		return httpRequest(req)
	}
}

// httpRequest sends http requests to cmdbd server endpoints for other functions.
func httpRequest(req *http.Request) (*httpResult, error) {

	req.Header.Add(`Accept`, `application/json; charset=UTF8`)
	req.Header.Add(`X-Custom-Header`, `cmdbc`)

	sl.Printf(`API call %s %s`, req.Method, req.URL)

	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	stat := httpStatus(resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)

	return &httpResult{stat, body}, err
}
