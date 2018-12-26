package fulfilment

var FF *FFClient

type FFClient struct {
}

func InitFF() {
	FF = &FFClient{}
}
