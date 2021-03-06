(function(){
    var account;
    var owner;
    files = listFiles()
    log("Importing payload " + files.length + " files.");

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
                content = StringReplace(getContent(file), "window.YTD.account.part0 =", "");
                account = JSON.parse(content)[0]["account"];

                content = StringReplace(getContent("profile.js"), "window.YTD.profile.part0 =", "");
                profile = JSON.parse(content)[0]["profile"];
                if(profile.avatarMediaUrl) {
                    profileImage = profile.avatarMediaUrl.replace("https://pbs.twimg.com/profile_images/", "").replace("/", "-");
                    if(!StringEndsWith(profileImage, "default_profile.png")) {
                        name = splitToLastBy(profileImage, '-');
                        imagePath = saveFile("profile_media/" + account["accountId"] + "-" + name)
                        profile["profileImage"] = imagePath;
                    } else {
                        profile["profileImage"] =  _defaultProfile;
                    }
                }

                if(profile.headerMediaUrl) {
                    backgroundImage = profile.headerMediaUrl.replace("https://pbs.twimg.com/profile_banners/", "").replace("/", "-");
                    imagePath = saveFile("profile_media/" + backgroundImage)
                    profile["backgroundImage"] = imagePath;
                }

                owner = getFileChecksum(file);
                account["profile"] = profile;
                account["provider_version"] = _version;
                data = {
                    "id": getFileChecksum(file),
                    "owner": owner,
                    "table": "accounts",
                    "provider": _provider,
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
                        log("importing conversation " + conversationId);
                        for(mi in conversation.messages) {
                            if(conversation.messages[mi].messageCreate) {
                                message = conversation.messages[mi].messageCreate;
                                checksum = getChecksum(JSON.stringify(message));
                                message["conversationId"] = conversationId;
                                message["account"] = account;
                                message["provider_version"] = _version;
                                messageData = {
                                    "id": checksum,
                                    "owner": owner,
                                    "created": message.createdAt,
                                    "table": "messages",
                                    "provider": _provider,
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
                    like["account"] = account;
                    like["provider_version"] = _version;
                    likeData = {
                        "id": checksum,
                        "owner": owner,
                        "table": "likes",
                        "provider": _provider,
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
                    if (tweet["tweet"] !== undefined) {
                        tweet = tweet["tweet"];
                    }
                    checksum = getChecksum(JSON.stringify(tweet));
                    createdDate = stringToDate(tweet.created_at);
                    created = ISODateString(createdDate);
                    tweet["account"] = account;
                    tweet["provider_version"] = _version;
                    tweetData = {
                        "id": checksum,
                        "owner": owner,
                        "created": created,
                        "table": "posts",
                        "provider": _provider,
                        "name": "tweet-" + tweet.id,
                        "content-type": "application/json",
                        "content": tweet
                    }
                    saveEntry(tweetData);
                    if(tweet.geo) {
                        geo = tweet.geo;
                        geo["tweet_id"] = checksum;
                        geo["tweet_name"] = "tweet-" + tweet.id;
                        geo["provider_version"] = _version;
                        geoData = {
                            "id": getChecksum(JSON.stringify(hashTag)),
                            "owner": owner,
                            "created": created,
                            "table": "locations",
                            "provider": _provider,
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
                        hashTag["provider_version"] = _version;
                        hashTagData = {
                            "id": getChecksum(JSON.stringify(hashTag)),
                            "owner": owner,
                            "created": created,
                            "table": "tags",
                            "provider": _provider,
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
                        mention["provider_version"] = _version;
                        mentionData = {
                            "id": getChecksum(JSON.stringify(mention)),
                            "owner": owner,
                            "created": created,
                            "table": "mentions",
                            "provider": _provider,
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
                        "owner": owner,
                        "table": "images",
                        "provider": _provider,
                        "file": file,
                        "name": name,
                        "content-type": contentType,
                        "filePath": saveFile(file)
                    }
                } else if (StringStartsWith(contentType, "video")) {
                    name = splitToLastBy(file, '-');
                    data = {
                        "id": checksum,
                        "owner": owner,
                        "table": "videos",
                        "provider": _provider,
                        "file": file,
                        "name": name,
                        "content-type": contentType,
                        "filePath": saveFile(file)
                    }
                } else {
                    data = {
                        "id": checksum,
                        "owner": owner,
                        "table": "files",
                        "provider": _provider,
                        "name": file,
                        "content-type": contentType,
                        "filePath": saveFile(file)
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
                log("linking mediaFiles to tweets.");
                for (i in tweets) {
                    tweet = tweets[i];
                    if (tweet["tweet"] !== undefined) {
                        tweet = tweet["tweet"];
                    }
                    if(tweet.entities && tweet.entities.media) {
                        linkMedia(tweet.entities.media, tweet, account);
                    } else if(tweet.extended_entities && tweet.extended_entities.media) {
                        linkMedia(tweet.extended_entities.media, tweet, account);
                    }
                }
            break;
        }
    }

    function linkMedia(media, tweet, account) {
        checksum = getChecksum(JSON.stringify(tweet));
        createdDate = stringToDate(tweet.created_at);
        created = ISODateString(createdDate);
        for(mi in media) {
            query = {};
            mediaFile = tweet.extended_entities.media[mi];
            if(mediaFile.video_info && mediaFile.video_info.variants) {
                variant = mediaFile.video_info.variants[0];
                filename = splitToLastBy(variant.url,"/");
                query = {"table":"videos", "content": {"name": filename}};
            } else {
                filename = splitToLastBy(mediaFile.media_url,"/");
                query = {"table":"images", "content": {"name": filename}};
            }
            elements = readEntry(query);
            for (i in elements) {
                element = elements[i];
                var content = JSON.parse(element["content"]);
                content["media"] = mediaFile;
                content["account"] = account;
                content["created"] = created;
                content["post_id"] = checksum;
                content["favorite_count"] = tweet.favorite_count;
                content["retweet_count"] = tweet.retweet_count;
                saveEntry(content);
            }
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