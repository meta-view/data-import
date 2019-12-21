(function(){
    console.log("[twitter.com] Import payload " + payloadPath);
    files = getFiles()
    console.log("[twitter.com] loading " + files.length + " files");
    for (i in files) {
        data = {
            "table": "files",
            "provider":"[twitter.com]",
            "name": files[i],
            "content-type": getContentType(files[i]), 
            "content": getContent(files[i])
        }
        saveEntry(data)
    }
})();