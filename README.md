# CMDBc
The _**Configuration Management Database Client**_ is a utility that manages information about devices attached to end-user workstations and reports that information to a cetralized repository over a RESTful JSON API provided by the complementary server component, the _**Configuration Management Database Daemon**_ or [**CMDBd**](https://github.com/jscherff/cmdbd/blob/master/README.md). **CMDBc** can register or _"check-in"_ attached devices with the server, obtain unique serial numbers from the server for devices that support serial number configuration, perform audits against previous device configurations, and report configuration changes found during the audit to the server for logging and analysis.

### System Requirements
**CMDBc** is written in **Go** and can be compiled for any operating system and architecture. It is intended to be installed on end-user workstations running **Microsoft Windows 7** or higher and should be invoked by a centralized management solution like **IBM BigFix**.

### Installation
Pre-compiled Windows binaries are available for both 32- and 64-bit systems and can be installed in any folder along with the required JSON configuration file:

* [`cmdbc.exe`](https://github.com/jscherff/cmdbc/raw/master/x86_32/cmdbc.exe) (x86_32)
* [`cmdbc.exe`](https://github.com/jscherff/cmdbc/raw/master/x86_64/cmdbc.exe) (x86_64)
* [`config.json`](https://raw.githubusercontent.com/jscherff/cmdbc/master/config.json)

### Configuration
The JSON configuration file, [`config.json`](https://github.com/jscherff/cmdbd/blob/master/config.json), is mostly self-explanatory. The default settings are sane and you should not have to change them in most use cases.

---
**WIP WIP WIP WIP WIP WIP WIP WIP WIP WIP WIP WIP**

---
Client operation is controlled through command-line actions and options. There are six top-level ‘actions,’ some of which require (or offer) additional sub-options:
 
* **`-audit`** performs a device change audit
    * **`-local`** audits against JSON state files stored on the local machine
    * **`-server`**	audits against the last device check-in stored in the database
-checkin	Checks devices in with the server, which stores device information in the database along with the check-in date
-legacy	Legacy mode, which produces the same output in the same file as the current inventory script
-report	Generates reports
      -console	Writes report output to the console (can be combined with -folder)
      -folder <path>	Writes report output files to <path>
      -format <format>	Specifies which report <format> to use
            csv	Comma-Separated Value format (default)
            nvp	Name-Value Pair format
            xml	Extensible Markup Language format
            json	JavaScript Object Notation format
-reset	Resets the device
-serial	Performs serial number operations
      -copy	Copies the factory serial number to the active serial number
      -erase	Erases the current serial number (can be combined with other options)
      -fetch	Fetches a unique serial number from the server
      -force	Forces a serial number change (overrides safety mechanism which prevents setting a serial number when one is already present)
      -set <value>	Sets serial number to the specified <value>
-help	Top-level actions help menu
-audit -help	Audit options help menu
-report -help	Report options help menu
-serial -help	Serial number options help menu
 
Actions and events are logged to system.log, errors are logged in error.log, and changes detected during audits are recorded in change.log. The log directory is configurable; the default is the log subdirectory under the application directory. All three kinds of logs can also be written to a local or remote Syslog server.
 
Device state is stored in JSON files in the state subdirectory directory (configurable)
 
Report files are written to the report subdirectory directory (configurable)
 
Serial number requests, check-ins, and audits record the following information in the database:
•	Hostname
•	Vendor ID
•	Product ID
•	Serial Number
•	Vendor Name
•	Product Name
•	Product Version
•	Software ID
•	Bus Number
•	Bus Address
•	Port Number
•	Buffer Size
•	Max Packet Size
•	USB Specification
•	USB Class
•	USB Subclass
•	USB Protocol
•	Device Speed
•	Device Version
•	Factory Serial Number
•	Date/Time
 
Audits also record the following information for each change detected:
•	Property
•	Old Value
•	New Value
•	Date/Time
 
The configuration file is self-explanatory and probably won’t need modification:

### Path Settings
```json
"Paths": {
    "LogDir": "log",
    "StateDir": "state",
    "ReportDir": "report"
}
```
* **`LogDir`** is the directory where log files are written. When a relative path like `log` is provided, the directory will be created below the appliation directory. 
* **`StateDir`** is where device state files are stored. State files are used in performing local audits.
* **`ReportDIr`** is where device reports are written.

### File Settings
```json
"Files": {
    "SystemLog": "system.log",
    "ChangeLog": "change.log",
    "ErrorLog": "error.log",
    "Legacy": "usb_serial.txt"
}
```

### Server Settings
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
### Logging Settings
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

### Syslog Settings
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
### Include Settings
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
        }
    },
    "Default": true
},
```
* **`VendorID`** specifies which vendors to include (`true`) or exclude (`false`). This setting applies to all of the vendor's products and overrides both _Product ID_ and _Default_ settings. Here, all devices with **Microsoft** _Vendor IDs_ `043d` and `045e` will be excluded.
* **`ProductID`** specifies which products to include (`true`) or exlude (`false`). This setting applies to specific _Product IDs_ under a given _Vendor ID_. Here, devices with a **MagTek** _Vendor ID_ `0801` and Card Reader _Product IDs_ `0001`, `0002`, `0011`, `0012`, and `0013` will be included, and it overrides the `Default` settings.
* **`Default`** specifies the default behavior for products that are not specifically included or excluded by _Vendor ID_ or _Product ID_. Here the default is to include (`true`), which effectively renders previous inclusions redundant.
```json
"Format": {
    "Report": "csv",
    "object": "json",
    "Legacy": "csv"
}
```

---
---

**Server Settings**
```json
"Server": {
    "Addr": ":8080",
    "ReadTimeout": 10,
    "WriteTimeout": 10,
    "MaxHeaderBytes": 1048576,
    "HttpBodySizeLimit": 1048576,
    "AllowedContentTypes": ["application/json"]
}
```
* **`Addr`** is the hostname or IP address and port of the listener, separated by a colon. If blank, the daemon will listen on all network interfaces.
* **`ReadTimeout`** is the maximum duration in seconds for reading the entire HTTP request, including the body.
* **`WriteTimeout`** is the maximum duration in seconds before timing out writes of the response.
* **`MaxHeaderBytes`** is the maximum size in bytes of the request header.
* **`HttpBodySizeLimit`** is the maximum size in bytes of the request body.
* **`AllowedContentTypes`** is a comma-separated list of allowed media types.

**Database Settings**
```json
"Database": {
    "Driver": "mysql",
    "Config": {
        "User": "cmdbd",
        "Passwd": "K2Cvg3NeyR",
        "Net": "",
        "Addr": "localhost",
        "DBName": "gocmdb",
        "Params": null
    },
    ...
}
```
* **`Driver`** is the database driver. Only `mysql` is supported.
* **`User`** is the database user the daemon uses to access the database.
* **`Passwd`** is the database user password. The default, shown, should be changed in production.
* **`Net`** is the port on which the database is listening. If blank, the daemon will use the MySQL default port, 3306.
* **`Addr`** is the database hostname or IP address.
* **`DBName`** is the database schema used by the application.
* **`Params`** are additional parameters to pass to the driver (advanced).

**Logger Settings**
```json
"Loggers": {
    "system": {
        "LogFile": "system.log",
        "LogFlags": ["date","time","shortfile"],
        "Stdout": false,
        "Stderr": false,
        "Syslog": false
    },
    "access": {
        "LogFile": "access.log",
        "LogFlags": [],
        "Stdout": false,
        "Stderr": false,
        "Syslog": true
    },
    "error": {
        "LogFile": "error.log",
        "LogFlags": ["date","time","shortfile"],
        "Stdout": false,
        "Stderr": false,
        "Syslog": false
    }
}
```
* **`LogFile`** is the filename of the log file.
* **`LogFlags`** specifies information to include in the prefix of each log entry. The following [case-sensitive] flags are supported:
  * **`date`** -- date of the event in `YYYY/MM/DD` format.
  * **`time`** -- local time of the event in `HH:MM:SS` 24-hour clock format.
  * **`utc`** -- time in UTC rather than local time.
  * **`standard`** -- shorthand for `date` and `time`.
  * **`longfile`** -- long filename of the source code file that generated the event.
  * **`shortfile`** -- short filename of the source code file that generated the event.
* **`Stdout`** causes the daemon to write log entries to standard output (console) in addition to other destinations.
* **`Stderr`** causes the daemon to write log entries to standard error in addition to other destinations.
* **`Syslog`** causes the daemon to write log entries to a local or remote syslog daemon using the `Syslog` configuration settings, below.

**Syslog Settings**
```json
"Syslog": {
    "Protocol": "tcp",
    "Port": "1468",
    "Host": "localhost",
    "Tag": "cmdbd",
    "Facility": "LOG_LOCAL7",
    "Severity": "LOG_INFO"
}
```

* **`Tag`** is an arbitrary string to add to the event.
* **`Facility`** specifies the type of program that is logging the message:
  * **`LOG_KERN`** -- kernel messages
  * **`LOG_USER`** -- user-level messages
  * **`LOG_MAIL`** -- mail system
  * **`LOG_DAEMON`** -- system daemons
  * **`LOG_AUTH`** -- security/authorization messages
  * **`LOG_SYSLOG`** -- messages generated internally by syslogd
  * **`LOG_LPR`** -- line printer subsystem
  * **`LOG_NEWS`** -- network news subsystem
  * **`LOG_UUCP`** -- UUCP subsystem
  * **`LOG_CRON`** -- security/authorization messages
  * **`LOG_AUTHPRIV`** -- FTP daemon
  * **`LOG_FTP`** -- scheduling daemon
  * **`LOG_LOCAL0`** -- local use 0
  * **`LOG_LOCAL1`** -- local use 1
  * **`LOG_LOCAL2`** -- local use 2
  * **`LOG_LOCAL3`** -- local use 3
  * **`LOG_LOCAL4`** -- local use 4
  * **`LOG_LOCAL5`** -- local use 5
  * **`LOG_LOCAL6`** -- local use 6
  * **`LOG_LOCAL7`** -- local use 7
* **`Severity`** specifies the severity of the event:
  * **`LOG_EMERG`** -- system is unusable
  * **`LOG_ALERT`** -- action must be taken immediately
  * **`LOG_CRIT`** -- critical conditions
  * **`LOG_ERR`** -- error conditions
  * **`LOG_WARNING`** -- warning conditions
  * **`LOG_NOTICE`** -- normal but significant conditions
  * **`LOG_INFO`** -- informational messages
  * **`LOG_DEBUG`** -- debug-level messages

**Log Directory Settings**
```json
"LogDir": {
    "Windows": "log",
    "Linux": "/var/log/cmdbd"
}
```
* **`Windows`** is the log directory to use for Windows installations.
* **`Linux`** is the log directory to use for Linux installations.

**Global Options**
```json
"Options": {
    "Stdout": false,
    "Stderr": false,
    "Syslog": false,
    "RecoveryStack": false
}
```
* **`Stdout`** causes _all logs_ to be written to standard output; it overrides `Stdout` setting for individual logs.
* **`Stderr`** causes all logs to be written to standard error; it overrides `Stderr` setting for individual logs.
* **`Syslog`** causes all logs to be written to the configured syslog daemon; it overrides `Syslog` settings for individual logs.
* **`RecoveryStack`** enables or suppresses writing of the stack track to the error log on panic conditions.

### Startup
Once all configuration tasks are complete, the daemon can be started with the following command:
```sh
systemctl start cmdbd
```
Service access, system events, and errors are written to the following log files:
* **`system.log`** records significant, non-error events.
* **`access.log`** records client activity in Apache Combined Log Format.
* **`error.log`** records service and database errors.

The daemon can also be started from the command line. The following command-line options are available:
* **`-config`** specifies an alternate JSON configuration file; the default is `/etc/cmdbd/config.json`.
* **`-stdout`** causes _all logs_ to be written to standard output; it overrides `Stdout` setting for individual logs.
* **`-stderr`** causes _all logs_ to be written to standard error; it overrides `Stderr` setting for individual logs.
* **`-syslog`** causes _all logs_ to be written to the configured syslog daemon; it overrides `Syslog` setting for individual logs.
* **`-help`** displays the above options with a short description.

You will need to become `root` or use the `sudo` command to start the daemon or it will not be able to write to its log files. (For security reasons, the daemon should never run as `root` in production; it should always run in the context of a nonprivileged account.) Manual startup example:
```sh
[root@sysadm-dev-01 ~]# /usr/sbin/cmdbd -help
Usage of /usr/sbin/cmdbd:
  -config file
        Web server configuration file (default "/etc/cmdbd/config.json")
  -stderr
        Enable logging to stderr
  -stdout
        Enable logging to stdout
  -syslog
        Enable logging to syslog

[root@sysadm-dev-01 ~]# /usr/sbin/cmdbd -stdout
system 2017/09/30 09:55:38 main.go:62: Database "10.2.9-MariaDB" (cmdbd@localhost/gocmdb) using "mysql" driver
system 2017/09/30 09:55:38 main.go:63: Server started and listening on ":8080"
```

### API Endpoints
| Endpoint | Method | Purpose
| :------ | :------ | :------ |
| **`/usbci/checkin/{host}/{vid}/{pid}`** | POST | Submit configuration information for a new device or update information for an existing device. |
| **`/usbci/checkout/{host}/{vid}/{pid}/{sn}`** | GET | Obtain configuration information for a previously-registered, serialized device in order to perform a change audit. |
| **`/usbci/audit/{host}/{vid}/{pid}/{sn}`** | POST | Submit the results of a change audit on a serialized device. Results include the attribute name, previous value, and new value for each modified attribute.
| **`/usbci/newsn/{host}/{vid}/{pid}`** | POST | Obtain a new unique serial number from the server for assignment to the attached device. |

### API Parameters
* **`host`** is the _hostname_ of the workstation to which the device is attached.
* **`vid`** is the _vendor ID_ of the device.
* **`pid`** is the _product ID_ of the device.
* **`sn`** is the _serial number_ of the device.
