package schema

type RegisterReq struct{
	Email 		string `json:"email" binding:"required"`
	Password 	string `json:"password" binding:"required"`
}

type RegisterRes struct{
	Message string `json:"message"`
}

type LoginReq struct{
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

type LoginRes struct{
	Message string `json:"message"`
	AccessToken string `json:"access_token"`
}