(function(){
    console.log("[twitter.com] Import payload " + _payloadPath);
    files = listFiles()
    console.log("[twitter.com] loading " + files.length + " files");
    for (i in files) {
        data = getData(files[i]);
        saveEntry(data);
    }

    function getData(file) {
        checksum = getSha1Checksum(file);
        contentType = getContentType(file);
        switch (file) {
            case "account.js":
                content = StringReplace(getContent(file), "window.YTD.account.part0 =", "");
                return {
                    "id": checksum,
                    "table": "accounts",
                    "provider": _provider,
                    "name": file,
                    "content-type": contentType,
                    "content": JSON.parse(content)
                }
            default:
                return {
                    "id": checksum,
                    "table": "files",
                    "provider": _provider,
                    "name": file,
                    "content-type": contentType,
                    "content": getBase64(file)
                }
        }
    }
})();