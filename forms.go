package main

type ApplyForm struct {
	Nickname string `binding:"Required"`
	Email    string `binding:"Email"`
	Country  string `binding:"Required"`
	City     string `binding:"Required"`
}

type SignInForm struct {
	Password string `binding:"Required"`
}

// func (cf ApplyForm) Error(ctx *macaron.Context, errs binding.Errors) {
// 	ctx.Data["Errors"] = errs
// }
