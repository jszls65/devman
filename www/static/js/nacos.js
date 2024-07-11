$(document).ready(function () {
  $('#myInput').on('keydown', function (event) {
    // keyCode 13 是回车键的键码
    if (event.keyCode == 13) {
      selectNamespace4Discovery(event);
    }
  })
})

layui.use('table', function () {
  var table = layui.table;
  var namespace = $('input[name="namespace"]:checked').val() || '';
  var serviceName = $("#serviceName").val()
  //第一个实例
  table.render({
    elem: '#demo'
    , height: 600
    , url: '/nacos_discovery/list'//数据接口
    , where: {
      "namespace": namespace,
      "serviceName": serviceName
    }
    , page: false //开启分页
    , cols: [[ //表头
      { field: 'no', title: '序号', width: 50, fixed: 'left' }
      , { field: 'service', title: '服务', width: 200, fixed: 'left' }
      , { field: 'ip', title: 'IP', width: 200 }
      , { field: 'weight', title: '权重', width: 80 }
      , { field: 'enable', title: '在线', width: 80 }
      , { field: 'healthy', title: '健康', width: 100 }
      // ,{field: 'ephemeral', title: 'Ephemeral', width: 100}
      , { field: 'metadata', title: 'Metadata', width: 150 }
      , {
        field: '', title: '操作', width: 150, templet: function (row) {
          var checkstr = row.enable ? "checked" : "";
          var ip = row.ip.split(":")[0]
          var port = row.ip.split(":")[1]
          var paramStr = "&ip=" + ip + "&port=" + port + "&serviceName=" + row.serviceName
          return ' <input type="checkbox" paramStr=' + paramStr + ' onclick="serviceEnable(this, row)"  lay-filter="switchTest"  name="switch" ' + checkstr + ' lay-skin="switch" lay-text="上线|下线">'
        }
      }
    ]]
  });

});


function selectNamespace4Discovery(event) {
  var namespace = $('input[name="namespace"]:checked').val() || '';
  var serviceName = $("#serviceName").val()
  if (event) {
    event.preventDefault(); // 阻止默认行为，比如表单提交
  }
  layui.table.reload("demo", {
    where: {
      "namespace": namespace,
      "serviceName": serviceName
    }
  });
}




// 开关监听
layui.use(['form'], function () {
  layui.form.on('switch(switchTest)', function (data) {

    var namespace = $('input[name="namespace"]:checked').val() || '';
    var serviceName = $("#serviceName").val()
    var state = this.checked ? 1 : 0;
    var paramStr = this.getAttribute("paramStr");
    $.get('/nacos_discovery/enable?v=' + state + paramStr + "&namespace=" + namespace + "&serviceName=" + serviceName, {}, function (str) {
      layer.msg(str.msg)

    });
  });
});
