package pixiv

type PixivArtworks struct {
	Error bool `json:"error"`
	Body  struct {
		Illustmanga struct {
			Data           []interface{} `json:"data"`
			Total          int           `json:"total"`
			Bookmarkranges []struct {
				Min interface{} `json:"min"`
				Max interface{} `json:"max"`
			} `json:"bookmarkRanges"`
		} `json:"illustManga"`
		Popular struct {
			Recent    []interface{} `json:"recent"`
			Permanent []interface{} `json:"permanent"`
		} `json:"popular"`
		Relatedtags []interface{} `json:"relatedTags"`
		Zoneconfig  struct {
			Header struct {
				URL string `json:"url"`
			} `json:"header"`
			Footer struct {
				URL string `json:"url"`
			} `json:"footer"`
			Infeed struct {
				URL string `json:"url"`
			} `json:"infeed"`
		} `json:"zoneConfig"`
		Extradata struct {
			Meta struct {
				Title              string `json:"title"`
				Description        string `json:"description"`
				Canonical          string `json:"canonical"`
				Alternatelanguages struct {
					Ja string `json:"ja"`
					En string `json:"en"`
				} `json:"alternateLanguages"`
				Descriptionheader string `json:"descriptionHeader"`
			} `json:"meta"`
		} `json:"extraData"`
	} `json:"body"`
}
