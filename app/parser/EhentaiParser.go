package parser

// EhentaiParser support Ehentai Album download
type EhentaiParser struct{ BaseParser }

func (parser EhentaiParser) GetMeta() Meta {
	return Meta{
		URLRgx:   `^https://e-hentai.org/g/\w+/\w+/?$`,
		Priority: 9,

		Name:        "Ehentai下载",
		Version:     "0.1",
		Description: "可解析下载Ehentai相册",
		Author:      "",
		Link:        "",
	}
}
