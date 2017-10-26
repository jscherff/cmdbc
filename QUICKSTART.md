# CMDBc QuickStart
The _**Configuration Management Database Client**_ is a utility that manages information about devices attached to end-user workstations and reports that information to a centralized repository over a RESTful JSON API provided by the complementary server component, the _**Configuration Management Database Daemon**_. See the [**CMDBc README.md**](https://github.com/jscherff/cmdbc/blob/master/README.md) and [**CMDBd README.md**](https://github.com/jscherff/cmdbd/blob/master/README.md) documents for more detail.

### System Requirements
**CMDBc** runs on end-user workstations running **Microsoft Windows 7** or higher and should be invoked by a centralized management solution like **IBM BigFix**.

### Installation
Save the appropriate executable file and the JSON configuration file to the desired installation folder, such as C:\CMDBc.

* [**`cmdbc.exe`**](https://github.com/jscherff/cmdbc/raw/master/i686/cmdbc.exe) (32-bit Windows 7 or higher)
* [**`cmdbc.exe`**](https://github.com/jscherff/cmdbc/raw/master/x86_64/cmdbc.exe) (64-bit Windows 7 or higher)
* [**`config.json`**](https://github.com/jscherff/cmdbc/raw/master/config.json) (Configuration file)

### Configuration
1. Obtain the hostname (or IP address) and listener port of the server.
1. In the **API** section of the JSON configuration file, change the hostname and port components of the **Server** to the correct values.
    ```json
    "API": {
        "Server": "http://cmdbsvcs-prd-01.24hourfit.com:8080",
        "Endpoint": {
        }
    }
    ```
1. Ensure a firewall rule is in place allowing communication from managed workstations to the IP address and port of the server.

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
    cmdbc.exe -audit
    ```
1. Periodically parse the contents of `error.log` for issues. This file is located in the "log" subdirectory beneath the installation folder.