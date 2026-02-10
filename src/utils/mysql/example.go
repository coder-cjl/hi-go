package mysql

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"gorm.io/gorm"
// )

// // ==================== 示例模型 ====================

// // User 用户模型示例
// type User struct {
// 	ID        uint           `gorm:"primaryKey" json:"id"`
// 	Name      string         `gorm:"size:100;not null" json:"name"`
// 	Email     string         `gorm:"size:200;uniqueIndex" json:"email"`
// 	Age       int            `gorm:"default:0" json:"age"`
// 	CreatedAt time.Time      `json:"created_at"`
// 	UpdatedAt time.Time      `json:"updated_at"`
// 	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
// }

// // TableName 指定表名
// func (User) TableName() string {
// 	return "users"
// }

// // ==================== 使用示例 ====================

// // ExampleAutoMigrate 自动迁移示例
// func ExampleAutoMigrate() {
// 	// 自动创建或更新表结构
// 	if err := AutoMigrate(&User{}); err != nil {
// 		fmt.Println("迁移失败:", err)
// 		return
// 	}
// 	fmt.Println("迁移成功")
// }

// // ExampleCreate 创建记录示例
// func ExampleCreate() {
// 	// 创建单条记录
// 	user := &User{
// 		Name:  "张三",
// 		Email: "zhangsan@example.com",
// 		Age:   25,
// 	}
// 	if err := Create(user); err != nil {
// 		fmt.Println("创建失败:", err)
// 		return
// 	}
// 	fmt.Println("创建成功, ID:", user.ID)
// }

// // ExampleCreateBatch 批量创建示例
// func ExampleCreateBatch() {
// 	users := []User{
// 		{Name: "用户1", Email: "user1@example.com", Age: 20},
// 		{Name: "用户2", Email: "user2@example.com", Age: 22},
// 		{Name: "用户3", Email: "user3@example.com", Age: 24},
// 	}
// 	// 每批次插入 100 条
// 	if err := CreateBatch(&users, 100); err != nil {
// 		fmt.Println("批量创建失败:", err)
// 		return
// 	}
// 	fmt.Println("批量创建成功")
// }

// // ExampleFindByID 根据 ID 查询示例
// func ExampleFindByID() {
// 	var user User
// 	if err := FindByID(&user, 1); err != nil {
// 		fmt.Println("查询失败:", err)
// 		return
// 	}
// 	fmt.Printf("用户: %+v\n", user)
// }

// // ExampleFindOne 条件查询单条示例
// func ExampleFindOne() {
// 	var user User
// 	if err := FindOne(&user, "email = ?", "zhangsan@example.com"); err != nil {
// 		fmt.Println("查询失败:", err)
// 		return
// 	}
// 	fmt.Printf("用户: %+v\n", user)
// }

// // ExampleFindAll 查询多条记录示例
// func ExampleFindAll() {
// 	var users []User
// 	// 查询年龄大于 20 的所有用户
// 	if err := FindAll(&users, "age > ?", 20); err != nil {
// 		fmt.Println("查询失败:", err)
// 		return
// 	}
// 	for _, u := range users {
// 		fmt.Printf("用户: %+v\n", u)
// 	}
// }

// // ExampleUpdate 更新记录示例
// func ExampleUpdate() {
// 	// 根据 ID 更新
// 	if err := Update(&User{}, 1, map[string]interface{}{
// 		"name": "李四",
// 		"age":  30,
// 	}); err != nil {
// 		fmt.Println("更新失败:", err)
// 		return
// 	}
// 	fmt.Println("更新成功")
// }

// // ExampleUpdateByCondition 条件更新示例
// func ExampleUpdateByCondition() {
// 	// 将年龄小于 18 的用户状态设为未成年
// 	if err := UpdateByCondition(&User{}, "age < ?", map[string]interface{}{
// 		"age": 18,
// 	}, 18); err != nil {
// 		fmt.Println("更新失败:", err)
// 		return
// 	}
// 	fmt.Println("条件更新成功")
// }

// // ExampleDelete 删除记录示例
// func ExampleDelete() {
// 	// 软删除（如果模型有 DeletedAt 字段）
// 	if err := Delete(&User{}, 1); err != nil {
// 		fmt.Println("删除失败:", err)
// 		return
// 	}
// 	fmt.Println("软删除成功")
// }

// // ExampleHardDelete 硬删除示例
// func ExampleHardDelete() {
// 	// 永久删除，不可恢复
// 	if err := HardDelete(&User{}, 1); err != nil {
// 		fmt.Println("硬删除失败:", err)
// 		return
// 	}
// 	fmt.Println("硬删除成功")
// }

// // ExampleExists 检查是否存在示例
// func ExampleExists() {
// 	exists, err := Exists(&User{}, "email = ?", "zhangsan@example.com")
// 	if err != nil {
// 		fmt.Println("检查失败:", err)
// 		return
// 	}
// 	if exists {
// 		fmt.Println("用户存在")
// 	} else {
// 		fmt.Println("用户不存在")
// 	}
// }

// // ExampleCount 统计数量示例
// func ExampleCount() {
// 	count, err := Count(&User{}, "age >= ?", 18)
// 	if err != nil {
// 		fmt.Println("统计失败:", err)
// 		return
// 	}
// 	fmt.Println("成年用户数量:", count)
// }

// // ExamplePaginate 分页查询示例
// func ExamplePaginate() {
// 	var users []User
// 	// 查询第 1 页，每页 10 条，年龄大于 18
// 	result, err := Paginate(&users, &User{}, 1, 10, "age > ?", 18)
// 	if err != nil {
// 		fmt.Println("分页查询失败:", err)
// 		return
// 	}
// 	fmt.Printf("总数: %d, 当前页: %d, 总页数: %d\n", result.Total, result.Page, result.Pages)
// 	for _, u := range users {
// 		fmt.Printf("用户: %+v\n", u)
// 	}
// }

// // ExamplePaginateWithOrder 带排序的分页查询示例
// func ExamplePaginateWithOrder() {
// 	var users []User
// 	// 按创建时间倒序分页
// 	result, err := PaginateWithOrder(&users, &User{}, 1, 10, "created_at DESC", "age > ?", 18)
// 	if err != nil {
// 		fmt.Println("分页查询失败:", err)
// 		return
// 	}
// 	fmt.Printf("总数: %d\n", result.Total)
// }

// // ExampleTransaction 事务示例
// func ExampleTransaction() {
// 	err := Transaction(func(tx *gorm.DB) error {
// 		// 在事务中执行多个操作
// 		user := &User{Name: "事务用户", Email: "tx@example.com", Age: 25}
// 		if err := tx.Create(user).Error; err != nil {
// 			return err // 返回错误会自动回滚
// 		}

// 		// 更新另一个用户
// 		if err := tx.Model(&User{}).Where("id = ?", 2).Update("age", 30).Error; err != nil {
// 			return err // 返回错误会自动回滚
// 		}

// 		return nil // 返回 nil 会自动提交
// 	})

// 	if err != nil {
// 		fmt.Println("事务执行失败:", err)
// 		return
// 	}
// 	fmt.Println("事务执行成功")
// }

// // ExampleManualTransaction 手动事务示例
// func ExampleManualTransaction() {
// 	tx := BeginTx()

// 	user := &User{Name: "手动事务用户", Email: "manual@example.com", Age: 28}
// 	if err := tx.Create(user).Error; err != nil {
// 		tx.Rollback() // 手动回滚
// 		fmt.Println("创建失败，已回滚:", err)
// 		return
// 	}

// 	// 可以继续执行其他操作...

// 	tx.Commit() // 手动提交
// 	fmt.Println("手动事务提交成功")
// }

// // ExampleRawQuery 原生 SQL 查询示例
// func ExampleRawQuery() {
// 	var users []User
// 	if err := RawQuery(&users, "SELECT * FROM users WHERE age > ? ORDER BY created_at DESC LIMIT ?", 18, 10); err != nil {
// 		fmt.Println("原生查询失败:", err)
// 		return
// 	}
// 	for _, u := range users {
// 		fmt.Printf("用户: %+v\n", u)
// 	}
// }

// // ExampleExec 原生 SQL 执行示例
// func ExampleExec() {
// 	// 执行原生 UPDATE
// 	if err := Exec("UPDATE users SET age = age + 1 WHERE age < ?", 30); err != nil {
// 		fmt.Println("执行失败:", err)
// 		return
// 	}
// 	fmt.Println("执行成功")
// }

// // ExampleDebug 调试模式示例
// func ExampleDebug() {
// 	var users []User
// 	// 开启调试模式会打印 SQL 语句
// 	Debug().Where("age > ?", 18).Find(&users)
// }

// // ExampleWithContext 使用 Context 示例
// func ExampleWithContext() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	var users []User
// 	// 带超时的查询
// 	if err := WithContext(ctx).Where("age > ?", 18).Find(&users).Error; err != nil {
// 		fmt.Println("查询失败:", err)
// 		return
// 	}
// }

// // ExamplePing 健康检查示例
// func ExamplePing() {
// 	if err := Ping(); err != nil {
// 		fmt.Println("数据库连接异常:", err)
// 		return
// 	}
// 	fmt.Println("数据库连接正常")
// }

// // ExampleConnectionPool 连接池配置示例
// func ExampleConnectionPool() {
// 	// 设置最大打开连接数
// 	SetMaxOpenConns(200)

// 	// 设置最大空闲连接数
// 	SetMaxIdleConns(20)

// 	// 设置连接最大存活时间
// 	SetConnMaxLifetime(2 * time.Hour)

// 	// 设置空闲连接最大存活时间
// 	SetConnMaxIdleTime(30 * time.Minute)

// 	fmt.Println("连接池配置完成")
// }
