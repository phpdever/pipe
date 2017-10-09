// Solo.go - A small and beautiful blogging platform written in golang.
// Copyright (C) 2017, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package console defines console controllers.
package console

import (
	"net/http"
	"strconv"

	"github.com/b3log/solo.go/model"
	"github.com/b3log/solo.go/service"
	"github.com/b3log/solo.go/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetBasicSettingsCtl(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	sessionData := util.GetSession(c)
	settings := service.Setting.GetAllSettings(sessionData.BID, model.SettingCategoryBasic)
	data := map[string]interface{}{}
	for _, setting := range settings {
		if model.SettingNameBasicCommentable == setting.Name {
			v, err := strconv.ParseBool(setting.Value)
			if nil != err {
				log.Errorf("value of basic setting [name=%s] must be \"true\" or \"false\"", setting.Name)
				data[setting.Name] = true
			} else {
				data[setting.Name] = v
			}
		} else {
			data[setting.Name] = setting.Value
		}
	}
	delete(data, model.SettingNameBasicPath)
	result.Data = data
}

func UpdateBasicSettingsCtl(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	args := map[string]interface{}{}
	if err := c.BindJSON(&args); nil != err {
		result.Code = -1
		result.Msg = "parses update basic settings request failed"

		return
	}

	sessionData := util.GetSession(c)
	basics := []*model.Setting{}
	for k, v := range args {
		var value interface{}
		switch v.(type) {
		case bool:
			value = strconv.FormatBool(v.(bool))
		default:
			value = v.(string)
		}

		basic := &model.Setting{
			Category: model.SettingCategoryBasic,
			BlogID:   sessionData.BID,
			Name:     k,
			Value:    value.(string),
		}

		basics = append(basics, basic)
	}

	if err := service.Setting.UpdateSettings(model.SettingCategoryBasic, basics); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}

func GetPreferenceSettingsCtl(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	sessionData := util.GetSession(c)
	settings := service.Setting.GetAllSettings(sessionData.BID, model.SettingCategoryPreference)
	data := map[string]interface{}{}
	for _, setting := range settings {
		if model.SettingNamePreferenceArticleListStyle != setting.Name {
			v, err := strconv.ParseInt(setting.Value, 10, 64)
			if nil != err {
				log.Errorf("value of preference setting [name=%s] must be an integer", setting.Name)
				data[setting.Name] = 10
			} else {
				data[setting.Name] = v
			}
		} else {
			data[setting.Name] = setting.Value
		}
	}
	delete(data, model.SettingNamePreferenceTheme)
	delete(data, model.SettingNamePreferenceVer)
	result.Data = data
}

func UpdatePreferenceSettingsCtl(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	args := map[string]interface{}{}
	if err := c.BindJSON(&args); nil != err {
		result.Code = -1
		result.Msg = "parses update preference settings request failed"

		return
	}

	sessionData := util.GetSession(c)
	prefs := []*model.Setting{}
	for k, v := range args {
		var value interface{}
		switch v.(type) {
		case float64:
			value = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		default:
			value = v.(string)
		}

		pref := &model.Setting{
			Category: model.SettingCategoryPreference,
			BlogID:   sessionData.BID,
			Name:     k,
			Value:    value.(string),
		}

		prefs = append(prefs, pref)
	}

	if err := service.Setting.UpdateSettings(model.SettingCategoryPreference, prefs); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}
