(function(){
    console.log("[twitter.com] Import payload " + _payloadPath);
    files = listFiles()
    console.log("[twitter.com] loading " + files.length + " files");
    for (i in files) {
        checksum = getSha1Checksum(files[i]);
        data = {
            "id": checksum,
            "table": "files",
            "provider": _provider,
            "name": files[i],
            "content-type": getContentType(files[i]), 
            "content": getContent(files[i])
        }
        saveEntry(data)
    }
})();