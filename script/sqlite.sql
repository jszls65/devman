-- ddl 

drop table if exists  request_log;
CREATE table `request_log`(  -- 访问日志表
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	ip text not null default '' , -- ip
	plateform text not null default '' ,  -- 平台
    path text not null default '',  -- 请求路径
	event text not null default '' ,  -- 事件
	params text not null default '',  -- 参数
	req_day text not null default '', -- 日期, yyyy-MM-ss
	create_time datetime not null,  -- 创建时间
    update_time datetime not null  -- 更新时间
);

drop table if exists request_summary ;
create table request_summary(  -- 访问次数统计表
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	ip text not null default '' , -- ip
	req_num int not null default 0,
	create_time datetime not null,  -- 创建时间
    update_time datetime not null  -- 更新时间
);



-- 接口配置
drop table if exists `interface_config`;
CREATE table `interface_config`(
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	`app_name` text not null default '',  -- 服务名称
	`type` text not null default '',  -- 类型: alive-存活, data-数据校验, autotest-自动化测试
	`http_method` text not null default '',  -- get post
	`url` text not null default '',  -- 请求的url
	`owmer` text not null default '',  -- 负责人
	`phone` text not null default '',  -- 负责人手机号
	`fail_num` INTEGER not null default 0,  -- 失败次数
	`call_num` INTEGER not null default 0,  -- 调用总次数
	create_time datetime not null,  -- 创建时间
    update_time datetime not null  -- 更新时间
);


-- 接口调用日志表
drop table if exists `interface_call_log`;
CREATE table `interface_call_log`(
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	`interface_config_id` INTEGER not null, 
	`result` integer not null default 1, -- 1-成功, 2-失败
	`cost_time` integer not null default 0, -- 耗时, 毫秒
	create_time datetime not null,  -- 创建时间
    update_time datetime not null  -- 更新时间
);