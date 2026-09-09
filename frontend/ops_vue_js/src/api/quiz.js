import { api } from './index'

export const quizApi = {
  /** 新增题库 */
  addBank(data) {
    return api.post('/quiz/bank/add', data)
  },
  /** 修改题库 */
  updateBank(data) {
    return api.post('/quiz/bank/update', data)
  },
  /** 删除题库 */
  deleteBank(id) {
    return api.post('/quiz/bank/delete', { id })
  },
  /** 题库列表(管理) */
  listBanks(params = {}) {
    return api.post('/quiz/banks/list', params)
  },
  /** 用户可用题库列表 */
  availableBanks(params = {}) {
    return api.post('/quiz/banks/available', params)
  },
  /** 学习模式：分页取题目(含答案与解析) */
  learnBank(bankId, page = 1, pageSize = 10) {
    return api.post('/quiz/bank/learn', { bankId, page, pageSize })
  },
  /** 新增题目 */
  addQuestion(data) {
    return api.post('/quiz/question/add', data)
  },
  /** 修改题目 */
  updateQuestion(data) {
    return api.post('/quiz/question/update', data)
  },
  /** 删除题目 */
  deleteQuestion(id) {
    return api.post('/quiz/question/delete', { id })
  },
  /** 题目列表(题库管理) */
  listQuestions(params = {}) {
    return api.post('/quiz/questions/list', params)
  },
  /** 开始答题 */
  startQuiz(bankId) {
    return api.post('/quiz/start', { bankId })
  },
  /** 提交答题 */
  submitQuiz(sessionId, durationSec, answers) {
    return api.post('/quiz/submit', { sessionId, durationSec, answers })
  },
  /** 我的成绩历史 */
  mySessions(params = {}) {
    return api.post('/quiz/sessions', params)
  },
  /** 答题回顾详情 */
  getSession(id) {
    return api.post('/quiz/session', { id })
  },
  /** 错题重做判分(仅判分不落库) */
  redoQuiz(sessionId, durationSec, answers) {
    return api.post('/quiz/redo', { sessionId, durationSec, answers })
  },
  /** 排行榜 */
  leaderboard(params = {}) {
    return api.post('/quiz/leaderboard', params)
  },
}
