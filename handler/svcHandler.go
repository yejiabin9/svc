package handler

import (
	"context"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/sirupsen/logrus"
	"github.com/yejiabin9/svc/domain/model"
	"github.com/yejiabin9/svc/domain/service"
	"github.com/yejiabin9/svc/proto/svc"
	"github.com/yejiabin9/svc/utils"
	"strconv"
)

type SvcHandler struct {
	//注意这里的类型是 ISvcDataService 接口类型
	SvcDataService service.ISvcDataService
}

// 添加服务
func (e *SvcHandler) AddSvc(ctx context.Context, info *svc.SvcInfo, rsp *svc.Response) error {
	log.Info("创建服务")
	svcModel := &model.Svc{}
	//数据类型转换
	if err := utils.SwapTo(info, svcModel); err != nil {
		logrus.Error(err)
		return err
	}

	//到 k8s 中创建服务
	if err := e.SvcDataService.CreateSvcToK8s(info); err != nil {
		logrus.Error(err)
		return err
	} else {
		svcID, err := e.SvcDataService.AddSvc(svcModel)
		if err != nil {
			//如果逻辑需要自行实现k8s中删除操作
			logrus.Error(err)
			return err
		}
		logrus.Info("Svc 添加数据成功ID号为：" + strconv.FormatInt(svcID, 10))
		rsp.Msg = "Svc 添加数据成功ID号为：" + strconv.FormatInt(svcID, 10)
	}
	return nil
}

// 删除服务
func (e *SvcHandler) DeleteSvc(ctx context.Context, req *svc.SvcId, rsp *svc.Response) error {
	log.Info("删除服务")
	service, err := e.SvcDataService.FIndSvcById(req.Id)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if err := e.SvcDataService.DeleteFromK8s(service); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// 更新svc
func (e *SvcHandler) UpdateSvc(ctx context.Context, req *svc.SvcInfo, rsp *svc.Response) error {
	log.Info("Received *svc.UpdateSvc request")
	//先更新k8s里面的数据
	if err := e.SvcDataService.UpdateSvcToK8s(req); err != nil {
		logrus.Error(err)
		return err
	}
	//查询数据库中的svc
	service, err := e.SvcDataService.FIndSvcById(req.Id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	//数据类型转换
	if err := utils.SwapTo(req, service); err != nil {
		logrus.Error(err)
		return err
	}
	//更新到数据中
	if err := e.SvcDataService.UpDateSvc(service); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// 服务查找
func (e *SvcHandler) FindSvcByID(ctx context.Context, req *svc.SvcId, rsp *svc.SvcInfo) error {
	log.Info("查找服务")
	svcModel, err := e.SvcDataService.FIndSvcById(req.Id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if err := utils.SwapTo(svcModel, rsp); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// 查找所有服务
func (e *SvcHandler) FindAllSvc(ctx context.Context, req *svc.FindAll, rsp *svc.AllSvc) error {
	log.Info("查询所有服务")
	allSvc, err := e.SvcDataService.FIndAllSvc()
	if err != nil {
		logrus.Error(err)
		return err
	}
	//整理格式
	for _, v := range allSvc {
		svcInfo := &svc.SvcInfo{}
		if err := utils.SwapTo(v, svcInfo); err != nil {
			logrus.Error(err)
			return err
		}
		rsp.SvcInfo = append(rsp.SvcInfo, svcInfo)
	}
	return nil
}
