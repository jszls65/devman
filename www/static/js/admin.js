
function menuFuncAlert(env, name){
    $("#iframe_body").attr("src", "/alert/load-list");
}


function menuFuncTool(env, name){
    $("#iframe_body").attr("src", "/tool");
}

function menuFunc(menuName, fn){
    $('#tabTitle').html(menuName);
    eval(fn)();
}

// 获取请求次数
function getRequestCount(){
    $.get("/log/sum", {}, function (data) {
        $('#sum').html(data.data)
    });
}