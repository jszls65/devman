/**
 * FastCode 快速代码工具
 *
 * @param data
 * @returns {string}
 */
function getFastCodeHtml(config) {
    var data = config.data;
    var tableName = config.tableName;
    var tableIndex = config.index;
    var formDivBegin = '<form class="layui-form" lay-filter="fastcodeform" action="">';
    var formDivBeginEnd = '</form>';
    // bean对象
    var javaBeanDiv = getFastCodeBaseDivHtml(
        "java Bean<br>实体类",
        "javaBean",
        getJavaBeanTxt(data)
    );
    // select语句
    var selectSqlDiv = getFastCodeRenameDivHtml(
        "select语句",
        "selectSql",
        getSelectSqlTxt(data, tableName, null),
        tableIndex,
        tableName
    );
    // MyBatis中insert语句
    var insertSql4MybatisDiv = getFastCodeRadioDivHtml(
        "insert语<br>MyBatis",
        "insertSql4MyBatis",
        getInsertSql4MybatisTxt(data, tableName, "1"),
        tableIndex,
        tableName
    );
    return formDivBegin
        + javaBeanDiv
        + selectSqlDiv
        + insertSql4MybatisDiv
        + formDivBeginEnd;
}

/**
 * FastCode 组装java实体类
 *
 * @param data
 * @returns {string|string}
 */
function getJavaBeanTxt(data) {
    var javaBeanTxt = "";
    for (var i in data) {
        var row = data[i];
        javaBeanTxt += "/** "+ row.comment +" */\n";
        javaBeanTxt += "private "+ getJavaClassType(row.tType) +" "+ camelCase(row.tField) +";\n";
    }
    return javaBeanTxt;
}

/**
 * FastCode 组装Mybatis新增语句
 *
 * @param data
 * @param tableName
 * @param type          1:单个，2:批量
 * @returns {string}
 */
function getInsertSql4MybatisTxt(data, tableName, type) {
    var batch = (type && type == "2");
    var insertSqlTxt = "<insert id=\"insert"+ (batch ? "Batch" : "") +"\" parameterType=\""+ (batch ? "java.util.List" : "bean") +"\">\nINSERT INTO "+ tableName +" ( \n\t";
    var insertParam = "";
    for (var i in data) {
        var row = data[i];
        insertSqlTxt += row.tField + ", ";
        insertParam += "#{" + (batch ? "item." : "" ) + camelCase(row.tField) + "}, ";
    }
    insertSqlTxt = insertSqlTxt.substring(0, insertSqlTxt.length -2);
    insertParam = insertParam.substring(0, insertParam.length -2);
    insertSqlTxt += "\n) VALUES \n";
    if (batch) {
        insertSqlTxt += "<foreach collection=\"list\" item=\"item\" separator=\",\">\n";
        insertSqlTxt += "(" + insertParam + ")\n";
        insertSqlTxt += "</foreach>\n";
    } else {
        insertSqlTxt += "(" + insertParam + ");\n";
    }
    insertSqlTxt += "</insert>";
    return insertSqlTxt;
}

/**
 * FastCode 组装查询语句
 *
 * @param data
 * @param tableName
 * @param shortName
 * @returns {string}
 */
function getSelectSqlTxt(data, tableName, shortName) {
    var selectSqlTxt = "SELECT \n\t";
    for (var i in data) {
        var row = data[i];
        selectSqlTxt += ((shortName ? shortName + "." : "") +row.tField + ", ");
    }
    selectSqlTxt = selectSqlTxt.substring(0, selectSqlTxt.length -2);
    selectSqlTxt += "\n" + "FROM " + tableName + (shortName ? " " + shortName : "") + ";";
    return selectSqlTxt;
}

/**
 * FastCode 给表加上别名
 *
 * @param t
 */
function renameTxt(t) {
    var shortName = $(t).val();
    var tableName = $(t).attr("tablename");
    var tableIndex = $(t).attr("tableindex");
    var textarea = $(t).parent().parent().find("textarea");
    if ("selectSql" === textarea.prop("id")) {
        var data = layui.table.cache[tableIndex];
        textarea.html(getSelectSqlTxt(data, tableName, shortName));
    }
}

/**
 * FastCode 切换单选
 *
 * @param t
 */
function radioClick(t, value) {
    var radioValue = value ? $(t).val() : value;
    var tableName = $(t).attr("tablename");
    var tableIndex = $(t).attr("tableindex");
    var textarea = $(t).parent().parent().find("textarea");
    if ("insertSql4MyBatis" === textarea.prop("id")) {
        var data = layui.table.cache[tableIndex];
        textarea.html(getInsertSql4MybatisTxt(data, tableName, radioValue));
    }
}

/**
 * FastCode 组装展示内容 基本模块
 *
 * @param labelInfo
 * @param txtId
 * @param textInfo
 * @returns {string}
 */
function getFastCodeBaseDivHtml(labelInfo, txtId, textInfo) {
    return '<div class="layui-form-item">\n' +
        '    <label class="layui-form-label">'+ labelInfo +'</label>\n' +
        '    <div class="layui-input-block">\n' +
        '      <textarea id="'+ txtId +'" class="layui-textarea">'+ textInfo +'</textarea>\n' +
        '    </div>\n' +
        '  </div>';
}

/**
 * FastCode 组装展示内容 别名模块
 *
 * @param labelInfo     内容名称
 * @param txtId         内容id
 * @param textInfo      内容
 * @param tableIndex    表格index
 * @param tableName     表格名称
 * @returns {string}
 */
function getFastCodeRenameDivHtml(labelInfo, txtId, textInfo, tableIndex, tableName) {
    return '<div class="layui-form-item">\n' +
        '    <label class="layui-form-label">'+ labelInfo +'</label>\n' +
        '<div class="layui-input-block">\n' +
        '<input type="text" placeholder="给表起个别名？" class="layui-input" tablename="'+ tableName +'" tableindex="'+ tableIndex +'" oninput="renameTxt(this)">\n' +
        '</div>'+
        '    <div class="layui-input-block">\n' +
        '      <textarea id="'+ txtId +'" class="layui-textarea">'+ textInfo +'</textarea>\n' +
        '    </div>\n' +
        '  </div>';
}

/**
 * FastCode 组装展示内容 单选模块
 *
 * @param labelInfo     内容名称
 * @param txtId         内容id
 * @param textInfo      内容
 * @param tableIndex    表格index
 * @param tableName     表格名称
 * @returns {string}
 */
function getFastCodeRadioDivHtml(labelInfo, txtId, textInfo, tableIndex, tableName) {
    return '<div class="layui-form-item">\n' +
        '    <label class="layui-form-label">'+ labelInfo +'</label>\n' +
        '<div class="layui-input-block">\n' +
        '      <input type="radio" name="sex" value="1" title="单个" lay-filter="'+ txtId +'" checked tablename="'+ tableName +'" tableindex="'+ tableIndex +'">\n' +
        '      <input type="radio" name="sex" value="2" title="批量" lay-filter="'+ txtId +'" tablename="'+ tableName +'" tableindex="'+ tableIndex +'">\n' +
        '    </div>' +
        '    <div class="layui-input-block">\n' +
        '      <textarea id="'+ txtId +'" class="layui-textarea">'+ textInfo +'</textarea>\n' +
        '    </div>\n' +
        '  </div>';
}

/**
 * FastCode 数据库type转java Type
 *
 * @param type
 * @returns {string}
 */
function getJavaClassType(type) {
    var lowCaseType = type.toLowerCase();
    if (lowCaseType.indexOf("bigint") >= 0) {
        return "Long";
    } else if (lowCaseType.indexOf("int") >= 0) {
        return "Integer";
    } else if (lowCaseType.indexOf("bit") >= 0) {
        return "Boolean";
    } else if (lowCaseType.indexOf("double") >= 0) {
        return "Double";
    } else if (lowCaseType.indexOf("float") >= 0) {
        return "Float";
    } else if (lowCaseType.indexOf("decimal") >= 0) {
        return "BigDecimal";
    } else if (lowCaseType.indexOf("date") >= 0 || lowCaseType.indexOf("time") >= 0 || lowCaseType.indexOf("year") >= 0) {
        return "Date";
    } else {
        return "String";
    }
}

/**
 * FastCode 数据库字段转驼峰命名
 *
 * @param name
 */
function camelCase(name) {
    if (!name) {
        return "";
    }
    var match = name.match(/_[\w]/g);
    if (!match || match.length <= 0) {
        return name;
    }
    match.forEach(function (value) {
        name = name.replace(value, value.split('')[1].toUpperCase());
    });
    return name;
}
