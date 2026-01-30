package xyOrm

import "errors"

//Define kModel
type KModel struct{
	baseModel
	FID string
}

func (km *KModel) CreateFields() {
	km.AssignFieldNames(km)
	km.baseModel.CreateFields()
	km.Fields = append(
		km.Fields, 
		TbField{
			Name: km.FID,
			Type: "string",
			Length: 50,
			IsKey: true,
		})
}

//query data
func (km *KModel) SelectByPk(value string) (*map[string]interface{}, error) {
	var rList *[]map[string]interface{}
	var err error
	if rList, err = km.baseModel.SelectByField(km.FID, value); err != nil {
		return nil, err
	}
	if len((*rList)) < 1{
		return nil, errors.New("No record found")
	}
	return &(*rList)[0], nil
}

//update data
func (km *KModel) UpdateByByPk(recordMap map[string]string, value string) error {
	return km.baseModel.UpdateByField(recordMap, km.FID, value)
}

//delete data
func (km *KModel) DeleteByByPk(value string) error {
	return km.baseModel.DeleteByField(km.FID, value)
}