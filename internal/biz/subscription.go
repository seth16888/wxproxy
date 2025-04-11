package biz

import (
	"context"
	"fmt"

	"github.com/seth16888/wxcommon/domain"
	wxError "github.com/seth16888/wxcommon/error"
	"github.com/seth16888/wxcommon/helpers"
	"github.com/seth16888/wxcommon/paths"
	v1 "github.com/seth16888/wxproxy/api/v1"

	. "github.com/seth16888/wxcommon/logger"
)

// SendSubscriptionMessageReq 发送订阅消息
type SendSubscriptionMessageReq struct {
	Touser      string                                   `json:"touser" binding:"required"`      // 接收者（用户）的 openid
	TemplateId  string                                   `json:"template_id" binding:"required"` // 所需下发的订阅模板id
	Page        string                                   `json:"page,omitempty"`                 // 点击模板卡片后的跳转页面
	Data        map[string]*SendSubscribeMessageDataItem `json:"data" binding:"required"`        // 模板内容，格式形如 { "key1": { "value": any }, "key2": { "
	MiniProgram struct {
		AppID    string `json:"appid"`    // 所需跳转到的小程序appid（该小程序appid必须与发模板消息的公众号是绑定关联关系）
		PagePath string `json:"pagepath"` // 所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar）
	} `json:"miniprogram,omitempty"` // 可选,跳转至小程序地址
}

type SendSubscribeMessageDataItem struct {
	Value string `json:"value"`
}

// GetPrivateTemplateListRes 获取当前帐号下的个人模板列表
type GetPrivateTemplateListRes struct {
	wxError.WXError
	Data []struct {
		PriTmplId string `json:"priTmplId"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		Example   string `json:"example"`
		Type      int    `json:"type"`
	} `json:"data"`
}

// GetPubTemplateTitlesRes 获取类目下的公共模板
type GetPubTemplateTitlesRes struct {
	wxError.WXError
	Count int `json:"count"`
	Data  []struct {
		Tid        string `json:"tid"`
		Title      string `json:"title"`
		Type       int    `json:"type"`
		CategoryId string `json:"categoryId"`
	} `json:"data"`
}

// GetPubTemplateKeyWordsRes 获取模板中的关键词
type GetPubTemplateKeyWordsRes struct {
	wxError.WXError
	Count int `json:"count"`
	Data  []struct {
		Kid     int    `json:"kid"`
		Name    string `json:"name"`
		Rule    string `json:"rule"`
		Example string `json:"example"`
	} `json:"data"`
}

type GetSubscribeCategoryRes struct {
	Data []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

func (m *MPProxyUsecase) AddSubscribeTpl(ctx context.Context, token string,
	tid string, sceneDesc string, kidList []int64) (string, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Add_Template, token)
	Debugf("url: %s", url)

	type AddSubscribeTplReq struct {
		Tid       string  `json:"tid"`
		SceneDesc string  `json:"sceneDesc"`
		KidList   []int64 `json:"kidList"`
	}
	req := &AddSubscribeTplReq{
		Tid:       tid,
		SceneDesc: sceneDesc,
		KidList:   kidList,
	}

	reader, err := helpers.BuildRequestBody[*AddSubscribeTplReq](req)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return "", err
	}

	resp, err := m.hc.Post(url, "application/json", reader)

	type resultT struct {
		wxError.WXError
		TmplID string `json:"priTmplId"`
	}
	rt, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("AddSubscribeTpl error: %d %s", wxErr.ErrCode, wxErr.Error())
		return "", fmt.Errorf("AddSubscribeTpl error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("AddSubscribeTpl error: %d %s", rt.ErrCode, rt.ErrMsg)
		return "", fmt.Errorf("AddSubscribeTpl error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt.TmplID, nil
}

// GetSubscribeCategory
func (m *MPProxyUsecase) GetSubscribeCategory(ctx context.Context, token string) (*GetSubscribeCategoryRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Get_Category, token)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[GetSubscribeCategoryRes](resp, err)
	if wxErr != nil {
		Errorf("GetSubscribeCategory error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetSubscribeCategory error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	return rt, nil
}

func (m *MPProxyUsecase) DelSubscribeTpl(ctx context.Context, token string, tplId string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Del_Subscription_Template, token)
	Debugf("url: %s", url)

	type DelSubscribeTplReq struct {
		PriTmplId string `json:"priTmplId"`
	}
	req := &DelSubscribeTplReq{
		PriTmplId: tplId,
	}

	reader, err := helpers.BuildRequestBody[*DelSubscribeTplReq](req)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("DelSubscribeTpl error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("DelSubscribeTpl error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("DelSubscribeTpl error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("DelSubscribeTpl error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}
func (m *MPProxyUsecase) GetSubscribeTplKeywords(ctx context.Context, token string, tplId string) (*GetPubTemplateKeyWordsRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s&tid=%s", domain.GetWXAPIDomain(), paths.Path_Get_PubTpl_KeyWorks, token, tplId)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[GetPubTemplateKeyWordsRes](resp, err)
	if wxErr != nil {
		Errorf("GetSubscribeTplKeywords error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetSubscribeTplKeywords error: %d %s", wxErr.ErrCode, wxErr.Error())
	}
	if rt.ErrCode != 0 {
		Errorf("GetSubscribeTplKeywords error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("GetSubscribeTplKeywords error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}
func (m *MPProxyUsecase) GetSubscribeTplTitles(ctx context.Context, token string,
  ids string, start int64, limit int64) (*GetPubTemplateTitlesRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s&ids=%s&start=%d&limit=%d",
    domain.GetWXAPIDomain(), paths.Path_Get_PubTpl_Titles, token, ids, start, limit)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[GetPubTemplateTitlesRes](resp, err)
	if wxErr != nil {
		Errorf("GetSubscribeTplTitles error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetSubscribeTplTitles error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("GetSubscribeTplTitles error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("GetSubscribeTplTitles error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}
func (m *MPProxyUsecase) GetSubscribePrivateTpl(ctx context.Context, token string) (*GetPrivateTemplateListRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Get_AllPrivateTmpl, token)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[GetPrivateTemplateListRes](resp, err)
	if wxErr != nil {
		Errorf("GetSubscribePrivateTpl error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetSubscribePrivateTpl error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("GetSubscribePrivateTpl error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("GetSubscribePrivateTpl error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}
func (m *MPProxyUsecase) SendSubscribeMessage(ctx context.Context, req *v1.SendSubscribeMessageRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Send_SubscribeMessage, req.AccessToken)
	Debugf("url: %s", url)

	var params = &SendSubscriptionMessageReq{
		Touser:     req.Touser,
		TemplateId: req.TemplateId,
		Page:       req.Page,
		Data:       map[string]*SendSubscribeMessageDataItem{},
	}
	for k, v := range req.Data {
		params.Data[k] = &SendSubscribeMessageDataItem{
			Value: v.Value,
		}
	}
	if req.Miniprogram != nil {
		params.MiniProgram.AppID = req.Miniprogram.Appid
		params.MiniProgram.PagePath = req.Miniprogram.PagePath
	}

	reader, err := helpers.BuildRequestBody[*SendSubscriptionMessageReq](params)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendSubscribeMessage error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendSubscribeMessage error: %d %s", wxErr.ErrCode, wxErr.Error())
	}
	if rt.ErrCode != 0 {
		Errorf("SendSubscribeMessage error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendSubscribeMessage error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}
