var table;

function dataTableInit(){
    layui.use('table', function(){
        table = layui.table;

        //第一个实例
        table.render({
            elem: '#alertlist'
            ,id: 'alertTableId'
            ,toolbar: "#tableToolbar"
            ,height: 516
            ,url: '/alert/list?env='+$('#env').val() //数据接口
            ,page: true //开启分页
            ,cols: [[ //表头
                {type:'checkbox'}
                ,{field: 'id', title: 'ID', width:80, sort: true}
                ,{field: 'app_name', title: '服务名称', width:90}
                ,{field: 'http_method', title: '方法', width:80, sort: false}
                ,{field: 'url', title: '请求', width: 250}
                ,{field: 'fail_num', title: '失败次数', width: 90, sort: false}
                ,{field: 'call_num', title: '调用次数', width: 90, sort: false}
                ,{field: 'owner', title: '负责人', width:80}
                ,{field: 'state', title: '状态', width:80, templet:function(row, data){
                    paramStr = "id="+row.id;
                    if(row.state == 0){
                        return '<input type="checkbox" paramStr="'+paramStr+'" name="open" lay-skin="switch" lay-filter="switchTest" title="">';
                    }else{
                        return '<input type="checkbox" paramStr="'+paramStr+'" checked name="open" lay-skin="switch" lay-filter="switchTest" title="">';
                    }
                }}
                ,{field: 'create_time', title: '创建日期', width: 160, sort: false}
                ,{field: 'update_time', title: '更新日期', width: 160}
            ]]
        });
    });
}

function search(){
    var appName = $('#appName').val() || "";
    var owner = $('#owner').val() || "";
    //上述方法等价于
    table.reload('alertTableId', {
        where: { //设定异步数据接口的额外参数，任意设
            appName: appName
            ,owner: owner
        }
        ,page: {
            curr: 1 //重新从第 1 页开始
        }
    }); //只重载数据

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
    var checkStatus = table.checkStatus('alertTableId');
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
    var checkStatus = table.checkStatus('alertTableId');
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

})

function openAddLayer(title, str){
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
