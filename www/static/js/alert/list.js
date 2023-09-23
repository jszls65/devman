var table;

function dataTableInit(){
    layui.use('table', function(){
        table = layui.table;

        //第一个实例
        table.render({
            elem: '#alertlist'
            ,toolbar: "#tableToolbar"
            ,height: 530
            ,url: '/alert/list?env='+$('#env').val() //数据接口
            ,page: true //开启分页
            ,cols: [[ //表头
                {field: 'id', title: 'ID', width:80, sort: true, fixed: 'left'}
                ,{field: 'app_name', title: '服务名称', width:90}
                ,{field: 'http_method', title: '方法', width:80, sort: false}
                ,{field: 'owner', title: '负责人', width:80}
                ,{field: 'url', title: '请求', width: 250}
                ,{field: 'fail_num', title: '失败次数', width: 90, sort: false}
                ,{field: 'call_num', title: '调用次数', width: 90, sort: false}
                ,{field: 'create_time', title: '创建日期', width: 180, sort: false}
                ,{field: 'update_time', title: '更新日期', width: 180}
                ,{field: '', title: '操作', width: 80}
            ]]
        });
    });
}

function search(){
    var appName = $('#appName').val() || "";
    var owner = $('#owner').val() || "";
    //上述方法等价于
    table.reload('alertlist', {
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
}

function loadAdd(){
    $.get('/alert/load-add', {}, function(str){
        layer.open({
            type: 1,
            content: str //注意，如果str是object，那么需要字符拼接。
        });
    });
}