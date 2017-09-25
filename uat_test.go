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

/* TODO: test each path through the application by setting flags and
   passing an object to the router for both magtek and generic devices.

	Actions:

	  -audit
	        Audit devices
	  -checkin
	        Check devices in
	  -legacy
	        Legacy operation
	  -report
	        Report actions
	  -reset
	        Reset device
	  -serial
	        Set serial number


	Audit Options:

	  -local
	        Audit against local state
	  -server
	        Audit against server state


	Report Options:

	  -console
	        Write reports to console
	  -folder <path>
	        Write reports to <path>
	  -format <format>
	        Report <format> {csv|nvp|xml|json}


	Serial Options:

	  -copy
	        Copy factory serial number
	  -erase
	        Erase current serial number
	  -fetch
	        Fetch serial number from server
	  -force
	        Force serial number change
	  -set <string>
	        Set serial number to <string>

 PATHS:

	  -audit -local
	  -audit -server
	  -checkin
	  -legacy
	  -report -folder -format csv
	  -report -folder -format nvp
	  -report -folder -format xml
	  -report -folder -format json
	  -reset
	  -serial -copy
	  -serial -erase
	  -serial -fetch
	  -serial -set <string>
	  -serial -fetch -force
	  -serial -set <string> -force

*/
