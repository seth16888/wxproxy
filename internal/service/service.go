package service

import (
	"context"
	"fmt"

	v1 "github.com/seth16888/wxproxy/api/v1"
	"github.com/seth16888/wxproxy/internal/biz"
	"go.uber.org/zap"
)

type MPProxyService struct {
	v1.UnimplementedMpproxyServer
	log *zap.Logger
	uc  *biz.MPProxyUsecase
}

func NewMPProxyService(uc *biz.MPProxyUsecase, logger *zap.Logger) *MPProxyService {
	return &MPProxyService{uc: uc, log: logger}
}

func (m *MPProxyService) BlockMember(ctx context.Context, req *v1.BlockMemberReq) (*v1.WXErrorReply, error) {
  wxErr := m.uc.BlockMember(ctx, req.AccessToken, req.OpenIds)

  return wxErr, nil
}

func (m *MPProxyService) UnBlockMember(ctx context.Context, req *v1.BlockMemberReq) (*v1.WXErrorReply, error) {
  wxErr := m.uc.UnBlockMember(ctx, req.AccessToken, req.OpenIds)

  return wxErr, nil
}

func (m *MPProxyService) GetMaterialCount(ctx context.Context, req *v1.AccessTokenParam) (*v1.GetMaterialCountReply, error) {
	res, err := m.uc.GetMaterialCoount(ctx, req.GetAccessToken())
	if err != nil {
		return nil, err
	}
	return &v1.GetMaterialCountReply{
		VoiceCount: res.VoiceCount,
		VideoCount: res.VideoCount,
		ImageCount: res.ImageCount,
		NewsCount:  res.NewsCount,
	}, nil
}

func (m *MPProxyService) GetMaterialList(ctx context.Context, req *v1.GetMaterialListRequest) (*v1.GetMaterialListReply, error) {
	res, err := m.uc.GetMaterialList(ctx, req.GetAccessToken(), req.GetType(), req.GetOffset(), req.GetCount())
	if err != nil {
		return nil, err
	}

	var items []*v1.MaterialItem
	for _, item := range res.Item {
		items = append(items, &v1.MaterialItem{
			MediaId:    item.MediaId,
			Name:       item.Name,
			UpdateTime: item.UpdateTime,
			Url:        item.Url,
		})
	}

	return &v1.GetMaterialListReply{
		TotalCount: res.TotalCount,
		ItemCount:  res.ItemCount,
		Item:       items,
	}, nil
}

func (m *MPProxyService) GetMemberList(ctx context.Context, req *v1.GetMemberListRequest) (*v1.GetMemberListReply, error) {
	res, err := m.uc.GetMemberList(ctx, req.AccessToken, req.NextOpenid)
	if err != nil {
		return nil, err
	}

	var items []*v1.OpenIdList
	for _, item := range res.Data.OpenidList {
		items = append(items, &v1.OpenIdList{
			Openid: item,
		})
	}

	return &v1.GetMemberListReply{
		Total:      res.Total,
		Count:      res.Count,
		NextOpenid: res.NextOpenid,
		Data: &v1.GetMemberListReply_IdList{
			Openid: items,
		},
	}, nil
}

func (m *MPProxyService) GetMemberInfo(ctx context.Context, req *v1.GetMemberInfoRequest) (*v1.GetMemberInfoReply, error) {
	res, err := m.uc.GetMemberInfo(ctx, req.AccessToken, req.Openid, req.Lang)
	if err != nil {
		return nil, err
	}

	return &v1.GetMemberInfoReply{
		Openid:         res.Openid,
		Subscribe:      res.Subscribe,
		SubscribeTime:  res.SubscribeTime,
		Unionid:        res.Unionid,
		Remark:         res.Remark,
		Groupid:        res.Groupid,
		TagidList:      res.TagidList,
		SubscribeScene: res.SubscribeScene,
		QrScene:        res.QrScene,
		QrSceneStr:     res.QrSceneStr,
		Language:       res.Language,
	}, nil
}

func (m *MPProxyService) BatchGetMemberInfo(ctx context.Context, req *v1.BatchGetMemberInfoRequest) (*v1.BatchGetMemberInfoReply, error) {
	var openidList []string
	for _, item := range req.UserList {
		openidList = append(openidList, item.Openid)
	}
	res, err := m.uc.BatchGetMemberInfo(ctx, req.AccessToken, openidList)
	if err != nil {
		return nil, err
	}

	var rt v1.BatchGetMemberInfoReply
	for _, item := range *res {
		rt.UserListInfo = append(rt.UserListInfo, &v1.GetMemberInfoReply{
			Openid:         item.Openid,
			Subscribe:      item.Subscribe,
			SubscribeTime:  item.SubscribeTime,
			Unionid:        item.Unionid,
			Remark:         item.Remark,
			Groupid:        item.Groupid,
			TagidList:      item.TagidList,
			SubscribeScene: item.SubscribeScene,
			QrScene:        item.QrScene,
			QrSceneStr:     item.QrSceneStr,
			Language:       item.Language,
		})
	}

	return &rt, nil
}

// GetMemberTags 获取用户身上的标签列表
func (m *MPProxyService) GetMemberTags(ctx context.Context, req *v1.GetMemberTagsRequest) (*v1.GetMemberTagsReply, error) {
	res, err := m.uc.GetMemberTags(ctx, req.AccessToken, req.Openid)
	if err != nil {
		return nil, err
	}

	return &v1.GetMemberTagsReply{
		TagidList: res,
	}, nil
}
// BatchTaggingMembers 批量为用户打标签
// TODO: biz层应该返回v1.WXErrorReply, 这里判断code是否为0=微信API错误
func (m *MPProxyService) UpdateMemberRemark(ctx context.Context, req *v1.UpdateMemberRemarkRequest) (*v1.WXErrorReply, error) {
	err := m.uc.UpdateMemberRemark(ctx, req.AccessToken, req.Openid, req.Remark)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{
		Errcode: 0,
		Errmsg:  "ok",
	}, nil
}

func (m *MPProxyService) GetTagList(ctx context.Context, req *v1.AccessTokenParam) (*v1.GetTagListReply, error) {
	res, err := m.uc.GetTagList(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	var tags []*v1.Tag
	for _, item := range res {
		tags = append(tags, &v1.Tag{
			Id:    item.Id,
			Name:  item.Name,
			Count: item.Count,
		})
	}

	return &v1.GetTagListReply{
		Tags: tags,
	}, nil
}

func (m *MPProxyService) CreateTag(ctx context.Context, req *v1.CreateTagRequest) (*v1.CreateTagReply, error) {
	res, err := m.uc.CreateTag(ctx, req.AccessToken, req.Name)
	if err != nil {
		return nil, err
	}

	return &v1.CreateTagReply{
		Tag: &v1.Tag{
			Id:    res.Id,
			Name:  res.Name,
			Count: res.Count,
		},
	}, nil
}

func (m *MPProxyService) UpdateTag(ctx context.Context, req *v1.UpdateTagRequest) (*v1.WXErrorReply, error) {
	err := m.uc.UpdateTag(ctx, req.AccessToken, req.Id, req.Name)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{
		Errcode: 0,
		Errmsg:  "ok",
	}, nil
}

func (m *MPProxyService) DeleteTag(ctx context.Context, req *v1.DeleteTagRequest) (*v1.WXErrorReply, error) {
	err := m.uc.DeleteTag(ctx, req.AccessToken, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.WXErrorReply{
		Errcode: 0,
		Errmsg:  "ok",
	}, nil
}

func (m *MPProxyService) GetTagMembers(ctx context.Context, req *v1.GetTagMembersRequest) (*v1.GetTagMembersReply, error) {
	res, err := m.uc.GetTagMembers(ctx, req.AccessToken, req.Id, req.NextOpenid)
	if err != nil {
		return nil, err
	}

	return &v1.GetTagMembersReply{
		Count:      res.Count,
		NextOpenid: res.NextOpenId,
		Data: &v1.GetTagMembersReply_DataT{
			Openid: res.Data.OpenIdList,
		},
	}, nil
}

func (m *MPProxyService) BatchTaggingMembers(ctx context.Context, req *v1.BatchTaggingMembersRequest) (*v1.WXErrorReply, error) {
	err := m.uc.BatchTaggingMembers(ctx, req.AccessToken, req.Id, req.OpenidList)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{
		Errcode: 0,
		Errmsg:  "ok",
	}, nil
}

func (m *MPProxyService) BatchUnTaggingMembers(ctx context.Context, req *v1.BatchUnTaggingMembersRequest) (*v1.WXErrorReply, error) {
	err := m.uc.BatchUntaggingMembers(ctx, req.AccessToken, req.Id, req.OpenidList)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{
		Errcode: 0,
		Errmsg:  "ok",
	}, nil
}

func (m *MPProxyService) CreateTemporaryQRCode(ctx context.Context,
	req *v1.CreateQRCodeRequest) (*v1.CreateQRCodeReply, error) {
	// TODO: scene根据数据类型同时支持数字、字符串两种类型
	res, err := m.uc.CreateTemporaryQRCode(ctx, req.AccessToken, req.Scene, req.ExpireSeconds)
	if err != nil {
		return nil, err
	}

	return &v1.CreateQRCodeReply{
		Ticket:        res.Ticket,
		ExpireSeconds: res.ExpireSeconds,
		URL:           res.Url,
	}, nil
}
func (m *MPProxyService) CreateLimitQRCode(ctx context.Context,
	req *v1.CreateQRCodeRequest) (*v1.CreateQRCodeReply, error) {
	// TODO: scene根据数据类型同时支持数字、字符串两种类型
	res, err := m.uc.CreateLimitQRCode(ctx, req.AccessToken, req.Scene, req.ExpireSeconds)
	if err != nil {
		return nil, err
	}

	return &v1.CreateQRCodeReply{
		Ticket:        res.Ticket,
		ExpireSeconds: res.ExpireSeconds,
		URL:           res.Url,
	}, nil
}
func (m *MPProxyService) GenShorten(ctx context.Context,
	req *v1.GenShortenRequest) (*v1.GenShortenReply, error) {
	res, err := m.uc.GenShorten(ctx, req.AccessToken, req.LongData, req.ExpireSeconds)
	if err != nil {
		return nil, err
	}

	return &v1.GenShortenReply{
		ShortKey: *res,
	}, nil
}
func (m *MPProxyService) FetchShorten(ctx context.Context,
	req *v1.FetchShortenRequest) (*v1.FetchShortenReply, error) {
	res, err := m.uc.FetchShorten(ctx, req.AccessToken, req.ShortKey)
	if err != nil {
		return nil, err
	}

	return &v1.FetchShortenReply{
		LongData:      res.LongData,
		CreateTime:    res.CreateTime,
		ExpireSeconds: res.ExpireSeconds,
	}, nil
}
func (m *MPProxyService) GetMenuInfo(ctx context.Context,
	req *v1.AccessTokenParam) (*v1.MenuInfoReply, error) {
	res, err := m.uc.GetMenuInfo(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	var rt = &v1.MenuInfoReply{
		Menu: &v1.MenuInfoReply_MenuType{
			Menuid: res.Menu.MenuID,
			Button: []*v1.MenuButton{},
		},
		Conditionalmenu: []*v1.ConditionalMenu{},
	}
	for _, item := range res.Menu.Button {
		var subBtns []*v1.MenuButton
		for _, subItem := range item.SubButtons {
			subBtns = append(subBtns, &v1.MenuButton{
				Type:      subItem.Type,
				Name:      subItem.Name,
				Key:       subItem.Key,
				Url:       subItem.URL,
				MediaId:   subItem.MediaID,
				AppId:     subItem.AppID,
				PagePath:  subItem.PagePath,
				SubButton: []*v1.MenuButton{},
			})
		}
		rt.Menu.Button = append(rt.Menu.Button, &v1.MenuButton{
			Type:      item.Type,
			Name:      item.Name,
			Key:       item.Key,
			Url:       item.URL,
			MediaId:   item.MediaID,
			AppId:     item.AppID,
			PagePath:  item.PagePath,
			SubButton: subBtns,
		})
	}
	var condition []*v1.ConditionalMenu
	for _, item := range res.Conditionalmenu {
		var cSubBtns []*v1.MenuButton
		for _, subItem := range item.Button {
			cSubBtns = append(cSubBtns, &v1.MenuButton{
				Type:      subItem.Type,
				Name:      subItem.Name,
				Key:       subItem.Key,
				Url:       subItem.URL,
				MediaId:   subItem.MediaID,
				AppId:     subItem.AppID,
				PagePath:  subItem.PagePath,
				SubButton: []*v1.MenuButton{},
			})
		}
		condition = append(condition, &v1.ConditionalMenu{
			Menuid: item.MenuID,
			Button: cSubBtns,
			Matchrule: &v1.ConditionalMatchRule{
				TagId:              item.MatchRule.TagId,
				ClientPlatformType: item.MatchRule.ClientPlatformType,
			},
		})
	}
	rt.Conditionalmenu = condition

	return rt, nil
}
func (m *MPProxyService) TryMatchMenu(ctx context.Context, req *v1.TryMatchMenuRequest) (*v1.TryMatchMenuReply, error) {
	return nil, nil
}
func (m *MPProxyService) PullMenu(ctx context.Context, req *v1.AccessTokenParam) (*v1.SelfMenuReply, error) {
	return nil, nil
}
func (m *MPProxyService) CreateMenu(ctx context.Context,
	req *v1.CreateMenuRequest) (*v1.WXErrorReply, error) {
	err := m.uc.CreateMenu(ctx, req.AccessToken, req.Button)
	if err != nil {
		m.log.Error("CreateMenu", zap.Error(err))
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}
func (m *MPProxyService) CreateConditionalMenu(ctx context.Context, req *v1.CreateMenuRequest) (*v1.WXErrorReply, error) {
	return nil, nil
}
func (m *MPProxyService) DeleteConditionalMenu(ctx context.Context, req *v1.DeleteConditionalMenuRequest) (*v1.WXErrorReply, error) {
	return nil, nil
}
func (m *MPProxyService) DeleteMenu(ctx context.Context, req *v1.AccessTokenParam) (*v1.WXErrorReply, error) {
	err := m.uc.DeleteMenu(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) GetIndustry(ctx context.Context, req *v1.AccessTokenParam) (*v1.GetIndustryReply, error) {
	res, err := m.uc.GetIndustry(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	return &v1.GetIndustryReply{
		PrimaryIndustry: &v1.GetIndustryReply_Industry{
			FirstClass:  res.PrimaryIndustry.FirstClass,
			SecondClass: res.PrimaryIndustry.SecondClass,
		},
		SecondaryIndustry: &v1.GetIndustryReply_Industry{
			FirstClass:  res.SecondaryIndustry.FirstClass,
			SecondClass: res.SecondaryIndustry.SecondClass,
		},
	}, nil
}
func (m *MPProxyService) GetAllPrivateTpl(ctx context.Context, req *v1.AccessTokenParam) (*v1.GetAllPrivateTplReply, error) {
	res, err := m.uc.GetAllPrivateTpl(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	var rt = &v1.GetAllPrivateTplReply{
		TemplateList: []*v1.GetAllPrivateTplReply_TplInfo{},
	}
	for _, item := range res.TemplateList {
		rt.TemplateList = append(rt.TemplateList, &v1.GetAllPrivateTplReply_TplInfo{
			TemplateId:        item.TemplateID,
			Title:             item.Title,
			Content:           item.Content,
			Example:           item.Example,
			PrimaryIndustry:   item.PrimaryIndustry,
			SecondaryIndustry: item.DeputyIndustry,
		})
	}

	return rt, nil
}
func (m *MPProxyService) SetIndustry(ctx context.Context, req *v1.SetIndustryRequest) (*v1.WXErrorReply, error) {
	res, err := m.uc.SetIndustry(ctx, req.AccessToken, req.IndustryId1, req.IndustryId2)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: int64(res.ErrCode), Errmsg: res.ErrMsg}, nil
}
func (m *MPProxyService) GetMessageTplId(ctx context.Context, req *v1.AddTemplateRequest) (*v1.AddMessageTplReply, error) {
	res, err := m.uc.GetMessageTplId(ctx, req.AccessToken, req.TemplateIdShort, req.KeywordNameList)
	if err != nil {
		return nil, err
	}

	return &v1.AddMessageTplReply{TemplateId: res.TemplateId}, nil
}
func (m *MPProxyService) DeleteMessageTpl(ctx context.Context, req *v1.DeleteMessageTplRequest) (*v1.WXErrorReply, error) {
	res, err := m.uc.DeleteMessageTpl(ctx, req.AccessToken, req.TemplateId)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: int64(res.ErrCode), Errmsg: res.ErrMsg}, nil
}
func (m *MPProxyService) SendTplMsg(ctx context.Context, req *v1.SendTplMsgRequest) (*v1.SendTplMsgReply, error) {
	res, err := m.uc.SendTplMsg(ctx, req.AccessToken, req)
	if err != nil {
		return nil, err
	}

	return &v1.SendTplMsgReply{Msgid: res.MsgID}, nil
}
func (m *MPProxyService) SendSubscribeMsg(ctx context.Context, req *v1.SendSubscribeMsgRequest) (*v1.WXErrorReply, error) {
	res, err := m.uc.SendSubscribeMsg(ctx, req.AccessToken, req)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: int64(res.ErrCode)}, nil
}
func (m *MPProxyService) GetBlockedTplMsg(ctx context.Context, req *v1.GetBlockedTplRequest) (*v1.GetBlockedTplMsgReply, error) {
	res, err := m.uc.GetBlockedTplMsg(ctx, req.AccessToken, req.TmplMsgId, req.LargestId, req.Limit)
	if err != nil {
		return nil, err
	}

	ret := &v1.GetBlockedTplMsgReply{Msginfo: []*v1.GetBlockedTplMsgReply_BlockedMsgInfo{}}
	for _, item := range res.MsgInfo {
		ret.Msginfo = append(ret.Msginfo, &v1.GetBlockedTplMsgReply_BlockedMsgInfo{
			Id:            item.Id,
			Openid:        item.OpenId,
			TmplMsgId:     item.TmplMsgId,
			Title:         item.Title,
			Content:       item.Content,
			SendTimestamp: item.SendTimestamp,
		})
	}

	return ret, nil
}

func (m *MPProxyService) AddSubscribeTpl(ctx context.Context, req *v1.AddSubscribeTplRequest) (*v1.AddSubscribeTplReply, error) {
	res, err := m.uc.AddSubscribeTpl(ctx, req.AccessToken, req.Tid, req.SceneDesc, req.KidList)
	if err != nil {
		return nil, err
	}

	return &v1.AddSubscribeTplReply{
		TemplateId: res,
	}, nil
}
func (m *MPProxyService) DelSubscribeTpl(ctx context.Context, req *v1.DelSubscribeTplRequest) (*v1.WXErrorReply, error) {
	err := m.uc.DelSubscribeTpl(ctx, req.AccessToken, req.TemplateId)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}
func (m *MPProxyService) GetSubscribeCategory(ctx context.Context, req *v1.AccessTokenParam) (*v1.GetSubscribeCategoryReply, error) {
	res, err := m.uc.GetSubscribeCategory(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	var rt = &v1.GetSubscribeCategoryReply{}
	for _, item := range res.Data {
		rt.Data = append(rt.Data, &v1.GetSubscribeCategoryReply_Category{
			Id:   item.Id,
			Name: item.Name,
		})
	}
	return rt, nil
}
func (m *MPProxyService) GetSubscribeTplKeywords(ctx context.Context, req *v1.GetSubscribeTplKeywordsRequest) (*v1.GetSubscribeTplKeywordsReply, error) {
	res, err := m.uc.GetSubscribeTplKeywords(ctx, req.AccessToken, req.TemplateId)
	if err != nil {
		return nil, err
	}
	var rt = &v1.GetSubscribeTplKeywordsReply{
		Count: int64(res.Count),
		Data:  []*v1.GetSubscribeTplKeywordsReply_Item{},
	}
	for _, item := range res.Data {
		rt.Data = append(rt.Data, &v1.GetSubscribeTplKeywordsReply_Item{
			Name:    item.Name,
			Example: item.Example,
			Kid:     int64(item.Kid),
			Rule:    item.Rule,
		})
	}

	return rt, nil
}
func (m *MPProxyService) GetSubscribeTplTitles(ctx context.Context, req *v1.GetSubscribeTplTitlesRequest) (*v1.GetSubscribeTplTitlesReply, error) {
	res, err := m.uc.GetSubscribeTplTitles(ctx, req.AccessToken, req.Ids, req.Start, req.Limit)
	if err != nil {
		return nil, err
	}

	var rt = &v1.GetSubscribeTplTitlesReply{
		Count: int64(res.Count),
		Data:  []*v1.GetSubscribeTplTitlesReply_Item{},
	}
	for _, item := range res.Data {
		rt.Data = append(rt.Data, &v1.GetSubscribeTplTitlesReply_Item{
			Tid:        item.Tid,
			Title:      item.Title,
			Type:       int64(item.Type),
			CategoryId: item.CategoryId,
		})
	}

	return rt, nil
}
func (m *MPProxyService) GetSubscribePrivateTpl(ctx context.Context, req *v1.AccessTokenParam) (*v1.GetSubscribePrivateTplReply, error) {
	res, err := m.uc.GetSubscribePrivateTpl(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}
	var rt = &v1.GetSubscribePrivateTplReply{
		Data: []*v1.GetSubscribePrivateTplReply_Item{},
	}
	for _, item := range res.Data {
		rt.Data = append(rt.Data, &v1.GetSubscribePrivateTplReply_Item{
			PriTmplId: item.PriTmplId,
			Title:     item.Title,
			Content:   item.Content,
			Example:   item.Example,
			Type:      int64(item.Type),
		})
	}

	return rt, nil
}
func (m *MPProxyService) SendSubscribeMessage(ctx context.Context, req *v1.SendSubscribeMessageRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendSubscribeMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

// GetKFList 获取客服列表
func (m *MPProxyService) GetKFList(ctx context.Context, req *v1.AccessTokenParam) (*v1.GetKFListReply, error) {
	res, err := m.uc.GetKFList(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	var rt = &v1.GetKFListReply{
		KfList: []*v1.KeFuInfo{},
	}
	for _, item := range res.KfList {
		rt.KfList = append(rt.KfList, &v1.KeFuInfo{
			KfAccount:        item.KfAccount,
			KfNick:           item.KfNick,
			KfId:             fmt.Sprintf("%d", item.KfID),
			KfWx:             item.KfWX,
			KfHeadImgUrl:     item.KfHeadImgURL,
			InviteWx:         item.InviteWX,
			InviteStatus:     item.InviteStatus,
			InviteExpireTime: int64(item.InviteExpTime),
		})
	}

	return rt, nil
}

// GetKFOnlineList 获取在线客服列表
// ctx 上下文对象，用于控制请求的生命周期
// req 包含访问令牌的参数
// 返回值：包含在线客服列表的回复对象和可能的错误信息
func (m *MPProxyService) GetKFOnlineList(ctx context.Context, req *v1.AccessTokenParam) (*v1.GetKFOnlineListReply, error) {
	res, err := m.uc.GetKFOnlineList(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	var rt = &v1.GetKFOnlineListReply{
		KfOnlineList: []*v1.KFOnlineInfo{},
	}
	for _, item := range res.KfOnlineList {
		rt.KfOnlineList = append(rt.KfOnlineList, &v1.KFOnlineInfo{
			AcceptedCase: int64(item.AcceptedCase),
			KfAccount:    item.KfAccount,
			KfId:         int64(item.KfID),
			Status:       int64(item.Status),
		})
	}

	return rt, nil
}

// GetKFMsgHistory 获取客服聊天记录
//   - OpCode操作码，2002（客服发送信息），2003（客服接收消息）
func (m *MPProxyService) GetKFMsgHistory(ctx context.Context, req *v1.GetKFMsgHistoryRequest) (*v1.GetKFMsgHistoryReply, error) {
	res, err := m.uc.GetKFMsgHistory(ctx, req.AccessToken, req.StartTime, req.EndTime, req.MsgId, req.Number)
	if err != nil {
		return nil, err
	}
	var rt = &v1.GetKFMsgHistoryReply{
		RecordList: []*v1.KFMsgHistory{},
		MsgId:      res.MsgId,
		Number:     res.Number,
	}
	for _, item := range res.MsgRecordList {
		rt.RecordList = append(rt.RecordList, &v1.KFMsgHistory{
			Worker: item.Worker,
			OpenId: item.OpenID,
			Text:   item.Text,
			Time:   item.Time,
			OpCode: item.OpCode,
		})
	}

	return rt, nil
}

// AddKFAccount 添加客服账号
//
//	整客服账号，格式为：账号前缀@公众号微信号，账号前缀最多10个字符，必须是英文、数字字符或者下划线，后缀为公众号微信号，长度不超过30个字符
func (m *MPProxyService) AddKFAccount(ctx context.Context, req *v1.AddKFAccountRequest) (*v1.WXErrorReply, error) {
	err := m.uc.AddKFAccount(ctx, req.AccessToken, req.KfAccount, req.Nickname, req.Password)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) UpdateKFAccount(ctx context.Context, req *v1.UpdateKFAccountRequest) (*v1.WXErrorReply, error) {
	err := m.uc.UpdateKFAccount(ctx, req.AccessToken, req.KfAccount, req.Nickname)
	if err != nil {
		return nil, err
	}
	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) DelKFAccount(ctx context.Context, req *v1.DelKFAccountRequest) (*v1.WXErrorReply, error) {
	err := m.uc.DelKFAccount(ctx, req.AccessToken, req.KfAccount)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) InviteKFWorker(ctx context.Context, req *v1.InviteKFWorkerRequest) (*v1.WXErrorReply, error) {
	err := m.uc.InviteKFWorker(ctx, req.AccessToken, req.KfAccount, req.InviteWx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m *MPProxyService) UpdateKFAvatar(context.Context, *v1.UpdateKFAvatarRequest) (*v1.WXErrorReply, error) {
	return nil, nil
}

// 下发客服输入状态
//
// "Typing"：正在输入"状态
// "CancelTyping"：取消对用户的"正在输入"状态
func (m *MPProxyService) UpdateKFTyping(ctx context.Context, req *v1.UpdateKFTypingRequest) (*v1.WXErrorReply, error) {
	err := m.uc.UpdateKFTyping(ctx, req.AccessToken, req.Touser, req.Command)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m *MPProxyService) GetKFSessionList(ctx context.Context, req *v1.GetKFSessionListRequest) (*v1.GetKFSessionListReply, error) {
	res, err := m.uc.GetKFSessionList(ctx, req.AccessToken, req.KfAccount)
	if err != nil {
		return nil, err
	}
	var rt = &v1.GetKFSessionListReply{
		SessionList: []*v1.KFSession{},
	}
	for _, item := range res.SessionList {
		rt.SessionList = append(rt.SessionList, &v1.KFSession{
			CreateTime: item.CreateTime,
			OpenId:     item.OpenID,
		})
	}
	return rt, nil
}

func (m *MPProxyService) GetKFSessionStatus(ctx context.Context, req *v1.GetKFSessionStatusRequest) (*v1.GetKFSessionStatusReply, error) {
	res, err := m.uc.GetKFSessionStatus(ctx, req.AccessToken, req.OpenId)
	if err != nil {
		return nil, err
	}
	var rt = &v1.GetKFSessionStatusReply{
		KfAccount:  res.KfAccount,
		CreateTime: res.CreateTime,
	}
	return rt, nil
}

func (m *MPProxyService) GetKFSessionUnaccepted(ctx context.Context, req *v1.AccessTokenParam) (*v1.GetKFSessionUnacceptedReply, error) {
	res, err := m.uc.GetKFSessionUnaccepted(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}
	var rt = &v1.GetKFSessionUnacceptedReply{
		Count:        int64(res.Count),
		WaitCaseList: []*v1.GetKFSessionUnacceptedReply_WaitCase{},
	}
	for _, item := range res.WaitCaseList {
		rt.WaitCaseList = append(rt.WaitCaseList, &v1.GetKFSessionUnacceptedReply_WaitCase{
			LatestTime: item.LatestTime,
			OpenId:     item.OpenID,
		})
	}
	return rt, nil
}

func (m *MPProxyService) CloseKFSession(ctx context.Context, req *v1.CloseKFSessionRequest) (*v1.WXErrorReply, error) {
	err := m.uc.CloseKFSession(ctx, req.AccessToken, req.KfAccount, req.OpenId)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) NewKFSession(ctx context.Context, req *v1.NewKFSessionRequest) (*v1.WXErrorReply, error) {
	err := m.uc.NewKFSession(ctx, req.AccessToken, req.KfAccount, req.OpenId)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) SendKFTextMsg(ctx context.Context, req *v1.SendKFTextMsgRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendKFTextMsg(ctx, req)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) SendKFImageMsg(ctx context.Context, req *v1.SendKFImageMsgRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendKFImageMsg(ctx, req)
	if err != nil {
		return nil, err
	}

	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) SendKFVoiceMsg(ctx context.Context, req *v1.SendKFVoiceMsgRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendKFVoiceMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) SendKFVideoMsg(ctx context.Context, req *v1.SendKFVideoMsgRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendKFVideoMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) SendKFMusicMsg(ctx context.Context, req *v1.SendKFMusicMsgRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendKFMusicMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

// SendKFNewsCardMsg 图文消息，点击跳转到外链
func (m *MPProxyService) SendKFNewsCardMsg(ctx context.Context, req *v1.SendKFNewsCardMsgRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendKFNewsCardMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) SendKFNewsPageMsg(ctx context.Context, req *v1.SendKFNewsPageMsgRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendKFNewsPageMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) SendKFToArticleMsg(ctx context.Context, req *v1.SendKFToArticleMsgRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendKFToArticleMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) SendKFMenuMsg(ctx context.Context, req *v1.SendKFMenuMsgRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendKFMenuMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) SendKFCardMsg(ctx context.Context, req *v1.SendKFCardMsgRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendKFCardMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) SendKFMiniProgramMsg(ctx context.Context, req *v1.SendKFMiniProgramMsgRequest) (*v1.WXErrorReply, error) {
	err := m.uc.SendKFMiniProgramMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.WXErrorReply{Errcode: 0, Errmsg: "ok"}, nil
}

func (m *MPProxyService) GetBlacklist(ctx context.Context, req *v1.GetBlacklistReq) (*v1.GetBlacklistReply, error) {
	reply,err:= m.uc.GetBlacklist(ctx, req.AccessToken, req.NextOpenid)
	if err!= nil {
		return nil, err
	}

	return reply, err
}
