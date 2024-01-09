module short

require mypackage v0.0.0

replace mypackage => ../test

//引用只需要精确到文件夹即可，和文件夹内的package无关，一个文件夹内只有一个package类型
//一个项目只能有一个go.mod进行资源管理  一个项目中可以有多个main文件夹运行不同的程序
go 1.19
