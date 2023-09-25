
function closeDialog(){
    layer.closeAll()
}

var submitting = false
layui.use(['form'], function () {
    var form = layui.form;
    //自定义验证规则
    form.verify({
        phone: [/^1[3|4|5|7|8]\d{9}$/, '手机必须11位，只能是数字！']
        , email: [/^[a-z0-9._%-]+@([a-z0-9-]+\.)+[a-z]{2,4}$|^1[3|4|5|7|8]\d{9}$/, '邮箱格式不对']
    });
    //监听提交
    form.on('submit(formDemo)', function(data){
        if(submitting){
           return false;
        }
        submitting = true;
        console.log(data) //当前容器的全部表单字段，名值对形式：{name: value};获取单个值 data.field["title"]
        $.post(
            "/alert/add",
            {
                "appName": data.field["app_name"],
                "httpMethod": data.field["http_method"],
                "url": data.field["url"],
                "owner": data.field["owner"],
                "state": data.field["state"] == 'on' ? 1 : 0,
                "type": 'alive'
            },
            function (data) {
                layer.msg("提交成功")
                submitting = false;
            }
        );

        return false; //阻止表单跳转。如果需要表单跳转，去掉这段即可。
    });
});