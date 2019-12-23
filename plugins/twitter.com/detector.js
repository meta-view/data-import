(function(){
    output = 0.0;
    console.log("[twitter.com] Checking payload " + payloadPath);
    files = getFiles()
    for (i in files) {
        console.log("scanning " + files[i])
    }
    return output.toFixed(2);
})();
