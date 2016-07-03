package models


type ComicInfo struct {

	ID int			`json:"id"`
	Title string		`json:"title"`
	Series string		`json:"series"`
	IssueNumber float32	`json:"issue_number"`
	PageCount int		`json:"page_count"`
	Credits credit		`json:"credits"`
	Volume string		`json:"volume"`
	Genres []string		`json:"genres"`
	DateAdded int		`json:"date_added"`
	PublishDate int		`json:"publish_date"`
	Synopsis string		`json:"synopsis"`
	Characters []string	`json:"characters"`
	Rating float32		`json:"rating"`
	Status string		`json:"status"`
	Bookmarks []int		`json:"bookmarks"`
	Other []string		`json:"other"` ////Other tags are formatted "[tagname]:[tag]" semicolon delimited
}

type ComicFile struct {
	ID int			`json:"id"`
	RelativePath string	`json:"relative_path"`
	AbsolutePath string	`json:"absolute_path"`
	Hash string		`json:"hash"` //MD5 hash
	FileSize int		`json:"filesize"`
}

type credit struct {
	Author string		`json:"author"`
	Artist string		`json:"artist"`
	Publisher string	`json:"publisher"`
	Other string		`json:"other"`

}

type User struct {
	ID int			`json:"id"`
	Name string		`json:"name"`
	Password string		`json:"-"`
	APIKeys []string	`json:"api_keys"`
}

type UserProgress struct {
	ID int			`json:"id"`
	ComicInfoID int		`json:"comic_info_id"`
	Read bool		`json:"read"`
	Completed bool		`json:"completed"`
	DateLastRead int	`json:"date_last_read"`
	DateCompleted int	`json:"date_completed"`
	LastReadPage int	`json:"last_read_page"`
}