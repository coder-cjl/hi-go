package aiservice

import "context"

// Skill AI技能接口
type Skill interface {
	// Name 返回技能名称
	Name() string

	// Description 返回技能描述（用于AI理解）
	Description() string

	// Parameters 返回技能参数定义（用于Function Calling）
	Parameters() map[string]interface{}

	// Execute 执行技能
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)

	// IsEnabled 是否启用
	IsEnabled() bool
}

// SkillRegistry 技能注册表
type SkillRegistry struct {
	skills map[string]Skill
}

// NewSkillRegistry 创建技能注册表
func NewSkillRegistry() *SkillRegistry {
	return &SkillRegistry{
		skills: make(map[string]Skill),
	}
}

// Register 注册技能
func (r *SkillRegistry) Register(skill Skill) {
	r.skills[skill.Name()] = skill
}

// Get 获取技能
func (r *SkillRegistry) Get(name string) (Skill, bool) {
	skill, ok := r.skills[name]
	return skill, ok
}

// GetAll 获取所有已启用的技能
func (r *SkillRegistry) GetAll() []Skill {
	skills := make([]Skill, 0, len(r.skills))
	for _, skill := range r.skills {
		if skill.IsEnabled() {
			skills = append(skills, skill)
		}
	}
	return skills
}
