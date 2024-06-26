

// 初始化页面, 光标默认键入
$('#inputContent').select();
var alias = "";

$("#leftContent").focus(function(){
    this.select();
});

$("#inputContent").focus(function(){
    this.select();
});

/**
 * 转换方法
 */
function convertSql(sence){
    
    // 获取表名
    var inputContent = $('#inputContent').val() || "";
    if(($('#inputContent').val() || "") === ''){
        
        $('#inputContent').attr("placeholder", "请输入表名和别名: table_name a").select();
        $('#leftContent').attr("placeholder", "字段01, 字段02, 字段03, 字段04 ...").select();
        return
    }
    if(($('#leftContent').val() || "") === ''){
        $('#leftContent').attr("placeholder", "请输入: 字段01, 字段02, 字段03, 字段04 ...").select();
        return
    }

    // 表名
    var tableName = inputContent.split(" ")[0];
    // 别名
    alias = inputContent.split(" ")[1] || "";
    alias = alias === '' ? tableName+'.' : alias+'.';

    // 获取左侧输入框的内容
    var fieldNameList = getInputContentFieldList();
    var newLines = [];
    var keyValList = [];

    switch(sence){
        case 1: // insert
        case 3:
            // 转换后的变量行列表
            var camelNameList = [];
            // 遍历每一行
            fieldNameList.forEach(function(fieldName){
                if(isBlankLine(fieldName)){
                    return;
                }
                camelNameList.push("#{"+ toCamelWord(fieldName) +"}");
            });

            keyValList = getDuplicateKeyItems(fieldNameList, false);
            // 加 insert
            newLines.push(' ----------------  insert ----------------  ');
            newLines.push(' ');
            newLines.push("insert into " + tableName + "");
            newLines.push('    ('+fieldNameList.join(', ')+')')
            newLines.push('values');
            newLines.push('    ('+camelNameList.join(', ')+')');
            newLines.push('    on duplicate key update ');
            newLines.push(' ' + keyValList.join(', '));

          // 批量insert on duplicate key


            // 转换后的变量行列表
            camelNameList = [];
            // 遍历每一行
            fieldNameList.forEach(function(fieldName){
                if(isBlankLine(fieldName)){
                    return;
                }
                camelNameList.push("#{item." + toCamelWord(fieldName)+"}");
            });

            keyValList = getDuplicateKeyItems(fieldNameList, true);
            // 加 insert
            newLines.push(' ');
            newLines.push(' ----------------  批量 insert ----------------  ');
            newLines.push(' ');
            newLines.push('<foreach collection="list" item="item" separator=";">');
            newLines.push("insert into " + tableName + "");
            newLines.push('    ('+fieldNameList.join(', ')+')')
            newLines.push('values');
            
            newLines.push('    ('+camelNameList.join(', ')+')');
            newLines.push('    on duplicate key update ');
            newLines.push(' ' + keyValList.join(', '));
            newLines.push('</foreach>');

            reqlog("insert")
            break;
        case 4:
        case 5:
            keyValList = [];
            // 遍历每一个字段名
            fieldNameList.forEach(function(fieldName){
                if(isBlankLine(fieldName)){
                    return;
                }
                keyValList.push(fieldName + " = #{"+ toCamelWord(fieldName)+"}");
            });
            newLines.push(' ----------------  update ----------------  ');
            newLines.push(' ');
            newLines.push("    update " + tableName + " set ");
            newLines.push('    '+keyValList.join(', '));
            newLines.push('    where id = #{id}')
        // 批量插入逻辑 开始
            // cid = #{cid}
            keyValList = [];
            // 遍历每一个字段名
            fieldNameList.forEach(function(fieldName){
                if(isBlankLine(fieldName)){
                    return;
                }
                keyValList.push(fieldName + " = #{item."+ toCamelWord(fieldName)+"}");
            });
            newLines.push(' ');
            newLines.push(' ----------------  批量update ----------------  ');
            newLines.push(' ');
            newLines.push('<foreach collection="list" item="item" separator=";">');
            newLines.push("    update " + tableName + " set ");
            newLines.push('    '+keyValList.join(', '));
            newLines.push('    where id = #{item.id}')
            newLines.push('</foreach>');

            reqlog("update")
            break;
        case 6:
            newLines.push(alias + fieldNameList.join(', ' + alias));
            reqlog("table-alias")
            break;
    }
    // 将转换后的内容写入右侧输入框
    $('#rightContent').html(newLines.join('\n'));
    $('#rightContent').attr("class","prettyprint")
    prettyPrint();
}

// 将下划线 转成 驼峰
function toCamelWord(word){
    var newOneWord = "";
    // 前一个字母
    var preLetter = '';
    word.split('').forEach(function(i){
        if(i != '_'){
            if(preLetter == '_'){
                i = i.toUpperCase()
            }
            newOneWord += i;
        }
        preLetter = i;
    });
    return newOneWord;
}

// 判断是否是空行
function isBlankLine(line){
    line = (line || "").trim();
    return line == "";
}

// 获取格式化后的输入的字段列表 例如 : id, name, age
function getInputContentFieldList(){
    var content = $("#leftContent").val() || "";
    if( content === ""){
        return;
    }
    content = content.replace(new RegExp(" |\n|`","g"), '');
    content = content.toLowerCase().trim();
    // 去掉最后一个逗号
    if(content.endsWith(",")){
        content = content.substring(0, content.length-1)
    }
    // 忽略的字段
    var ignoreFieldNames = ['id', 'mod_time', 'add_time', 'create_time', 'update_time'];
    var newLines = [];
    content.split(',').forEach(function(i){
        if($.inArray(i, ignoreFieldNames) === -1){
            newLines.push(i);
        }
    })
    return newLines;
}


/**
 * 转换方法
 */
 function javaGetSet(){
    // 获取目标类名
    var targetName = $('#inputContent').val() || "className";

    // 获取左侧输入框的内容
    var lines = getInputLines();
    var newLines = [];
    // 加3行判空
    newLines.push("if(null == "+ targetName +"){");
    newLines.push("    return null;");
    newLines.push("}");
    // 遍历每一行
    lines.forEach(function(line){
        if(isBlankLine(line)){
            return;
        }
        line = line.trim();
        var getMethodName = "g"+line.substring(line.indexOf(".") + 2, line.lastIndexOf(")") + 1);
        var replace = line.replace(")", "");
        var newLine = replace + targetName + "." + getMethodName + ");";
        newLines.push(newLine);
    });
    // 最后的return
    newLines.push('return ' + newLines[3].split('.')[0] + ";");
    // 将转换后的内容写入右侧输入框
    $('#rightContent').html(newLines.join('\n'));

    reqlog("java-get-set")
}

/**
 * java 转成 json
 */
 function java2Json(){
    // 获取左侧输入框的内容
    var lines = getInputLines();
    if(lines.length == 0){
        return;
    }
    
    // 转换后的json行
    var jsonLineList = ["{"];
    var fieldObjList = getFieldObjList(lines);
    var total = fieldObjList.length;
    fieldObjList.forEach(function(i){
        total--;
        jsonLineList.push(getJsonLine(i.field, i.type, i.desc, total === 0));
    });

    // 将转换后的内容写入右侧输入框
    jsonLineList.push('}')
    $('#rightContent').html(jsonLineList.join('\n'));

    reqlog('java-json')
}

//对字符串扩展
String.prototype.resetBlank=function(){
    var regEx = /\s+/g; 
    return this.replace(regEx, ' '); 
};


// 获取json的行
function getJsonLine( field, dataType, desc, isLatsed){
    var defaultVal = "\"\"";
    if($.inArray(dataType.toLowerCase(), ['integer', 'int', 'float', 'double', 'long', 'bigdecimal']) != -1){
        defaultVal = 1;
        dataType = "Number";
    }
    ['list', 'set', 'arraylist'].forEach(function(i){
        if(dataType.toLowerCase().indexOf(i) != -1){
            defaultVal = '["string01", "string02"]';
            dataType = "数组";    
        }
    });
    
    return "    \""+field+"\": " + defaultVal + (isLatsed ? "": ",") + " // "+ dataType +" "+ desc;
}


// 获取输入的行, 去掉空行
function getInputLines(){
    var content = $("#leftContent").val() || "";
    if( content == ""){
        return [];
    }
    var lines = content.trim().split('\n');
    var newLines = [];
    lines.forEach(function(line){
        line = (line || "").trim();
        if(line == ''){
            return;
        }
        line = line.replace(';', '');
        line = line.resetBlank();
        newLines.push(line);
        
    });
    return newLines;
}


// 获取字段属性对象列表[(field, type, desc)]
function getFieldObjList(lines){
    // 标记方法是否开始
    var methodFlag = false;
    var jsonObjList = [];
    // 注释
    var desc = "";
    // 遍历每一行
    lines.forEach(function(line){
        
        // 跳过包, 导包, 注释, 注解, 常量, 类
        if(line.indexOf('/*') != -1 || line.indexOf('*/') != -1 || line.startsWith('package') 
            || line.startsWith('import') || line.indexOf('@') != -1  
            || line.indexOf('static') != -1 || line.indexOf('final') != -1
            || line.indexOf('public class') != -1) {
            return;
        }
        
        // 跳过方法
        if(line.indexOf('{') !== -1){
            methodFlag = true;
            return;
        }
        if(line.indexOf('}') !== -1){
            methodFlag = false;
            return;
        }
        if(methodFlag){
            return;
        }
        
        if(line.startsWith('*') || line.startsWith('//')){
            // 得到注释
            line = line.replaceAll('*','');
            line = line.replaceAll('//','');
            line = line.trim() ;
            desc = line;
            return;
        }
        // 只处理 两个关键词的行
        if(line.indexOf('private') == -1 && line.indexOf('public') == -1){
            return;
        }

        // 变量名, 数据类型, 注释
        jsonObjList.push({
            'field': line.split(' ')[2],
            'type': line.split(' ')[1],
            'desc': desc
        });
        desc = "";
     
    });

    return jsonObjList;
}


function mergeRequst() {
    // 获取分支名称, 个人分支, 去掉后缀就是目标分支
    var privateBranch = $("#inputContent").val() || "";
    if(privateBranch == ''){
        $('#inputContent').attr("placeholder", "请输入个人分支");
        $('#inputContent').focus();
        return;
    }
    var splis = privateBranch.split(' ');
    var targetBranch = "master"; // 默认master
    if(splis.length >= 2){
        privateBranch = splis[0] || ""
        targetBranch = splis[1] || ""
    }else{
        // 只有一个元素
        if(privateBranch.lastIndexOf('_') !== -1){
            targetBranch = privateBranch.substring(0, privateBranch.lastIndexOf('_'))
        }else{
            targetBranch = 'master';
        }
    }
   
    // 最后一个中划线
    // http://101.37.39.148:20080/smartgroup/smartadjava/-/merge_requests/new?merge_request%5Bsource_branch%5D=dev_1.0.0_0721-zls&merge_request%5Btarget_branch%5D=dev_1.0.0_0721
    
     // 被指派人的id
    var assignee_ids = $("#leftContent").val() || "50";
    privateBranch = privateBranch.replaceAll('/', '%2F');
    targetBranch = targetBranch.replaceAll('/', '%2F');
    var url = "http://101.37.39.148:20080/smartgroup/smartadjava/-/merge_requests/new?merge_request%5Bsource_branch%5D="+ privateBranch +"&merge_request%5Btarget_branch%5D=" + targetBranch+"&merge_request[assignee_ids][]="+assignee_ids+"&merge_request[force_remove_source_branch]=1";
    window.open(url);
    reqlog("MR")
}


function myReplace(){
    var suffix = "_hd"
    $("#rightContent").val("")
    // 多个单词
    var privateBranch = $("#inputContent").val() || "";
    if(privateBranch == ''){
        alert('请输入要替换的内容');
        $('#inputContent').focus()
        return;
    }
    privateBranch = privateBranch.replaceAll(suffix, "");
    // 被替换内容
    var leftContent = $("#leftContent").val() || "";
    if(leftContent == ''){
        alert('请输入被替换的内容');
        $('#leftContent').focus()
        return;
    }
    var results = [];
    var split = privateBranch.split(",")
    split.forEach(function(i){
        i = i.trim()
        if(i!= ""){
            var originStr = " " + i + " ";
            if(leftContent.indexOf(originStr) != -1){
                results.push(i+suffix);
                leftContent = leftContent.replaceAll(originStr, " " + i + suffix +" ")
            }
        }
        
    })
    $("#rightContent").val(results.join(", "))
}

function clearAll(){
    $("#inputContent").val("").select();
    $("#leftContent").val("");
    $("#rightContent").html("");

    $('#inputContent').attr("placeholder", "请输入...");
    $('#leftContent').attr("placeholder", "");
    window.localStorage.inputContent = "";
    window.localStorage.leftContent = "";
}

function getDuplicateKeyItems(fieldNameList, isBatch){
    var keyValList = [];
    var pre = isBatch ? "item." : "";
    // 遍历每一个字段名
    fieldNameList.forEach(function(fieldName){
        if(isBlankLine(fieldName)){
            return;
        }
        keyValList.push(fieldName + " = #{"+ pre + toCamelWord(fieldName)+"}");
    });
    return keyValList;
}


function reqlog(event){
    $.get("/log/save", {event: event}, function (data) {
    });
}

function cacheInput(obj) {
    window.localStorage.inputContent = $(obj).val();
}

function cacheLeft(obj) {
    window.localStorage.leftContent = $(obj).val();
}