
function menuFuncAlert(env, name){
    $("#iframe_body").attr("src", "/alert/load-list");
}


function menuFuncTool(env, name){
    $("#iframe_body").attr("src", "/tool");
}

function menuFuncDatamap(env, name){
        $("#iframe_body").attr("src", "/datamap?env="+env);
    
}

function menuFunc(menuName, fn){
    // 记录上一次点击的菜单
    window.localStorage.menuName = menuName
    window.localStorage.menuFunc = fn
    $('#tabTitle').html(menuName);
    eval(fn)(menuName);
}

// 获取请求次数
function getRequestCount(){
    $.get("/log/sum", {}, function (data) {
        $('#sum').html(data.data)
    });
}

// 获取sqlite数据库是否开启
function getSqliteDbOpen(){
    $.get("/alert/sdb-open", {}, function(data){
        console.log(data)
        if(!data){
            // 把菜单隐藏
            $("#menuFuncAlertDl").hide();
        }
    })
}