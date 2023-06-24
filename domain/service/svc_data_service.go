package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/yejiabin9/svc/domain/model"
	"github.com/yejiabin9/svc/domain/repository"
	"github.com/yejiabin9/svc/proto/svc"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

type ISvcDataService interface {
	AddSvc(svc *model.Svc) (int64, error)
	DeleteSvc(int64) error
	UpDateSvc(svc *model.Svc) error
	FIndSvcById(int64) (*model.Svc, error)
	FIndAllSvc() ([]model.Svc, error)
	CreateSvcToK8s(info *svc.SvcInfo) error
	UpdateSvcToK8s(info *svc.SvcInfo) error
	DeleteFromK8s(svc2 *model.Svc) error
}

func NewSvcDataService(svcRepository repository.ISvcRepository, clientSet *kubernetes.Clientset) ISvcDataService {
	return &SvcDataService{
		SvcRepository: svcRepository,
		K8sClientSet:  clientSet,
	}

}

type SvcDataService struct {
	SvcRepository repository.ISvcRepository
	K8sClientSet  *kubernetes.Clientset
}

func (s *SvcDataService) CreateSvcToK8s(info *svc.SvcInfo) (err error) {
	service := s.setService(info)
	if _, err = s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Get(context.TODO(),
		info.SvcName, v12.GetOptions{}); err != nil {
		//create
		if _, err = s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Create(context.TODO(),
			service, v12.CreateOptions{}); err != nil {
			logrus.Error("create service error")
			return err
		}
		return nil
	} else {
		logrus.Error("Service" + info.SvcName + "已经存在")
		return errors.New("Service" + info.SvcName + "已经存在")
	}
}

func (s *SvcDataService) setService(svcInfo *svc.SvcInfo) *v1.Service {
	service := &v1.Service{}

	service.TypeMeta = v12.TypeMeta{
		Kind:       "v1",
		APIVersion: "Service",
	}
	service.ObjectMeta = v12.ObjectMeta{
		Name:         svcInfo.SvcName,
		GenerateName: svcInfo.SvcNamespace,
		Labels: map[string]string{
			"app-name": svcInfo.SvcPodName,
			"author":   "jiabin",
		},
		Annotations: map[string]string{
			"k8s.generated-by-jiabin": "auther by jiabin",
		},
	}
	//clusterIp
	service.Spec = v1.ServiceSpec{
		Ports: s.getSvcPort(svcInfo),
		Selector: map[string]string{
			"app-name": svcInfo.SvcPodName,
		},

		Type: "ClusterIP",
	}
	return service

}

func (s *SvcDataService) getSvcPort(info *svc.SvcInfo) (servicePort []v1.ServicePort) {
	for _, v := range info.SvcPort {
		servicePort = append(servicePort, v1.ServicePort{
			Name:       "port-" + strconv.FormatInt(int64(v.SvcPort), 10),
			Protocol:   v1.Protocol(v.SvcPortProtocol),
			Port:       v.SvcPort,
			TargetPort: intstr.FromInt(int(v.SvcTargetPort)),
		})
	}
	return servicePort
}

func (s *SvcDataService) UpdateSvcToK8s(info *svc.SvcInfo) (err error) {
	service := s.setService(info)
	if _, err = s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Get(context.TODO(),
		info.SvcName, v12.GetOptions{}); err != nil {

		return nil
	} else {
		//create
		if _, err = s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Update(context.TODO(),
			service, v12.UpdateOptions{}); err != nil {
			logrus.Error("create service error")
			return err
		}

		logrus.Info("Service " + info.SvcName + "update success")
		return nil
	}

}

func (s *SvcDataService) DeleteFromK8s(svc2 *model.Svc) (err error) {
	if err = s.K8sClientSet.CoreV1().Services(svc2.SvcNamespace).Delete(context.TODO(), svc2.SvcName, v12.DeleteOptions{}); err != nil {
		logrus.Error(err)
		return err
	} else {
		if err := s.DeleteSvc(svc2.ID); err != nil {
			logrus.Error(err)
			return err
		}
		logrus.Info("delete service ID " + strconv.FormatInt(svc2.ID, 10) + "success")
		return nil
	}
}

func (s *SvcDataService) AddSvc(svc *model.Svc) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s SvcDataService) DeleteSvc(i int64) error {
	//TODO implement me
	panic("implement me")
}

func (s SvcDataService) UpDateSvc(svc *model.Svc) error {
	//TODO implement me
	panic("implement me")
}

func (s SvcDataService) FIndSvcById(i int64) (*model.Svc, error) {
	//TODO implement me
	panic("implement me")
}

func (s SvcDataService) FIndAllSvc() ([]model.Svc, error) {
	//TODO implement me
	panic("implement me")
}
