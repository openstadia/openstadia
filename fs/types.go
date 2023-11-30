package fs

type Command string

const (
	Ls Command = "LS"
)

type File struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"is_dir"`
	IsHidden bool   `json:"is_hidden"`
}

type LsReq struct {
	Path []string `json:"path"`
}

type LsRes struct {
	Files []File   `json:"files"`
	Path  []string `json:"path"`
}
