var binwrap = require("binwrap");
var path = require("path");

var packageInfo = require(path.join(__dirname, "package.json"));
var version = packageInfo.version;
var root = "https://shift72-sites.s3.amazonaws.com/s72-web/kibble/" + version;

module.exports = binwrap({
  dirname: __dirname,
  binaries: ["kibble"],
  urls: {
    "darwin-x64": root + "/kibble_" + version + "_macOS_64-bit.zip",
    "darwin-arm64": root + "/kibble_" + version + "_macOS_arm64-bit.zip",
    "linux-x64":  root + "/kibble_" + version + "_Tux_64-bit.tar.gz",
    "win32-x64":  root + "/kibble_" + version + "_windows_64-bit.zip"
  }
});
