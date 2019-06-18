package model

// Instance dummy struct
type Instance struct{}

// InstanceQuery01 raw query
var InstanceQuery01 string = `
select straight_join
  alert_instance.id as instance_id
  ,alert_instance.node as node_name
  ,alert_instance.name as instance_name
  ,ifnull(m01.num_value, '-') as agent_linux
  ,ifnull(m02.num_value, '-') as agent_mysql
  ,ifnull(m03.num_value, '-') as mysql_readonly
  ,ifnull(m04.str_value, '') as repl_master
  ,ifnull(m05.num_value, '-') as mysql_uptime
  ,ifnull(datediff(now(),from_unixtime(m06.num_value)), '-') as node_uptime
  ,ifnull(m07.str_value, '') as mysql_version
  ,ifnull(m08.str_value, '') as node_kernel
  ,ifnull(m09.num_value, '-') as mysql_innodb_buffer_pool_size
from (
  select 
    alert_instance.*, 
    substring_index(name, '::', 1) node
  from alert_instance
) alert_instance
left join snapshot_metric m01 on m01.instance = alert_instance.node and m01.name = 'agent_linux'
left join snapshot_metric m02 on m02.instance = alert_instance.name and m02.name = 'agent_mysql'
left join snapshot_metric m03 on m03.instance = alert_instance.name and m03.name = 'mysql_readonly'
left join snapshot_metric m04 on m04.instance = alert_instance.name and m04.name = 'mysql_repl_io_stus'
left join snapshot_metric m05 on m05.instance = alert_instance.name and m05.name = 'mysql_uptime'
left join snapshot_metric m06 on m06.instance = alert_instance.node and m06.name = 'node_uptime'
left join snapshot_metric m07 on m07.instance = alert_instance.name and m07.name = 'mysql_version'
left join snapshot_metric m08 on m08.instance = alert_instance.node and m08.name = 'node_kernel'
left join snapshot_metric m09 on m09.instance = alert_instance.name and m09.name = 'mysql_innodb_buffer_pool_size'
order by alert_instance.name
`

// GetInstanceList get current snapshot instance information
func (o *Instance) GetInstanceList() []map[string]string {
	results, _ := orm.QueryString(InstanceQuery01)
	//fmt.Println(results)
	return results
}
