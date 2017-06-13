var binwrap = require("binwrap");
var path = require("path");

var packageInfo = require(path.join(__dirname, "package.json"));
var version = packageInfo.version;
var root = "https://github.com/indiereign/shift72-kibble/releases/download";

module.exports = binwrap({
  binaries: ["kibble"],
  urls: {
    "darwin-x64": root + "/v" + version + "/kibble_" + version + "_darwin_amd64.tar.gz",
    "linux-x64": root + "/v" + version + "/kibble_" + version + "_linux_amd64.tar.gz",
    "win32-x64": root + "/v" + version + "/kibble_" + version + "_windows_amd64.tar.gz"
  }
});