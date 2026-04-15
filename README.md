# (Ext)ended Golang Library

Extended Golang Library is an extension for standard golang library.

<img src="docs/assets/diablo2_hell.jpg" style="zoom:50%">

<div style="text-align: center"><b>HELL STARTS HERE</b></div>

## Conventions

* Separate module is used for each Go version that is significant for the library implementation (e.g. `go1.25/`).
* Backward compatibility is not a primary concern - anything can be refactored once a better implementation is recognized, even if it breaks compatibility. But the major version will be incremented accordingly.
* Library contains some highly experimental features which can be identified by source file path (e.g. `_exp` file name suffix) or by comments at the top of the source file.

