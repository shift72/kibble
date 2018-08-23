var binwrap = require("binwrap");
var path = require("path");

var packageInfo = require(path.join(__dirname, "package.json"));
var version = packageInfo.version;
var root = "https://s3-ap-southeast-2.amazonaws.com/shift72-sites/s72-web/kibble/" + version;

module.exports = binwrap({
  dirname: __dirname,
  binaries: ["kibble"],
  urls: {
    "darwin-x64": root + "/shift72-kibble_" + version + "_macOS_64-bit.zip",
    "linux-x64":  root + "/shift72-kibble_" + version + "_Tux_64-bit.tar.gz",
    "win32-x64":  root + "/shift72-kibble_" + version + "_windows_64-bit.zip"
  }
});
