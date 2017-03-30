package models

/*
MGA database structures used for internal stuff
*/
type ComicInfo struct {
	ID          int      `json:"id"`
	Hash        string   `json:"hash"`
	Title       string   `json:"title"`
	Series      string   `json:"series"`
	IssueNumber float32  `json:"issue_number"`
	PageCount   int      `json:"page_count"`
	Credits     Credit   `json:"credits"`
	Volume      string   `json:"volume"`
	Genres      []string `json:"genres"`
	DateAdded   int64    `json:"date_added"`
	PublishDate int64    `json:"publish_date"`
	Synopsis    string   `json:"synopsis"`
	Characters  []string `json:"characters"`
	Rating      float32  `json:"rating"`
	Status      string   `json:"status"`
	Bookmarks   []int    `json:"bookmarks"`
	Other       []string `json:"other"` //Other tags are formatted "[tagname]:[tag]" semicolon delimited
}

type ComicFile struct {
	ID           int    `json:"id"`
	RelativePath string `json:"relative_path"`
	AbsolutePath string `json:"absolute_path"`
	FileName     string `json:"file_name"`
	Hash         string `json:"hash"` //MD5 hash
	FileSize     int64  `json:"filesize"`
}

type ComicWrapper struct {
	ComicInfo ComicInfo `json:"comic_info"`
	ComicFile ComicFile `json:"comic_file"`
}

type Credit struct {
	Author    string `json:"author"`
	Artist    string `json:"artist"`
	Publisher string `json:"publisher"`
	Other     string `json:"other"`
}

type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Password string   `json:"-"`
	APIKeys  []string `json:"api_keys"`
}

type UserProgress struct {
	ID            int   `json:"id"`
	ComicInfoID   int   `json:"comic_info_id"`
	Read          bool  `json:"read"`
	Completed     bool  `json:"completed"`
	DateLastRead  int64 `json:"date_last_read"`
	DateCompleted int64 `json:"date_completed"`
	LastReadPage  int   `json:"last_read_page"`
}

//Used to create a pseudo directory structure.
//All comics must belong to a category. The default behavior
//is to make each folder that is walked a category
type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Parent   string `json:"parent"`
	ParentId int    `json:"parent_id"`
	IsRoot   bool   `json:"is_root"`
	Full     string `json:"full"`
}

/*
ComicStreamer compatibility structures
*/
type CSComicWrapper struct {
	CSComicInfo CSComic




}
type CSComic struct {
	AddedTs      string   `json:"added_ts"`
	Series       string   `json:"series"`
	PageCount    int      `json:"page_count"`
	Locations    []string `json:"locations"`
	Month        string   `json:"month"`
	Imprint      string   `json:"imprint"`
	Year         string   `json:"year"`
	Id           int      `json:"id"`
	DeletedTs    string   `json:"deleted_ts"`
	Genres       []string `json:"genres"`
	Title        string   `json:"title"`
	Comments     string   `json:"comments"`
	Filesize     int      `json:"filesize"`
	Issue        string   `json:"issue"`
	Hash         string   `json:"hash"`
	StoryArcs    []string `json:"storyarcs"`
	ModTs        string   `json:"mod_ts"`
	LastReadPage int      `json:"last_read_page"`
	Weblink      string   `json:"weblink"`
	Volume       string   `json:"volume"`
	Credits      Credit   `json:"credits,omitempty"`
	GenericTags  []string `json:"generictags"`
	Characters   []string `json:"characters"`
	LastReadTs   string   `json:"lastread_ts"`
	Date         string   `json:"date"`
	Path         string   `json:"path"`
	Day          string   `json:"day"`
	Publisher    string   `json:"publisher"`
	Teams        []string `json:"teams"`
}

type CSComicResult struct {
	Comics     []CSComic `json:"comics"`
	TotalCount int       `json:"total_count"`
	PageCount  int       `json:"page_count"`
}

type CSFolderResponse struct {
	Current string               `json:"current"`
	Folders []CSFolder           `json:"folders"`
	Comics  CSComicCountResponse `json:"comics"`
}

type CSComicCountResponse struct {
	Count    int    `json:"count"`
	URL_Path string `json:"url_path"`
}

type CSFolder struct {
	URL_Path string `json:"url_path"`
	Name     string `json:"name"`
}
