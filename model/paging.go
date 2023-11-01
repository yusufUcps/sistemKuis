package model

type Paging struct {
    Page    int `json:"page"`
	Page_size int `json:"page_size"`
    Max_page int64 `json:"max_page"`	
}

func ConvertPaging(page int , pageSize int, maxData int64) *Paging {

    maxPage:= maxData / int64(pageSize)
    if maxData%int64(pageSize) != 0{
        maxPage = maxPage + 1
    }

    paging := Paging{
       Page: page,
       Page_size: pageSize,
       Max_page: maxPage,
    }
    return &paging
}
