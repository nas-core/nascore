package checkpsw

import (
	"math/rand"
	"time"
)

// recordFailedLogin 记录失败登录IP并计数
func recordFailedLogin(ip string, maxMap int) {
	failedLoginMutex.Lock()
	defer failedLoginMutex.Unlock()
	failedLoginAttempts[ip]++
	// 如果列表满了，随机删除一个
	if len(failedLoginAttempts) > maxMap {
		var keys []string
		for k := range failedLoginAttempts {
			keys = append(keys, k)
		}
		if len(keys) > 0 {
			randomIndex := rand.Intn(len(keys))
			delete(failedLoginAttempts, keys[randomIndex])
		}
	}
}

// getDelayTime 根据失败次数计算延迟时间
func getDelayTime(attempts int, maxTime int) time.Duration {
	delay := time.Duration(attempts) * time.Second
	delay = min(delay, time.Duration(maxTime)*time.Second)
	return delay
}

func get_sleep_time(user_ip string, maxTime int) time.Duration {
	// 获取延迟时间并延迟
	attempts := failedLoginAttempts[user_ip]
	return getDelayTime(attempts, maxTime)
}
