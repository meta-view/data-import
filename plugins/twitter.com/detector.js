(function(){
    console.log("\tJS\t[twitter.com] Checking payload " + payloadPath);
    files = getFiles()
    console.log("\tJS\t[twitter.com] loading " + files.length + " files");
    for (i in files) {
        console.log("\tJS\t[twitter.com] ContentType of " + files[i] + " is " + getContentType(files[i]));
    }
    return "0.50";
})();
