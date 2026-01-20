package xyOrm

//Define kModel
type KModel struct{
	baseModel
	FID string
}

func (km *KModel) CreateFields() {
	AssignFieldNames(km)
	km.Fields = append(
		km.Fields, 
		TbField{
			Name: km.FID,
			Type: "string",
			Length: 50,
			IsKey: true,
		})
}