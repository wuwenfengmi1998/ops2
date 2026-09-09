package routers

import (
	"encoding/json"
	"ops/models"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	mrand "math/rand/v2"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 题型常量
const (
	QuizTypeSingle   = "single"   // 单选题
	QuizTypeMultiple = "multiple" // 多选题
	QuizTypeBlank    = "blank"    // 填空题
	QuizTypeJudge    = "judge"    // 判断题
)

var (
	quizUserGroup TabUserGroups
	quizAdmins    []uint
)

// quizAdminCheck 题库管理权限（quiz_admin 用户组成员或系统管理员）
func quizAdminCheck(userID uint) bool {
	if slices.Contains(quizAdmins, userID) {
		return true
	}
	if SysAdminCheck(userID) {
		return true
	}
	return false
}

// QuizUpdateAdminsCash 刷新题库管理员缓存
func QuizUpdateAdminsCash() {
	quizAdmins = nil
	quizAdmins = append(quizAdmins, 1)
	var binds []TabUserGroupBinds
	models.DB.Where("group_id = ?", quizUserGroup.ID).Find(&binds)
	for _, item := range binds {
		if !slices.Contains(quizAdmins, item.UserID) {
			quizAdmins = append(quizAdmins, item.UserID)
		}
	}
}

// TabQuizBank 题库表
type TabQuizBank struct {
	ID            uint           `gorm:"primarykey"`
	CreatorID     uint           `gorm:"not null;comment:创建人ID"`
	Name          string         `gorm:"size:200;not null;comment:题库名称"`
	Description   string         `gorm:"size:500;comment:题库描述"`
	CountSingle   int            `gorm:"default:20;comment:单选题抽题数"`
	CountMultiple int            `gorm:"default:10;comment:多选题抽题数"`
	CountJudge    int            `gorm:"default:20;comment:判断题抽题数"`
	CountBlank    int            `gorm:"default:10;comment:填空题抽题数"`
	DurationSec   int            `gorm:"comment:整卷答题时限(秒),0=不限时"`
	Active        bool           `gorm:"default:true;comment:是否启用"`
	CreatedAt     *time.Time     `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt     *time.Time     `gorm:"type:datetime;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// TabQuizQuestion 题目表
type TabQuizQuestion struct {
	ID        uint           `gorm:"primarykey"`
	BankID    uint           `gorm:"not null;index;comment:关联题库ID"`
	CreatorID uint           `gorm:"not null;comment:创建人ID"`
	Type      string         `gorm:"size:20;not null;comment:题型: single-单选 multiple-多选 blank-填空 judge-判断"`
	Title     string         `gorm:"type:text;not null;comment:题目内容"`
	Options   string         `gorm:"type:text;comment:选项JSON数组(单选/多选)"`
	Answer    string         `gorm:"type:text;not null;comment:答案JSON(single-索引 multiple-索引数组 blank-字符串数组 judge-布尔)"`
	Explain   string         `gorm:"type:text;comment:答案解析"`
	Score     int            `gorm:"default:5;comment:分值"`
	Active    bool           `gorm:"default:true;comment:是否启用"`
	CreatedAt *time.Time     `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt *time.Time     `gorm:"type:datetime;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TabQuizSession 答题会话表
type TabQuizSession struct {
	ID            uint       `gorm:"primarykey"`
	UserID        uint       `gorm:"not null;index;comment:答题用户ID"`
	BankID        uint       `gorm:"not null;index;comment:关联题库ID"`
	DurationLimit int        `gorm:"default:0;comment:时段快照(秒),0=不限时"`
	Count         int        `gorm:"not null;comment:题目数"`
	CorrectCount  int        `gorm:"default:0;comment:答对数"`
	WrongCount    int        `gorm:"default:0;comment:答错数"`
	Score         int        `gorm:"default:0;comment:得分"`
	TotalScore    int        `gorm:"default:0;comment:卷面总分"`
	DurationSec   int        `gorm:"default:0;comment:耗时(秒)"`
	Key           string     `gorm:"type:text;comment:抽题顺序配置JSON"`
	Status        int        `gorm:"default:0;comment:状态: 0-进行中 1-已完成"`
	FinishedAt    *time.Time `gorm:"type:datetime;comment:完成时间"`
	CreatedAt     *time.Time `gorm:"type:datetime;autoCreateTime"`
}

// TabQuizAnswer 答题明细表
type TabQuizAnswer struct {
	ID         uint       `gorm:"primarykey"`
	SessionID  uint       `gorm:"not null;index;comment:关联SessionID"`
	QuestionID uint       `gorm:"not null;index;comment:关联题目ID"`
	UserAnswer string     `gorm:"type:text;comment:用户答案JSON(原始索引空间)"`
	IsCorrect  bool       `gorm:"default:false;comment:是否答对"`
	GotScore   int        `gorm:"default:0;comment:得分"`
	CreatedAt  *time.Time `gorm:"type:datetime;autoCreateTime"`
}

// quizItemKey 单题的抽题配置：perm[打乱后的位置] = 原始位置
type quizItemKey struct {
	QuestionID uint  `json:"qid"`
	Perm       []int `json:"perm"`
}

// QuizPlayQuestion 答题时下发的题目(不含答案)
type QuizPlayQuestion struct {
	ID      uint     `json:"id"`
	Type    string   `json:"type"`
	Title   string   `json:"title"`
	Options []string `json:"options"`
	Score   int      `json:"score"`
	Blanks  int      `json:"blanks"` // 填空题空格数
}

// QuizReviewItem 判分后的题目回顾
type QuizReviewItem struct {
	Index             int      `json:"index"` // 1-based 顺序
	QuestionID        uint     `json:"qid"`
	Type              string   `json:"type"`
	Title             string   `json:"title"`
	Options           []string `json:"options"` // 答题时的顺序
	Score             int      `json:"score"`
	Explain           string   `json:"explain"` // 答案解析
	YourAnswerText    []string `json:"yourAnswerText"`
	CorrectAnswerText []string `json:"correctAnswerText"`
	IsCorrect         bool     `json:"isCorrect"`
	GotScore          int      `json:"gotScore"`
}

// normalizeBlank 填空判分归一化：忽略前后空格、大小写不敏感
func normalizeBlank(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func parseIndex(raw string) (int, bool) {
	var v interface{}
	if json.Unmarshal([]byte(raw), &v) != nil {
		return 0, false
	}
	switch t := v.(type) {
	case float64:
		return int(t), true
	case json.Number:
		if i, err := t.Int64(); err == nil {
			return int(i), true
		}
	case string:
		if i, err := strconv.Atoi(t); err == nil {
			return i, true
		}
	}
	return 0, false
}

func parseIntArray(raw string) ([]int, bool) {
	var arr []interface{}
	if json.Unmarshal([]byte(raw), &arr) != nil {
		return nil, false
	}
	result := []int{}
	for _, item := range arr {
		switch t := item.(type) {
		case float64:
			result = append(result, int(t))
		case string:
			if i, err := strconv.Atoi(t); err == nil {
				result = append(result, i)
			} else {
				return nil, false
			}
		default:
			return nil, false
		}
	}
	return result, true
}

// gradeBlank 填空中每空支持 "答案1|答案2" 多种合法写法，全部空答对才判对
func gradeBlank(userRaw, answerRaw string) bool {
	var ua, aa []string
	if json.Unmarshal([]byte(userRaw), &ua) != nil || json.Unmarshal([]byte(answerRaw), &aa) != nil {
		return false
	}
	if len(ua) != len(aa) {
		return false
	}
	for i := range aa {
		matched := false
		for _, v := range strings.Split(aa[i], "|") {
			if normalizeBlank(ua[i]) == normalizeBlank(v) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

// gradeQuestion 按题型宽松判分
func gradeQuestion(q TabQuizQuestion, userRaw string) bool {
	if userRaw == "" {
		return false
	}
	switch q.Type {
	case QuizTypeSingle:
		ui, ok1 := parseIndex(userRaw)
		ai, ok2 := parseIndex(q.Answer)
		return ok1 && ok2 && ui == ai
	case QuizTypeMultiple:
		us, ok1 := parseIntArray(userRaw)
		as, ok2 := parseIntArray(q.Answer)
		if !ok1 || !ok2 || len(us) == 0 || len(us) != len(as) {
			return false
		}
		sort.Ints(us)
		sort.Ints(as)
		return slices.Equal(us, as)
	case QuizTypeBlank:
		return gradeBlank(userRaw, q.Answer)
	case QuizTypeJudge:
		return strings.EqualFold(strings.TrimSpace(userRaw), strings.TrimSpace(q.Answer))
	}
	return false
}

// decodeQuizOptions 解析题目选项为字符串数组
func decodeQuizOptions(raw string) []string {
	if raw == "" || raw == "null" {
		return []string{}
	}
	var opts []string
	if json.Unmarshal([]byte(raw), &opts) != nil {
		return []string{}
	}
	return opts
}

// quizShuffleOptions 选项打乱：返回 perm[打乱后位置]=原始位置 及打乱后的选项
func quizShuffleOptions(options []string) ([]int, []string) {
	n := len(options)
	if n <= 1 {
		perm := make([]int, n)
		for i := range perm {
			perm[i] = i
		}
		return perm, options
	}
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i
	}
	mrand.Shuffle(n, func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	shuffled := make([]string, n)
	for i := range perm {
		shuffled[i] = options[perm[i]]
	}
	return perm, shuffled
}

// buildReview 由会话及答题记录构建回顾数据（顺序沿用答题时）
func buildReview(session TabQuizSession, qs map[uint]TabQuizQuestion, answers []TabQuizAnswer, keys []quizItemKey) []QuizReviewItem {
	answerByQid := make(map[uint]TabQuizAnswer, len(answers))
	for _, a := range answers {
		answerByQid[a.QuestionID] = a
	}
	// invPerm[原始位置]=打乱后的位置
	invByQid := make(map[uint]map[int]int, len(keys))
	orderByQid := make(map[uint]int, len(keys))
	for i, k := range keys {
		orderByQid[k.QuestionID] = i
		inv := make(map[int]int, len(k.Perm))
		for shuffled, orig := range k.Perm {
			inv[orig] = shuffled
		}
		invByQid[k.QuestionID] = inv
	}
	items := []QuizReviewItem{}
	for i, k := range keys {
		q, ok := qs[k.QuestionID]
		if !ok {
			continue
		}
		ans := answerByQid[k.QuestionID]
		item := QuizReviewItem{
			Index:      i + 1,
			QuestionID: q.ID,
			Type:       q.Type,
			Title:      q.Title,
			Options:    []string{},
			Score:      q.Score,
			Explain:    q.Explain,
			IsCorrect:  ans.IsCorrect,
			GotScore:   ans.GotScore,
		}
		// 选项按答题时的顺序还原
		origOptions := decodeQuizOptions(q.Options)
		for _, orig := range k.Perm {
			if orig >= 0 && orig < len(origOptions) {
				item.Options = append(item.Options, origOptions[orig])
			}
		}
		inv := invByQid[k.QuestionID]
		switch q.Type {
		case QuizTypeSingle, QuizTypeMultiple:
			// 用户答案与正确答案的索引映射回"打乱后空间"展示
			if q.Type == QuizTypeSingle {
				if idx, ok := parseIndex(ans.UserAnswer); ok {
					if shuffledIdx, exists := inv[idx]; exists && shuffledIdx < len(item.Options) {
						item.YourAnswerText = append(item.YourAnswerText, item.Options[shuffledIdx])
					}
				}
				if ai, ok := parseIndex(q.Answer); ok {
					if shuffledIdx, exists := inv[ai]; exists && shuffledIdx < len(item.Options) {
						item.CorrectAnswerText = append(item.CorrectAnswerText, item.Options[shuffledIdx])
					}
				}
			} else {
				if uIdxs, ok := parseIntArray(ans.UserAnswer); ok {
					for _, origIdx := range uIdxs {
						if shuffledIdx, exists := inv[origIdx]; exists && shuffledIdx < len(item.Options) {
							item.YourAnswerText = append(item.YourAnswerText, item.Options[shuffledIdx])
						}
					}
				}
				if correctIdxs, ok := parseIntArray(q.Answer); ok {
					for _, origIdx := range correctIdxs {
						if shuffledIdx, exists := inv[origIdx]; exists && shuffledIdx < len(item.Options) {
							item.CorrectAnswerText = append(item.CorrectAnswerText, item.Options[shuffledIdx])
						}
					}
				}
			}
		case QuizTypeBlank:
			var ua []string
			if json.Unmarshal([]byte(ans.UserAnswer), &ua) == nil {
				item.YourAnswerText = ua
			}
			var aa []string
			if json.Unmarshal([]byte(q.Answer), &aa) == nil {
				for _, one := range aa {
					item.CorrectAnswerText = append(item.CorrectAnswerText, strings.Split(one, "|")[0])
				}
			}
		case QuizTypeJudge:
			if strings.TrimSpace(ans.UserAnswer) != "" {
				item.YourAnswerText = []string{strings.TrimSpace(ans.UserAnswer)}
			}
			item.CorrectAnswerText = []string{strings.TrimSpace(q.Answer)}
		}
		items = append(items, item)
	}
	return items
}

// quizUserMaps 批量获取用户名称与头像映射
func quizUserMaps(userIDs []uint) (map[uint]string, map[uint]string) {
	nameMap := map[uint]string{}
	avatarMap := map[uint]string{}
	ids := []uint{}
	seen := map[uint]bool{}
	for _, id := range userIDs {
		if id > 0 && !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nameMap, avatarMap
	}
	var users []TabUser
	models.DB.Where("id IN ?", ids).Find(&users)
	for _, u := range users {
		nameMap[u.ID] = u.Name
	}
	var infos []TabUserInfo
	models.DB.Where("user_id IN ?", ids).Find(&infos)
	for _, info := range infos {
		avatarMap[info.UserID] = info.AvatarPath
	}
	return nameMap, avatarMap
}

// quizCreatorMaps 批量获取题库创建者的名称与头像映射
func quizCreatorMaps(banks []TabQuizBank) (map[uint]string, map[uint]string) {
	ids := []uint{}
	for _, b := range banks {
		ids = append(ids, b.CreatorID)
	}
	return quizUserMaps(ids)
}

// quizBankNameMap 批量获取题库名称映射
func quizBankNameMap(sessions []TabQuizSession) map[uint]string {	ids := []uint{}
	seen := map[uint]bool{}
	for _, s := range sessions {
		if s.BankID > 0 && !seen[s.BankID] {
			seen[s.BankID] = true
			ids = append(ids, s.BankID)
		}
	}
	m := map[uint]string{}
	if len(ids) == 0 {
		return m
	}
	var banks []TabQuizBank
	models.DB.Where("id IN ?", ids).Find(&banks)
	for _, b := range banks {
		m[b.ID] = b.Name
	}
	return m
}

// quizDefaultCount 各题型默认抽题数
func quizDefaultCount(qType string) int {
	switch qType {
	case QuizTypeSingle:
		return 20
	case QuizTypeMultiple:
		return 10
	case QuizTypeJudge:
		return 20
	case QuizTypeBlank:
		return 10
	}
	return 10
}

// ApiQuizInit 初始化答题模块
func ApiQuizInit() {
	models.DB.AutoMigrate(&TabQuizBank{})
	models.DB.AutoMigrate(&TabQuizQuestion{})
	models.DB.AutoMigrate(&TabQuizSession{})
	models.DB.AutoMigrate(&TabQuizAnswer{})

	// 自动创建 quiz_admin 用户组
	quizUserGroup.Name = "quiz_admin"
	if models.DB.Where(&quizUserGroup).First(&quizUserGroup).Error == nil {
		QuizUpdateAdminsCash()
	} else {
		quizUserGroup.Type = "usergroup"
		models.DB.Create(&quizUserGroup)
	}

	// 兼容旧库：未设置过题型抽题数的题库补默认值(单选20 多选10 判断20 填空10)
	models.DB.Model(&TabQuizBank{}).
		Where("count_single = 0 AND count_multiple = 0 AND count_judge = 0 AND count_blank = 0").
		Updates(map[string]interface{}{
			"count_single":   20,
			"count_multiple": 10,
			"count_judge":    20,
			"count_blank":    10,
		})
}

// ApiQuiz 题库路由
func ApiQuiz(r *gin.RouterGroup) {
	// ── 题库管理（quiz_admin） ──
	r.POST("/bank/add", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		if !quizAdminCheck(user.ID) {
			ReturnJson(ctx, "quiz_permission_denied", nil)
			return
		}
		type FromAdd struct {
			Name          string `json:"name"`
			Description   string `json:"description"`
			CountSingle   int    `json:"countSingle"`
			CountMultiple int    `json:"countMultiple"`
			CountJudge    int    `json:"countJudge"`
			CountBlank    int    `json:"countBlank"`
			DurationSec   int    `json:"durationSec"`
		}
		var from FromAdd
		if err := decodeJSON(data, &from); err != nil || strings.TrimSpace(from.Name) == "" {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		if from.CountSingle <= 0 {
			from.CountSingle = quizDefaultCount(QuizTypeSingle)
		}
		if from.CountMultiple <= 0 {
			from.CountMultiple = quizDefaultCount(QuizTypeMultiple)
		}
		if from.CountJudge <= 0 {
			from.CountJudge = quizDefaultCount(QuizTypeJudge)
		}
		if from.CountBlank <= 0 {
			from.CountBlank = quizDefaultCount(QuizTypeBlank)
		}
		if from.DurationSec < 0 {
			from.DurationSec = 0
		}
		bank := TabQuizBank{
			CreatorID:     user.ID,
			Name:          strings.TrimSpace(from.Name),
			Description:   strings.TrimSpace(from.Description),
			CountSingle:   from.CountSingle,
			CountMultiple: from.CountMultiple,
			CountJudge:    from.CountJudge,
			CountBlank:    from.CountBlank,
			DurationSec:   from.DurationSec,
			Active:        true,
		}
		if err := models.DB.Create(&bank).Error; err != nil {
			ReturnJson(ctx, "dbErr", nil)
			return
		}
		ReturnJson(ctx, "apiOK", gin.H{"id": bank.ID})
	})

	r.POST("/bank/update", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		if !quizAdminCheck(user.ID) {
			ReturnJson(ctx, "quiz_permission_denied", nil)
			return
		}
		type FromUpdate struct {
			ID            uint   `json:"id"`
			Name          string `json:"name"`
			Description   string `json:"description"`
			CountSingle   int    `json:"countSingle"`
			CountMultiple int    `json:"countMultiple"`
			CountJudge    int    `json:"countJudge"`
			CountBlank    int    `json:"countBlank"`
			DurationSec   int    `json:"durationSec"`
			Active        *bool  `json:"active"`
		}
		var from FromUpdate
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		var bank TabQuizBank
		if err := models.DB.First(&bank, from.ID).Error; err != nil {
			ReturnJson(ctx, "quiz_bank_not_found", nil)
			return
		}
		if strings.TrimSpace(from.Name) != "" {
			bank.Name = strings.TrimSpace(from.Name)
		}
		if from.Description != "" {
			bank.Description = from.Description
		}
		if from.CountSingle > 0 {
			bank.CountSingle = from.CountSingle
		}
		if from.CountMultiple > 0 {
			bank.CountMultiple = from.CountMultiple
		}
		if from.CountJudge > 0 {
			bank.CountJudge = from.CountJudge
		}
		if from.CountBlank > 0 {
			bank.CountBlank = from.CountBlank
		}
		if from.DurationSec >= 0 {
			bank.DurationSec = from.DurationSec
		}
		if from.Active != nil {
			bank.Active = *from.Active
		}
		if err := models.DB.Save(&bank).Error; err != nil {
			ReturnJson(ctx, "dbErr", nil)
			return
		}
		ReturnJson(ctx, "apiOK", nil)
	})

	r.POST("/bank/delete", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		if !quizAdminCheck(user.ID) {
			ReturnJson(ctx, "quiz_permission_denied", nil)
			return
		}
		type FromDelete struct {
			ID uint `json:"id"`
		}
		var from FromDelete
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		var bank TabQuizBank
		if err := models.DB.First(&bank, from.ID).Error; err != nil {
			ReturnJson(ctx, "quiz_bank_not_found", nil)
			return
		}
		// 级联软删除题库下的题目
		models.DB.Delete(&TabQuizQuestion{}, "bank_id = ?", bank.ID)
		if err := models.DB.Delete(&bank).Error; err != nil {
			ReturnJson(ctx, "dbErr", nil)
			return
		}
		ReturnJson(ctx, "apiOK", nil)
	})

	r.POST("/banks/list", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		if !quizAdminCheck(user.ID) {
			ReturnJson(ctx, "quiz_permission_denied", nil)
			return
		}
		type FromList struct {
			Keyword  string `json:"keyword"`
			Page     int    `json:"page"`
			PageSize int    `json:"pageSize"`
		}
		var from FromList
		decodeJSON(data, &from)
		if from.Page <= 0 {
			from.Page = 1
		}
		if from.PageSize <= 0 || from.PageSize > 100 {
			from.PageSize = 10
		}
		query := models.DB.Model(&TabQuizBank{})
		if strings.TrimSpace(from.Keyword) != "" {
			query = query.Where("name LIKE ?", "%"+strings.TrimSpace(from.Keyword)+"%")
		}
		var total int64
		query.Count(&total)
		var banks []TabQuizBank
		query.Order("id DESC").Offset((from.Page - 1) * from.PageSize).Limit(from.PageSize).Find(&banks)
		type countRow struct {
			BankID uint `gorm:"column:bank_id"`
			N      int  `gorm:"column:n"`
		}
		var counts []countRow
		models.DB.Model(&TabQuizQuestion{}).
			Select("bank_id, COUNT(*) AS n").
			Group("bank_id").
			Scan(&counts)
		countMap := map[uint]int{}
		for _, c := range counts {
			countMap[c.BankID] = c.N
		}
		items := []gin.H{}
		creatorNames, creatorAvatars := quizCreatorMaps(banks)
		for _, b := range banks {
			items = append(items, gin.H{
				"id":            b.ID,
				"name":          b.Name,
				"description":   b.Description,
				"countSingle":   b.CountSingle,
				"countMultiple": b.CountMultiple,
				"countJudge":    b.CountJudge,
				"countBlank":    b.CountBlank,
				"durationSec":   b.DurationSec,
				"active":        b.Active,
				"questionCnt":   countMap[b.ID],
				"createdAt":     b.CreatedAt,
				"creatorName":   creatorNames[b.CreatorID],
				"creatorAvatar": creatorAvatars[b.CreatorID],
			})
		}
		ReturnJson(ctx, "apiOK", gin.H{
			"items":    items,
			"total":    total,
			"page":     from.Page,
			"pageSize": from.PageSize,
		})
	})

	// ── 用户可用题库 ──
	r.POST("/banks/available", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		type FromList struct {
			Keyword string `json:"keyword"`
		}
		var from FromList
		decodeJSON(data, &from)
		query := models.DB.Model(&TabQuizBank{}).Where("active = ?", true)
		if strings.TrimSpace(from.Keyword) != "" {
			query = query.Where("name LIKE ?", "%"+strings.TrimSpace(from.Keyword)+"%")
		}
		var banks []TabQuizBank
		query.Order("id ASC").Find(&banks)
		type countRow struct {
			BankID uint `gorm:"column:bank_id"`
			N      int  `gorm:"column:n"`
		}
		type scoreRow struct {
			BankID uint `gorm:"column:bank_id"`
			Best   int  `gorm:"column:best"`
			Times  int  `gorm:"column:times"`
		}
		var counts []countRow
		models.DB.Model(&TabQuizQuestion{}).
			Select("bank_id, COUNT(*) AS n").
			Where("active = ?", true).
			Group("bank_id").
			Scan(&counts)
		countMap := map[uint]int{}
		for _, c := range counts {
			countMap[c.BankID] = c.N
		}
		var scores []scoreRow
		models.DB.Model(&TabQuizSession{}).
			Select("bank_id, MAX(score) AS best, COUNT(*) AS times").
			Where("user_id = ? AND status = 1", user.ID).
			Group("bank_id").
			Scan(&scores)
		bestMap := map[uint]int{}
		timesMap := map[uint]int{}
		for _, s := range scores {
			bestMap[s.BankID] = s.Best
			timesMap[s.BankID] = s.Times
		}
		items := []gin.H{}
		creatorNames, creatorAvatars := quizCreatorMaps(banks)
		for _, b := range banks {
			items = append(items, gin.H{
				"id":            b.ID,
				"name":          b.Name,
				"description":   b.Description,
				"countSingle":   b.CountSingle,
				"countMultiple": b.CountMultiple,
				"countJudge":    b.CountJudge,
				"countBlank":    b.CountBlank,
				"durationSec":   b.DurationSec,
				"questionCnt":   countMap[b.ID],
				"myBest":        bestMap[b.ID],
				"myTimes":       timesMap[b.ID],
				"creatorName":   creatorNames[b.CreatorID],
				"creatorAvatar": creatorAvatars[b.CreatorID],
			})
		}
		ReturnJson(ctx, "apiOK", gin.H{"items": items})
	})

	// ── 学习模式（登录用户）：分页返回题目含答案与解析 ──
	r.POST("/bank/learn", func(ctx *gin.Context) {
		isAuth, _, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		type FromLearn struct {
			BankID   uint `json:"bankId"`
			Page     int  `json:"page"`
			PageSize int  `json:"pageSize"`
		}
		var from FromLearn
		if err := decodeJSON(data, &from); err != nil || from.BankID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		var bank TabQuizBank
		if err := models.DB.First(&bank, from.BankID).Error; err != nil || !bank.Active {
			ReturnJson(ctx, "quiz_bank_not_found", nil)
			return
		}
		if from.Page <= 0 {
			from.Page = 1
		}
		if from.PageSize <= 0 || from.PageSize > 20 {
			from.PageSize = 10
		}
		query := models.DB.Model(&TabQuizQuestion{}).Where("bank_id = ? AND active = ?", bank.ID, true)
		var total int64
		query.Count(&total)
		var questions []TabQuizQuestion
		query.Order("id ASC").Offset((from.Page - 1) * from.PageSize).Limit(from.PageSize).Find(&questions)

		creatorIDs := []uint{}
		for _, q := range questions {
			creatorIDs = append(creatorIDs, q.CreatorID)
		}
		nameMap, avatarMap := quizUserMaps(creatorIDs)

		items := []gin.H{}
		for _, q := range questions {
			items = append(items, gin.H{
				"id":            q.ID,
				"type":          q.Type,
				"title":         q.Title,
				"options":       decodeQuizOptions(q.Options),
				"answer":        json.RawMessage(q.Answer),
				"explain":       q.Explain,
				"score":         q.Score,
				"creatorName":   nameMap[q.CreatorID],
				"creatorAvatar": avatarMap[q.CreatorID],
			})
		}
		ReturnJson(ctx, "apiOK", gin.H{
			"items":    items,
			"total":    total,
			"page":     from.Page,
			"pageSize": from.PageSize,
			"bankName": bank.Name,
		})
	})

	// ── 题目管理（quiz_admin） ──
	r.POST("/question/add", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		if !quizAdminCheck(user.ID) {
			ReturnJson(ctx, "quiz_permission_denied", nil)
			return
		}
		type FromAdd struct {
			BankID  uint            `json:"bankId"`
			Type    string          `json:"type"`
			Title   string          `json:"title"`
			Options []string        `json:"options"`
			Answer  json.RawMessage `json:"answer"`
			Explain string          `json:"explain"`
			Score   int             `json:"score"`
		}
		var from FromAdd
		if err := decodeJSON(data, &from); err != nil || from.BankID == 0 || from.Type == "" || strings.TrimSpace(from.Title) == "" {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		var bank TabQuizBank
		if err := models.DB.First(&bank, from.BankID).Error; err != nil {
			ReturnJson(ctx, "quiz_bank_not_found", nil)
			return
		}
		if from.Score <= 0 {
			from.Score = 5
		}
		if (from.Type == QuizTypeSingle || from.Type == QuizTypeMultiple) && len(from.Options) < 2 {
			ReturnJson(ctx, "quiz_answer_format_err", nil)
			return
		}
		answerBytes, err := json.Marshal(from.Answer)
		if err != nil {
			ReturnJson(ctx, "quiz_answer_format_err", nil)
			return
		}
		question := TabQuizQuestion{
			BankID:    from.BankID,
			CreatorID: user.ID,
			Type:      from.Type,
			Title:     from.Title,
			Explain:   strings.TrimSpace(from.Explain),
			Score:     from.Score,
			Active:    true,
		}
		if len(from.Options) > 0 {
			optsBytes, _ := json.Marshal(from.Options)
			question.Options = string(optsBytes)
		}
		question.Answer = string(answerBytes)
		if err := models.DB.Create(&question).Error; err != nil {
			ReturnJson(ctx, "dbErr", nil)
			return
		}
		ReturnJson(ctx, "apiOK", gin.H{"id": question.ID})
	})

	r.POST("/question/update", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		if !quizAdminCheck(user.ID) {
			ReturnJson(ctx, "quiz_permission_denied", nil)
			return
		}
		type FromUpdate struct {
			ID      uint            `json:"id"`
			BankID  uint            `json:"bankId"`
			Type    string          `json:"type"`
			Title   string          `json:"title"`
			Options []string        `json:"options"`
			Answer  json.RawMessage `json:"answer"`
			Explain *string         `json:"explain"`
			Score   int             `json:"score"`
			Active  *bool           `json:"active"`
		}
		var from FromUpdate
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		var question TabQuizQuestion
		if err := models.DB.First(&question, from.ID).Error; err != nil {
			ReturnJson(ctx, "quiz_question_not_found", nil)
			return
		}
		if from.BankID > 0 {
			question.BankID = from.BankID
		}
		if from.Type != "" {
			question.Type = from.Type
		}
		if strings.TrimSpace(from.Title) != "" {
			question.Title = from.Title
		}
		if from.Score > 0 {
			question.Score = from.Score
		}
		if len(from.Options) > 0 {
			optsBytes, _ := json.Marshal(from.Options)
			question.Options = string(optsBytes)
		}
		if len(from.Answer) > 0 && string(from.Answer) != "null" {
			answerBytes, _ := json.Marshal(from.Answer)
			question.Answer = string(answerBytes)
		}
		if from.Explain != nil {
			question.Explain = strings.TrimSpace(*from.Explain)
		}
		if from.Active != nil {
			question.Active = *from.Active
		}
		if err := models.DB.Save(&question).Error; err != nil {
			ReturnJson(ctx, "dbErr", nil)
			return
		}
		ReturnJson(ctx, "apiOK", nil)
	})

	r.POST("/question/delete", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		if !quizAdminCheck(user.ID) {
			ReturnJson(ctx, "quiz_permission_denied", nil)
			return
		}
		type FromDelete struct {
			ID uint `json:"id"`
		}
		var from FromDelete
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		if err := models.DB.Delete(&TabQuizQuestion{}, from.ID).Error; err != nil {
			ReturnJson(ctx, "dbErr", nil)
			return
		}
		ReturnJson(ctx, "apiOK", nil)
	})

	r.POST("/questions/list", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		if !quizAdminCheck(user.ID) {
			ReturnJson(ctx, "quiz_permission_denied", nil)
			return
		}
		type FromList struct {
			BankID   uint   `json:"bankId"`
			Type     string `json:"type"`
			Keyword  string `json:"keyword"`
			Page     int    `json:"page"`
			PageSize int    `json:"pageSize"`
		}
		var from FromList
		decodeJSON(data, &from)
		if from.Page <= 0 {
			from.Page = 1
		}
		if from.PageSize <= 0 || from.PageSize > 100 {
			from.PageSize = 10
		}
		query := models.DB.Model(&TabQuizQuestion{})
		if from.BankID > 0 {
			query = query.Where("bank_id = ?", from.BankID)
		}
		if from.Type != "" {
			query = query.Where("type = ?", from.Type)
		}
		if strings.TrimSpace(from.Keyword) != "" {
			query = query.Where("title LIKE ?", "%"+strings.TrimSpace(from.Keyword)+"%")
		}
		var total int64
		query.Count(&total)
		var questions []TabQuizQuestion
		query.Order("id DESC").Offset((from.Page - 1) * from.PageSize).Limit(from.PageSize).Find(&questions)

		creatorIDs := []uint{}
		for _, q := range questions {
			creatorIDs = append(creatorIDs, q.CreatorID)
		}
		nameMap := map[uint]string{}
		if len(creatorIDs) > 0 {
			var users []TabUser
			models.DB.Where("id IN ?", creatorIDs).Find(&users)
			for _, u := range users {
				nameMap[u.ID] = u.Name
			}
		}
		items := []gin.H{}
		for _, q := range questions {
			items = append(items, gin.H{
				"id":          q.ID,
				"bankId":      q.BankID,
				"type":        q.Type,
				"title":       q.Title,
				"options":     decodeQuizOptions(q.Options),
				"answer":      json.RawMessage(q.Answer),
				"explain":     q.Explain,
				"score":       q.Score,
				"active":      q.Active,
				"creatorName": nameMap[q.CreatorID],
				"createdAt":   q.CreatedAt,
			})
		}
		ReturnJson(ctx, "apiOK", gin.H{
			"items":    items,
			"total":    total,
			"page":     from.Page,
			"pageSize": from.PageSize,
		})
	})

	// ── 答题（登录用户） ──
	r.POST("/start", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		type FromStart struct {
			BankID uint `json:"bankId"`
		}
		var from FromStart
		if err := decodeJSON(data, &from); err != nil || from.BankID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		var bank TabQuizBank
		if err := models.DB.First(&bank, from.BankID).Error; err != nil || !bank.Active {
			ReturnJson(ctx, "quiz_bank_not_found", nil)
			return
		}
		// 按题型分别随机抽取，不足则该题型全部抽取
		typeDraws := []struct {
			qType string
			count int
		}{
			{QuizTypeSingle, bank.CountSingle},
			{QuizTypeMultiple, bank.CountMultiple},
			{QuizTypeJudge, bank.CountJudge},
			{QuizTypeBlank, bank.CountBlank},
		}
		questions := []TabQuizQuestion{}
		for _, d := range typeDraws {
			if d.count <= 0 {
				d.count = quizDefaultCount(d.qType)
			}
			var pool []TabQuizQuestion
			models.DB.Where("bank_id = ? AND active = ? AND type = ?", bank.ID, true, d.qType).Find(&pool)
			if len(pool) == 0 {
				continue
			}
			mrand.Shuffle(len(pool), func(i, j int) {
				pool[i], pool[j] = pool[j], pool[i]
			})
			n := d.count
			if n > len(pool) {
				n = len(pool)
			}
			questions = append(questions, pool[:n]...)
		}
		if len(questions) == 0 {
			ReturnJson(ctx, "quiz_no_questions", nil)
			return
		}

		keys := []quizItemKey{}
		playQues := []QuizPlayQuestion{}
		totalScore := 0
		for _, q := range questions {
			origOptions := decodeQuizOptions(q.Options)
			var perm []int
			shuffled := origOptions
			if q.Type == QuizTypeSingle || q.Type == QuizTypeMultiple {
				perm, shuffled = quizShuffleOptions(origOptions)
			} else {
				for i := range origOptions {
					perm = append(perm, i)
				}
				if len(origOptions) == 0 {
					perm = []int{}
				}
			}
			blankCount := 0
			if q.Type == QuizTypeBlank {
				var aa []string
				if json.Unmarshal([]byte(q.Answer), &aa) == nil {
					blankCount = len(aa)
				}
			}
			keys = append(keys, quizItemKey{QuestionID: q.ID, Perm: perm})
			playQues = append(playQues, QuizPlayQuestion{
				ID:      q.ID,
				Type:    q.Type,
				Title:   q.Title,
				Options: shuffled,
				Score:   q.Score,
				Blanks:  blankCount,
			})
			totalScore += q.Score
		}
		keyBytes, _ := json.Marshal(keys)

		// 清理该用户未完成的会话
		models.DB.Where("user_id = ? AND status = 0", user.ID).Delete(&TabQuizSession{})

		session := TabQuizSession{
			UserID:        user.ID,
			BankID:        bank.ID,
			DurationLimit: bank.DurationSec,
			Count:         len(questions),
			TotalScore:    totalScore,
			Key:           string(keyBytes),
			Status:        0,
		}
		if err := models.DB.Create(&session).Error; err != nil {
			ReturnJson(ctx, "dbErr", nil)
			return
		}
		ReturnJson(ctx, "apiOK", gin.H{
			"sessionId":     session.ID,
			"bankName":      bank.Name,
			"questions":     playQues,
			"totalScore":    totalScore,
			"durationLimit": bank.DurationSec,
		})
	})

	r.POST("/submit", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		type FromAnswer struct {
			QuestionID uint            `json:"qid"`
			UserAnswer json.RawMessage `json:"answer"`
		}
		type FromSubmit struct {
			SessionID   uint         `json:"sessionId"`
			DurationSec int          `json:"durationSec"`
			Answers     []FromAnswer `json:"answers"`
		}
		var from FromSubmit
		if err := decodeJSON(data, &from); err != nil || from.SessionID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		var session TabQuizSession
		if err := models.DB.First(&session, from.SessionID).Error; err != nil || session.UserID != user.ID {
			ReturnJson(ctx, "quiz_session_not_found", nil)
			return
		}
		if session.Status != 0 {
			ReturnJson(ctx, "quiz_session_finished", nil)
			return
		}
		var keys []quizItemKey
		if err := json.Unmarshal([]byte(session.Key), &keys); err != nil {
			ReturnJson(ctx, "quiz_session_not_found", nil)
			return
		}
		userAnswerByQid := map[uint]json.RawMessage{}
		for _, a := range from.Answers {
			if a.QuestionID != 0 {
				userAnswerByQid[a.QuestionID] = a.UserAnswer
			}
		}
		// 取题目
		qids := []uint{}
		for _, k := range keys {
			qids = append(qids, k.QuestionID)
		}
		qMap := map[uint]TabQuizQuestion{}
		var qs []TabQuizQuestion
		models.DB.Where("id IN ?", qids).Find(&qs)
		for _, q := range qs {
			qMap[q.ID] = q
		}

		correctCount := 0
		wrongCount := 0
		score := 0
		answers := []TabQuizAnswer{}
		for _, k := range keys {
			q, ok := qMap[k.QuestionID]
			if !ok {
				continue
			}
			raw, has := userAnswerByQid[k.QuestionID]
			normalized := ""
			if has && len(raw) > 0 && string(raw) != "null" {
				normalized = string(raw)
				// 选项类题目把用户答案从"打乱后空间"映射回"原始空间"
				if q.Type == QuizTypeSingle || q.Type == QuizTypeMultiple {
					mapped := true
					if q.Type == QuizTypeSingle {
						if idx, ok2 := parseIndex(normalized); ok2 {
							if idx < 0 || idx >= len(k.Perm) {
								mapped = false
							} else {
								normalized = strconv.Itoa(k.Perm[idx])
							}
						} else {
							mapped = false
						}
					} else {
						if idxArr, ok2 := parseIntArray(normalized); ok2 {
							newArr := []int{}
							for _, idx := range idxArr {
								if idx < 0 || idx >= len(k.Perm) {
									mapped = false
									break
								}
								newArr = append(newArr, k.Perm[idx])
							}
							if mapped {
								out, _ := json.Marshal(newArr)
								normalized = string(out)
							}
						} else {
							mapped = false
						}
					}
					if !mapped {
						normalized = ""
					}
				}
			}
			isCorrect := gradeQuestion(q, normalized)
			gotScore := 0
			if isCorrect {
				gotScore = q.Score
			}
			if isCorrect {
				correctCount++
			} else {
				wrongCount++
			}
			score += gotScore
			answers = append(answers, TabQuizAnswer{
				SessionID:  session.ID,
				QuestionID: q.ID,
				UserAnswer: normalized,
				IsCorrect:  isCorrect,
				GotScore:   gotScore,
			})
		}
		if len(keys) > 0 && (correctCount+wrongCount > 0) {
			models.DB.Create(&answers)
		}
		now := time.Now()
		session.CorrectCount = correctCount
		session.WrongCount = wrongCount
		session.Score = score
		if from.DurationSec > 0 {
			session.DurationSec = from.DurationSec
		}
		session.Status = 1
		session.FinishedAt = &now
		models.DB.Save(&session)

		ReturnJson(ctx, "apiOK", gin.H{
			"session": gin.H{
				"id":            session.ID,
				"bankId":        session.BankID,
				"durationLimit": session.DurationLimit,
				"count":         session.Count,
				"correctCount":  session.CorrectCount,
				"wrongCount":    session.WrongCount,
				"score":         session.Score,
				"totalScore":    session.TotalScore,
				"durationSec":   session.DurationSec,
			},
			"review": buildReview(session, qMap, answers, keys),
		})
	})

	r.POST("/sessions", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		type FromList struct {
			Page     int `json:"page"`
			PageSize int `json:"pageSize"`
		}
		var from FromList
		decodeJSON(data, &from)
		if from.Page <= 0 {
			from.Page = 1
		}
		if from.PageSize <= 0 || from.PageSize > 100 {
			from.PageSize = 10
		}
		query := models.DB.Model(&TabQuizSession{}).Where("user_id = ? AND status = 1", user.ID)
		var total int64
		query.Count(&total)
		var sessions []TabQuizSession
		query.Order("id DESC").Offset((from.Page - 1) * from.PageSize).Limit(from.PageSize).Find(&sessions)
		bankNames := quizBankNameMap(sessions)
		items := []gin.H{}
		for _, s := range sessions {
			items = append(items, gin.H{
				"id":            s.ID,
				"bankId":        s.BankID,
				"bankName":      bankNames[s.BankID],
				"count":         s.Count,
				"correctCount":  s.CorrectCount,
				"wrongCount":    s.WrongCount,
				"score":         s.Score,
				"totalScore":    s.TotalScore,
				"durationSec":   s.DurationSec,
				"durationLimit": s.DurationLimit,
				"createdAt":     s.CreatedAt,
			})
		}
		ReturnJson(ctx, "apiOK", gin.H{
			"items":    items,
			"total":    total,
			"page":     from.Page,
			"pageSize": from.PageSize,
		})
	})

	r.POST("/session", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		type FromGet struct {
			ID uint `json:"id"`
		}
		var from FromGet
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		var session TabQuizSession
		if err := models.DB.First(&session, from.ID).Error; err != nil || session.UserID != user.ID {
			ReturnJson(ctx, "quiz_session_not_found", nil)
			return
		}
		var keys []quizItemKey
		if err := json.Unmarshal([]byte(session.Key), &keys); err != nil {
			ReturnJson(ctx, "quiz_session_not_found", nil)
			return
		}
		qids := []uint{}
		for _, k := range keys {
			qids = append(qids, k.QuestionID)
		}
		qMap := map[uint]TabQuizQuestion{}
		var qs []TabQuizQuestion
		models.DB.Where("id IN ?", qids).Find(&qs)
		for _, q := range qs {
			qMap[q.ID] = q
		}
		var answers []TabQuizAnswer
		models.DB.Where("session_id = ?", session.ID).Find(&answers)
		bankNames := quizBankNameMap([]TabQuizSession{session})
		ReturnJson(ctx, "apiOK", gin.H{
			"session": gin.H{
				"id":            session.ID,
				"bankId":        session.BankID,
				"bankName":      bankNames[session.BankID],
				"durationLimit": session.DurationLimit,
				"count":         session.Count,
				"correctCount":  session.CorrectCount,
				"wrongCount":    session.WrongCount,
				"score":         session.Score,
				"totalScore":    session.TotalScore,
				"durationSec":   session.DurationSec,
			},
			"review": buildReview(session, qMap, answers, keys),
		})
	})

	// 删除本人的成绩记录（连带答题明细）
	r.POST("/session/delete", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		type FromDelete struct {
			ID uint `json:"id"`
		}
		var from FromDelete
		if err := decodeJSON(data, &from); err != nil || from.ID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		var session TabQuizSession
		if err := models.DB.First(&session, from.ID).Error; err != nil || session.UserID != user.ID {
			ReturnJson(ctx, "quiz_session_not_found", nil)
			return
		}
		models.DB.Delete(&TabQuizAnswer{}, "session_id = ?", session.ID)
		models.DB.Delete(&TabQuizSession{}, session.ID)
		ReturnJson(ctx, "apiOK", nil)
	})

	// 错题重做：仅判分，不写入任何数据
	r.POST("/redo", func(ctx *gin.Context) {
		isAuth, user, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		type FromAnswer struct {
			QuestionID uint            `json:"qid"`
			UserAnswer json.RawMessage `json:"answer"`
		}
		type FromRedo struct {
			SessionID   uint         `json:"sessionId"`
			DurationSec int          `json:"durationSec"`
			Answers     []FromAnswer `json:"answers"`
		}
		var from FromRedo
		if err := decodeJSON(data, &from); err != nil || from.SessionID == 0 {
			ReturnJson(ctx, "jsonErr", nil)
			return
		}
		var session TabQuizSession
		if err := models.DB.First(&session, from.SessionID).Error; err != nil || session.UserID != user.ID {
			ReturnJson(ctx, "quiz_session_not_found", nil)
			return
		}
		if session.Status != 1 {
			ReturnJson(ctx, "quiz_session_finished", nil)
			return
		}
		var keys []quizItemKey
		if err := json.Unmarshal([]byte(session.Key), &keys); err != nil {
			ReturnJson(ctx, "quiz_session_not_found", nil)
			return
		}
		// 原答卷错题集
		var wrongSet []TabQuizAnswer
		models.DB.Where("session_id = ? AND is_correct = ?", session.ID, false).Find(&wrongSet)
		wrongQids := map[uint]bool{}
		for _, w := range wrongSet {
			wrongQids[w.QuestionID] = true
		}
		userAnswerByQid := map[uint]json.RawMessage{}
		for _, a := range from.Answers {
			if a.QuestionID != 0 {
				userAnswerByQid[a.QuestionID] = a.UserAnswer
			}
		}
		// 仅用于判分的 qid 集合
		qids := []uint{}
		wrongKeys := []quizItemKey{}
		for _, k := range keys {
			if wrongQids[k.QuestionID] {
				qids = append(qids, k.QuestionID)
				wrongKeys = append(wrongKeys, k)
			}
		}
		qMap := map[uint]TabQuizQuestion{}
		var qs []TabQuizQuestion
		models.DB.Where("id IN ?", qids).Find(&qs)
		for _, q := range qs {
			qMap[q.ID] = q
		}

		correctCount := 0
		wrongCount := 0
		score := 0
		totalScore := 0
		answers := []TabQuizAnswer{}
		for _, k := range wrongKeys {
			q, ok := qMap[k.QuestionID]
			if !ok {
				continue
			}
			totalScore += q.Score
			raw, has := userAnswerByQid[k.QuestionID]
			normalized := ""
			if has && len(raw) > 0 && string(raw) != "null" {
				normalized = string(raw)
				if q.Type == QuizTypeSingle || q.Type == QuizTypeMultiple {
					mapped := true
					if q.Type == QuizTypeSingle {
						if idx, ok2 := parseIndex(normalized); ok2 {
							if idx < 0 || idx >= len(k.Perm) {
								mapped = false
							} else {
								normalized = strconv.Itoa(k.Perm[idx])
							}
						} else {
							mapped = false
						}
					} else {
						if idxArr, ok2 := parseIntArray(normalized); ok2 {
							newArr := []int{}
							for _, idx := range idxArr {
								if idx < 0 || idx >= len(k.Perm) {
									mapped = false
									break
								}
								newArr = append(newArr, k.Perm[idx])
							}
							if mapped {
								out, _ := json.Marshal(newArr)
								normalized = string(out)
							}
						} else {
							mapped = false
						}
					}
					if !mapped {
						normalized = ""
					}
				}
			}
			isCorrect := gradeQuestion(q, normalized)
			gotScore := 0
			if isCorrect {
				gotScore = q.Score
			}
			if isCorrect {
				correctCount++
			} else {
				wrongCount++
			}
			score += gotScore
			answers = append(answers, TabQuizAnswer{
				SessionID:  session.ID,
				QuestionID: q.ID,
				UserAnswer: normalized,
				IsCorrect:  isCorrect,
				GotScore:   gotScore,
			})
		}
		ReturnJson(ctx, "apiOK", gin.H{
			"score":       score,
			"totalScore":  totalScore,
			"correct":     correctCount,
			"wrong":       wrongCount,
			"durationSec": from.DurationSec,
			"review":      buildReview(session, qMap, answers, wrongKeys),
		})
	})

	r.POST("/leaderboard", func(ctx *gin.Context) {
		isAuth, _, data := AuthenticationAuthority(ctx)
		if !isAuth {
			ReturnJson(ctx, "userCookieError", nil)
			return
		}
		type FromLB struct {
			BankID uint `json:"bankId"`
			Limit  int  `json:"limit"`
		}
		var from FromLB
		decodeJSON(data, &from)
		if from.Limit <= 0 || from.Limit > 50 {
			from.Limit = 10
		}
		type lbRow struct {
			UserID   uint `gorm:"column:user_id"`
			Best     int  `gorm:"column:best"`
			Sessions int  `gorm:"column:sessions"`
		}
		var rows []lbRow
		query := models.DB.Model(&TabQuizSession{}).Where("status = 1")
		if from.BankID > 0 {
			query = query.Where("bank_id = ?", from.BankID)
		}
		query.Select("user_id, MAX(score) AS best, COUNT(*) AS sessions").
			Group("user_id").
			Order("best DESC").
			Limit(from.Limit).
			Scan(&rows)
		userIDs := []uint{}
		for _, r := range rows {
			userIDs = append(userIDs, r.UserID)
		}
		nameMap, avatarMap := quizUserMaps(userIDs)
		items := []gin.H{}
		for i, r := range rows {
			lastAt := ""
			lastQuery := models.DB.Where("user_id = ? AND status = 1", r.UserID)
			if from.BankID > 0 {
				lastQuery = lastQuery.Where("bank_id = ?", from.BankID)
			}
			var last TabQuizSession
			if lastQuery.Order("id DESC").First(&last).Error == nil && last.CreatedAt != nil {
				lastAt = last.CreatedAt.Format("2006-01-02 15:04:05")
			}
			items = append(items, gin.H{
				"rank":     i + 1,
				"userId":   r.UserID,
				"name":     nameMap[r.UserID],
				"avatar":   avatarMap[r.UserID],
				"best":     r.Best,
				"sessions": r.Sessions,
				"lastAt":   lastAt,
			})
		}
		ReturnJson(ctx, "apiOK", gin.H{"items": items})
	})
}
