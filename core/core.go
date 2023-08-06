package core

/*
Init 伪初始化，为了不使用副作用导入来执行此包内的init函数而指定，不做任何操作
*/
func Init() {
	//just empty
}

func init() {
	//加载配置
	loadConfiguration()
	//创建logger
	initZap()
	//创建mysql连接
	initGormMysql()
	//创建router
	initGin()
}
