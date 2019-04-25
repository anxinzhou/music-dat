package http

type ChainBalanceController struct {
	ContractController
}

func (this *ChainBalanceController) Get() {
	kind:=this.Ctx.Input.Param(":kind")
	if kind == "avatar" {

	} else {

	}
	//this.
	//this.ServeJSON()
}
