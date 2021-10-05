package datastruct


// Jawaban No 2
type DetailLanguage struct {
	Language string `json:"language"`
	Appeared int `json:"appeared"`
	Created []string `json:"created"`
	Functional bool `json:"functional"`
	Objectorient bool `json:"object-oriented"`
	Relation Influen `json:"relation"`
}

type Influen struct {
	InfluencedBy []string `json:"influenced-by"`
	Influences []string `json:"influences"`
}

type ResponseProcess struct{
	// Status int `json:"status"`
	Desc string `json:"desc"`
}

type ResponsePalindrome struct{
	Status int `json:"status"`
	Desc string `json:"desc"`
}

type Languages struct {
	Languages []DetailLanguage
}

type ResponseStore struct{
	Status int `json:"status"`
	Desc string `json:"desc"`
}

type StoredData struct{
	ID int
	ListDetailLang DetailLanguage
}

type RequestStore struct{
	Language string `json:"language"`
	Appeared int `json:"appeared"`
	Created []string `json:"created"`
	Functional bool `json:"functional"`
	Objectorient bool `json:"object-oriented"`
	Relation Influen `json:"relation"`
}