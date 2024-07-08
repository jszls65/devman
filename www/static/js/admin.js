
function menuFuncAlert(env, name) {
    $("#iframe_body").attr("src", "/alert/load-list");
}


function menuFuncNacosConfig(menuName) {
    // namespace从本地缓存中获取
    var namespace = window.localStorage.namespace || "k8s-test";
    var proid = menuName.split(',')[1]  ;  // git 项目id
    $("#iframe_body").attr("src", "/nacos_config?namespace=" + namespace + "&proid=" + proid);
}


function menuFuncTool(env, name) {
    $("#iframe_body").attr("src", "/tool");
}

function menuFuncDatamap(configId) {
    $("#iframe_body").attr("src", "/datamap?configId=" + configId);

}

function menuFunc(configId, fn) {
    // 记录上一次点击的菜单
    window.localStorage.menuName = configId
    window.localStorage.menuFunc = fn
    var split = configId.split(",")

    if (fn == 'menuFuncNacosConfig') {
        split[0] = (window.localStorage.namespace || '') == '' ? 'dev' : window.localStorage.namespace;
        $('#indexPageTitle').html(split[0] + " > " + (split[2]||''));
    }else{
        if (split.length >=2){
            $('#indexPageTitle').html(split[0] + " > " + split[1]);
        }else{
            $('#indexPageTitle').html(split[0]);
        }
    }
    
    eval(fn)(configId);
}

// 获取sqlite数据库是否开启
function getSqliteDbOpen() {
    $.get("/alert/sdb-open", {}, function (data) {
        if (!data) {
            // 把菜单隐藏
            $("#menuFuncAlertDl").hide();
        }
    })
}

// 切换nacos的namespace单选框
function selectNamespace() {
    var namespace = $('input[name="namespace"]:checked').val() || '';
    
    // alert(namespace)
    window.localStorage.namespace = namespace;

    // 刷新父页面
    window.parent.location.reload();
}

var errorMsg = $("#errorMsg") || "";
if (errorMsg.length > 0){
    layui.layer.msg(errorMsg);
}
