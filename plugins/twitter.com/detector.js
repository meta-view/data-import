(function(){
    console.log("\tJS\t[twitter.com] Checking payload " + payloadPath);
    files = getFiles()
    console.log("loading " + files.length + " files");
    for (i in files) {
        console.log("ContentType of " + files[i] + " is " + getContentType(files[i]));
    }
    return "0.50";
})();
