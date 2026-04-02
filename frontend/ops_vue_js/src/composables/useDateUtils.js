// composables/useDateUtils.js
import { useI18n } from "vue-i18n";

/**
 * 日期字符串 转 日期对象 (纯日期，无时间)
 * @param {string} dateStr - 格式：YYYY-MM-DD 例如 "2026-04-02"
 * @returns {Date} 日期对象（时间自动设为 00:00:00）
 */
export function strToDate(dateStr) {
  return new Date(dateStr);
}

/**
 * 日期对象 转回 日期字符串 (YYYY-MM-DD)
 * @param {Date} date - 日期对象
 * @returns {string} 格式：YYYY-MM-DD
 */
export function dateToStr(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

// 封装成组合式函数，方便在组件中按需引入
export const useDateUtils = () => {
  const { t } = useI18n();
  return {
    strToDate,
    dateToStr,
    // 可以直接把 FullCalendar 常用的日期处理也封装进来
    /**
     * FullCalendar 显示用：给结束日期 +1 天（包含当天）
     * @param {string | Date} endDate - 真实结束日期
     * @returns {Date} FullCalendar 可用的 end 日期
     */
    toCalendarEnd: (endDate) => {
      const date =
        typeof endDate === "string" ? strToDate(endDate) : new Date(endDate);
      date.setDate(date.getDate() + 1);
      return date;
    },
    /**
     * 提交数据库用：给 FullCalendar end 日期 -1 天（还原真实结束日）
     * @param {Date} calendarEnd - FullCalendar 返回的 end 日期
     * @returns {string} 可存库的 YYYY-MM-DD 字符串
     */
    toRealEnd: (calendarEnd) => {
      const date = new Date(calendarEnd);
      date.setDate(date.getDate() - 1);
      return dateToStr(date);
    },
    /**
     * 获取日期是星期几
     * @param {string | Date} date - YYYY-MM-DD 或 Date 对象
     * @returns {string} 星期一、星期二...星期日
     */
    getI18nWeekday(date) {
      const d = typeof date === "string" ? strToDate(date) : new Date(date);
      const weekKeys = ["sun", "mon", "tue", "wed", "thu", "fri", "sat"];
      const key = weekKeys[d.getDay()];
      return t(`week.${key}`); // 自动切换中英文
    },
  };
};
