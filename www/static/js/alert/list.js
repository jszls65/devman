var tableVar;

function dataTableInit(){
    layui.use('table', function(){
        //第一个实例
        tableVar = layui.table.render({
            elem: '#alertlist'
            ,id: 'alertTableId'
            ,toolbar: "#tableToolbar"
            // ,height: 516
            ,url: '/alert/list'
            ,page: false //开启分页
            ,cols: [[ //表头
                {type:'checkbox'}
                ,{field: 'app_name', title: '服务名称', width:160}
                ,{field: 'http_method', title: '方法', width:80, sort: false}
                ,{field: 'url', title: '请求', width: 250}

                ,{field: 'heath_state', title: '健康状态', width:90, templet:function(row, data){

                    switch (row.heath_state){
                        case 0:
                            return ' <span class="layui-badge layui-bg-green">健康</span>'
                        case 1:
                            return ' <span class="layui-badge layui-bg-orange">警告</span>'
                        case 2:
                            return ' <span class="layui-badge layui-bg-red">不可用</span>'

                    }
                }}
                ,{field: 'state', title: '运行状态', width:90, templet:function(row, data){
                    paramStr = "id="+row.id;
                    if(row.state == 0){
                        return '<input type="checkbox" paramStr="'+paramStr+'" name="open" lay-skin="switch" lay-filter="switchTest" title="">';
                    }else{
                        return '<input type="checkbox" paramStr="'+paramStr+'" checked name="open" lay-skin="switch" lay-filter="switchTest" title="">';
                    }
                }}
                ,{field: 'last_fail_time', title: '上次失败', width: 120,templet: function (row){
                     
                        return formatDate(row.last_fail_time);
                    }}
                ,{field: 'owner', title: '负责人', width:80}
                ,{field: 'create_time', title: '创建日期', width: 120, sort: false,templet: function (row){
                        return formatDate(row.create_time)
                    }}
                ,{field: 'update_time', title: '更新日期', width: 120, sort: false,templet: function (row){
                    return formatDate(row.update_time)
                }}

            ]]
        });
    });
}

function search(){

    var appName = $('#appName').val() || "";
    var owner = $('#owner').val() || "";
    //上述方法等价于
    tableVar.reload({
    // tableVar.reload('alertTableId', {
        where: { //设定异步数据接口的额外参数，任意设
            appName: appName
            ,owner: owner
        }
        ,page: false
    }); //只重载数据
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
    layer.close(addLayerId);
    search();
}

// 删除
function delRow(){
    reloadSwitchRefresh(false)
    var checkStatus = layui.table.checkStatus('alertTableId');
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
    var checkStatus = layui.table.checkStatus('alertTableId');
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

    layui.$('#searchBtn').click(function(event){
        // event.preventDefault(); // 阻止默认跳转行为
        search();
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
    }, 10*1000);
    // 10s
}


function reloadSwitchRefresh(f){
    $('#autoRefresh').prop('checked',f)
    // layui.form.render("checkbox");
}

/**
 * 格式化时间为人类可读的字符串格式
 * @param {number|string|Date} time_value 13位时间戳
 * @returns {string}
 *
 * 时间格式化为：
 * 刚刚
 * 1分钟前-56分钟前
 * 1小时前-23小时前
 * 1天前-7天前
 * 2022-10-09 13:33
 */
function formatDate(time_value) {
    // 兼容其他类型的参数
    if (typeof time_value != 'number') {
        time_value = Date.parse(time_value)
    }

    // 进制转换
    let millisecond = 1
    let second = millisecond * 1000
    let minute = second * 60
    let hour = minute * 60
    let day = hour * 24
    let day_8 = day * 8

    now_time = Date.now()
    duration = now_time - time_value

    if (duration < minute) {
        return '刚刚'
    } else if (duration < hour) {
        return Math.floor(duration / minute) + '分钟前'
    } else if (duration < day) {
        return Math.floor(duration / hour) + '小时前'
    } else if (duration < day_8) {
        return Math.floor(duration / day) + '天前'
    } else {
        let date = new Date(time_value)

        return [
            [
                date.getFullYear(),
                ('0' + (date.getMonth() + 1)).slice(-2),
                ('0' + date.getDate()).slice(-2),
            ].join('-'),
            [
                ('0' + date.getHours()).slice(-2),
                ('0' + date.getMinutes()).slice(-2),
            ].join(':'),
        ].join(' ')
    }
}


