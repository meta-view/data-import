(function(){
    var provider = "twitter.com";
    var account;
    files = listFiles()
    console.log("[" + provider + "] Importing payload " + files.length + " files");

    // Importing files into the database
    for (i in files) {
        saveData(files[i]);
    }

    // Linking files, and adding additional info defails.
    for (i in files) {
        linkFiles(files[i]);
    }

    function saveData(file) {
        switch (file) {
            case "account.js":
            content = StringReplace(getContent("profile.js"), "window.YTD.profile.part0 =", "");
            profile = JSON.parse(content)[0]["profile"];
            if(profile.avatarMediaUrl) {
                profileImage = profile.avatarMediaUrl.replace("https://pbs.twimg.com/profile_images/", "").replace("/", "-");
                if(StringEndsWith(profileImage, "default_profile.png")) {
                    profile["profileImage"] = "default_profile.png"
                } else {
                    profile["profileImage"] = profileImage;
                }
            }
            if(profile.headerMediaUrl) {
                backgroundImage = profile.headerMediaUrl.replace("https://pbs.twimg.com/profile_banners/", "").replace("/", "-");
                profile["backgroundImage"] = backgroundImage;
            }
            content = StringReplace(getContent(file), "window.YTD.account.part0 =", "");
                account = JSON.parse(content)[0]["account"];
                account["profile"] = profile;
                data = {
                    "id": getFileChecksum(file),
                    "table": "accounts",
                    "provider": provider,
                    "name": file,
                    "content-type": "application/json",
                    "content": account
                }
                saveEntry(data);
                break;
            case "direct-message.js":
                content = StringReplace(getContent(file), "window.YTD.direct_message.part0 = ", "");
                conversations = JSON.parse(content);
                for (i in conversations) {
                    if(conversations[i].dmConversation) {
                        conversation = conversations[i].dmConversation;
                        conversationId = conversation.conversationId;
                        console.log("importing conversation " + conversationId);
                        for(mi in conversation.messages) {
                            if(conversation.messages[mi].messageCreate) {
                                message = conversation.messages[mi].messageCreate;
                                checksum = getChecksum(JSON.stringify(message));
                                message["conversationId"] = conversationId;
                                messageData = {
                                    "id": checksum,
                                    "created": message.createdAt,
                                    "table": "messages",
                                    "provider": provider,
                                    "name": "message-" + message.id,
                                    "content-type": "application/json",
                                    "content": message
                                }
                                saveEntry(messageData);
                            }
                        }
                    }
                }
                break;
            case "like.js":
                content = StringReplace(getContent(file), "window.YTD.like.part0 = ", "");
                likes = JSON.parse(content);
                for (i in likes) {
                    like = likes[i];
                    checksum = getChecksum(JSON.stringify(like));
                    likeData = {
                        "id": checksum,
                        "table": "likes",
                        "provider": provider,
                        "name": "like-"+like.tweetId,
                        "content-type": "application/json",
                        "content": like
                    }
                    saveEntry(likeData);
                }
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
                        "provider": provider,
                        "name": "tweet-" + tweet.id,
                        "content-type": "application/json",
                        "content": tweet
                    }
                    saveEntry(tweetData);
                    if(tweet.geo) {
                        geo = tweet.geo;
                        geo["tweet_id"] = checksum;
                        geo["tweet_name"] = "tweet-" + tweet.id;
                        geoData = {
                            "id": getChecksum(JSON.stringify(hashTag)),
                            "created": created,
                            "table": "locations",
                            "provider": provider,
                            "name": "geo-" + tweet.id,
                            "content-type": "application/json",
                            "content": geo
                        }
                        saveEntry(geoData);
                    }
                    for(ti in tweet.entities.hashtags) {
                        hashTag = tweet.entities.hashtags[ti];
                        hashTag["tweet_id"] = checksum;
                        hashTag["tweet_name"] = "tweet-" + tweet.id;
                        hashTagData = {
                            "id": getChecksum(JSON.stringify(hashTag)),
                            "created": created,
                            "table": "tags",
                            "provider": provider,
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
                            "provider": provider,
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
                        "provider": provider,
                        "file": file,
                        "name": name,
                        "content-type": contentType,
                        "content": getBase64(file)
                    }
                } else if (StringStartsWith(contentType, "video")) {
                    name = splitToLastBy(file, '-');
                    data = {
                        "id": checksum,
                        "table": "videos",
                        "provider": provider,
                        "file": file,
                        "name": name,
                        "content-type": contentType,
                        "content": getBase64(file)
                    }
                } else {
                    data = {
                        "id": checksum,
                        "table": "files",
                        "provider": provider,
                        "name": file,
                        "content-type": contentType,
                        "content": getBase64(file)
                    }
                }
                saveEntry(data);
        }
    }

    function linkFiles(file) {
        switch (file) {
            case "tweet.js":
                content = StringReplace(getContent(file), "window.YTD.tweet.part0 = ", "");
                tweets = JSON.parse(content);
                console.log("linking images to tweets.");
                for (i in tweets) {
                    tweet = tweets[i];
                    checksum = getChecksum(JSON.stringify(tweet));
                    createdDate = stringToDate(tweet.created_at);
                    created = ISODateString(createdDate);
                    if(tweet.extended_entities && tweet.extended_entities.media) {
                        for(mi in tweet.extended_entities.media) {
                            mediaFile = tweet.extended_entities.media[mi];
                            filename = mediaFile.media_url.replace("http://pbs.twimg.com/media/", "");
                            query = {"table":"images", "content": {"name": filename}};
                            images = readEntry(query);
                            for (i in images) {
                                image = images[i];
                                var content = JSON.parse(image["content"]);
                                content["media"] = mediaFile;
                                content["account"] = account;
                                content["created"] = created;
                                content["favorite_count"] = tweet.favorite_count;
                                content["retweet_count"] = tweet.retweet_count;
                                saveEntry(content);
                            }
                        }
                    }
                }
            break;
        }
    }
    function splitToLastBy(element, divider) {
        parts = element.split(divider);
        return parts[parts.length - 1];
    }

    // see https://stackoverflow.com/a/13133124
    function stringToDate(s) {
        var b = s.split(/[: ]/g);
        var m = {   jan:0, feb:1, mar:2, apr:3, may:4, jun:5, 
                    jul:6, aug:7, sep:8, oct:9, nov:10, dec:11};
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