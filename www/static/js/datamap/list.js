// 复制链接对象
var copyBtn;
// 当前页layer对象
var listLayer;
// form对象
var layuiForm;
var loadOnId;
//JavaScript代码区域
layui.use(['element', 'util', 'table', 'layer', 'form','code'], function(){
    var util = layui.util;
    listLayer = layui.layer;
    layuiForm = layui.form;
    //执行top块
    util.fixbar({});

    initRefreshCache();
});

/**
 * 打开目录
 */

function openTableSearch(){

    var dataKey = $('#dataKey').val();

    var tableNamesStr = $('#tableNames').html().trim();
    if (tableNamesStr === "") {
        return;
    }
    var tableNameList = tableNamesStr.split(",");

    $.ajax({
        url: '/datamap/table-search?env='+dataKey
        ,type: 'GET'
        ,async : true
        ,headers: {'Content-Type': 'application/json'}
        ,success: function (str){
            layui.layer.open({
            // listLayer.open({
                // title: '表名检索'
                title: '<a href="#topM" title="回到顶部">表名检索 <span class="layui-badge">'+tableNameList.length+'</span> ' +
                    '<i class="layui-icon layui-icon-up"></i></a> '
                ,content: str
                ,id: 'catalogueBox'
                ,shade:0
                ,offset:'rt'
                ,area: ['300px', '585px']
                ,btn:[]
                ,type:0
                ,maxmin:false
                ,move: false
                ,anim:1
                ,closeBtn: 0
                ,restore: function () {
                }
                ,success: function(layero, index){

                }
            });
        }
        ,complete: function (){
            // showCreateTableRunning = false;
        }
        ,error: function (){

        }
    })
}

/**
 * 搜索框绑定按钮事件 支持 上 下 回车
 *
 * @param e
 */
function dealWithKeyEvent(e) {
    var hiddenList = $('div[class="catalogue-list"] div:hidden');
    if (hiddenList.length >= 0) {
        hiddenList.each(function () {
            $(this).removeClass('selected-div');
        });
    }
    var showList = $('div[class="catalogue-list"] div:visible');
    if (showList.length <= 0) {
        return;
    }
    if (!e) {
        return;
    }
    var keyCode = e.keyCode;
    if (13 !== keyCode && 38 !== keyCode && 40 !== keyCode) {
        return;
    }
    var showListSelected = $('div[class="catalogue-detail selected-div"]:visible');
    if (showListSelected.length === 1) {
        if (13 === keyCode) {
            $(showListSelected[0]).children().first()[0].click();
        } else if (38 === keyCode) {
            $(showListSelected[0]).removeClass('selected-div');
            var index = showList.index(showListSelected[0]);
            if (0 === index) {
                $(showList[showList.length - 1]).addClass('selected-div');
            } else {
                $(showList[index - 1]).addClass('selected-div');
            }
        } else if (40 === keyCode) {
            $(showListSelected[0]).removeClass('selected-div');
            var index = showList.index(showListSelected[0]);
            if ((showList.length - 1) === index) {
                $(showList[0]).addClass('selected-div');
            } else {
                $(showList[index + 1]).addClass('selected-div');
            }
        }
    } else if (showListSelected.length === 0) {
        if (38 === keyCode) {
            $(showList[showList.length - 1]).addClass('selected-div');
        } else if (40 === keyCode) {
            $(showList[0]).addClass('selected-div');
        }
    }
}

/**
 * 查询弹窗中的a标签点击事件
 *
 * @param name
 */
function selectTable(name) {
    dealWithMaxLayer();
    dealWithH3Selected(name);
    dealWithCatalogueDetail(name);
}

/**
 * 处理目录弹窗中的选中事件
 *
 * @param name
 */
function dealWithCatalogueDetail(name) {
    if (!name) {
        return;
    }
    var selectedDetailDiv = $('div[class="catalogue-detail selected-div"]');
    if (selectedDetailDiv.length > 0) {
        selectedDetailDiv.each(function () {
            $(this).removeClass('selected-div');
        })
    }
    $('a[href="#' + name + '"]').parent().addClass('selected-div');
}

/**
 * 查询后选中的表格添加样式 其余的移除样式
 *
 * @param name
 */
function dealWithH3Selected(name) {
    var h3List = $("h3");
    if (h3List.length <= 0) {
        return;
    }
    h3List.each(function () {
       var h3 = $(this);
       if (h3.html() == name) {
       // if (h3.attr("title") == name) {
           h3.addClass("selected-h3");
       } else {
           h3.removeClass("selected-h3");
       }
    });
}

/**
 * 目录最大化时 若触发定位 则还原窗口
 */
function dealWithMaxLayer() {
    var maxButton = $('a[class="layui-layer-ico layui-layer-max layui-layer-maxmin"]');
    if (maxButton.length <= 0) {
        return;
    }
    $(maxButton[0]).click();
}

/**
 * 搜索目录
 */
function searchCatalogue() {
    var text = $('#searchBox').val();
    text = text.toLowerCase();
    var catalogues = $('div[class="catalogue-list"] a');
    if (catalogues.length <= 0) {
        return;
    }
    for (var i = 0; i < catalogues.length; i++) {
        var tableName = $(catalogues[i]).html();
        tableName = tableName.toLowerCase();
        var tableNameNew = tableName.replaceAll('_','');
        if (tableName.indexOf(text) != -1 || tableNameNew.indexOf(text) != -1) {
            $(catalogues[i]).parent().show();
        } else {
            $(catalogues[i]).parent().hide();
        }
    }
}

/**
 * 初始化复制链接事件
 */
function initCopyUrl() {
    copyBtn = new ClipboardJS('#copyUrl');
    copyBtn.on('success',function(e) {
        listLayer.alert('链接已复制到剪切板，内容如下：<br>' + e.text, {title: '链接',offset: 't' });
        e.clearSelection();
    });
    copyBtn.on('error',function(e) {
        //复制失败；
        listLayer.alert('链接复制失败，需要您手动复制如下链接：<br>' + e.text, {title: '链接', offset: 't' });
    });
}

function copyUrl(tableName){
    var dataKey = $('#dataKey').val();
    var url = window.location.host + '/datamap/share?env=' + dataKey + '&tableName=' + tableName;
    $('#copyUrl').attr('data-clipboard-text', url);
    $('#copyUrl').click();
}

/**
 * 刷新缓存
 */
var refreshRunning = false;
function initRefreshCache() {

    if ($('#refreshCache').length <= 0) {
        return;
    }

    // 点击刷新按钮
    $('#refreshCache').click(function () {
        if( refreshRunning){
            return;
        }
        refreshRunning = true;

        $('#refreshCache').html("加载中...")
        var dataKey = $('#dataKey').val();
        if (!dataKey || '' === dataKey) {
            return;
        }

        $.ajax({
            type : 'GET',
            url : '/datamap/refreshCache?env='+dataKey,
            data : {},
            success : function(result) {
                refreshRunning = false;
                result ? window.location.reload() : layer.alert('刷新当前数据库缓存失败！');
            },
            error : function () {
                refreshRunning = false;
                layer.alert('刷新当前数据库缓存失败！');
            }
        });
    });
}

var showCreateTableRunning = false;
function showCreateTable(tableName){
    if(showCreateTableRunning){
        return;
    }
    // loading
    // var currLoadIndex = listLayer.load(4, {
    //     shade: [0.7,'#fff'] //0.1透明度的白色背景
    // });
    var currLoadIndex = layer.msg('加载中', {
        icon: 16
        , offset: '250px'
        ,shade: 0.01
    });
    showCreateTableRunning = true;
    var env = $('#dataKey').val()

    $.ajax({
        url: '/datamap/load-code?env='+env+'&tableName='+tableName
        ,type: 'GET'
        ,async : true
        ,headers: {'Content-Type': 'application/json'}
        ,success: function (str){
            listLayer.open({
                type: 1
                , title: "建表语句"
                , area: ['700px', '550px']
                , offset: '10px'
                , id: "showCreateTable"
                , shade: 0.6 // 遮罩透明度
                , shadeClose: true
                , maxmin: true
                , scrollbar: false // 屏蔽浏览器滚动条
                , content: str //注意，如果str是object，那么需要字符拼接。
            });
        }
        ,complete: function (){
            listLayer.close(currLoadIndex);
            showCreateTableRunning = false;
        }
        ,error: function (){

        }
    })

}

