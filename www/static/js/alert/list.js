var tableVar;

moment.locale('zh-cn')

function dataTableInit(){
    layui.use('table', function(){
        tableVar = layui.table;

        //第一个实例
        tableVar.render({
            elem: '#alertlist'
            ,id: 'alertTableId'
            ,toolbar: "#tableToolbar"
            // ,height: 516
            ,url: '/alert/list'
            ,page: false //开启分页
            ,limits: [10000]
            ,cols: [[ //表头
                {type:'checkbox'}
                ,{field: 'app_name', title: '服务名称', width:150}
                ,{field: 'http_method', title: '方法', width:80, sort: false}
                ,{field: 'url', title: '请求', width: 250}

                ,{field: 'heath_state', title: '健康状态', width:80, templet:function(row, data){

                    switch (row.heath_state){
                        case 0:
                            return ' <span class="layui-badge layui-bg-green">健康</span>'
                        case 1:
                            return ' <span class="layui-badge layui-bg-yellow">警告</span>'
                        case 2:
                            return ' <span class="layui-badge layui-bg-red">不可用</span>'

                    }
                }}
                ,{field: 'state', title: '运行状态', width:80, templet:function(row, data){
                    paramStr = "id="+row.id;
                    if(row.state == 0){
                        return '<input type="checkbox" paramStr="'+paramStr+'" name="open" lay-skin="switch" lay-filter="switchTest" title="">';
                    }else{
                        return '<input type="checkbox" paramStr="'+paramStr+'" checked name="open" lay-skin="switch" lay-filter="switchTest" title="">';
                    }
                }}
                ,{field: 'last_fail_time', title: '上次失败', width: 120,templet: function (row){
                        var calendar = moment(row.last_fail_time, "YYYY-MM-DD HH:mm").calendar();
                        if(calendar.startsWith("0001")){
                            return "";
                        }
                        calendar = calendar.replaceAll("/", "-")

                        return calendar;
                    }}
                ,{field: 'owner', title: '负责人', width:80}
                ,{field: 'create_time', title: '创建日期', width: 120, sort: false,templet: function (row){
                        return moment(row.create_time, "YYYY-MM-DD HH:mm").calendar();
                    }}
                ,{field: 'update_time', title: '更新日期', width: 120, sort: false,templet: function (row){
                    return moment(row.update_time, "YYYY-MM-DD HH:mm").calendar();
                }}

            ]]
        });
    });
}

function search(){

    var appName = $('#appName').val() || "";
    var owner = $('#owner').val() || "";
    //上述方法等价于
    layui.table.reload('alertTableId', {
        where: { //设定异步数据接口的额外参数，任意设
            appName: appName
            ,owner: owner
        }
        ,page: false
    }); //只重载数据
    moment.locale('zh-cn');
    reloadSwitchRefresh(taskRunFlag)
}

function clearForm(){
    $('#appName').val("")
    $('#owner').val("")
    search();
}

// 打开添加弹窗
var addLayerId;
function loadAdd(){
    $.get('/alert/load-add', {}, function(str){
        openAddLayer('添加',str);
    });
}

function closeAddLayer(){
    search();
    layer.close(addLayerId);
}

// 删除
function delRow(){
    reloadSwitchRefresh(false)
    var checkStatus = tableVar.checkStatus('alertTableId');
    if (checkStatus.data.length == 0) {
        layer.msg("请选择要删除的行")
        return
    }
    var ids = []
    checkStatus.data.forEach(function(i){
        ids.push(i.id)
    })
    layer.confirm('确认要删除吗？', {
        btn : [ '确定', '取消' ]//按钮
    }, function(index) {
        $.ajax({
            type: 'POST',
            url: '/alert/del',
            data: JSON.stringify({"ids": ids}),
            contentType: 'application/json; charset=UTF-8',
            dataType: 'json',
            success: function (data) {
                layer.msg(data.msg)
                search();
            },
            error: function () { }
        });
    }, function (){
        // 取消
    });
}



// 编辑
function loadEdit(){
    reloadSwitchRefresh(false)
    var checkStatus = tableVar.checkStatus('alertTableId');
    if (checkStatus.data.length == 0) {
        layer.msg("请选择要删除的行")
        return
    }
    var ids = []
    checkStatus.data.forEach(function(i){
        ids.push(i.id)
    })
    if(ids.length != 1){
        layer.msg("只能选中一条记录编辑");
        return;
    }

    $.get('/alert/load-edit?id='+ids[0], {}, function(str){
        openAddLayer('编辑',str);
    });
}

// 开关监听
layui.use(['form'], function () {
    layui.form.on('switch(switchTest)', function(data){
        layer.msg('开关checked：'+ (this.checked ? 'true' : 'false'), {
        });
        var paramStr = this.getAttribute("paramStr");
        var state = this.checked ? 1 : 0;
        $.get('/alert/update-state?'+ paramStr+"&state="+state, {}, function(str){
            layer.msg(str.msg)
            search()
        });
    });

    layui.form.on('switch(autoRefresh)', function(data){
        layer.msg((this.checked ? '开启自动刷新' : '关闭自动刷新'));
        taskRunFlag = this.checked
        autoRefreshFunc()
    });
})

function openAddLayer(title, str){
    reloadSwitchRefresh(false)
    addLayerId = layer.open({
        type: 1
        , title: title
        /*, area: ['700px', '450px']*/
        , offset: 'auto'
        , shade: 0.6 // 遮罩透明度
        , shadeClose: false
        , maxmin: true
        , content: str //注意，如果str是object，那么需要字符拼接。
    });
}

var taskRunFlag = false;
var timerId
autoRefreshFunc()
function autoRefreshFunc(){
    if(timerId){
        clearInterval(timerId)
    }
    timerId = setInterval(function () {
        if(taskRunFlag){
            search()
        }
    }, 10000);
    // 10s
}


function reloadSwitchRefresh(f){
    $('#autoRefresh').prop('checked',f)
    // layui.form.render("checkbox");
}