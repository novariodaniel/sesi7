package datastruct

type Persegi struct{
	Luas int
	Keliling int
}

type ParamRequest struct {
	Sisi int `json:"sisi"`
}

type ResponseMath struct {
	Status int
	Desc string
}