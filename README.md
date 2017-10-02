# CMDBc
The _**Configuration Management Database Client**_ is a utility that manages information about devices attached to end-user workstations and reports that information to a centralized repository over a RESTful JSON API provided by the complementary server component, the _**Configuration Management Database Daemon**_ or [**CMDBd**](https://github.com/jscherff/cmdbd/blob/master/README.md). **CMDBc** can register or _"check-in"_ attached devices with the server, obtain unique serial numbers from the server for devices that support serial number configuration, perform audits against previous device configurations, and report configuration changes found during the audit to the server for logging and analysis.

### System Requirements
**CMDBc** is written in **Go** and can be compiled for any operating system and architecture. It is intended to be installed on end-user workstations running **Microsoft Windows 7** or higher and should be invoked by a centralized management solution like **IBM BigFix**.

### Installation
To implement **CMDBc** quickly with minimal configuration, please see the [**QuickStart** document.](https://github.com/jscherff/cmdbc/blob/master/QUICKSTART.md)

Pre-compiled binaries are also available for both 32- and 64-bit Windows systems and can be installed in any folder along with the required JSON configuration file:

* [**`cmdbc.exe`**](https://github.com/jscherff/cmdbc/raw/master/i686/cmdbc.exe) (32-bit Windows 7 or higher)
* [**`cmdbc.exe`**](https://github.com/jscherff/cmdbc/raw/master/x86_64/cmdbc.exe) (64-bit Windows 7 or higher)
* [**`config.json`**](https://github.com/jscherff/cmdbc/raw/master/config.json) (Configuration file)

### Configuration
The JSON configuration file, [`config.json`](https://github.com/jscherff/cmdbd/blob/master/config.json), is mostly self-explanatory. The default settings are sane and you should not have to change them in most use cases.

#### Server Settings
Parameters for communicating with the **CMDBd** server:
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

#### Path Settings
Directories where logs, state files, and reports will be written:
```json
"Paths": {
    "LogDir": "log",
    "StateDir": "state",
    "ReportDir": "report"
}
```
* **`LogDir`** is the directory where log files are written. When a relative path like `log` is provided, the directory will be created below the application directory. 
* **`StateDir`** is where device state files are stored. State files are used in performing local audits.
* **`ReportDir`** is where device reports are written.

#### File Settings
Filenames for logs and the legacy report file:
```json
"Files": {
    "SystemLog": "system.log",
    "ChangeLog": "change.log",
    "ErrorLog": "error.log",
    "Legacy": "usb_serial.txt"
}
```
* **`SystemLog`** is the name of the file where **CMDBc** records significant, non-error events.
* **`ChangeLog`** is the name of the file where **CMDBc** records changes found during audits. (It also reports changes to the **CMDBd** server.)
* **`ErrorLog`** is the name of the file where **CMDBc** records errors.
* **`Legacy`** is the name of the file where **CMDBc** writes the legacy inventory report.

#### Logging Settings
Granular logging options for the system, change, and error log:
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

#### Syslog Settings
Parameters for communicating with a local or remote syslog server:
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

#### Include Settings
Vendors and products to include (_true_) or exclude (_false_) when inventorying devices:
```json
"Include": {
    "VendorID": {
        "043d": false,
        "045e": false
    },
    "ProductID": {
        "0801": {
            "0001": true,
            "0002": true,
            "0011": true,
            "0012": true,
            "0013": true
        },
        "0acd": {
            "2030": true
        }
    },
    "Default": true
},
```
* **`VendorID`** specifies which vendors to include or exclude. This setting applies to all of the vendor's products and overrides both the _ProductID_ and _Default_ configuration settings; that is, if a vendor is excluded under _VendorID_, that vendor's products cannot be included under _ProductID_. Here, all devices with **Microsoft** _Vendor IDs_ `043d` and `045e` will be excluded.
* **`ProductID`** specifies which products to include or exclude. This setting applies to specific _ProductIDs_ under a given _VendorID_ and overrides the _Default_ configuration setting. Here, **MagTek** (_VendorID_ `0801`) card readers with  _ProductIDs_ `0001`, `0002`, `0011`, `0012`, and `0013` will be included, as will **ID TECH** (_VendorID_ `0acd`) Card 
* **`Default`** specifies the default behavior for products that are not specifically included or excluded by _Vendor ID_ or _Product ID_. Here the default is to include, which effectively renders previous inclusions redundant; however, specific _VendorID_ and _ProductID_ inclusions ensure that those devices will be inventoried even if the _Default_ setting is changed to 'exclude' (_false_).

#### Format Settings
Default file formats for various use cases:
```json
"Format": {
    "Report": "csv",
    "Default": "json"
}
```
* **`Report`** is the default output format for inventory reports.
* **`Default`** is the default output format for other use cases.

### Command-Line Flags
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

### Serial Number Configuration
Configure serial numbers on attached devices with the `serial` _action flag_.

The `set`, `copy`, and `fetch` _option flags_ are mutually-exclusive. You assign a specific serial number string with the `set` _option flag_, copy the immutable factory serial number (if one exists) to the configurable serial number with the `copy` _option flag_, or request a new, unique serial number from the server with the `fetch` _option flag_.

The `copy`, `fetch`, and `set` _option flags_ can each be combined with `erase` and `force`. By default, **CMDBc** ignores serial number changes for devices that already have serial numbers. The `erase` _option flag_ bypasses this by erasing the existing serial number before attempting to assign a new one, effectively removing the constraint. The `force` _option flag_ simply overrides the safeguard feature.

**Examples**:
```sh
cmdbc.exe -serial -fetch -force
```
The preceding command will, for each compatible device, fetch a new serial number from the server and configure the device with it, overriding the safety mechanism that prevents overwriting existing serial numbers.
```sh
cmdbc.exe -serial -erase -fetch
```
The preceding command will, for each compatible device, erase the existing serial number, fetch a new, unique serial number from the server, and configure the device with it.

While the previous two examples would normally produce the same result, a subtle difference is that, if **CMDBc** were unable to obtain a new serial number, `force` would leave existing serial numbers in place whereas `erase` would leave devices without serial numbers.

You can also use the `erase` _option flag_ by itself to erase device serial numbers, although this is an unusual use case.

**Caution**: action and option flags apply to _all attached devices_; if you use the `serial` _action flag_ with the `fetch` _option flag_, **CMDBc** will only configure new serial numbers on compatible devices that don't already have one. If all attached devices already have serial numbers or are not configurable, nothing will happen. However, if you add the `force` flag, it will overwrite the serial number on all compatible devices -- even those that already have a serial number. If you use the `set` and `force` _option flags_ and there is more than one configurable device attached, you will end up having multiple devices with the same serial number.

Refer to the _Database Structure_ section in the documentation for [**CMDBd**](https://github.com/jscherff/cmdbd/blob/master/README.md) for details on device information transferred to the server and tables/columns affected on serial number requests.
 
### Device Registration
Register attached devices with the server using the `checkin` _action flag_. This will create a new object in the device repository.

Refer to the _Database Structure_ section in the documentation for [**CMDBd**](https://github.com/jscherff/cmdbd/blob/master/README.md) for details on device information transferred to the server and tables/columns affected on device registrations.

### Device Audits
Perform a configuration change audit for attached devices using the `audit` _action flag._

You can audit against device state files saved on the local workstation with the `local` _option flag_, or you can audit against device information stored in the database with the `server` _option flag_. The latter is preferred. By default, device state for local audits is stored in JSON files in the `state` subdirectory under the utility installation directory (configurable). Changes detected during an audit are written to the local change log and are also reported to the server.

Audits are only supported on serialized devices.

Refer to the _Database Structure_ section in the documentation for [**CMDBd**](https://github.com/jscherff/cmdbd/blob/master/README.md) for details on device information transferred to the server and tables/columns affected on device audits.

### Device Reports
Generate device reports for attached devices using the `report` _action flag._

Select the report format with the `format` _option flag_. Four formats are currently supported: _comma-separated value_ (CSV), _name-value pairs_ (NVP), _extensible markup language_ (XML), and _JavaScript object notation_ (JSON). 

By default, report files are written to the `report` subdirectory under the utility installation directory (configurable). A separate report file is generated for each device. The report filename is `{bn}-{ba}-{pn}-{vid}-{pid}.{fmt}`, where
* `bn` is a three-digit decimal value representing _bus number_,
* `ba` is a three-digit decimal value representing _bus address_,
* `pn` is a three-digit decimal value representing _port number_,
* `vid` is a four-digit hexadecimal value representing _vendor ID_,
* `pid` is a four-digit hexadecimal value representing _product ID_, and
* `fmt` is the report format (csv, nvp, xml, or json)

Change the report destination folder with the `folder` _option flag_.

Write the report to the console with the `console` _option flag_. If you use the `console` _option flag_ without `folder`, the report will only be written to the console. If you use the `console` _option flag_ after `folder`, the report will be written to the specified folder _and_ to the console. If you use the `console` _option flag_ before `folder`, the report will only be written to the console and `folder` will be ignored.

**Examples**:
```sh
cmdbc.exe -report -format csv
```
The preceding command writes the device reports in CSV format to the 'reports' subdirectory.
```sh
cmdbc.exe -report -format json -console
cmdbc.exe -report -format json -console -folder c:\reports
```
Both of the preceding commands write the device reports in JSON format to the console. The `folder` _option flag_ in the second command is ignored.
```sh
cmdbc.exe -report -format xml -folder c:\reports
cmdbc.exe -report -format xml -folder c:\reports -console
```
Both of the preceding commands write the device reports in XML format to the c:\reports folder. The second command also writes the reports to the console.

### Device Resets
Reset attached devices using the `reset` _action flag_.

Depending on the device, this either does a host-side reset, refreshing the USB device descriptor, or a low-level hardware reset on the device.

### Legacy Reports
Write a legacy device report using the `legacy` _action flag_.

This feature mimics the behavior of previous device inventory utilities for integration backward compatibility. It simply writes the hostname and device serial number in CSV format to a file named `usb_serials.txt` in the utility installation directory, then exits. It filters all but MagTek card readers, and if there is more than one card reader attached, it arbitrarily chooses one.

Renaming the utility from **cmdbd.exe** to **magtek_inventory.exe** forces this behavior without command-line flags.
