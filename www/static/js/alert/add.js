
function closeDialog(){
    layer.closeAll()
}

var submitting = false
var _form
layui.use(['form'], function () {
    _form = layui.form;
    //自定义验证规则
    _form.verify({
        phone: [/^1[3|4|5|7|8]\d{9}$/, '手机必须11位，只能是数字！']
        , email: [/^[a-z0-9._%-]+@([a-z0-9-]+\.)+[a-z]{2,4}$|^1[3|4|5|7|8]\d{9}$/, '邮箱格式不对']
    });
    //监听提交
    _form.on('submit(formDemo)', function(data){
        if(submitting){
           return false;
        }
        submitting = true;
        $.ajax({
            type: 'POST',
            url: "/alert/add",
            data: JSON.stringify({
                "id": data.field["id"] == '' || data.field["id"] == '0' ? 0 : Number(data.field["id"]) ,
                "appName": data.field["app_name"],
                "httpMethod": data.field["http_method"],
                "url": data.field["url"],
                "owner": data.field["owner"],
                "state": data.field["state"] == 'on' ? 1 : 0,
                "note": data.field["note"],
                "body": data.field["body"],
                "type": 'alive'
            }),
            contentType: 'application/json;charset=UTF-8',
            dataType: 'json',
            success: function (data) {
                layer.msg(data.msg)
                search();
            },
            error: function (data) {
                layer.msg(data.responseJSON.msg)
            },
            complete: function (){
                submitting = false;

            }
        });

        return false; //阻止表单跳转。如果需要表单跳转，去掉这段即可。
    });
});


function copyData(){
    $('#id').val(0);
    $('#submitBtn').click();
}