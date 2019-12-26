(function(){
    console.log("[twitter.com] Import payload " + _payloadPath);
    files = listFiles()
    console.log("[twitter.com] loading " + files.length + " files");
    for (i in files) {
        saveData(files[i]);
    }

    function saveData(file) {
        switch (file) {
            case "account.js":
                content = StringReplace(getContent(file), "window.YTD.account.part0 =", "");
                data = {
                    "id": getFileChecksum(file),
                    "table": "accounts",
                    "name": file,
                    "content-type": "application/json",
                    "content": JSON.parse(content)
                }
                saveEntry(data);
                break;
            case "tweet.js":
                content = StringReplace(getContent(file), "window.YTD.tweet.part0 = ", "");
                tweets = JSON.parse(content);
                for (i in tweets) {
                    tweet = tweets[i];
                    checksum = getChecksum(JSON.stringify(tweet));
                    createdDate = stringToDate(tweet.created_at);
                    created = ISODateString(createdDate);
                    tweetData = {
                        "id": checksum,
                        "created": created,
                        "table": "posts",
                        "name": "tweet-" + tweet.id,
                        "content-type": "application/json",
                        "content": tweet
                    }
                    saveEntry(tweetData);
                    for(ti in tweet.entities.hashtags) {
                        hashTag = tweet.entities.hashtags[ti];
                        hashTag["tweet_id"] = checksum;
                        hashTag["tweet_name"] = "tweet-" + tweet.id;
                        hashTagData = {
                            "id": getChecksum(JSON.stringify(hashTag)),
                            "created": created,
                            "table": "tags",
                            "name": hashTag.text,
                            "content-type": "application/json",
                            "content": hashTag
                        }
                        saveEntry(hashTagData);
                    }
                    for(mi in tweet.entities.user_mentions) {
                        mention = tweet.entities.user_mentions[mi];
                        mention["tweet_id"] = checksum;
                        mention["tweet_name"] = "tweet-" + tweet.id;
                        mentionData = {
                            "id": getChecksum(JSON.stringify(mention)),
                            "created": created,
                            "table": "mentions",
                            "name": mention.screen_name,
                            "content-type": "application/json",
                            "content": mention
                        }
                        saveEntry(mentionData);
                    }
                }
            default:
                contentType = getContentType(file);
                checksum = getFileChecksum(file);
                if (StringStartsWith(contentType, "image")) {
                    name = splitToLastBy(file, '-');
                    data = {
                        "id": checksum,
                        "table": "images",
                        "file": file,
                        "name": name,
                        "content-type": contentType,
                        "content": getBase64(file)
                    }
                } else {
                    data = {
                        "id": checksum,
                        "table": "files",
                        "name": file,
                        "content-type": contentType,
                        "content": getBase64(file)
                    }
                }
                saveEntry(data);
        }
    }

    function splitToLastBy(element, divider) {
        parts = element.split(divider);
        return parts[parts.length - 1];
    }

    // see https://stackoverflow.com/a/13133124
    function stringToDate(s) {
        var b = s.split(/[: ]/g);
        var m = {jan:0, feb:1, mar:2, apr:3, may:4, jun:5, jul:6,
                    aug:7, sep:8, oct:9, nov:10, dec:11};
        return new Date(Date.UTC(b[7], m[b[1].toLowerCase()], b[2], b[3], b[4], b[5]));
    }

    // see https://stackoverflow.com/a/7244288
    function ISODateString(d){
        function pad(n){return n<10 ? '0'+n : n}
        return d.getUTCFullYear()+'-'
             + pad(d.getUTCMonth()+1)+'-'
             + pad(d.getUTCDate())+'T'
             + pad(d.getUTCHours())+':'
             + pad(d.getUTCMinutes())+':'
             + pad(d.getUTCSeconds())+'Z'}
})();