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

    function renderDetail(entry) {
        out = '';
        out += '<div class="card">';
        out += '    <div class="card-body d-flex flex-column">';
        out += '        <div>';
        out += '            <a href="/?provider={{provider}}" class="text-default">{{provider}}</a>/';
        out += '            <a href="/?table={{table}}" class="text-default">{{table}}</a>';
        out += '        </div>';
        out += '        <div class="text-muted"><small>{{post-content}}</small></div>';
        out += '        <div class="d-flex align-items-center pt-5 mt-auto">';
        out += '            <div>';
        out += '                <small class="d-block">imported <span class="text-muted">{{imported}}</span></small>';
        out += '                <small class="d-block">created <span class="text-muted">{{created}}</span></small>';
        out += '                <small class="d-block">updated <span class="text-muted">{{updated}}</span></small>';
        out += '            </div>';
        out += '       </div>';
        out += '    </div>';
        out += '</div>';
        out = replaceAll(out, "{{table}}", entry["table"]);
        out = replaceAll(out, "{{provider}}", entry["provider"]);
        out = out.replace("{{imported}}", entry["imported"]);
        out = out.replace("{{updated}}", entry["updated"]);
        out = out.replace("{{created}}", entry["created"]);
        out = out.replace("{{post-content}}", 
            replaceAll(entry["content"],",",",\n"));
        return out;
    }

    function replaceAll(target, search, replacement) {
        return target.split(search).join(replacement);
    };
})();