package request

// PageQuery 请求一个页面数据的dto对象
type PageQuery struct {
	Query        interface{}
	Page         int
	PageSize     int
	SortProperty string
	SortRule     string
}
