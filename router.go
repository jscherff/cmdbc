// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/jscherff/gocmdb"
)

func magtekRouter(musb gocmdb.MagtekUSB) (err error) {

	switch {

	case *fActionSerial:
		err = serialAction(musb)
		if err == nil {defer musb.Reset()}

	default:
		err = genericRouter(musb)
	}

	return err
}

func genericRouter(gusb gocmdb.GenericUSB) (err error) {

	switch {

	case *fActionAudit:
		err = auditRequest(gusb)

	case *fActionCheckin:
		err = checkinRequest(gusb)

	case *fActionLegacy:
		err = legacyAction(gusb)

	case *fActionReport:
		err = reportAction(gusb)

	case *fActionReset:
		defer gusb.Reset()
	}

	return err
}
