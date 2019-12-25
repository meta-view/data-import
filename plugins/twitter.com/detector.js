(function(){
    var output = 0.0;
    var markers = ["account.js", "block.js", "profile.js", "tweet.js", "like.js", "direct-message.js", "README.txt"];
    var part = 100 / markers.length;

    console.log("[" + _provider + "] Checking payload " + _payloadPath);

    function calcOutput(filename) {
        for (i in markers) {
            if (StringEndsWith(filename, markers[i])) {
                return part;
            }
        }
        return 0.0;
    }

    files = readDir()
    for (i in files) {
        output += calcOutput(files[i]);
    }
    console.log("[" + _provider + "] calculated values: " + output);
    return output;
})();