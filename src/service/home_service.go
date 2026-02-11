package service

import (
	"fmt"
	"hi-go/src/config"
	"hi-go/src/model"
	"hi-go/src/repository"
)

// 首页业务逻辑层
type HomeService struct {
	homeRepo *repository.HomeRepository
}

// 创建首页服务实例
func NewHomeService() *HomeService {
	return &HomeService{
		homeRepo: repository.NewHomeRepository(),
	}
}

// 获取首页列表
func (s *HomeService) GetList(req *model.HomeListRequest) (*model.HomeListResponse, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = config.Config.Business.DefaultPageSize
	}
	if req.PageSize > config.Config.Business.MaxPageSize {
		req.PageSize = config.Config.Business.MaxPageSize
	}

	// 查询数据
	list, total, err := s.homeRepo.List(req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &model.HomeListResponse{
		List:  list,
		Total: total,
	}, nil
}

// 创建模拟数据（30条）
func (s *HomeService) CreateMockData() error {
	// 生成30条模拟数据
	homes := make([]model.Home, 30)
	for i := 0; i < 30; i++ {
		homes[i] = model.Home{
			Title:       fmt.Sprintf("首页内容标题 %d", i+1),
			Description: fmt.Sprintf("这是第 %d 条首页内容的描述信息，用于展示在首页列表中", i+1),
			ImageURL:    fmt.Sprintf("https://picsum.photos/400/300?random=%d", i+1),
			Link:        fmt.Sprintf("https://example.com/detail/%d", i+1),
			Sort:        i + 1,
			Status:      1,
		}
	}

	// 批量插入数据库
	return s.homeRepo.BatchCreate(homes)
}
