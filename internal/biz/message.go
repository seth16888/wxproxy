package biz

import (
	"context"
	"fmt"

	wxError "github.com/seth16888/wxcommon/error"

	"github.com/seth16888/wxcommon/domain"
	"github.com/seth16888/wxcommon/helpers"
	. "github.com/seth16888/wxcommon/logger"
	"github.com/seth16888/wxcommon/paths"

	v1 "github.com/seth16888/wxproxy/api/v1"
)

// 发送订阅消息
type SubscribeMessage struct {
	TemplateMessage

	Scene string `json:"scene"`
	Title string `json:"title"`
}

// SendTemplateMessageRes 发送模版消息返回结果
type SendTemplateMessageRes struct {
	wxError.WXError

	MsgID int64 `json:"msgid"`
}

// GetBlockedMessagesReq 获取被屏蔽的消息请求
type GetBlockedMessagesReq struct {
	// 模板消息ID
	TmplMsgId string `json:"tmpl_msg_id" binding:"required" msg:"tmpl_msg_id required"`
	// 上一页查询结果最大的id，用于翻页，第一次传0
	LargestId int64 `json:"largest_id,omitempty"`
	// 单页查询的大小，最大100
	Limit int64 `json:"limit" binding:"required,gt=0,lte=100" msg:"limit required"`
}

// GetBlockedMessagesRes 获取被屏蔽的消息返回
type GetBlockedMessagesRes struct {
	wxError.WXError

	MsgInfo []BlockedMessageInfo `json:"msginfo"`
}

// BlockedMessageInfo 被屏蔽的消息
type BlockedMessageInfo struct {
	Id            int64  `json:"id"`
	TmplMsgId     string `json:"tmpl_msg_id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	SendTimestamp int64  `json:"send_timestamp"`
	OpenId        string `json:"openid"`
}

// TemplateMessage 发送的模板消息内容
type TemplateMessage struct {
	ToUser      string                       `json:"touser"`                  // 必须, 接受者OpenID
	TemplateID  string                       `json:"template_id"`             // 必须, 模版ID
	URL         string                       `json:"url,omitempty"`           // 可选, 用户点击后跳转的URL, 该URL必须处于开发者在公众平台网站中设置的域中
	Color       string                       `json:"color,omitempty"`         // 可选, 整个消息的颜色, 可以不设置
	Data        map[string]*TemplateDataItem `json:"data"`                    // 必须, 模板数据
	ClientMsgID string                       `json:"client_msg_id,omitempty"` // 可选, 防重入ID
	MiniProgram MiniProgram                  `json:"miniprogram,omitempty"`   // 可选,跳转至小程序地址
}

type MiniProgram struct {
	AppID    string `json:"appid"`    // 所需跳转到的小程序appid（该小程序appid必须与发模板消息的公众号是绑定关联关系）
	PagePath string `json:"pagepath"` // 所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar）
}

// TemplateDataItem 模版内某个 .DATA 的值
type TemplateDataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

// GetTemplateIndustryResp 获取设置的行业信息
type GetTemplateIndustryResp struct {
	PrimaryIndustry struct {
		FirstClass  string `json:"first_class"`
		SecondClass string `json:"second_class"`
	} `json:"primary_industry"`
	SecondaryIndustry struct {
		FirstClass  string `json:"first_class"`
		SecondClass string `json:"second_class"`
	} `json:"secondary_industry"`
}

// SetTemplateIndustryReq 设置所属行业
type SetTemplateIndustryReq struct {
	// 公众号模板消息所属行业编号
	IndustryId1 string `json:"industry_id1"`
	// 公众号模板消息所属行业编号
	IndustryId2 string `json:"industry_id2"`
}

// GetAllPrivateTemplateRes 获取模板列表返回
type GetAllPrivateTemplateRes struct {
	TemplateList []struct {
		TemplateID      string `json:"template_id"`
		Title           string `json:"title"`
		Content         string `json:"content"`
		Example         string `json:"example"`
		PrimaryIndustry string `json:"primary_industry"`
		DeputyIndustry  string `json:"deputy_industry"`
	} `json:"template_list"`
}

// GetTemplateIdReq 获取模板Id请求
type GetTemplateIdReq struct {
	TemplateIdShort string   `json:"template_id_short" binding:"required" msg:"template_id_short required"`
	KeywordNameList []string `json:"keyword_name_list" binding:"required,gt=0" msg:"keyword_name_list required"`
}

// GetTemplateIdRes 获取模板Id返回
type GetTemplateIdRes struct {
	wxError.WXError
	TemplateId string `json:"template_id"`
}

// GetIndustry
func (m *MPProxyUsecase) GetIndustry(ctx context.Context, token string) (*GetTemplateIndustryResp, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_Industry,
		token,
	)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[GetTemplateIndustryResp](resp, err)
	if wxErr != nil {
		Errorf("GetIndustry error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetIndustry error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	return rt, nil
}

// SetIndustry
func (m *MPProxyUsecase) SetIndustry(ctx context.Context, token string, id1 string, id2 string) (*wxError.WXError, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Set_Industry,
		token,
	)
	Debugf("url: %s", url)

	req := SetTemplateIndustryReq{IndustryId1: id1, IndustryId2: id2}
	body, err := helpers.BuildRequestBody[SetTemplateIndustryReq](req)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return nil, err
	}

	resp, err := m.hc.Post(url, "application/json", body)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SetIndustry error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("SetIndustry error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendSubscribeMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("SendSubscribeMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}

// GetAllPrivateTpl
func (m *MPProxyUsecase) GetAllPrivateTpl(ctx context.Context, token string) (*GetAllPrivateTemplateRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_AllPrivateTmpl,
		token,
	)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[GetAllPrivateTemplateRes](resp, err)
	if wxErr != nil {
		Errorf("GetAllPrivateTpl error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetAllPrivateTpl error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	return rt, nil
}

// GetMessageTplId
func (m *MPProxyUsecase) GetMessageTplId(ctx context.Context, token string, tmpId string, keywords []string) (*GetTemplateIdRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_TemplateId,
		token,
	)
	Debugf("url: %s", url)
	type GetTemplateIdReq struct {
		TemplateIDShort string   `json:"template_id_short"`
		KeywordList     []string `json:"keyword_name_list"`
	}
	req := GetTemplateIdReq{TemplateIDShort: tmpId, KeywordList: keywords}

	body, err := helpers.BuildRequestBody[GetTemplateIdReq](req)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return nil, err
	}

	resp, err := m.hc.Post(url, "application/json", body)
	rt, wxErr := helpers.BuildHttpResponse[GetTemplateIdRes](resp, err)
	if wxErr != nil {
		Errorf("GetMessageTplId error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetMessageTplId error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	return rt, nil
}

// DeleteMessageTpl
func (m *MPProxyUsecase) DeleteMessageTpl(ctx context.Context, token string, tmpId string) (*wxError.WXError, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Del_Template,
		token,
	)
	Debugf("url: %s", url)

	type DelTemplateReq struct {
		TemplateID string `json:"template_id"`
	}
	req := DelTemplateReq{TemplateID: tmpId}
	body, err := helpers.BuildRequestBody[DelTemplateReq](req)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return nil, err
	}

	resp, err := m.hc.Post(url, "application/json", body)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("DeleteMessageTpl error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("DeleteMessageTpl error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendSubscribeMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("SendSubscribeMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}

// SendTplMsg
func (m *MPProxyUsecase) SendTplMsg(ctx context.Context, token string, req *v1.SendTplMsgRequest) (*SendTemplateMessageRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Send_TemplateMessage,
		token,
	)
	Debugf("url: %s", url)

	var body = &TemplateMessage{
		ToUser:      req.Touser,
		TemplateID:  req.TemplateId,
		URL:         req.Url,
		Color:       "",
		Data:        map[string]*TemplateDataItem{},
		ClientMsgID: req.ClientMsgId,
	}
	if req.Miniprogram != nil {
		body.MiniProgram = MiniProgram{
			AppID:    req.Miniprogram.Appid,
			PagePath: req.Miniprogram.PagePath,
		}
	}
	for name, item := range req.Data {
		body.Data[name] = &TemplateDataItem{
			Value: item.Value,
			Color: item.Color,
		}
	}
	Debugf("body: %+v", body)

	bodyReader, err := helpers.BuildRequestBody[*TemplateMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return nil, err
	}

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	rt, wxErr := helpers.BuildHttpResponse[SendTemplateMessageRes](resp, err)
	if wxErr != nil {
		Errorf("SendTplMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("SendTplMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendSubscribeMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("SendSubscribeMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}

// SendSubscribeMsg
func (m *MPProxyUsecase) SendSubscribeMsg(ctx context.Context, token string, req *v1.SendSubscribeMsgRequest) (*wxError.WXError, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Send_SubscribeMessage,
		token,
	)
	Debugf("url: %s", url)

	var body = &TemplateMessage{
		ToUser:      req.Touser,
		TemplateID:  req.TemplateId,
		URL:         req.Url,
		Color:       "",
		Data:        map[string]*TemplateDataItem{},
		ClientMsgID: req.ClientMsgId,
	}
	if req.Miniprogram != nil {
		body.MiniProgram = MiniProgram{
			AppID:    req.Miniprogram.Appid,
			PagePath: req.Miniprogram.PagePath,
		}
	}
	for name, item := range req.Data {
		body.Data[name] = &TemplateDataItem{
			Value: item.Value,
			Color: item.Color,
		}
	}
	Debugf("body: %+v", body)

	bodyReader, err := helpers.BuildRequestBody[*TemplateMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return nil, err
	}

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendSubscribeMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("SendSubscribeMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendSubscribeMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("SendSubscribeMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}

// GetBlockedTplMsg
func (m *MPProxyUsecase) GetBlockedTplMsg(ctx context.Context, token string,
	msgId string, largest int64, limit int64) (*GetBlockedMessagesRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_BlockedMsg,
		token,
	)
	Debugf("url: %s", url)

	body := &GetBlockedMessagesReq{TmplMsgId: msgId, LargestId: largest, Limit: limit}
	bodyReader, err := helpers.BuildRequestBody[*GetBlockedMessagesReq](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return nil, err
	}

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	rt, wxErr := helpers.BuildHttpResponse[GetBlockedMessagesRes](resp, err)
	if wxErr != nil {
		Errorf("GetBlockedTplMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetBlockedTplMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendSubscribeMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("SendSubscribeMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}
