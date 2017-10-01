# CMDBc
The _**Configuration Management Database Client**_ is a utility that manages information about devices attached to end-user workstations and reports that information to a cetralized repository over a RESTful JSON API provided by the complementary server component, the _**Configuration Management Database Daemon**_ or [**CMDBd**](https://github.com/jscherff/cmdbd/blob/master/README.md). **CMDBc** can register or _"check-in"_ attached devices with the server, obtain unique serial numbers from the server for devices that support serial number configuration, perform audits against previous device configurations, and report configuration changes found during the audit to the server for logging and analysis.

### System Requirements
**CMDBc** is written in **Go** and can be compiled for any operating system and architecture. It is intended to be installed on end-user workstations running **Microsoft Windows 7** or higher and should be invoked by a centralized management solution like **IBM BigFix**.

### Installation
Pre-compiled Windows binaries are available for both 32- and 64-bit systems and can be installed in any folder along with the required JSON configuration file:

* [`cmdbc.exe`](https://github.com/jscherff/cmdbc/raw/master/i686/cmdbc.exe) (x86_32)
* [`cmdbc.exe`](https://github.com/jscherff/cmdbc/raw/master/x86_64/cmdbc.exe) (x86_64)
* [`config.json`](https://raw.githubusercontent.com/jscherff/cmdbc/master/config.json)


### Configuration
The JSON configuration file, [`config.json`](https://github.com/jscherff/cmdbd/blob/master/config.json), is mostly self-explanatory. The default settings are sane and you should not have to change them in most use cases.

##### Server Settings
Parameters for communicating with the **CMDBd** server.
```json
"Server": {
    "URL": "http://sysadm-dev-01.24hourfit.com:8080",
    "CheckinPath": "usbci/checkin",
    "CheckoutPath": "usbci/checkout",
    "NewSNPath": "usbci/newsn",
    "AuditPath": "usbci/audit"
}
```
* **`URL`** is the base URL for the **CMDBd** server hosting the REST API.
* **`CheckinPath`** is the path below the server URL for registration or "_check-in_" actions.
* **`CheckoutPath`** is the path below the server URL for obtaining previously submitted device confignuration information for the purpose of conducting an audit.
* **`NewSNPath`** is the path below the server URL for obtaining a new, unique serial number for assignment to devices that support serial number configuration.
* **`AuditPath`** is the path below the server URL for submitting device configuration changes discovered during an audit.

##### Path Settings
Directories where logs, state files, and reports will be written. 
```json
"Paths": {
    "LogDir": "log",
    "StateDir": "state",
    "ReportDir": "report"
}
```
* **`LogDir`** is the directory where log files are written. When a relative path like `log` is provided, the directory will be created below the appliation directory. 
* **`StateDir`** is where device state files are stored. State files are used in performing local audits.
* **`ReportDir`** is where device reports are written.

##### File Settings
Filenames for logs and the legacy report file.
```json
"Files": {
    "SystemLog": "system.log",
    "ChangeLog": "change.log",
    "ErrorLog": "error.log",
    "Legacy": "usb_serial.txt"
}
```
* **`SystemLog`** is the name of the file where **CMDBc** records significant, non-error events.
* **`ChangeLog`** is the name of the file where **CMDBc** records changes found during audits. It also reports changes to the **CMDBd** server.
* **`ErrorLog`** is the name of the file where **CMDBc** records errors.
* **`Legacy`** is the name of the file where **CMDBc** writes the legacy inventory report.

##### Logging Settings
Granular logging options for the system, change, and error log.
```json
"Logging": {
    "System": {
        "Logfile": true,
        "Console": false,
        "Syslog": false
    },
    "Change": {
        "Logfile": true,
        "Console": false,
        "Syslog": false
    },
    "Error": {
        "Logfile": true,
        "Console": true,
        "Syslog": false
    }
},
```
* **`Logfile`** specifies whether or not events are written to log files.
* **`Console`** specifies whether or not events are written to the console (stdout).
* **`Syslog`** causes the utility to write events to a local or remote syslog daemon using the `Syslog` configuration settings, below.

##### Syslog Settings
Parameters for communicating with a local or remote syslog server.
```json
"Syslog": {
    "Protocol": "tcp",
    "Port": "1468",
    "Host": "localhost"
},
```
* **`Protocol`** is the transport-layer protocol used by the syslog daemon (blank for local).
* **`Port`** is the port used by the syslog daemon (blank for local).
* **`Host`** is the hostname or IP address of the syslog daemon (blank for local).

##### Include Settings
Vendors and products to include (_true_) or exclude (_false_) when inventorying devices.
```json
"Include": {
    "VendorID": {
        "043d": false,
        "045e": false
    },
    "ProductID": {
        "0801": {
            "0001": true
            "0002": true
            "0011": true
            "0012": true
            "0013": true
        },
        "0acd": {
            "2030": true
        }
    },
    "Default": true
},
```
* **`VendorID`** specifies which vendors to include (_true_) or exclude (_false_). This setting applies to all of the vendor's products and overrides both the _ProductID_ and _Default_ configuration settings; that is, if a vendor is excluded under _VendorID_, that vendor's products cannot be included under _ProductID_. Here, all devices with **Microsoft** _Vendor IDs_ `043d` and `045e` will be excluded.
* **`ProductID`** specifies which products to include (_true_) or exclude (_false_). This setting applies to specific _ProductIDs_ under a given _VendorID_ and overrides the _Default_ configuration setting. Here, **MagTek** (_VendorID_ `0801`) card readers with  _ProductIDs_ `0001`, `0002`, `0011`, `0012`, and `0013` will be included, as will **ID TECH** (_VendorID_ `0acd`) Card 
* **`Default`** specifies the default behavior for products that are not specifically included or excluded by _Vendor ID_ or _Product ID_. Here the default is to include, which effectively renders previous inclusions redundant; however, specific _VendorID_ and _ProductID_ inclusions ensure that those devices will be inventoried even if the _Default_ setting is changed to 'exclude' (_false_).

##### Format Settings
Default file formats for various use cases.
```json
"Format": {
    "Report": "csv",
    "Default": json
}
```
* **`Report`** is the default output format for inventory reports.
* **`Default`** is the default output format for other use cases.

### Operation
##### Command-Line Flags
Client operation is controlled through command-line _flags_. There are seven top-level _action flags_ -- `audit`, `checkin`, `legacy`, `report`, `reset`, `serial`, and `help`.  Some of these require (or offer) additional _option flags_.
 
* **`-audit`** performs a device configuration change audit.
    * **`-local`** audits against JSON state files stored on the local machine
    * **`-server`**	audits against the last device check-in stored in the database.
    * **`-help`** lists _audit option flags_ and their descriptions.
* **`-checkin`** checks devices in with the server, which stores device information in the database along with the check-in date.
* **`-legacy`** specifies _legacy mode_, which produces the same output to the same filename, `usb_serials.txt`, as the legacy inventory utility. The utility will also operate in legacy mode if the executable is renamed from **cmdbc.exe** to **magtek_inventory.exe**, the name of the legacy inventory utility executable.
* **`-report`** generates device configuration reports.
    * **`-console`** writes report output to the console.
    * **`-folder`** _`<path>`_ writes report output files to _`<path>`_. It defaults to the `report` folder beneath the installation directory.
    * **`-format`** _`<format>`_ specifies which report _`<format>`_ to use:
        * **`csv`** specifies comma-separated value format (default).
        * **`nvp`** specifies name-value pair format.
        * **`xml`** specifies extensible markup language format.
        * **`json`** specifies JavaScript object notation format.
    * **`-help`** lists _report option flags_ and their descriptions.
* **`-reset`** resets the device.
* **`-serial`** performs serial number operations. (By default, **CMDBc** will not configure a serial number on a device that already has one.)
    * **`-copy`** copies the factory serial number (if present) to the active serial number.
    * **`-erase`** erases the current serial number.
    * **`-fetch`** fetches a unique serial number from the server.
    * **`-force`** forces a serial number change, even if the device already has one.
    * **`-set`** _`<value>`_ sets serial number to the specified _`<value>`_.
    * **`-help`** lists _serial option flags_ and their descriptions.
* **`-help`** lists top-level _action flags_ and their descriptions.

##### Flag Combinations
Some _action flags_ can take multiple options.
* The **`-report`** _option flags_ can be used together in any combination. Example:
    * **`-report -format`** `json` **`-folder`** `C:\Reports` **`-console`** will write device configuration reports in JSON format to `C:\Reports` and will also display the reports on the screen.
* The **`-serial`** _option flags_ `-copy`, `-fetch`, and `-set` are mutually-exclusive, but each can be combined with `-erase` and `-force`. Examples:
    * **`-serial -fetch -force`** will fetch a new, unique serial number from the **CMDBd** server and will configure the device with it, overriding the safety mechanism that prevents overwriting existing serial numbers.
    * **`-serial -erase -fetch`** will erase the existing serial number, fetch a new, unique serial number from the **CMDBd** server, and will configure the device with it (same end-result as `-serial -fetch -force`).

**Caution**: action and option flags apply to _all attached devices_; if you run the utility with the `-serial -fetch` flags, it will only configure new serial numbers on compatible devices that don't already have a serial number. If all attached devices already have serial numbers or are not configurable, nothing will happen. However, if you add the `-force` flag, it will overwrite the serial number on all compatible devices -- even those that already have a serial number. If you run the utility with the `-serial -set -force` and there is more than one configurable device attached, you will end up having multiple devices with the same serial number.


Actions and events are recorded in `system.log`, errors are recorded in `error.log`, and changes detected during audits are recorded in `change.log`. The log directory is configurable; the default is the `log` subdirectory under the folder in which the utility is installed. All three logs can also be written to the console (stdout) and/or to a local or remote syslog server.
 
Device state is stored in JSON files in the state subdirectory directory (configurable)
 
Report files are written to the report subdirectory directory (configurable)
 
Serial number requests, check-ins, and audits record the following information in the database:
* Hostname
* Vendor ID
* Product ID
* Serial Number
* Vendor Name
* Product Name
* Product Version
* Software ID
* Bus Number
* Bus Address
* Port Number
* Buffer Size
* Max Packet Size
* USB Specification
* USB Class
* USB Subclass
* USB Protocol
* Device Speed
* Device Version
* Factory Serial Number
* Date/Time
 
Audits also record the following information for each change detected:
* Property
* Old Value
* New Value
* Date/Time
