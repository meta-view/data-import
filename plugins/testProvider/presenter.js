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
        out += '        <div class="text-muted">';
        out += '            <pre>{{post-content}}</pre>';
        out += '        </div>';
        out += '    </div>';
        out += '</div>';
        out = out.replace("{{post-content}}", entry["content"].replace(",", ",\n"));
        return out;
    }
})();