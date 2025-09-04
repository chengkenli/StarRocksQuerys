// 全局变量
let dashboardInterval = null;

/**
 * 更新仪表盘数据
 * @param {number} cpuValue - CPU使用率 (0-1)
 * @param {number} memoryValue - 内存使用率 (0-1)
 */
function updateDataDashboard(cpuValue, memoryValue, ioValue, diskValue) {
  console.log("dashboard:", { cpuValue, memoryValue, ioValue, diskValue });

  // 获取DOM元素
  const cpuGauge = document.querySelector(".gauge-cpu");
  const cpuFill = cpuGauge.querySelector(".gauge-fill");
  const cpuText = cpuGauge.querySelector(".gauge-value");

  const memoryGauge = document.querySelector(".gauge-memory");
  const memoryFill = memoryGauge.querySelector(".gauge-fill");
  const memoryText = memoryGauge.querySelector(".gauge-value");

  const ioGauge = document.querySelector(".gauge-io");
  const ioFill = ioGauge.querySelector(".gauge-fill");
  const ioText = ioGauge.querySelector(".gauge-value");

  const diskGauge = document.querySelector(".gauge-disk");
  const diskFill = diskGauge.querySelector(".gauge-fill");
  const diskText = diskGauge.querySelector(".gauge-value");

  // 计算百分比
  const cpuPercent = Math.min(100, Math.round(cpuValue * 1000) / 10);
  const memoryPercent = Math.min(100, Math.round(memoryValue * 1000) / 10);
  const ioPercent = Math.min(100, Math.round(ioValue * 1000) / 10);
  const diskPercent = Math.min(100, Math.round(diskValue * 1000) / 10);

  // 更新CPU
  cpuFill.style.background = `
                conic-gradient(
                    ${cpuPercent >= 80 ? "#F44336" : "#4CAF50"} 0% ${cpuPercent}%,
                    rgba(76, 175, 80, 0.15) ${cpuPercent}% 100%
                )`;
  cpuText.textContent = `${cpuPercent}%`;
  cpuGauge.classList.toggle("gauge-warning", cpuPercent >= 80);

  // 更新内存
  memoryFill.style.background = `
                conic-gradient(
                    ${memoryPercent >= 80 ? "#F44336" : "#2196F3"} 0% ${memoryPercent}%,
                    rgba(33, 150, 243, 0.15) ${memoryPercent}% 100%
                )`;
  memoryText.textContent = `${memoryPercent}%`;
  memoryGauge.classList.toggle("gauge-warning", memoryPercent >= 80);

  // 更新IO
  ioFill.style.background = `
                conic-gradient(
                    ${ioPercent >= 80 ? "#F44336" : "#2196F3"} 0% ${ioPercent}%,
                    rgba(33, 150, 243, 0.15) ${ioPercent}% 100%
                )`;
  ioText.textContent = `${ioPercent}%`;
  ioGauge.classList.toggle("gauge-warning", ioPercent >= 80);

  // 更新存储
  diskFill.style.background = `
                conic-gradient(
                    ${diskPercent >= 80 ? "#F44336" : "#2196F3"} 0% ${diskPercent}%,
                    rgba(33, 150, 243, 0.15) ${diskPercent}% 100%
                )`;
  diskText.textContent = `${diskPercent}%`;
  diskGauge.classList.toggle("gauge-warning", diskPercent >= 80);
}

/**
 * 初始化仪表盘
 */
async function initDashboard() {
  try {
    const data = await apiFetch("/getgrafan");
    updateDataDashboard(data.cpu, data.memory, data.io, data.disk);
    return true;
  } catch (error) {
    console.error("dashboard failed:", error);
    return false;
  }
}

/**
 * 设置定时刷新
 * @param {number} interval - 刷新间隔(毫秒)
 */
function setupDashboard(interval = 5000) {
  // 清除现有定时器
  if (dashboardInterval) clearInterval(dashboardInterval);

  // 立即加载一次
  initDashboard();

  // 设置定时刷新
  dashboardInterval = setInterval(async () => {
    try {
      const data = await apiFetch("/getgrafan");
      updateDataDashboard(data.cpu, data.memory, data.io, data.disk);
    } catch (error) {
      console.error("dashboard failed:", error);
    }
  }, interval);
}

/**
 * 清理资源
 */
function cleanupDashboard() {
  if (dashboardInterval) {
    clearInterval(dashboardInterval);
    dashboardInterval = null;
  }
}
