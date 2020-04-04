package table

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/encoding/gparser"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"strings"
	tableModel "yj-app/app/model/tool/table"
	tableColumnModel "yj-app/app/model/tool/table_column"
	userService "yj-app/app/service/system/user"
	"yj-app/app/utils/convert"
	"yj-app/app/utils/page"
)

//根据主键查询数据
func SelectRecordById(id int64) (*tableModel.EntityExtend, error) {
	entity, err := tableModel.SelectRecordById(id)
	if err != nil {
		return nil, err
	}
	//表附加属性
	SetTableFromOptions(entity)
	return entity, nil
}

//根据主键删除数据
func DeleteRecordById(id int64) bool {
	result, err := tableModel.Delete("table_id", id)
	if err == nil {
		affected, _ := result.RowsAffected()
		if affected > 0 {
			return true
		}
	}

	return false
}

//批量删除数据记录
func DeleteRecordByIds(ids string) int64 {
	idarr := convert.ToInt64Array(ids, ",")
	result, err := tableModel.Delete("table_id in (?)", idarr)
	if err != nil {
		return 0
	}

	nums, _ := result.RowsAffected()

	if nums > 0 {
		tableColumnModel.Delete("table_id in (?)", idarr)
	}

	return nums
}

//保存修改数据
func SaveEdit(req *tableModel.EditReq, session *ghttp.Session) (int64, error) {
	if req == nil {
		return 0, gerror.New("参数错误")
	}

	table, err := tableModel.FindOne("table_id=?", req.TableId)
	if err != nil || table == nil {
		return 0, gerror.New("数据不存在")
	}

	if req.TableName != "" {
		table.TableName = req.TableName
	}

	if req.TableComment != "" {
		table.TableComment = req.TableComment
	}

	if req.BusinessName != "" {
		table.BusinessName = req.BusinessName
	}

	if req.ClassName != "" {
		table.ClassName = req.ClassName
	}

	if req.FunctionAuthor != "" {
		table.FunctionAuthor = req.FunctionAuthor
	}

	if req.FunctionName != "" {
		table.FunctionName = req.FunctionName
	}

	if req.ModuleName != "" {
		table.ModuleName = req.ModuleName
	}

	if req.PackageName != "" {
		table.PackageName = req.PackageName
	}

	if req.Remark != "" {
		table.Remark = req.Remark
	}

	if req.TplCategory != "" {
		table.TplCategory = req.TplCategory
	}

	if req.Params != "" {
		table.Options = req.Params
	}

	table.UpdateTime = gtime.Now()

	user := userService.GetProfile(session)

	if user != nil {
		table.UpdateBy = user.LoginName
	}

	tx, err := g.DB().Begin()
	if err != nil {
		return 0, err
	}

	_, err = tx.Table("gen_table").Update(table, "table_id="+gconv.String(table.TableId))

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	//保存列数据
	if req.Columns != "" {
		if j, err := gjson.DecodeToJson([]byte(req.Columns)); err != nil {
			glog.Error(err)
		} else {
			var columnList []tableColumnModel.Entity
			err = j.ToStructs(&columnList)
			if err == nil && columnList != nil && len(columnList) > 0 {
				for _, column := range columnList {
					if column.ColumnId > 0 {
						tmp, _ := tableColumnModel.FindOne("column_id=?", column.ColumnId)
						if tmp != nil {
							tmp.ColumnComment = column.ColumnComment
							tmp.GoType = column.GoType
							tmp.HtmlType = column.HtmlType
							tmp.QueryType = column.QueryType
							tmp.GoField = column.GoField
							tmp.DictType = column.DictType
							tmp.IsInsert = column.IsInsert
							tmp.IsEdit = column.IsEdit
							tmp.IsList = column.IsList
							tmp.IsQuery = column.IsQuery

							_, err = tx.Table("gen_table_column").Update(tmp, "column_id="+gconv.String(tmp.ColumnId))

							if err != nil {
								tx.Rollback()
								return 0, err
							}
						}
					}
				}
			}
		}
	}

	return 1, tx.Commit()
}

//设置代码生成其他选项值
func SetTableFromOptions(entity *tableModel.EntityExtend) {
	if entity != nil && entity.Options != "" {
		if p, e := gparser.LoadJson([]byte(entity.Options)); e != nil {
			glog.Error(e)
		} else {
			treeCode := p.GetString("treeCode")
			treeParentCode := p.GetString("treeParentCode")
			treeName := p.GetString("treeName")
			entity.TreeCode = treeCode
			entity.TreeParentCode = treeParentCode
			entity.TreeName = treeName
		}
	}

}

//设置主键列信息
func SetPkColumn(table *tableModel.EntityExtend, columns []tableColumnModel.Entity) {
	for _, column := range columns {
		if column.IsPk == "1" {
			table.PkColumn = column
			break
		}
	}
	if &(table.PkColumn) == nil {
		table.PkColumn = columns[0]
	}
}

//根据条件分页查询数据
func SelectListByPage(param *tableModel.SelectPageReq) ([]tableModel.Entity, *page.Paging, error) {
	return tableModel.SelectListByPage(param)
}

//查询据库列表
func SelectDbTableList(param *tableModel.SelectPageReq) ([]tableModel.Entity, *page.Paging, error) {
	return tableModel.SelectDbTableList(param)
}

//查询据库列表
func SelectDbTableListByNames(tableNames []string) ([]tableModel.Entity, error) {
	return tableModel.SelectDbTableListByNames(tableNames)
}

//根据table_id查询表列数据
func SelectGenTableColumnListByTableId(tableId int64) ([]tableColumnModel.Entity, error) {
	return tableColumnModel.SelectGenTableColumnListByTableId(tableId)
}

//查询据库列表
func SelectTableByName(tableName string) (*tableModel.Entity, error) {
	return tableModel.SelectTableByName(tableName)
}

//查询表ID业务信息
func SelectGenTableById(tableId int64) (*tableModel.Entity, error) {
	return tableModel.SelectGenTableById(tableId)
}

//查询表名称业务信息
func SelectGenTableByName(tableName string) (*tableModel.Entity, error) {
	return tableModel.SelectGenTableByName(tableName)
}

//导入表结构
func ImportGenTable(tableList []tableModel.Entity, operName string) error {
	if tableList != nil && operName != "" {
		tx, err := g.DB().Begin()
		if err != nil {
			return err
		}

		for _, table := range tableList {
			tableName := table.TableName
			InitTable(&table, operName)
			result, err := tx.Table("gen_table").Insert(table)
			if err != nil {
				return err
			}

			tmpid, err := result.LastInsertId()

			if err != nil || tmpid <= 0 {
				tx.Rollback()
				return gerror.New("保存数据失败")
			}

			table.TableId = tmpid

			// 保存列信息
			genTableColumns, err := tableColumnModel.SelectDbTableColumnsByName(tableName)

			if err != nil || len(genTableColumns) <= 0 {
				tx.Rollback()
				return gerror.New("获取列数据失败")
			}

			for _, column := range genTableColumns {
				InitColumnField(&column, &table)
				_, err = tx.Table("gen_table_column").Insert(column)
				if err != nil {
					tx.Rollback()
					return gerror.New("保存列数据失败")
				}
			}
		}
		return tx.Commit()
	} else {
		return gerror.New("参数错误")
	}
}

//初始化表信息
func InitTable(table *tableModel.Entity, operName string) {
	table.ClassName = ConvertClassName(table.TableName)
	table.PackageName = g.Cfg().GetString("gen.packageName")
	table.ModuleName = g.Cfg().GetString("gen.moduleName")
	table.BusinessName = GetBusinessName(table.TableName)
	table.FunctionName = strings.ReplaceAll(table.TableComment, "表", "")
	table.FunctionAuthor = g.Cfg().GetString("gen.author")
	table.CreateBy = operName
	table.TplCategory = "crud"
	table.CreateTime = gtime.Now()
}

//初始化列属性字段
func InitColumnField(column *tableColumnModel.Entity, table *tableModel.Entity) {
	dataType := GetDbType(column.ColumnType)
	columnName := column.ColumnName
	column.TableId = table.TableId
	column.CreateBy = table.CreateBy
	//设置字段名
	column.GoField = ConvertToCamelCase(columnName)
	column.HtmlField = ConvertToCamelCase1(columnName)

	if tableColumnModel.IsStringObject(dataType) {
		//字段为字符串类型
		column.GoType = "string"
		if strings.EqualFold(dataType, "text") || strings.EqualFold(dataType, "tinytext") || strings.EqualFold(dataType, "mediumtext") || strings.EqualFold(dataType, "longtext") {
			column.HtmlType = "textarea"
		} else {
			columnLength := GetColumnLength(column.ColumnType)
			if columnLength >= 500 {
				column.HtmlType = "textarea"
			} else {
				column.HtmlType = "input"
			}
		}
	} else if tableColumnModel.IsTimeObject(dataType) {
		//字段为时间类型
		column.GoType = "Time"
		column.HtmlType = "datatime"
	} else if tableColumnModel.IsNumberObject(dataType) {
		//字段为数字类型
		column.HtmlType = "input"
		// 如果是浮点型
		tmp := column.ColumnType
		if tmp == "float" || tmp == "double" {
			column.GoType = "float64"
		} else {
			start := strings.Index(tmp, "(")
			end := strings.Index(tmp, ")")
			result := tmp[start+1 : end]
			arr := strings.Split(result, ",")
			if len(arr) == 2 && gconv.Int(arr[1]) > 0 {
				column.GoType = "float64"
			} else if len(arr) == 1 && gconv.Int(arr[0]) <= 10 {
				column.GoType = "int"
			} else {
				column.GoType = "int64"
			}
		}
	}
	//新增字段
	if columnName == "create_by" || columnName == "create_time" || columnName == "update_by" || columnName == "update_time" {
		column.IsRequired = "0"
		column.IsInsert = "0"
	} else {
		column.IsRequired = "0"
		column.IsInsert = "1"
		if strings.Index(columnName, "name") >= 0 || strings.Index(columnName, "status") >= 0 {
			column.IsRequired = "1"
		}
	}

	// 编辑字段
	if tableColumnModel.IsNotEdit(columnName) {
		if column.IsPk == "1" {
			column.IsEdit = "0"
		} else {
			column.IsEdit = "1"
		}
	} else {
		column.IsEdit = "0"
	}
	// 列表字段
	if tableColumnModel.IsNotList(columnName) {
		column.IsList = "1"
	} else {
		column.IsList = "0"
	}
	// 查询字段
	if tableColumnModel.IsNotQuery(columnName) {
		column.IsQuery = "1"
	} else {
		column.IsQuery = "0"
	}

	// 查询字段类型
	if CheckNameColumn(columnName) {
		column.QueryType = "LIKE"
	} else {
		column.QueryType = "EQ"
	}

	// 状态字段设置单选框
	if CheckStatusColumn(columnName) {
		column.HtmlType = "radio"
	} else if CheckTypeColumn(columnName) || CheckSexColumn(columnName) {
		// 类型&性别字段设置下拉框
		column.HtmlType = "select"
	}
}

//检查字段名后3位是否是sex
func CheckSexColumn(columnName string) bool {
	if len(columnName) >= 3 {
		end := len(columnName)
		start := end - 3

		if start <= 0 {
			start = 0
		}

		if columnName[start:end] == "sex" {
			return true
		}
	}
	return false
}

//检查字段名后4位是否是type
func CheckTypeColumn(columnName string) bool {
	if len(columnName) >= 4 {
		end := len(columnName)
		start := end - 4

		if start <= 0 {
			start = 0
		}

		if columnName[start:end] == "type" {
			return true
		}
	}
	return false
}

//检查字段名后4位是否是name
func CheckNameColumn(columnName string) bool {
	if len(columnName) >= 4 {
		end := len(columnName)
		start := end - 4

		if start <= 0 {
			start = 0
		}

		tmp := columnName[start:end]

		if tmp == "name" {
			return true
		}
	}
	return false
}

//检查字段名后6位是否是status
func CheckStatusColumn(columnName string) bool {
	if len(columnName) >= 6 {
		end := len(columnName)
		start := end - 6

		if start <= 0 {
			start = 0
		}
		tmp := columnName[start:end]

		if tmp == "status" {
			return true
		}
	}

	return false
}

//获取数据库类型字段
func GetDbType(columnType string) string {
	if strings.Index(columnType, "(") > 0 {
		return columnType[0:strings.Index(columnType, "(")]
	} else {
		return columnType
	}
}

//表名转换成类名
func ConvertClassName(tableName string) string {
	autoRemovePre := g.Cfg().GetBool("gen.autoRemovePre")
	tablePrefix := g.Cfg().GetString("gen.tablePrefix")
	if autoRemovePre && tablePrefix != "" {
		searchList := strings.Split(tablePrefix, ",")
		for _, str := range searchList {
			tableName = strings.ReplaceAll(tableName, str, "")
		}
	}
	return tableName
}

//获取业务名
func GetBusinessName(tableName string) string {
	lastIndex := strings.LastIndex(tableName, "_")
	nameLength := len(tableName)
	businessName := tableName[lastIndex+1 : nameLength]
	return businessName
}

//将下划线大写方式命名的字符串转换为驼峰式。如果转换前的下划线大写方式命名的字符串为空，则返回空字符串。 例如：HELLO_WORLD->HelloWorld
func ConvertToCamelCase(name string) string {
	if name == "" {
		return ""
	} else if !strings.Contains(name, "_") {
		// 不含下划线，仅将首字母大写
		return strings.ToUpper(name[0:1]) + name[1:len(name)]
	}
	var result string = ""
	camels := strings.Split(name, "_")
	for index := range camels {
		if camels[index] == "" {
			continue
		}
		camel := camels[index]
		result = result + strings.ToUpper(camel[0:1]) + strings.ToLower(camel[1:len(camel)])
	}
	return result
}

////将下划线大写方式命名的字符串转换为驼峰式,首字母小写。如果转换前的下划线大写方式命名的字符串为空，则返回空字符串。 例如：HELLO_WORLD->helloWorld
func ConvertToCamelCase1(name string) string {
	if name == "" {
		return ""
	} else if !strings.Contains(name, "_") {
		// 不含下划线，原值返回
		return name
	}
	var result string = ""
	camels := strings.Split(name, "_")
	for index := range camels {
		if camels[index] == "" {
			continue
		}
		camel := camels[index]
		if result == "" {
			result = strings.ToLower(camel[0:1]) + strings.ToLower(camel[1:len(camel)])
		} else {
			result = result + strings.ToUpper(camel[0:1]) + strings.ToLower(camel[1:len(camel)])
		}
	}
	return result
}

//获取字段长度
func GetColumnLength(columnType string) int {
	start := strings.Index(columnType, "(")
	end := strings.Index(columnType, ")")
	result := columnType[start+1 : end-1]
	return gconv.Int(result)
}

//获取Go类别下拉框
func GoTypeTpl() string {
	return `<script id="goTypeTpl" type="text/x-jquery-tmpl">
<div>
<select class='form-control' name='columns[${index}].goType'>
    <option value="int64" {{if goType==="int64"}}selected{{/if}}>int64</option>
    <option value="int" {{if goType==="int"}}selected{{/if}}>int</option>
    <option value="string" {{if goType==="string"}}selected{{/if}}>string</option>
    <option value="Time" {{if goType==="Time"}}selected{{/if}}>Time</option>
    <option value="float64" {{if goType==="float64"}}selected{{/if}}>float64</option>
    <option value="byte" {{if goType==="byte"}}selected{{/if}}>byte</option>
</select>
</div>
</script>`
}

//获取查询方式下拉框
func QueryTypeTpl() string {
	return `<script id="queryTypeTpl" type="text/x-jquery-tmpl">
<div>
<select class='form-control' name='columns[${index}].queryType'>
    <option value="EQ" {{if queryType==="EQ"}}selected{{/if}}>=</option>
    <option value="NE" {{if queryType==="NE"}}selected{{/if}}>!=</option>
    <option value="GT" {{if queryType==="GT"}}selected{{/if}}>></option>
    <option value="GTE" {{if queryType==="GTE"}}selected{{/if}}>>=</option>
    <option value="LT" {{if queryType==="LT"}}selected{{/if}}><</option>
    <option value="LTE" {{if queryType==="LTE"}}selected{{/if}}><=</option>
    <option value="LIKE" {{if queryType==="LIKE"}}selected{{/if}}>Like</option>
    <option value="BETWEEN" {{if queryType==="BETWEEN"}}selected{{/if}}>Between</option>
</select>
</div>
</script>`
}

// 获取显示类型下拉框
func HtmlTypeTpl() string {
	return `<script id="htmlTypeTpl" type="text/x-jquery-tmpl">
<div>
<select class='form-control' name='columns[${index}].htmlType'>
    <option value="input" {{if htmlType==="input"}}selected{{/if}}>文本框</option>
    <option value="textarea" {{if htmlType==="textarea"}}selected{{/if}}>文本域</option>
    <option value="select" {{if htmlType==="select"}}selected{{/if}}>下拉框</option>
    <option value="radio" {{if htmlType==="radio"}}selected{{/if}}>单选框</option>
    <option value="checkbox" {{if htmlType==="checkbox"}}selected{{/if}}>复选框</option>
    <option value="datetime" {{if htmlType==="datetime"}}selected{{/if}}>日期控件</option>
</select>
</div>
</script>`
}
