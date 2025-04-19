package biz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/seth16888/wxcommon/domain"
	wxError "github.com/seth16888/wxcommon/error"
	"github.com/seth16888/wxcommon/hc"
	"github.com/seth16888/wxcommon/helpers"
	. "github.com/seth16888/wxcommon/logger"
	"github.com/seth16888/wxcommon/paths"
	v1 "github.com/seth16888/wxproxy/api/v1"
	"go.uber.org/zap"
)

const (
	ActionId  = "QR_SCENE"
	ActionStr = "QR_STR_SCENE"

	ActionLimitId  = "QR_LIMIT_SCENE"
	ActionLimitStr = "QR_LIMIT_STR_SCENE"
)

// MenuRes 获取自定义菜单的返回数据
type MenuRes struct {
	Menu struct {
		Button []Button `json:"button,omitempty"`
		MenuID int64    `json:"menuid,omitempty"`
	} `json:"menu,omitempty"`
	Conditionalmenu []conditionalMenuRes `json:"conditionalmenu,omitempty"`
}

// ConditionalMenu 个性化菜单返回结果
type conditionalMenuRes struct {
	Button    []Button  `json:"button"`
	MatchRule MatchRule `json:"matchrule"`
	MenuID    int64     `json:"menuid"`
}

// Button 菜单按钮
type Button struct {
	Type       string    `json:"type,omitempty"`
	Name       string    `json:"name,omitempty"`
	Key        string    `json:"key,omitempty"`
	URL        string    `json:"url,omitempty"`
	MediaID    string    `json:"media_id,omitempty"`
	AppID      string    `json:"appid,omitempty"`
	PagePath   string    `json:"pagepath,omitempty"`
	SubButtons []*Button `json:"sub_button,omitempty"`
}

// MatchRule 个性化菜单规则
type MatchRule struct {
	TagId              string `json:"tag_id,omitempty"`
	ClientPlatformType string `json:"client_platform_type,omitempty"`
}

type FetchShortenRes struct {
	LongData      string `json:"long_data"`
	ExpireSeconds int64  `json:"expire_seconds"`
	CreateTime    int64  `json:"create_time"`
}

type CreateQRCodeReq struct {
	ActionInfo struct {
		Scene struct {
			SceneId  int64  `json:"scene_id,omitempty"`
			SceneStr string `json:"scene_str,omitempty"`
		} `json:"scene"`
	} `json:"action_info"`

	ActionName    string `json:"action_name"`
	ExpireSeconds int64  `json:"expire_seconds"`
}

type Ticket struct {
	Ticket        string `json:"ticket"`
	Url           string `json:"url"`
	ExpireSeconds int64  `json:"expire_seconds"`
}

// TagMembersRes 从WX获取标签下的用户列表返回结果
type TagMembersRes struct {
	Count int64 `json:"count"`
	Data  struct {
		OpenIdList []string `json:"openid"`
	} `json:"data"`
	NextOpenId string `json:"next_openid"`
}

type Tag struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Count int64  `json:"count,omitempty"`
}

type GetMaterialCountReply struct {
	VoiceCount int64 `json:"voice_count"`
	VideoCount int64 `json:"video_count"`
	ImageCount int64 `json:"image_count"`
	NewsCount  int64 `json:"news_count"`
}

type GetMaterialListRes struct {
	TotalCount int64          `json:"total_count"`
	ItemCount  int64          `json:"item_count"`
	Item       []MaterialItem `json:"item"`
}

type MaterialItem struct {
	MediaId    string `json:"media_id"`
	Name       string `json:"name"`
	UpdateTime int64  `json:"update_time"`
	Url        string `json:"url"`
}

type GetMemberListRes struct {
	Total      int64  `json:"total"`
	Count      int64  `json:"count"`
	NextOpenid string `json:"next_openid"`
	Data       struct {
		OpenidList []string `json:"openid"`
	} `json:"data"`
}

type GetMemberInfoRes struct {
	Subscribe      int64
	Openid         string
	SubscribeTime  int64
	Unionid        string
	Remark         string
	Groupid        int64
	TagidList      []int64
	SubscribeScene string
	QrScene        int64
	QrSceneStr     string
	Language       string
}

type MPProxyUsecase struct {
	log *zap.Logger
	hc  *hc.Client
}

func (m *MPProxyUsecase) UnBlockMember(ctx context.Context, token string, ids []string) *v1.WXErrorReply {
  url := fmt.Sprintf("https://%s%s?access_token=%s",
  domain.GetWXAPIDomain(),
  paths.Path_Batch_Remove_Black_List,
  token,
)
m.log.Debug("url", zap.String("url", url))

type blockMemberReq struct {
  OpenIds []string `json:"openid_list"`
}
req:= &blockMemberReq{
  OpenIds: ids,
}

reader,err:= helpers.BuildRequestBody(req)
if err != nil {
  m.log.Error("build request body error", zap.Error(err))
  return &v1.WXErrorReply{Errcode: 500, Errmsg: "build request body error"}
}

resp, err := m.hc.Post(url, "application/json", reader)
rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
if wxErr != nil {
  m.log.Error("UpdateKFTyping error", zap.Error(err))
  return &v1.WXErrorReply{Errcode: 500, Errmsg: "UpdateKFTyping error"}
}

return &v1.WXErrorReply{Errcode: int64(rt.ErrCode), Errmsg: rt.ErrMsg}
}

func (m *MPProxyUsecase) BlockMember(ctx context.Context, token string, ids []string) *v1.WXErrorReply {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Batch_Add_Black_List,
		token,
	)
	m.log.Debug("url", zap.String("url", url))

  type blockMemberReq struct {
    OpenIds []string `json:"openid_list"`
  }
  req:= &blockMemberReq{
    OpenIds: ids,
  }

  reader,err:= helpers.BuildRequestBody(req)
  if err != nil {
    m.log.Error("build request body error", zap.Error(err))
    return &v1.WXErrorReply{Errcode: 500, Errmsg: "build request body error"}
  }

  resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		m.log.Error("UpdateKFTyping error", zap.Error(err))
		return &v1.WXErrorReply{Errcode: 500, Errmsg: "UpdateKFTyping error"}
	}

	return &v1.WXErrorReply{Errcode: int64(rt.ErrCode), Errmsg: rt.ErrMsg}
}

func NewMPProxyUsecase(hc *hc.Client, logger *zap.Logger) *MPProxyUsecase {
	return &MPProxyUsecase{
		hc:  hc,
		log: logger,
	}
}

// GetMaterialCoount
func (m *MPProxyUsecase) GetMaterialCoount(ctx context.Context, token string) (*GetMaterialCountReply, error) {
	// url
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_Material_Count,
		token,
	)
	m.log.Debug("url", zap.String("url", url))

	resp, err := m.hc.Get(url)
	if err != nil {
		m.log.Error("GetMaterialCoount error", zap.Error(err))
		return nil, err
	}

	result, wxErr := helpers.BuildHttpResponse[GetMaterialCountReply](resp, err)
	if wxErr != nil {
		m.log.Error("GetMaterialCoount error", zap.Error(err))
		return nil, fmt.Errorf("GetMaterialCoount error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	return result, nil
}

func (m *MPProxyUsecase) GetMaterialList(ctx context.Context, token string, mediaType string,
	offset int64, count int64,
) (*GetMaterialListRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Batch_Get_Material,
		token,
	)
	m.log.Debug("url", zap.String("url", url))

	body := map[string]interface{}{
		"type":   mediaType,
		"offset": offset,
		"count":  count,
	}
	m.log.Debug("body", zap.Any("body", body))
	bodyJson, err := json.Marshal(body)
	if err != nil {
		m.log.Error("GetMaterialList error", zap.Error(err))
		return nil, err
	}
	bodyReader := bytes.NewReader(bodyJson)

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		m.log.Error("GetMaterialList error", zap.Error(err))
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		GetMaterialListRes
	}

	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		m.log.Error("GetMaterialList error", zap.Error(err))
		return nil, fmt.Errorf("GetMaterialList error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		m.log.Error("GetMaterialList error", zap.Error(err))
		return nil, fmt.Errorf("GetMaterialList error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return &result.GetMaterialListRes, nil
}

func (m *MPProxyUsecase) GetMemberList(ctx context.Context, token string, openid string) (*GetMemberListRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s&next_openid=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_Member_List,
		token,
		openid,
	)
	m.log.Debug("url", zap.String("url", url))

	resp, err := m.hc.Get(url)
	if err != nil {
		m.log.Error("GetMemberList error", zap.Error(err))
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		GetMemberListRes
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		m.log.Error("GetMemberList error:", zap.Error(err))
		return nil, fmt.Errorf("GetMemberList error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		m.log.Error("GetMemberList error", zap.Error(err))
		return nil, fmt.Errorf("GetMemberList error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return &result.GetMemberListRes, nil
}

func (m *MPProxyUsecase) GetMemberInfo(ctx context.Context, token string, openid string, lang string) (*GetMemberInfoRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s&openid=%s&lang=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_Member_Info,
		token,
		openid,
		lang,
	)
	m.log.Debug("url", zap.String("url", url))

	resp, err := m.hc.Get(url)
	if err != nil {
		m.log.Error("GetMemberInfo error", zap.Error(err))
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		GetMemberInfoRes
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		m.log.Error("GetMemberInfo error", zap.Error(err))
		return nil, fmt.Errorf("GetMemberInfo error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		m.log.Error("GetMemberInfo error", zap.Error(err))
		return nil, fmt.Errorf("GetMemberInfo error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return &result.GetMemberInfoRes, nil
}

func (m *MPProxyUsecase) BatchGetMemberInfo(ctx context.Context, token string, openidList []string) (*[]GetMemberInfoRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Batch_Get_Member_Info,
		token,
	)
	m.log.Debug("url", zap.String("url", url))

	type batchGetMemberInfoReq struct {
		UserList []struct {
			Openid string `json:"openid"`
		} `json:"user_list"`
	}
	body := batchGetMemberInfoReq{}
	for _, openid := range openidList {
		body.UserList = append(body.UserList, struct {
			Openid string `json:"openid"`
		}{
			Openid: openid,
		})
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		m.log.Error("BatchGetMemberInfo error", zap.Error(err))
		return nil, err
	}
	bodyReader := bytes.NewReader(bodyJson)

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		m.log.Error("BatchGetMemberInfo error", zap.Error(err))
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		UserInfoList []GetMemberInfoRes `json:"user_info_list"`
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		m.log.Error("BatchGetMemberInfo error", zap.Error(wxErr))
		return nil, fmt.Errorf("BatchGetMemberInfo error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		m.log.Error("BatchGetMemberInfo error", zap.Any("result", result))
		return nil, fmt.Errorf("BatchGetMemberInfo error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return &result.UserInfoList, nil
}

func (m *MPProxyUsecase) GetMemberTags(ctx context.Context, token string, openid string) ([]int64, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_Member_Tags,
		token,
	)
	Debugf("url: %s", url)

	body := map[string]string{
		"openid": openid,
	}
	bodyJson, err := json.Marshal(body) // TODO: helpers.BuildRequestBody
	if err != nil {
		Errorf("GetMemberTags error: %s", err.Error())
		return nil, err
	}
	bodyReader := bytes.NewReader(bodyJson)

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		Errorf("GetMemberTags error: %s", err.Error())
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		TagidList []int64 `json:"tagid_list"`
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("GetMemberTags error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetMemberTags error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		Errorf("GetMemberTags error: %d %s", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("GetMemberTags error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return result.TagidList, nil
}

func (m *MPProxyUsecase) UpdateMemberRemark(ctx context.Context, token string, openid string, remark string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Update_Member_Remark,
		token,
	)
	Debugf("url: %s", url)

	body := map[string]string{
		"openid": openid,
		"remark": remark,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		Errorf("UpdateMemberRemark error: %s", err.Error())
		return err
	}
	bodyReader := bytes.NewReader(bodyJson)

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		Errorf("UpdateMemberRemark error: %s", err.Error())
		return err
	}

	type resultT struct {
		wxError.WXError
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("UpdateMemberRemark error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("UpdateMemberRemark error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		Errorf("UpdateMemberRemark error: %d %s", result.ErrCode, result.ErrMsg)
		return fmt.Errorf("UpdateMemberRemark error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

func (m *MPProxyUsecase) GetTagList(ctx context.Context, token string) ([]Tag, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_Tags,
		token,
	)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	if err != nil {
		Errorf("GetTagList error: %s", err.Error())
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		Tags []Tag `json:"tags"`
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("GetTagList error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetTagList error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		Errorf("GetTagList error: %d %s", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("GetTagList error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return result.Tags, nil
}

func (m *MPProxyUsecase) CreateTag(ctx context.Context, token string, name string) (*Tag, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Create_Tag,
		token,
	)
	Debugf("url: %s", url)

	body := map[string]interface{}{
		"tag": map[string]interface{}{
			"name": name,
		},
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		Errorf("CreateTag error: %s", err.Error())
		return nil, err
	}
	bodyReader := bytes.NewReader(bodyJson)

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		Errorf("CreateTag error: %s", err.Error())
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		Tag Tag
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("CreateTag error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("CreateTag error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		Errorf("CreateTag error: %d %s", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("CreateTag error: %d %s", result.ErrCode, result.ErrMsg)
	}
	return &result.Tag, nil
}

func (m *MPProxyUsecase) UpdateTag(ctx context.Context, token string, id int64, name string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Update_Tag,
		token,
	)

	body := map[string]interface{}{
		"tag": map[string]interface{}{
			"id":   id,
			"name": name,
		},
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		Errorf("UpdateTag error: %s", err.Error())
		return err
	}
	bodyReader := bytes.NewReader(bodyJson)

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		Errorf("UpdateTag error: %s", err.Error())
		return err
	}

	type resultT struct {
		wxError.WXError
	}

	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("UpdateTag error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("UpdateTag error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		Errorf("UpdateTag error: %d %s", result.ErrCode, result.ErrMsg)
		return fmt.Errorf("UpdateTag error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

func (m *MPProxyUsecase) DeleteTag(ctx context.Context, token string, id int64) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Delete_Tag,
		token,
	)

	body := map[string]interface{}{
		"tag": map[string]interface{}{
			"id": id,
		},
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(bodyJson)
	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		Errorf("DeleteTag error: %s", err.Error())
		return err
	}

	type resultT struct {
		wxError.WXError
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("DeleteTag error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("DeleteTag error: %d %s", wxErr.ErrCode, wxErr.Error())
	}
	if result.ErrCode != 0 {
		Errorf("DeleteTag error: %d %s", result.ErrCode, result.ErrMsg)
		return fmt.Errorf("DeleteTag error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

func (m *MPProxyUsecase) GetTagMembers(ctx context.Context, token string, id int64, nextOpenid string) (*TagMembersRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_Tag_Members,
		token,
	)
	Debugf("url: %s", url)

	body := map[string]interface{}{
		"tagid":       id,
		"next_openid": nextOpenid,
	}
	Debugf("body: %+v", body)
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(bodyJson)
	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		Errorf("GetTagMembers error: %s", err.Error())
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		TagMembersRes
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("GetTagMembers error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetTagMembers error: %d %s", wxErr.ErrCode, wxErr.Error())
	}
	if result.ErrCode != 0 {
		Errorf("GetTagMembers error: %d %s", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("GetTagMembers error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return &result.TagMembersRes, nil
}

// BatchTaggingMembers
func (m *MPProxyUsecase) BatchTaggingMembers(ctx context.Context, token string, tagid int64, openids []string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Batch_Tagging,
		token,
	)
	Debugf("url: %s", url)

	body := map[string]interface{}{
		"tagid":       tagid,
		"openid_list": openids,
	}
	Debugf("body: %+v", body)
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(bodyJson)
	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		Errorf("BatchTaggingMembers error: %s", err.Error())
		return err
	}

	type resultT struct {
		wxError.WXError
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("BatchTaggingMembers error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("BatchTaggingMembers error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		Errorf("BatchTaggingMembers error: %d %s", result.ErrCode, result.ErrMsg)
		return fmt.Errorf("BatchTaggingMembers error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

// BatchUntaggingMembers
func (m *MPProxyUsecase) BatchUntaggingMembers(ctx context.Context, token string, tagid int64, openids []string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Batch_Untagging,
		token,
	)
	Debugf("url: %s", url)

	body := map[string]interface{}{
		"tagid":       tagid,
		"openid_list": openids,
	}
	Debugf("body: %+v", body)
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(bodyJson)
	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		Errorf("BatchUntaggingMembers error: %s", err.Error())
		return err
	}

	type resultT struct {
		wxError.WXError
	}

	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("BatchUntaggingMembers error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("BatchUntaggingMembers error: %d %s", wxErr.ErrCode, wxErr.Error())
	}
	if result.ErrCode != 0 {
		Errorf("BatchUntaggingMembers error: %d %s", result.ErrCode, result.ErrMsg)
		return fmt.Errorf("BatchUntaggingMembers error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

// CreateLimitQRCode
func (m *MPProxyUsecase) CreateLimitQRCode(ctx context.Context, token string,
	scene interface{}, expireSeconds int64,
) (*Ticket, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Create_QRCode,
		token,
	)
	Debugf("url: %s", url)

	tq := CreateQRCodeReq{
		ExpireSeconds: expireSeconds,
	}
	switch reflect.ValueOf(scene).Kind() {
	case reflect.String:
		Debugf("create qrcode scene: %s", scene)
		tq.ActionName = ActionStr
		tq.ActionInfo.Scene.SceneStr = scene.(string)
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		Debugf("create qrcode scene: %d", scene)
		tq.ActionName = ActionId
		tq.ActionInfo.Scene.SceneId = scene.(int64)
	default:
		Errorf("scene not supported: %v", reflect.ValueOf(scene).Kind())
		return nil, fmt.Errorf("scene not supported: %v", reflect.ValueOf(scene).Kind())
	}
	Debugf("body: %+v", tq)

	body, err := helpers.BuildRequestBody[CreateQRCodeReq](tq)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return nil, err
	}

	resp, err := m.hc.Post(url, "application/json", body)
	if err != nil {
		Errorf("create limit qrcode error: %s", err.Error())
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		Ticket
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("create limit qrcode error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("create limit qrcode error: %d %s", wxErr.ErrCode, wxErr.Error())
	}
	if result.ErrCode != 0 {
		Errorf("create limit qrcode error: %d %s", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("create limit qrcode error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return &result.Ticket, nil
}

// CreateTemporaryQRCode
func (m *MPProxyUsecase) CreateTemporaryQRCode(ctx context.Context, token string,
	scene interface{}, expireSeconds int64,
) (*Ticket, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Create_QRCode,
		token,
	)
	Debugf("url: %s", url)

	tq := CreateQRCodeReq{
		ExpireSeconds: expireSeconds,
	}
	switch reflect.ValueOf(scene).Kind() {
	case reflect.String:
		Debugf("create qrcode scene: %s", scene)
		tq.ActionName = ActionStr
		tq.ActionInfo.Scene.SceneStr = scene.(string)
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		Debugf("create qrcode scene: %d", scene)
		tq.ActionName = ActionId
		tq.ActionInfo.Scene.SceneId = scene.(int64)
	default:
		Errorf("scene not supported: %v", reflect.ValueOf(scene).Kind())
		return nil, fmt.Errorf("scene not supported: %v", reflect.ValueOf(scene).Kind())
	}
	Debugf("body: %+v", tq)

	body, err := helpers.BuildRequestBody[CreateQRCodeReq](tq)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return nil, err
	}

	resp, err := m.hc.Post(url, "application/json", body)
	if err != nil {
		Errorf("create temporary qrcode error: %s", err.Error())
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		Ticket
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("create temporary qrcode error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("create temporary qrcode error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		Errorf("create temporary qrcode error: %d %s", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("create temporary qrcode error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return &result.Ticket, nil
}

// GenShorten
func (m *MPProxyUsecase) GenShorten(ctx context.Context, token string, longData string,
	expireSeconds int64,
) (*string, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Gen_Shorten,
		token,
	)
	Debugf("url: %s", url)

	body := map[string]interface{}{
		"long_data":      longData,
		"expire_seconds": expireSeconds,
	}
	Debugf("body: %+v", body)

	bodyReader, err := helpers.BuildRequestBody[map[string]interface{}](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return nil, err
	}

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		Errorf("gen shorten error: %s", err.Error())
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		ShortKey string `json:"short_key"`
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("gen shorten error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("gen shorten error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		Errorf("gen shorten error: %d %s", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("gen shorten error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return &result.ShortKey, nil
}

// FetchShorten
func (m *MPProxyUsecase) FetchShorten(ctx context.Context, token string, shortKey string) (*FetchShortenRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Fetch_Shorten,
		token,
	)
	Debugf("url: %s", url)

	body := map[string]interface{}{
		"short_key": shortKey,
	}
	Debugf("body: %+v", body)
	bodyReader, err := helpers.BuildRequestBody[map[string]interface{}](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return nil, err
	}

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		Errorf("fetch shorten error: %s", err.Error())
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		FetchShortenRes
	}

	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("fetch shorten error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("fetch shorten error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		Errorf("fetch shorten error: %d %s", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("fetch shorten error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return &result.FetchShortenRes, nil
}

// GetMenuInfo
func (m *MPProxyUsecase) GetMenuInfo(ctx context.Context, token string) (*MenuRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_Menu,
		token,
	)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	if err != nil {
		Errorf("get menu info error: %s", err.Error())
		return nil, err
	}

	type resultT struct {
		wxError.WXError
		MenuRes
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("get menu info error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("get menu info error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		Errorf("get menu info error: %d %s", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("get menu info error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return &result.MenuRes, nil
}

// CreateMenu
func (m *MPProxyUsecase) CreateMenu(ctx context.Context, token string,
	buttons []*v1.MenuButton,
) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Create_Menu,
		token,
	)
	Debugf("url: %s", url)
	Debugf("buttons: %+v", buttons)

	type CreateMenuReq struct {
		Button []Button `json:"button"`
	}
	body := CreateMenuReq{
		Button: make([]Button, 0),
	}
	for _, btn := range buttons {
		var subBtns []*Button
		if len(btn.SubButton) > 0 {
			for _, subBtn := range btn.SubButton {
				subBtns = append(subBtns, &Button{
					Type:       subBtn.Type,
					Name:       subBtn.Name,
					Key:        subBtn.Key,
					URL:        subBtn.Url,
					MediaID:    subBtn.MediaId,
					AppID:      subBtn.AppId,
					PagePath:   subBtn.PagePath,
					SubButtons: []*Button{},
				})
			}
		}

		body.Button = append(body.Button, Button{
			Type:       btn.Type,
			Name:       btn.Name,
			Key:        btn.Key,
			URL:        btn.Url,
			MediaID:    btn.MediaId,
			AppID:      btn.AppId,
			PagePath:   btn.PagePath,
			SubButtons: subBtns,
		})
	}
	Debugf("body: %+v", body)

	bodyReader, err := helpers.BuildRequestBody[CreateMenuReq](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	if err != nil {
		Errorf("create menu error: %s", err.Error())
		return err
	}

	type resultT struct {
		wxError.WXError
	}
	result, wxErr := helpers.BuildHttpResponse[resultT](resp, err)
	if wxErr != nil {
		Errorf("create menu error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("create menu error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if result.ErrCode != 0 {
		Errorf("create menu error: %d %s", result.ErrCode, result.ErrMsg)
		return fmt.Errorf("create menu error: %d %s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

// DeleteMenu
func (m *MPProxyUsecase) DeleteMenu(ctx context.Context, token string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Del_Menu,
		token,
	)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	if err != nil {
		Errorf("delete menu error: %s", err.Error())
		return err
	}

	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("delete menu error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("delete menu error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("delete menu error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("delete menu error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}
