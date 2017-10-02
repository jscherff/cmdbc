# CMDBc QuickStart
The _**Configuration Management Database Client**_ is a utility that manages information about devices attached to end-user workstations and reports that information to a centralized repository over a RESTful JSON API provided by the complementary server component, the _**Configuration Management Database Daemon**_ or [**CMDBd**](https://github.com/jscherff/cmdbd/blob/master/README.md). **CMDBc** can register or _"check-in"_ attached devices with the server, obtain unique serial numbers from the server for devices that support serial number configuration, perform audits against previous device configurations, and report configuration changes found during the audit to the server for logging and analysis.

### System Requirements
**CMDBc** is written in **Go** and can be compiled for any operating system and architecture. It is intended to be installed on end-user workstations running **Microsoft Windows 7** or higher and should be invoked by a centralized management solution like **IBM BigFix**.

### Installation
Save the appropriate executable file and the JSON configuration file to the desired installation folder, such as C:\CMDBc

* [**`cmdbc.exe`**](https://github.com/jscherff/cmdbc/raw/master/i686/cmdbc.exe) (32-bit Windows 7 or higher)
* [**`cmdbc.exe`**](https://github.com/jscherff/cmdbc/raw/master/x86_64/cmdbc.exe) (64-bit Windows 7 or higher)
* [**`config.json`**](https://raw.githubusercontent.com/jscherff/cmdbc/master/config.json)

### Configuration
In the JSON configuration file, find the **Server** section and change the **URL** setting to the URL of the **CMDBd** server. Do not modify any other settings.
```json
"Server": {
    "URL": "http://sysadm-dev-01.24hourfit.com:8080",
    "CheckinPath": "usbci/checkin",
    "CheckoutPath": "usbci/checkout",
    "NewSNPath": "usbci/newsn",
    "AuditPath": "usbci/audit"
}
```
### Operation
Using an _enterprise endpoint managment solution_ like **IBM BigFix**:
1. Schedule the following command to run once per month initially, then once per quarter or as necessary:
    ```sh
    cmdbc.exe -serial -fetch
    ```
1. Schedule the following command to run once per week:
    ```sh
    cmdbc.exe -checkin
    ```
1. Schedule the following command to run once per month:
    ```sh
    cmdbc.exe -audit -server
    ```
1. If legacy operation is required, schedule the following command to run once per week:
    ```sh
    cmdbc.exe -legacy
    ```
1. Periodically parse the contents of `error.log` for issues.
