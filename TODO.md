General Enhancements
--------------------
- [ ] Add option to inventory devices that could not be opened.
	* For devices that cannot be opened, enable only audit capability.
- [ ] Rethink how devices that cannot be opened are instantiated.
	* Force underlying device to be nil and add nil checks on all methods?
	* Create a special, non-device type/package (like 'nil' or 'null') with no methods?
