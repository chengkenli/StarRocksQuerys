// 全局变量
let begintime, endtime;
let timeUpdateInterval;
let isModalOpen = false;

// 设置默认值 - 当前时间前30分钟和当前时间
function setDefaultTimes() {
  if (isModalOpen) return; // 如果模态框打开，不更新

  const now = new Date();
  endtime = formatDateTime(now);
  begintime = formatDateTime(new Date(now.getTime() - 30 * 60 * 1000)); // 30分钟前

  console.log("时间已更新:", { begintime, endtime });
}

// 格式化日期为 yyyy-mm-dd 时:分:秒
function formatDateTime(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  const seconds = String(date.getSeconds()).padStart(2, "0");

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

// 解析日期时间字符串为Date对象
function parseDateTime(datetimeStr) {
  const [datePart, timePart] = datetimeStr.split(" ");
  const [year, month, day] = datePart.split("-").map(Number);
  const [hours, minutes, seconds] = timePart.split(":").map(Number);

  return new Date(year, month - 1, day, hours, minutes, seconds);
}

// 验证日期时间格式
function isValidDateTime(datetimeStr) {
  const regex = /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/;
  if (!regex.test(datetimeStr)) return false;

  try {
    const date = parseDateTime(datetimeStr);
    return !isNaN(date.getTime());
  } catch (e) {
    return false;
  }
}

// 初始化默认时间并设置定时器
function initTimeUpdater() {
  setDefaultTimes(); // 立即设置一次
  timeUpdateInterval = setInterval(setDefaultTimes, 5000); // 每5秒更新一次
}

// 停止时间更新
function stopTimeUpdater() {
  if (timeUpdateInterval) {
    clearInterval(timeUpdateInterval);
    timeUpdateInterval = null;
  }
}

// 页面加载时初始化
document.addEventListener("DOMContentLoaded", function () {
  initTimeUpdater();

  document.getElementById("datetime-info").addEventListener("click", function () {
    // 停止自动更新时间
    stopTimeUpdater();
    isModalOpen = true;

    // 创建模态框
    const modal = document.createElement("div");
    modal.style.position = "fixed";
    modal.style.top = "0";
    modal.style.left = "0";
    modal.style.width = "100%";
    modal.style.height = "100%";
    modal.style.backgroundColor = "rgba(0,0,0,0.5)";
    modal.style.display = "flex";
    modal.style.justifyContent = "center";
    modal.style.alignItems = "center";
    modal.style.zIndex = "1000";

    // 模态框内容
    modal.innerHTML = `
      <div style="background: white; padding: 20px; border-radius: 5px; width: 300px;">
        <h4>选择时间范围</h4>
        <div style="margin-bottom: 15px;">
          <label>开始时间 (yyyy-mm-dd 时:分:秒):</label>
          <input type="text" id="start-time" style="width: 100%;" placeholder="例如: 2023-01-01 08:00:00">
        </div>
        <div style="margin-bottom: 15px;">
          <label>结束时间 (yyyy-mm-dd 时:分:秒):</label>
          <input type="text" id="end-time" style="width: 100%;" placeholder="例如: 2023-01-01 18:00:00">
        </div>
        <div style="display: flex; justify-content: space-between;">
          <button id="cancel-btn" style="padding: 5px 10px;">取消</button>
          <button id="confirm-btn" style="padding: 5px 10px;">确认</button>
        </div>
      </div>
    `;

    document.body.appendChild(modal);

    // 设置默认值
    const startTimeInput = document.getElementById("start-time");
    const endTimeInput = document.getElementById("end-time");

    startTimeInput.value = begintime;
    endTimeInput.value = endtime;

    // 取消按钮事件
    document.getElementById("cancel-btn").addEventListener("click", function () {
      document.body.removeChild(modal);
      isModalOpen = false;
      initTimeUpdater(); // 重新开始自动更新
    });

    // 确认按钮事件
    document.getElementById("confirm-btn").addEventListener("click", function () {
      const startValue = startTimeInput.value.trim();
      const endValue = endTimeInput.value.trim();

      if (!startValue || !endValue) {
        alert("请填写开始和结束时间");
        return;
      }

      if (!isValidDateTime(startValue) || !isValidDateTime(endValue)) {
        alert("时间格式不正确，请使用 yyyy-mm-dd 时:分:秒 格式");
        return;
      }

      const startDate = parseDateTime(startValue);
      const endDate = parseDateTime(endValue);

      if (endDate <= startDate) {
        alert("结束时间必须晚于开始时间");
        return;
      }

      begintime = startValue;
      endtime = endValue;

      document.body.removeChild(modal);
      isModalOpen = false;
      // 用户选择后不再自动更新时间

      console.log("用户选择的时间:", { begintime, endtime });
    });
  });
});
