package xyOrm

//Define kModel
type kModel struct{
	baseModel
	FID string
}

func (km *kModel) CreateFields() {
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