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
)

drop table if exists request_summary ;
create table request_summary(  -- 访问次数统计表
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	ip text not null default '' , -- ip
	req_num int not null default 0,
	create_time datetime not null,  -- 创建时间
    update_time datetime not null  -- 更新时间
)