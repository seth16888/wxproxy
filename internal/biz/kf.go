package biz

import (
	"context"
	"fmt"

	"github.com/seth16888/wxcommon/domain"
	wxError "github.com/seth16888/wxcommon/error"
	"github.com/seth16888/wxcommon/helpers"
	. "github.com/seth16888/wxcommon/logger"
	"github.com/seth16888/wxcommon/paths"

	v1 "github.com/seth16888/wxproxy/api/v1"
)

// KfMessageComm 客服发送消息公共参数
type KfMessageComm struct {
	ToUser  string `json:"touser" binding:"required" msg:"touser required"`
	MsgType string `json:"msgtype" binding:"required" msg:"msgtype required"`

	// KF 客服帐号:如果要指定发送的客服账号，则必须填写此字段
	KF struct {
		KFAccount string `json:"kf_account"`
	} `json:"customservice,omitempty"`
}

// KfTextMessage 文本消息
type KfTextMessage struct {
	KfMessageComm
	Text struct {
		Content string `json:"content" binding:"required" msg:"content required"`
	} `json:"text" binding:"required" msg:"text required"` // 发送文本消息时，支持插入跳小程序的文字链
}

// KfImageMessage 图片消息
type KfImageMessage struct {
	KfMessageComm
	Image struct {
		MediaID string `json:"media_id" binding:"required" msg:"media_id required"`
	} `json:"image" binding:"required" msg:"image required"`
}

// KfVoiceMessage 语音消息
type KfVoiceMessage struct {
	KfMessageComm
	Voice struct {
		MediaID string `json:"media_id" binding:"required" msg:"media_id required"`
	} `json:"voice" binding:"required" msg:"voice required"`
}

// KfVideoMessage 视频消息
type KfVideoMessage struct {
	KfMessageComm
	Video struct {
		MediaID      string `json:"media_id" binding:"required" msg:"media_id required"`
		ThumbMediaID string `json:"thumb_media_id" binding:"required" msg:"thumb_media_id required"`
		Title        string `json:"title"`
		Description  string `json:"description"`
	} `json:"video" binding:"required" msg:"video required"`
}

// KfMusicMessage 音乐消息
type KfMusicMessage struct {
	KfMessageComm
	Music struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		MusicURL     string `json:"musicurl" binding:"required" msg:"musicurl required"`
		HQMusicURL   string `json:"hqmusicurl" binding:"required" msg:"hqmusicurl required"`
		ThumbMediaID string `json:"thumb_media_id" binding:"required" msg:"thumb_media_id required"`
	} `json:"music" binding:"required" msg:"music required"`
}

// KfNewsLinkToURLMessage 图文消息，点击跳转到外链
type KfNewsLinkToURLMessage struct {
	KfMessageComm
	News struct { // 图文消息，限制在1条以内，注意，如果图文数超过1，则将会返回错误码45008。
		Articles []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			PicURL      string `json:"picurl"`
		} `json:"articles"`
	} `json:"news" binding:"required" msg:"news required"`
}

// KfNewsLinkToPageMessage 图文消息，点击跳转到图文消息页面
//
//	草稿接口灰度完成后，将不再支持此前客服接口中带 media_id 的 mpnews 类型的图文消息
type KfNewsLinkToPageMessage struct {
	KfMessageComm
	MPNews struct {
		MediaID string `json:"media_id" binding:"required" msg:"media_id required"`
	} `json:"mpnews" binding:"required" msg:"mpnews required"`
}

// KfNewsLinkToArticleMessage 图文消息，点击跳转到图文消息页面
type KfNewsLinkToArticleMessage struct {
	KfMessageComm
	MPNewsArticle struct {
		ArticleId string `json:"article_id" binding:"required" msg:"article_id required"`
	} `json:"mpnewsarticle" binding:"required" msg:"mpnewsarticle required"`
}

// KfMenuMessage 菜单消息
//
// 当用户点击后，微信会发送一条XML消息到开发者服务器，格式如下：
//
//	<xml>
//	<ToUserName><![CDATA[ToUser]]></ToUserName>
//	<FromUserName><![CDATA[FromUser]]></FromUserName>
//	<CreateTime>1500000000</CreateTime>
//	<MsgType><![CDATA[text]]></MsgType>
//	<Content><![CDATA[满意]]></Content>
//	<MsgId>1234567890123456</MsgId>
//	<bizmsgmenuid>101</bizmsgmenuid>
//	</xml>
type KfMenuMessage struct {
	KfMessageComm
	MsgMenu struct {
		HeadContent string `json:"head_content" binding:"required" msg:"head_content required"`
		List        []struct {
			Id      string `json:"id" binding:"required" msg:"id required"`
			Content string `json:"content" binding:"required" msg:"content required"`
		} `json:"list" binding:"required" msg:"list required"`
		TailContent string `json:"tail_content"`
	} `json:"msgmenu" binding:"required" msg:"msgmenu required"`
}

// KfCardMessage 卡券消息
//
//	特别注意客服消息接口投放卡券仅支持非自定义Code码和导入code模式的卡券的卡券，详情请见：创建卡券。
type KfCardMessage struct {
	KfMessageComm
	WXCard struct {
		CardId string `json:"card_id" binding:"required" msg:"card_id required"`
	} `json:"wxcard" binding:"required" msg:"wxcard required"`
}

// KfMiniProgramMessage 小程序卡片消息
type KfMiniProgramMessage struct {
	KfMessageComm
	MiniProgramPage struct {
		AppID        string `json:"appid" binding:"required" msg:"appid required"`
		PagePath     string `json:"pagepath" binding:"required" msg:"pagepath required"`
		Title        string `json:"title"`
		ThumbMediaID string `json:"thumb_media_id" binding:"required" msg:"thumb_media_id required"`
	} `json:"miniprogrampage" binding:"required" msg:"miniprogrampage required"`
}

// KeFuInfo 客服基本信息
type KeFuInfo struct {
	KfAccount     string `json:"kf_account"`         // 完整客服帐号，格式为：帐号前缀@公众号微信号
	KfNick        string `json:"kf_nick"`            // 客服昵称
	KfID          int    `json:"kf_id"`              // 客服编号
	KfHeadImgURL  string `json:"kf_headimgurl"`      // 客服头像
	KfWX          string `json:"kf_wx"`              // 如果客服帐号已绑定了客服人员微信号， 则此处显示微信号
	InviteWX      string `json:"invite_wx"`          // 如果客服帐号尚未绑定微信号，但是已经发起了一个绑定邀请， 则此处显示绑定邀请的微信号
	InviteExpTime int    `json:"invite_expire_time"` // 如果客服帐号尚未绑定微信号，但是已经发起过一个绑定邀请， 邀请的过期时间，为unix 时间戳
	InviteStatus  string `json:"invite_status"`      // 邀请的状态，有等待确认"waiting,rejected,expired
}

// KeFuInfoListRes 客服信息列表
type KeFuInfoListRes struct {
	wxError.WXError
	KfList []*KeFuInfo `json:"kf_list"` // 完整客服信息
}

// KeFuOnlineInfo 客服在线信息
type KeFuOnlineInfo struct {
	KfAccount    string `json:"kf_account"`
	Status       int    `json:"status"`
	KfID         int    `json:"kf_id"`
	AcceptedCase int    `json:"accepted_case"`
}

type KeFuOnlineListRes struct {
	wxError.WXError
	KfOnlineList []*KeFuOnlineInfo `json:"kf_online_list"`
}

// KeFuMsgRecordRes 客服消息记录
type KeFuMsgRecordRes struct {
	wxError.WXError
	Number        int64            `json:"number"`
	MsgId         int64            `json:"msgid"`
	MsgRecordList []*KeFuMsgRecord `json:"recordlist"`
}

// KeFuMsgRecord 客服消息记录
type KeFuMsgRecord struct {
	OpenID string `json:"openid"`   // 用户openid
	Worker string `json:"worker"`   // 客服账号
	OpCode int64  `json:"opercode"` // 操作码，1001表示文本，1002表示图片，1003表示语音，1004表示视频，1006表示第三方小程序，1040表示链接
	Text   string `json:"text"`     // 操作内容
	Time   int64  `json:"time"`     // 操作时间戳
}

// GetKFSessionStatusRes 获取会话状态返回结果
type GetKFSessionStatusRes struct {
	wxError.WXError
	KfAccount  string `json:"kf_account"`
	CreateTime int64  `json:"createtime"`
}

// GetKFSessionListRes 获取客服会话列表返回结果
type GetKFSessionListRes struct {
	wxError.WXError
	SessionList []KFSession `json:"sessionlist"`
}

// KFSession 客服会话
type KFSession struct {
	OpenID     string `json:"openid"`
	CreateTime int64  `json:"createtime"`
}

// GetUnacceptedSessionListRes 获取未接入会话列表返回结果
type GetUnacceptedSessionListRes struct {
	wxError.WXError
	Count        int64 `json:"count"` // 未接入会话数量
	WaitCaseList []struct {
		OpenID     string `json:"openid"`      // 粉丝的openid
		LatestTime int64  `json:"latest_time"` // 粉丝的最后一条消息的时间，unix时间戳
	} `json:"waitcaselist"` // 未接入会话列表，最多返回100条数据，按照来访顺序
}

// AddKFAccount
func (m *MPProxyUsecase) AddKFAccount(ctx context.Context, token string, account string, nick string, pwd string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Add_KfAccount, token)
	Debugf("url: %s", url)

	type AddKFAccountReq struct {
		Account string `json:"kf_account"`
		Nick    string `json:"nickname"`
	}
	req := &AddKFAccountReq{
		Account: account,
		Nick:    nick,
	}
	reader, err := helpers.BuildRequestBody[*AddKFAccountReq](req)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("AddKFAccount error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("AddKFAccount error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("AddKFAccount error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("AddKFAccount error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// InviteKFWorker
func (m *MPProxyUsecase) InviteKFWorker(ctx context.Context, token string, account string, inviteWx string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Invite_Worker, token)
	Debugf("url: %s", url)

	type InviteKFWorkerReq struct {
		Account string `json:"kf_account"`
		Wx      string `json:"invite_wx"`
	}
	req := &InviteKFWorkerReq{
		Account: account,
		Wx:      inviteWx,
	}
	reader, err := helpers.BuildRequestBody[*InviteKFWorkerReq](req)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("InviteKFWorker error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("InviteKFWorker error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("InviteKFWorker error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("InviteKFWorker error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// SendKFTextMsg
func (m *MPProxyUsecase) SendKFTextMsg(ctx context.Context, req *v1.SendKFTextMsgRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_KF_Send_Message, req.AccessToken)
	Debugf("url: %s", url)

	body := &KfTextMessage{
		KfMessageComm: KfMessageComm{
			ToUser:  req.Common.ToUser,
			MsgType: "text",
			KF: struct {
				KFAccount string `json:"kf_account"`
			}{
				KFAccount: req.Common.CustomerService.KfAccount,
			},
		},
		Text: struct {
			Content string `json:"content" binding:"required" msg:"content required"`
		}{
			Content: req.Text.Content,
		},
	}

	reader, err := helpers.BuildRequestBody[*KfTextMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendKFTextMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendKFTextMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendKFTextMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendKFTextMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// GetKFSessionList
func (m *MPProxyUsecase) GetKFSessionList(ctx context.Context, token string, kf string) (*GetKFSessionListRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s&kf_account=%s", domain.GetWXAPIDomain(), paths.Path_Get_KFSessionList, token, kf)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[GetKFSessionListRes](resp, err)
	if wxErr != nil {
		Errorf("GetKFSessionList error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetKFSessionList error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("GetKFSessionList error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("GetKFSessionList error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}

// NewKFSession
func (m *MPProxyUsecase) NewKFSession(ctx context.Context, token string, kf string, openid string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Create_Session, token)
	Debugf("url: %s", url)

	type NewKFSessionReq struct {
		OpenID    string `json:"openid"`
		KfAccount string `json:"kf_account"`
	}
	req := &NewKFSessionReq{
		OpenID:    openid,
		KfAccount: kf,
	}
	reader, err := helpers.BuildRequestBody[*NewKFSessionReq](req)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("NewKFSession error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("NewKFSession error: %d %s", wxErr.ErrCode, wxErr.Error())
	}
	if rt.ErrCode != 0 {
		Errorf("NewKFSession error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("NewKFSession error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// UpdateKFAccount
func (m *MPProxyUsecase) UpdateKFAccount(ctx context.Context, token string, account string, nick string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Update_KfAccount, token)
	Debugf("url: %s", url)

	type UpdateKFAccountReq struct {
		Account string `json:"kf_account"`
		Nick    string `json:"nickname"`
	}
	req := &UpdateKFAccountReq{
		Account: account,
		Nick:    nick,
	}

	reader, err := helpers.BuildRequestBody[*UpdateKFAccountReq](req)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("UpdateKFAccount error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("UpdateKFAccount error: %d %s", wxErr.ErrCode, wxErr.Error())
	}
	if rt.ErrCode != 0 {
		Errorf("UpdateKFAccount error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("UpdateKFAccount error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// DelKFAccount
func (m *MPProxyUsecase) DelKFAccount(ctx context.Context, token string, account string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s&kf_account=%s", domain.GetWXAPIDomain(), paths.Path_Del_KfAccount, token, account)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("DelKFAccount error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("DelKFAccount error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("DelKFAccount error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("DelKFAccount error: %d %s", rt.ErrCode, rt.ErrMsg)
	}
	return nil
}

// UpdateKFTyping
func (m *MPProxyUsecase) UpdateKFTyping(ctx context.Context, token string, touser string, command string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Update_KfStatus, token)
	Debugf("url: %s", url)

	type UpdateKFTypingReq struct {
		ToUser  string `json:"touser"`
		Command string `json:"command"`
	}
	req := &UpdateKFTypingReq{
		ToUser:  touser,
		Command: command,
	}

	reader, err := helpers.BuildRequestBody[*UpdateKFTypingReq](req)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("UpdateKFTyping error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("UpdateKFTyping error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("UpdateKFTyping error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("UpdateKFTyping error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// GetKFSessionStatus
func (m *MPProxyUsecase) GetKFSessionStatus(ctx context.Context, token string, openid string) (*GetKFSessionStatusRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s&openid=%s", domain.GetWXAPIDomain(), paths.Path_Get_SessionStatus, token, openid)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[GetKFSessionStatusRes](resp, err)
	if wxErr != nil {
		Errorf("GetKFSessionStatus error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetKFSessionStatus error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("GetKFSessionStatus error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("GetKFSessionStatus error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}

// GetKFSessionUnaccepted
func (m *MPProxyUsecase) GetKFSessionUnaccepted(ctx context.Context, token string) (*GetUnacceptedSessionListRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Get_UnacceptedSessionList, token)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[GetUnacceptedSessionListRes](resp, err)
	if wxErr != nil {
		Errorf("GetKFSessionUnaccepted error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("GetKFSessionUnaccepted error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("GetKFSessionUnaccepted error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("GetKFSessionUnaccepted error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}

// CloseKFSession
func (m *MPProxyUsecase) CloseKFSession(ctx context.Context, token string, kf string, openid string) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Close_Session, token)
	Debugf("url: %s", url)

	type CloseKFSessionReq struct {
		KfAccount string `json:"kf_account"`
		OpenID    string `json:"openid"`
	}
	req := &CloseKFSessionReq{
		KfAccount: kf,
		OpenID:    openid,
	}

	reader, err := helpers.BuildRequestBody[*CloseKFSessionReq](req)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("CloseKFSession error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("CloseKFSession error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("CloseKFSession error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("CloseKFSession error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// SendKFImageMsg
func (m *MPProxyUsecase) SendKFImageMsg(ctx context.Context, req *v1.SendKFImageMsgRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_KF_Send_Message, req.AccessToken)
	Debugf("url: %s", url)

	body := &KfImageMessage{
		KfMessageComm: KfMessageComm{
			ToUser:  req.Common.ToUser,
			MsgType: "image",
			KF: struct {
				KFAccount string `json:"kf_account"`
			}{
				KFAccount: req.Common.CustomerService.KfAccount,
			},
		},
		Image: struct {
			MediaID string `json:"media_id" binding:"required" msg:"media_id required"`
		}{
			MediaID: req.Image.MediaId,
		},
	}

	reader, err := helpers.BuildRequestBody[*KfImageMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendKFImageMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendKFImageMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}
	if rt.ErrCode != 0 {
		Errorf("SendKFImageMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendKFImageMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// SendKFVoiceMsg
func (m *MPProxyUsecase) SendKFVoiceMsg(ctx context.Context, req *v1.SendKFVoiceMsgRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_KF_Send_Message, req.AccessToken)
	Debugf("url: %s", url)

	body := &KfVoiceMessage{
		KfMessageComm: KfMessageComm{
			ToUser:  req.Common.ToUser,
			MsgType: "voice",
			KF: struct {
				KFAccount string `json:"kf_account"`
			}{
				KFAccount: req.Common.CustomerService.KfAccount,
			},
		},
		Voice: struct {
			MediaID string `json:"media_id" binding:"required" msg:"media_id required"`
		}{
			MediaID: req.Voice.MediaId,
		},
	}

	reader, err := helpers.BuildRequestBody[*KfVoiceMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendKFVoiceMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendKFVoiceMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendKFVoiceMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendKFVoiceMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// SendKFVideoMsg
func (m *MPProxyUsecase) SendKFVideoMsg(ctx context.Context, req *v1.SendKFVideoMsgRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_KF_Send_Message, req.AccessToken)
	Debugf("url: %s", url)

	body := &KfVideoMessage{
		KfMessageComm: KfMessageComm{
			ToUser:  req.Common.ToUser,
			MsgType: "video",
			KF: struct {
				KFAccount string `json:"kf_account"`
			}{
				KFAccount: req.Common.CustomerService.KfAccount,
			},
		},
		Video: struct {
			MediaID      string `json:"media_id" binding:"required" msg:"media_id required"`
			ThumbMediaID string `json:"thumb_media_id" binding:"required" msg:"thumb_media_id required"`
			Title        string `json:"title"`
			Description  string `json:"description"`
		}{
			MediaID:      req.Video.MediaId,
			ThumbMediaID: req.Video.ThumbMediaId,
			Title:        req.Video.Title,
			Description:  req.Video.Description,
		},
	}

	reader, err := helpers.BuildRequestBody[*KfVideoMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendKFVideoMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendKFVideoMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendKFVideoMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendKFVideoMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// SendKFMusicMsg
func (m *MPProxyUsecase) SendKFMusicMsg(ctx context.Context, req *v1.SendKFMusicMsgRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_KF_Send_Message, req.AccessToken)
	Debugf("url: %s", url)

	body := &KfMusicMessage{
		KfMessageComm: KfMessageComm{
			ToUser:  req.Common.ToUser,
			MsgType: "music",
			KF: struct {
				KFAccount string `json:"kf_account"`
			}{
				KFAccount: req.Common.CustomerService.KfAccount,
			},
		},
		Music: struct {
			Title        string `json:"title"`
			Description  string `json:"description"`
			MusicURL     string `json:"musicurl" binding:"required" msg:"musicurl required"`
			HQMusicURL   string `json:"hqmusicurl" binding:"required" msg:"hqmusicurl required"`
			ThumbMediaID string `json:"thumb_media_id" binding:"required" msg:"thumb_media_id required"`
		}{
			Title:        req.Music.Title,
			Description:  req.Music.Description,
			MusicURL:     req.Music.MusicUrl,
			HQMusicURL:   req.Music.MusicUrl,
			ThumbMediaID: req.Music.ThumbMediaId,
		},
	}

	reader, err := helpers.BuildRequestBody[*KfMusicMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendKFMusicMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendKFMusicMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendKFMusicMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendKFMusicMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// SendKFNewsCardMsg
func (m *MPProxyUsecase) SendKFNewsCardMsg(ctx context.Context, req *v1.SendKFNewsCardMsgRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_KF_Send_Message, req.AccessToken)
	Debugf("url: %s", url)

	body := &KfNewsLinkToURLMessage{
		KfMessageComm: KfMessageComm{
			ToUser:  req.Common.ToUser,
			MsgType: "news",
			KF: struct {
				KFAccount string `json:"kf_account"`
			}{
				KFAccount: req.Common.CustomerService.KfAccount,
			},
		},
		News: struct {
			Articles []struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				URL         string `json:"url"`
				PicURL      string `json:"picurl"`
			} `json:"articles"`
		}{},
	}
	type Article struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		PicURL      string `json:"picurl"`
	}
	body.News.Articles = append(body.News.Articles, Article{
		Title:       req.News.Title,
		Description: req.News.Description,
		URL:         req.News.Url,
		PicURL:      req.News.PicUrl,
	})

	reader, err := helpers.BuildRequestBody[*KfNewsLinkToURLMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendKFNewsCardMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendKFNewsCardMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendKFNewsCardMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendKFNewsCardMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// SendKFNewsPageMsg
func (m *MPProxyUsecase) SendKFNewsPageMsg(ctx context.Context, req *v1.SendKFNewsPageMsgRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_KF_Send_Message, req.AccessToken)
	Debugf("url: %s", url)

	body := &KfNewsLinkToPageMessage{}
	body.KfMessageComm.KF.KFAccount = req.Common.CustomerService.KfAccount
	body.KfMessageComm.MsgType = "mpnews"
	body.KfMessageComm.ToUser = req.Common.ToUser
	body.MPNews.MediaID = req.MpNews.MediaId

	reader, err := helpers.BuildRequestBody[*KfNewsLinkToPageMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendKFNewsPageMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendKFNewsPageMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendKFNewsPageMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendKFNewsPageMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// SendKFToArticleMsg
func (m *MPProxyUsecase) SendKFToArticleMsg(ctx context.Context, req *v1.SendKFToArticleMsgRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_KF_Send_Message, req.AccessToken)
	Debugf("url: %s", url)

	body := &KfNewsLinkToArticleMessage{}
	body.KfMessageComm.KF.KFAccount = req.Common.CustomerService.KfAccount
	body.KfMessageComm.MsgType = "mpnewsarticle"
	body.KfMessageComm.ToUser = req.Common.ToUser
	body.MPNewsArticle.ArticleId = req.MpNewsArticle.ArticleId

	reader, err := helpers.BuildRequestBody[*KfNewsLinkToArticleMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendKFToArticleMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendKFToArticleMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendKFToArticleMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendKFToArticleMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// SendKFMenuMsg
func (m *MPProxyUsecase) SendKFMenuMsg(ctx context.Context, req *v1.SendKFMenuMsgRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_KF_Send_Message, req.AccessToken)
	Debugf("url: %s", url)

	body := &KfMenuMessage{}
	body.KfMessageComm.KF.KFAccount = req.Common.CustomerService.KfAccount
	body.KfMessageComm.MsgType = "msgmenu"
	body.KfMessageComm.ToUser = req.Common.ToUser
	body.MsgMenu.HeadContent = req.MsgMenu.HeadContent
	body.MsgMenu.TailContent = req.MsgMenu.TailContent
	for _, item := range req.MsgMenu.List {
		body.MsgMenu.List = append(body.MsgMenu.List, struct {
			Id      string `json:"id" binding:"required" msg:"id required"`
			Content string `json:"content" binding:"required" msg:"content required"`
		}{
			Id:      item.Id,
			Content: item.Content,
		})
	}

	reader, err := helpers.BuildRequestBody[*KfMenuMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendKFMenuMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendKFMenuMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendKFMenuMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendKFMenuMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// SendKFCardMsg
func (m *MPProxyUsecase) SendKFCardMsg(ctx context.Context, req *v1.SendKFCardMsgRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_KF_Send_Message, req.AccessToken)
	Debugf("url: %s", url)

	body := &KfCardMessage{}
	body.KfMessageComm.KF.KFAccount = req.Common.CustomerService.KfAccount
	body.KfMessageComm.MsgType = "wxcard"
	body.KfMessageComm.ToUser = req.Common.ToUser
	body.WXCard.CardId = req.WxCard.CardId

	reader, err := helpers.BuildRequestBody[*KfCardMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendKFCardMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendKFCardMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendKFCardMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendKFCardMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// SendKFMiniProgramMsg
func (m *MPProxyUsecase) SendKFMiniProgramMsg(ctx context.Context, req *v1.SendKFMiniProgramMsgRequest) error {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_KF_Send_Message, req.AccessToken)
	Debugf("url: %s", url)

	body := &KfMiniProgramMessage{}
	body.KfMessageComm.KF.KFAccount = req.Common.CustomerService.KfAccount
	body.KfMessageComm.MsgType = "miniprogrampage"
	body.KfMessageComm.ToUser = req.Common.ToUser
	body.MiniProgramPage.AppID = req.MiniProgramPage.AppId
	body.MiniProgramPage.Title = req.MiniProgramPage.Title
	body.MiniProgramPage.PagePath = req.MiniProgramPage.PagePath
	body.MiniProgramPage.ThumbMediaID = req.MiniProgramPage.ThumbMediaId

	reader, err := helpers.BuildRequestBody[*KfMiniProgramMessage](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return err
	}

	resp, err := m.hc.Post(url, "application/json", reader)
	rt, wxErr := helpers.BuildHttpResponse[wxError.WXError](resp, err)
	if wxErr != nil {
		Errorf("SendKFMiniProgramMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
		return fmt.Errorf("SendKFMiniProgramMsg error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("SendKFMiniProgramMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
		return fmt.Errorf("SendKFMiniProgramMsg error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return nil
}

// GetKFList
func (m *MPProxyUsecase) GetKFList(ctx context.Context, token string) (*KeFuInfoListRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_KfList,
		token,
	)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[KeFuInfoListRes](resp, err)
	if wxErr != nil {
		Errorf("get kf list error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("get kf list error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("get kf list error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("get kf list error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}

// GetKFOnlineList
func (m *MPProxyUsecase) GetKFOnlineList(ctx context.Context, token string) (*KeFuOnlineListRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Get_OnlineKfList,
		token,
	)
	Debugf("url: %s", url)

	resp, err := m.hc.Get(url)
	rt, wxErr := helpers.BuildHttpResponse[KeFuOnlineListRes](resp, err)
	if wxErr != nil {
		Errorf("get kf online list error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("get kf online list error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("get kf online list error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("get kf online list error: %d %s", rt.ErrCode, rt.ErrMsg)
	}
	return rt, nil
}

// GetKFMsgHistory
func (m *MPProxyUsecase) GetKFMsgHistory(ctx context.Context, token string, start int64,
	end int64, msgId int64, number int64) (*KeFuMsgRecordRes, error) {
	url := fmt.Sprintf("https://%s%s?access_token=%s", domain.GetWXAPIDomain(), paths.Path_Get_MsgRecord, token)
	Debugf("url: %s", url)

	type GetKFMsgHistoryReq struct {
		StartTime int64 `json:"starttime"`
		EndTime   int64 `json:"endtime"`
		MsgId     int64 `json:"msgid"`
		Number    int64 `json:"number"`
	}
	body := &GetKFMsgHistoryReq{
		StartTime: start,
		EndTime:   end,
		MsgId:     msgId,
		Number:    number,
	}
	bodyReader, err := helpers.BuildRequestBody[*GetKFMsgHistoryReq](body)
	if err != nil {
		Errorf("build request body error: %s", err.Error())
		return nil, err
	}

	resp, err := m.hc.Post(url, "application/json", bodyReader)
	rt, wxErr := helpers.BuildHttpResponse[KeFuMsgRecordRes](resp, err)
	if wxErr != nil {
		Errorf("get kf msg history error: %d %s", wxErr.ErrCode, wxErr.Error())
		return nil, fmt.Errorf("get kf msg history error: %d %s", wxErr.ErrCode, wxErr.Error())
	}

	if rt.ErrCode != 0 {
		Errorf("get kf msg history error: %d %s", rt.ErrCode, rt.ErrMsg)
		return nil, fmt.Errorf("get kf msg history error: %d %s", rt.ErrCode, rt.ErrMsg)
	}

	return rt, nil
}
