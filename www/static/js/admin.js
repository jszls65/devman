
function menuFuncAlert(env, name){
    $("#iframe_body").attr("src", "/alert/load-list");
}


function menuFuncTool(env, name){
    $("#iframe_body").attr("src", "/alert/load-list");
}

function menuFunc(menuName, fn){
    $('#tabTitle').html(menuName);
    eval(fn)();
}