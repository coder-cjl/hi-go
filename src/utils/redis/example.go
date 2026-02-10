package redis

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	goredis "github.com/redis/go-redis/v9"
// )

// // Example 展示Redis工具的基本使用方法
// func Example() {
// 	ctx := context.Background()

// 	fmt.Println("=== Redis 工具使用示例 ===\n")

// 	// ========== 1. String 操作 ==========
// 	fmt.Println("1. String 操作")

// 	// 设置键值对
// 	err := Set(ctx, "user:name", "张三", 10*time.Minute)
// 	if err != nil {
// 		fmt.Printf("   ❌ Set 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Println("   ✓ Set user:name = 张三")

// 	// 获取值
// 	name, err := Get(ctx, "user:name")
// 	if err != nil {
// 		fmt.Printf("   ❌ Get 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ Get user:name = %s\n", name)

// 	// 批量设置
// 	err = MSet(ctx, "user:age", "25", "user:city", "北京")
// 	if err != nil {
// 		fmt.Printf("   ❌ MSet 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Println("   ✓ MSet user:age=25, user:city=北京")

// 	// 批量获取
// 	values, err := MGet(ctx, "user:name", "user:age", "user:city")
// 	if err != nil {
// 		fmt.Printf("   ❌ MGet 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ MGet: %v\n", values)

// 	// 计数器操作
// 	count, err := Incr(ctx, "page:views")
// 	if err != nil {
// 		fmt.Printf("   ❌ Incr 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 页面访问量: %d\n\n", count)

// 	// ========== 2. Hash 操作 ==========
// 	fmt.Println("2. Hash 操作")

// 	// 设置哈希字段
// 	_, err = HSet(ctx, "user:1001", "name", "李四", "age", "30", "city", "上海")
// 	if err != nil {
// 		fmt.Printf("   ❌ HSet 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Println("   ✓ HSet user:1001")

// 	// 获取单个字段
// 	userName, err := HGet(ctx, "user:1001", "name")
// 	if err != nil {
// 		fmt.Printf("   ❌ HGet 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ HGet user:1001 name = %s\n", userName)

// 	// 获取所有字段
// 	userInfo, err := HGetAll(ctx, "user:1001")
// 	if err != nil {
// 		fmt.Printf("   ❌ HGetAll 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ HGetAll user:1001: %v\n", userInfo)

// 	// 字段计数器
// 	newAge, err := HIncrBy(ctx, "user:1001", "age", 1)
// 	if err != nil {
// 		fmt.Printf("   ❌ HIncrBy 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 年龄增加后: %d\n\n", newAge)

// 	// ========== 3. List 操作 ==========
// 	fmt.Println("3. List 操作")

// 	// 从左侧插入
// 	_, err = LPush(ctx, "messages", "消息3", "消息2", "消息1")
// 	if err != nil {
// 		fmt.Printf("   ❌ LPush 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Println("   ✓ LPush messages")

// 	// 从右侧插入
// 	_, err = RPush(ctx, "messages", "消息4", "消息5")
// 	if err != nil {
// 		fmt.Printf("   ❌ RPush 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Println("   ✓ RPush messages")

// 	// 获取列表长度
// 	length, err := LLen(ctx, "messages")
// 	if err != nil {
// 		fmt.Printf("   ❌ LLen 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 列表长度: %d\n", length)

// 	// 获取列表元素
// 	messages, err := LRange(ctx, "messages", 0, -1)
// 	if err != nil {
// 		fmt.Printf("   ❌ LRange 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 所有消息: %v\n", messages)

// 	// 弹出元素
// 	msg, err := LPop(ctx, "messages")
// 	if err != nil {
// 		fmt.Printf("   ❌ LPop 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ LPop: %s\n\n", msg)

// 	// ========== 4. Set 操作 ==========
// 	fmt.Println("4. Set 操作")

// 	// 添加成员
// 	_, err = SAdd(ctx, "tags", "Go", "Redis", "MySQL", "Docker")
// 	if err != nil {
// 		fmt.Printf("   ❌ SAdd 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Println("   ✓ SAdd tags")

// 	// 获取所有成员
// 	tags, err := SMembers(ctx, "tags")
// 	if err != nil {
// 		fmt.Printf("   ❌ SMembers 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 所有标签: %v\n", tags)

// 	// 检查成员
// 	isMember, err := SIsMember(ctx, "tags", "Go")
// 	if err != nil {
// 		fmt.Printf("   ❌ SIsMember 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 'Go' 是否存在: %v\n", isMember)

// 	// 获取成员数量
// 	tagCount, err := SCard(ctx, "tags")
// 	if err != nil {
// 		fmt.Printf("   ❌ SCard 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 标签数量: %d\n\n", tagCount)

// 	// ========== 5. Sorted Set 操作 ==========
// 	fmt.Println("5. Sorted Set 操作（排行榜）")

// 	// 添加成员（分数为得分）
// 	_, err = ZAdd(ctx, "leaderboard",
// 		goredis.Z{Score: 100, Member: "玩家A"},
// 		goredis.Z{Score: 95, Member: "玩家B"},
// 		goredis.Z{Score: 110, Member: "玩家C"},
// 		goredis.Z{Score: 88, Member: "玩家D"},
// 	)
// 	if err != nil {
// 		fmt.Printf("   ❌ ZAdd 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Println("   ✓ ZAdd leaderboard")

// 	// 获取排行榜（从高到低）
// 	topPlayers, err := ZRevRange(ctx, "leaderboard", 0, 2)
// 	if err != nil {
// 		fmt.Printf("   ❌ ZRevRange 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 前3名: %v\n", topPlayers)

// 	// 获取玩家分数
// 	score, err := ZScore(ctx, "leaderboard", "玩家C")
// 	if err != nil {
// 		fmt.Printf("   ❌ ZScore 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 玩家C 分数: %.0f\n", score)

// 	// 增加分数
// 	newScore, err := ZIncrBy(ctx, "leaderboard", 5, "玩家B")
// 	if err != nil {
// 		fmt.Printf("   ❌ ZIncrBy 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 玩家B 新分数: %.0f\n\n", newScore)

// 	// ========== 6. 过期时间操作 ==========
// 	fmt.Println("6. 过期时间操作")

// 	// 设置键的过期时间
// 	Set(ctx, "session:abc123", "user_data", 0)
// 	ok, err := Expire(ctx, "session:abc123", 30*time.Second)
// 	if err != nil {
// 		fmt.Printf("   ❌ Expire 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 设置过期时间: %v\n", ok)

// 	// 获取剩余时间
// 	ttl, err := TTL(ctx, "session:abc123")
// 	if err != nil {
// 		fmt.Printf("   ❌ TTL 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 剩余时间: %v\n\n", ttl)

// 	// ========== 7. 分布式锁示例 ==========
// 	fmt.Println("7. 分布式锁（SetNX）")

// 	lockKey := "lock:resource:1"
// 	locked, err := SetNX(ctx, lockKey, "locked", 10*time.Second)
// 	if err != nil {
// 		fmt.Printf("   ❌ SetNX 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ 获取锁: %v\n", locked)

// 	if locked {
// 		fmt.Println("   ✓ 执行业务逻辑...")
// 		time.Sleep(1 * time.Second)

// 		// 释放锁
// 		Del(ctx, lockKey)
// 		fmt.Println("   ✓ 释放锁\n")
// 	}

// 	// ========== 8. 管道操作 ==========
// 	fmt.Println("8. 管道操作（Pipeline）")

// 	cmds, err := Pipeline(ctx, func(pipe goredis.Pipeliner) error {
// 		pipe.Set(ctx, "key1", "value1", 0)
// 		pipe.Set(ctx, "key2", "value2", 0)
// 		pipe.Set(ctx, "key3", "value3", 0)
// 		pipe.Incr(ctx, "counter")
// 		return nil
// 	})
// 	if err != nil {
// 		fmt.Printf("   ❌ Pipeline 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ Pipeline 执行了 %d 条命令\n\n", len(cmds))

// 	// ========== 9. 键操作 ==========
// 	fmt.Println("9. 键操作")

// 	// 检查键是否存在
// 	exists, err := Exists(ctx, "user:name")
// 	if err != nil {
// 		fmt.Printf("   ❌ Exists 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ user:name 存在: %v\n", exists)

// 	// 获取键的类型
// 	keyType, err := Type(ctx, "user:1001")
// 	if err != nil {
// 		fmt.Printf("   ❌ Type 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ user:1001 类型: %s\n", keyType)

// 	// 查找所有符合模式的键
// 	allKeys, err := Keys(ctx, "user:*")
// 	if err != nil {
// 		fmt.Printf("   ❌ Keys 失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("   ✓ user:* 匹配的键: %v\n", allKeys)

// 	fmt.Println("\n=== 示例结束 ===")
// }

// // ExampleCache 缓存使用示例
// func ExampleCache() {
// 	ctx := context.Background()

// 	fmt.Println("\n=== 缓存使用场景示例 ===\n")

// 	// 模拟从数据库获取数据的函数
// 	getUserFromDB := func(userID string) map[string]string {
// 		// 模拟数据库查询
// 		time.Sleep(100 * time.Millisecond)
// 		return map[string]string{
// 			"id":    userID,
// 			"name":  "张三",
// 			"email": "zhangsan@example.com",
// 			"age":   "25",
// 		}
// 	}

// 	userID := "1001"
// 	cacheKey := "cache:user:" + userID

// 	fmt.Println("1. 尝试从缓存获取用户信息")

// 	// 先尝试从缓存获取
// 	cachedData, err := HGetAll(ctx, cacheKey)
// 	if err == nil && len(cachedData) > 0 {
// 		fmt.Printf("   ✓ 缓存命中: %v\n", cachedData)
// 	} else {
// 		fmt.Println("   ⚠ 缓存未命中，从数据库获取...")

// 		// 从数据库获取
// 		userData := getUserFromDB(userID)
// 		fmt.Printf("   ✓ 从数据库获取: %v\n", userData)

// 		// 存入缓存
// 		values := make([]interface{}, 0, len(userData)*2)
// 		for k, v := range userData {
// 			values = append(values, k, v)
// 		}
// 		HSet(ctx, cacheKey, values...)
// 		Expire(ctx, cacheKey, 5*time.Minute)
// 		fmt.Println("   ✓ 已缓存用户信息（5分钟过期）")
// 	}

// 	fmt.Println("\n2. 第二次获取（应该从缓存读取）")
// 	cachedData, err = HGetAll(ctx, cacheKey)
// 	if err == nil && len(cachedData) > 0 {
// 		fmt.Printf("   ✓ 缓存命中: %v\n", cachedData)
// 	}

// 	fmt.Println("\n=== 缓存示例结束 ===")
// }

// // ExampleRateLimiter 限流器示例
// func ExampleRateLimiter() {
// 	ctx := context.Background()

// 	fmt.Println("\n=== 限流器示例（滑动窗口）===\n")

// 	// 模拟API调用限流：每分钟最多10次请求
// 	checkRateLimit := func(userID string) bool {
// 		key := fmt.Sprintf("rate_limit:%s", userID)
// 		now := time.Now().Unix()
// 		windowStart := now - 60 // 60秒的时间窗口

// 		// 使用事务执行限流逻辑
// 		err := Watch(ctx, func(tx *goredis.Tx) error {
// 			// 移除窗口外的记录
// 			tx.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

// 			// 获取当前窗口内的请求数
// 			count, err := tx.ZCard(ctx, key).Result()
// 			if err != nil {
// 				return err
// 			}

// 			if count >= 10 {
// 				return fmt.Errorf("rate limit exceeded")
// 			}

// 			// 添加当前请求
// 			_, err = tx.ZAdd(ctx, key, goredis.Z{
// 				Score:  float64(now),
// 				Member: fmt.Sprintf("%d", now),
// 			}).Result()
// 			if err != nil {
// 				return err
// 			}

// 			// 设置过期时间
// 			tx.Expire(ctx, key, 60*time.Second)
// 			return nil
// 		}, key)

// 		return err == nil
// 	}

// 	userID := "user123"

// 	// 模拟多次API调用
// 	for i := 1; i <= 12; i++ {
// 		allowed := checkRateLimit(userID)
// 		if allowed {
// 			fmt.Printf("   ✓ 请求 %d: 允许\n", i)
// 		} else {
// 			fmt.Printf("   ❌ 请求 %d: 限流（超过限制）\n", i)
// 		}
// 		time.Sleep(50 * time.Millisecond)
// 	}

// 	fmt.Println("\n=== 限流器示例结束 ===")
// }
