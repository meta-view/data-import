# MetaView Service

This project has the goal to support each an every user in analysing and accessing the data dumps provided by internet services since the introduction of DSGVO and so on.

The user should be able to:

1. Find a direct entry on how to reveice the data.
2. Being able to upload and store the data dump itself.
3. Being able to perform various analysis with the saved data.
4. … let's see :-).


## Setup & Run

This application is written in Go. So the easiest way of running it is with a simple

```
% go run main.go
```

This should lead to the following log output:

```
2019/12/18 22:29:36 templates loading successful
2019/12/18 22:29:36 buffer allocation successful
2019/12/18 22:29:36 loading map[download_request:https://twitter.com/settings/your_twitter_data name:twitter.com version:0.0.1] from plugins/twitter.com
2019/12/18 22:29:36 Serving Application on port 9000
```

You can now go to `http://localhost:9000` with your browser.  
You should be greeted with an upload Page.  

If you upload an archive, you will see the following log output:

```
2019/12/18 22:29:51 uploaded: map[Content-Disposition:[form-data; name="files[]"; filename="twitter-2019-12-17-a7688c22d2722121907d0e818259c7ad0588752c3a3a6ba39c46f49031a6c42c.zip"] Content-Type:[application/zip]]
2019/12/18 22:29:51 importing data/zip/twitter-2019-12-17-a7688c22d2722121907d0e818259c7ad0588752c3a3a6ba39c46f49031a6c42c.zip to data/raw/twitter-2019-12-17-a7688c22d2722121907d0e818259c7ad0588752c3a3a6ba39c46f49031a6c42c
	extracting file data/raw/twitter-2019-12-17-a7688c22d2722121907d0e818259c7ad0588752c3a3a6ba39c46f49031a6c42c/README.txt
…
	extracting file data/raw/twitter-2019-12-17-a7688c22d2722121907d0e818259c7ad0588752c3a3a6ba39c46f49031a6c42c/tweet_media/1078942711678423040-DNooJU2U8AI3KUg.png
	extracting file data/raw/twitter-2019-12-17-a7688c22d2722121907d0e818259c7ad0588752c3a3a6ba39c46f49031a6c42c/verified.js
2019/12/18 22:29:51 Detect if data/raw/twitter-2019-12-17-a7688c22d2722121907d0e818259c7ad0588752c3a3a6ba39c46f49031a6c42c is for [twitter.com]
	JS	[twitter.com] Checking payload data/raw/twitter-2019-12-17-a7688c22d2722121907d0e818259c7ad0588752c3a3a6ba39c46f49031a6c42c
```

A list of Download URLs for the service that are planed to be supported can be find [here](providers.json).

## Todo

* creating a basic Application: ✅
* creating a basic abstraction for plugins: ✅
* creating Desktop Apps for Mac/Linux/Windows: ❌
* adding a wizard for step by step import of user data: ❌
* Creating analysis function (time-based, geographical, structural): ❌
* … (feel free to leave an idea in the issues)

