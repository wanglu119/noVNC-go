package api

import (
	"github.com/spf13/afero"

	"github.com/wanglu119/me-deps/webCommon"
)


type data struct {
	*webCommon.WebDataImplPart
}

func CreateWebData() webCommon.WebData {
	return &data{
		&webCommon.WebDataImplPart{},
	}
}

func (d *data) GetAuthToken() string {
	return ""
}

func (d *data) SetAuthData(aData interface{}) {
}
func (d *data) GetFs() afero.Fs {
	log.Error("not impl GetFs")
	return nil
}
func (d *data) GetAuthData() interface{} {
	return nil
}
