(function(){
    var provider = "twitter.com";
    var output = 0.0;
    var markers = ["tweet_media", "account.js", "block.js", "profile.js", "tweet.js", "like.js", "direct-message.js", "README.txt"];
    var part = 100 / markers.length;

    console.log("[" + provider + "] Checking payload " + _payloadPath);

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
    console.log("[" + provider + "] calculated values: " + output);
    return output > 100 ? 100 : output;
})();