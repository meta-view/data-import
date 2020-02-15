(function(){
    if(render == "row")
        return renderRow(entry);
    else if(render == "node")
        return renderNode(entry);
    else
        return renderDetail(entry);

    function renderRow(entry) {
        out = "";
        return out;
    }

    function renderNode(entry) {
        out = "";

        switch(entry["table"]) {
            default:
                out += ""
        }
        return out;
    }

    function replaceAll(target, search, replacement) {
        return target.split(search).join(replacement);
    };

    function renderDetail(entry) {
        var out = '';
        switch(entry["table"]) {
        case "images":
            var content = JSON.parse(entry["content"]);
            var type = content["content-type"].replace("image/", "");

            if(content["account"] !== undefined){ 
                out += '<div id={{id}}" class="card p-3 text-center">';
                out += '    <a href="" class="mb-3">';
                out += '        <img src="{{image}}" alt="Photo by {{author}}" class="rounded" />';
                out += '    </a>';
                out += '    <div class="d-flex align-items-center px-2">';
                out += '        <div class="avatar avatar-md mr-3" style="background-image:{{avatar-image}}"></div>';
                out += '        <div>';
                out += '            <div>{{author}}</div>';
                out += '            <small class="d-block text-muted">{{created-relative}}</small>';
                out += '        </div>';
                out += '        <div class="ml-auto text-muted">';
                if(content["retweet_count"] !== undefined){
                out += '            <span class="icon"><i class="fe fe-repeat mr-1"></i> {{retweets}}</span>';
                }
                if(content["favorite_count"] !== undefined){
                out += '            <spanclass="icon d-none d-md-inline-block ml-3"><i class="fe fe-heart mr-1"></i> {{likes}}</span>';
                }
                out += '        </div>';
                out += '    </div>';
                out += '</div>';
                out = replaceAll(out, "{{author}}", content["account"]["accountDisplayName"]);
                out = out.replace("{{avatar-image}}", "url(/files/" + content["account"]["profile"]["profileImage"] + ")");
                out = out.replace("{{retweets}}", content["retweet_count"]);
                out = out.replace("{{likes}}", content["favorite_count"]);
            } else {
                out += '<div id={{id}}" class="card p-3 text-center">';
                out += '    <a href="" class="mb-3">';
                out += '        <img src="{{image}}" alt="Photo by {{author}}" class="rounded" />';
                out += '    </a>';
                out += '    <div class="d-flex align-items-center px-2">';
                out += '        <div>';
                out += '            <small class="d-block text-muted">{{created-relative}}</small>';
                out += '        </div>';
                out += '        <div class="ml-auto text-muted">';
                out += '        </div>';
                out += '    </div>';
                out += '</div>';
            }
            out = replaceAll(out, "{{created-relative}}", content["created"]);
            out = out.replace("{{id}}", entry["id"]);
            out = out.replace("{{image}}", "/files/"+content["filePath"]);
            return out;
        case "videos":
            var content = JSON.parse(entry["content"]);
            var type = content["content-type"].replace("video/", "");
            if(content["account"] !== undefined){
                out += '<div id={{id}}" class="card p-3 text-center">';
                out += '    <a href="" class="mb-3">';
                out += '        <video controls="controls" src="{{video}}" class="card-video rounded">';
                out += '        Your browser does not support the HTML5 Video element.';
                out += '        </video>';
                out += '    </a>';
                out += '    <div class="d-flex align-items-center px-2">';
                out += '        <div class="avatar avatar-md mr-3" style="background-image:{{avatar-image}}"></div>';
                out += '        <div>';
                out += '            <div>{{author}}</div>';
                out += '            <small class="d-block text-muted">{{created-relative}}</small>';
                out += '        </div>';
                out += '        <div class="ml-auto text-muted">';
                if(content["retweet_count"] !== undefined){
                out += '            <span class="icon"><i class="fe fe-repeat mr-1"></i> {{retweets}}</span>';
                }
                if(content["favorite_count"] !== undefined){
                out += '            <spanclass="icon d-none d-md-inline-block ml-3"><i class="fe fe-heart mr-1"></i> {{likes}}</span>';
                }
                out += '        </div>';
                out += '    </div>';
                out += '</div>';
                out = replaceAll(out, "{{author}}", content["account"]["accountDisplayName"]);
                out = out.replace("{{avatar-image}}", "url(/files/" + content["account"]["profile"]["profileImage"] + ")");
                out = out.replace("{{retweets}}", content["retweet_count"]);
                out = out.replace("{{likes}}", content["favorite_count"]);
            } else {
                out += '<div id={{id}}" class="card p-3 text-center">';
                out += '    <a href="" class="mb-3">';
                out += '        <video controls="controls" src="{{video}}" class="card-video rounded">';
                out += '        Your browser does not support the HTML5 Video element.';
                out += '        </video>';
                out += '    </a>';
                out += '    <div class="d-flex align-items-center px-2">';
                out += '        <div>';
                out += '            <small class="d-block text-muted">{{created-relative}}</small>';
                out += '        </div>';
                out += '        <div class="ml-auto text-muted">';
                out += '        </div>';
                out += '    </div>';
                out += '</div>';
            }
            out = replaceAll(out, "{{created-relative}}", content["created"]);
            out = out.replace("{{id}}", entry["id"]);
            out = out.replace("{{video}}", "/files/"+content["filePath"]);
            return out;
        case "posts":
            out += '<div class="card">';
            out += '    <div class="card-body d-flex flex-column">';
            out += '        <div class="text-muted">{{post-content}}</div>';
            out += '        <div class="d-flex align-items-center pt-5 mt-auto">';
            out += '            <div class="avatar avatar-md mr-3" style="background-image:{{avatar-image}}"></div>';
            out += '            <div>';
            out += '                <a href="" class="text-default">{{author}}</a>';
            out += '                <small class="d-block text-muted">{{created-relative}}</small>';
            out += '            </div>';
            out += '        <div class="ml-auto text-muted">';
            out += '            <span class="icon"><i class="fe fe-repeat mr-1"></i> {{retweets}}</span>';
            out += '            <span class="icon d-none d-md-inline-block ml-3"><i class="fe fe-heart mr-1"></i> {{likes}}</span>';
            out += '        </div>';
            out += '       </div>';
            out += '    </div>';
            out += '</div>';
            var data = JSON.parse(entry["content"]);
            var content = data["content"]
            out = out.replace("{{post-content}}", content["full_text"]);
            if(content["account"] !== undefined){
                out = replaceAll(out, "{{author}}", content["account"]["accountDisplayName"]);
                out = out.replace("{{avatar-image}}", "url(/files/" + content["account"]["profile"]["profileImage"] + ")");
            }
            out = replaceAll(out, "{{created-relative}}", data["created"]);
            out = out.replace("{{retweets}}", content["retweet_count"]);
            out = out.replace("{{likes}}", content["favorite_count"]);
            return out;
        case "tags":
            var data = JSON.parse(entry["content"]);
            var content = data["content"]
            out += '<div class="card">';
            out += '    <div class="card-body d-flex flex-column">';
            out += '        <div class="h1 m-0">#{{post-content}}</div>';
            out += '        <div class="d-flex align-items-center pt-5 mt-auto">';
            out += '            <div>';
            out += '                <small class="d-block text-muted">{{created-relative}}</small>';
            out += '            </div>';
            out += '       </div>';
            out += '    </div>';
            out += '</div>';
            out = out.replace("{{post-content}}", content["text"]);
            out = replaceAll(out, "{{created-relative}}", data["created"]);
            return out;
        default:
            out += '<div class="card">';
            out += '    <div class="card-body d-flex flex-column">';
            out += '        <div class="text-muted">{{post-content}}</div>';
            out += '        <div class="d-flex align-items-center pt-5 mt-auto">';
            out += '            <div class="avatar avatar-md mr-3" style="background-image:{{avatar-image}}"></div>';
            out += '            <div>';
            out += '                <a href="" class="text-default">{{author}}</a>';
            out += '                <small class="d-block text-muted">{{created-relative}}</small>';
            out += '            </div>';
            out += '        <div class="ml-auto text-muted">';
            out += '            <span class="icon"><i class="fe fe-repeat mr-1"></i> {{retweets}}</span>';
            out += '            <span class="icon d-none d-md-inline-block ml-3"><i class="fe fe-heart mr-1"></i> {{likes}}</span>';
            out += '        </div>';
            out += '       </div>';
            out += '    </div>';
            out += '</div>';
            var data = JSON.parse(entry["content"]);
            var content = data["content"]
            out = out.replace("{{post-content}}", content["full_text"]);
            if(content["account"] !== undefined){
                out = replaceAll(out, "{{author}}", content["account"]["accountDisplayName"]);
                out = out.replace("{{avatar-image}}", "url(/files/" + content["account"]["profile"]["profileImage"] + ")");
            }
            out = replaceAll(out, "{{created-relative}}", data["created"]);
            out = out.replace("{{retweets}}", content["retweet_count"]);
            out = out.replace("{{likes}}", content["favorite_count"]);
            return out;
        }
        return out;
    }
})();