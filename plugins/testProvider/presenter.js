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
        out += '        <h4><a href="#">{{post-title}}</a></h4>';
        out += '        <div class="text-muted">{{post-content}}</div>';
        out += '        <div class="d-flex align-items-center pt-5 mt-auto">';
        out += '            <div class="avatar avatar-md mr-3" style="background-image: {{avatar-image}}"></div>';
        out += '            <div>';
        out += '                <a href="" class="text-default">{{author}}</a>';
        out += '                <small class="d-block text-muted">{{created-relative}}</small>';
        out += '            </div>';
        out += '            <div class="ml-auto text-muted">';
        out += '                <a href="" class="icon d-none d-md-inline-block ml-3"><i class="fe fe-heart mr-1"></i></a>';
        out += '            </div>';
        out += '       </div>';
        out += '    </div>';
        out += '</div>';
        return out;
    }
})();