# MetaView Service

This project has the goal to support each an every user in analysing and accessing the data dumps provided by internet services since the introduction of DSGVO and so on.

The user should be able to:

1. Find a direct entry on how to reveice the data.
2. Being able to upload and store the data dump itself.
3. Being able to perform various analysis with the saved data.
4. … let's see :-).

## RUN

You can use the `run-docker.sh` script to run the latest build of this Application within docker.
You can find a nightly Docker Build on Docker Hub: https://hub.docker.com/r/phaus/meta-view-service

## Build & Run

If you want to run from the source, you can do that as well!

This application is written in Go. So the easiest way of running it is with the `run.sh` script.

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

## Plugins

A Plugin should run in a sandbox.  
Every Plugin runs in its own JS-VM (based on [Otto](https://github.com/robertkrimen/otto)).
There are some functions injected into the VM Context to help the plugin itself to deal with the surroundings (see _Plugin API_)

For now the plugin has to fulfill three different use-cases:

* `detector.js` - detects if a payload can be handled by that specific plugin. The result is a percentage value (0.0 - 100.00).
* `importer.js` - the importer handles the conversion of the plain payload data into a normalized datastructure for the Database - for now that is (id,crated,update,json).
* `presenter.js` - the presenter will display a specific entry or a list of entries in a given context (list->detail, graph).

As you can also see, the default plugin `detector.js` is able to access and identify files and folder from within the uploaded directory:

```
loading 64 files
ContentType of README.txt is text/plain; charset=utf-8
ContentType of account-creation-ip.js is application/octet-stream
…
ContentType of tweet.js is text/plain; charset=utf-8
…
ContentType of tweet_media/368453822730862593-BR0CzLSCAAAfv1l.png is image/png
ContentType of tweet_media/878731487729811457-DBTIKevUQAAURjC.jpg is image/jpeg
ContentType of tweet_media/899359035618623489-DHr-Q58WsAILnJZ.jpg is image/jpeg
ContentType of tweet_media/949981183256989697-DS8CkUzV4AA8fWh.jpg is image/jpeg
ContentType of verified.js is application/octet-stream
```

## Plugin API

At the moment the plugin JS API has two methods on hand: 

* `getFiles()` - a method that lists all files rooted in the payload folder.
* `getContentType(…)` - a method to determine the content type of the given file.

You can see an example uses in the file [detector.js](plugins/twitter.com/detector.js) of the twitter plugin.

A list of Download URLs for the service that are planed to be supported can be find in the file [providers.json](providers.json).

## Todo

* creating a basic Application: ✅
* creating a basic abstraction for plugins: ✅
* creating Desktop Apps for Mac/Linux/Windows: ❌
* adding a wizard for step by step import of user data: ❌
* Creating analysis function (time-based, geographical, structural): ❌
* … (feel free to leave an idea in the issues)
