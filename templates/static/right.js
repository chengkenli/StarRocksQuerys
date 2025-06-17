/**
 * right.js - 数据库查询监控系统右侧面板功能模块
 *
 * 功能：
 * 1. 初始化右侧资源监控面板
 * 2. 定时获取并显示查询资源消耗数据
 * 3. 提供资源数据格式化显示
 * 4. 生命周期管理
 */

// 模块私有变量
let rightRefreshInterval = null;
let isRightPanelActive = false;

/**
 * 初始化右侧面板
 */
function initRightPanel() {
  if (!document.getElementById("resource-data")) {
    console.warn("右侧面板容器未找到，初始化中止");
    return;
  }

  isRightPanelActive = true;
  console.log("正在初始化右侧资源面板...");

  // 初始加载数据
  fetchResourceData()
    .then(() => console.log("右侧面板初始化完成"))
    .catch((error) => console.error("右侧面板初始化失败:", error));
}

/**
 * 设置右侧面板自动刷新
 * @param {number} interval - 刷新间隔(毫秒)，默认5000
 */
function setupRightAutoRefresh(interval = 5000) {
  if (!isRightPanelActive) {
    console.warn("右侧面板未激活，请先调用initRightPanel()");
    return;
  }

  // 清除现有定时器
  if (rightRefreshInterval) {
    clearInterval(rightRefreshInterval);
    console.log("已清除现有右侧刷新定时器");
  }

  // 设置新定时器
  rightRefreshInterval = setInterval(async () => {
    try {
      await fetchResourceData();
    } catch (error) {
      console.error("右侧自动刷新失败:", error);
      // 错误时尝试重新连接
      if (error.status === 401) {
        try {
          await verifyToken();
          await fetchResourceData();
        } catch (reauthError) {
          console.error("重新认证失败:", reauthError);
        }
      }
    }
  }, interval);

  console.log(`右侧面板自动刷新已设置: 每 ${interval / 1000} 秒刷新一次`);
}

/**
 * 获取资源数据并更新UI
 */
async function fetchResourceData() {
  if (!isRightPanelActive) return;

  try {
    console.debug("正在获取资源数据...");
    const startTime = performance.now();

    const data = await apiFetch("/queries");
    renderResourceTable(data);

    const duration = (performance.now() - startTime).toFixed(2);
    console.debug(`队列资源数据获取完成，耗时 ${duration}ms`);
    return data;
  } catch (error) {
    console.error("获取资源数据失败:", error);
    showErrorToast("获取资源数据失败");
    throw error;
  }
}

/**
 * 渲染资源表格
 * @param {Array} data - 资源数据数组
 */
function renderResourceTable(data) {
  const tbody = document.getElementById("resource-data");
  if (!tbody) return;

  // 使用文档片段提高性能
  const fragment = document.createDocumentFragment();

  if (!data || data.length === 0) {
    const row = document.createElement("tr");
    row.innerHTML = `<td colspan="8" class="text-center text-muted">没有可用的资源数据</td>`;
    fragment.appendChild(row);
  } else {
    data.forEach((item, index) => {
      const row = document.createElement("tr");

      let scanBytesDisplay = item.ScanBytes;
      let scanBytesNumber = 0;
      let unit = "B"; // 默认单位
      // 提取单位和数值
      if (typeof scanBytesDisplay === "string") {
        const match = scanBytesDisplay.match(/(\d+\.?\d*)\s*([KMGTP]?B)/i);
        if (match) {
          unit = match[2].toUpperCase(); // 提取单位并转为大写
          scanBytesNumber = parseInt(match[1], 10);
        } else {
          scanBytesNumber = parseInt(scanBytesDisplay, 10) || 0;
        }
      } else {
        scanBytesNumber = scanBytesDisplay;
      }
      scanBytesDisplay = scanBytesNumber;
      // 根据单位添加标签
      switch (unit) {
        case "TB":
          scanBytesDisplay += ' <span class="badge bg-danger">TB</span>';
          break;
        case "GB":
          scanBytesDisplay += ' <span class="badge bg-warning">GB</span>';
          break;
        case "MB":
          scanBytesDisplay += ' <span class="badge bg-info">MB</span>';
          break;
        case "KB":
          scanBytesDisplay += ' <span class="badge bg-primary">KB</span>';
          break;
        case "B":
          scanBytesDisplay += ' <span class="badge bg-secondary">B</span>';
          break;
        default:
          // 未知单位不添加标签
          break;
      }

      let scanRowsDisplay = item.ScanRows;
      let rowsNumber = 0;
      // 提取数字部分
      if (typeof scanRowsDisplay === "string") {
        const match = scanRowsDisplay.match(/^(\d+)\s*rows?$/i);
        if (match) {
          rowsNumber = parseInt(match[1], 10);
        } else {
          // 如果不是"数字 rows"格式，尝试直接转换为数字
          rowsNumber = parseInt(scanRowsDisplay, 10) || 0;
        }
      } else {
        // 如果不是字符串，直接使用数值
        rowsNumber = scanRowsDisplay;
      }
      scanRowsDisplay = rowsNumber;
      // 根据数值大小添加标签
      if (rowsNumber > 10000000000) {
        // >100亿
        scanRowsDisplay += ' <span class="badge bg-danger">100亿+</span>';
      } else if (rowsNumber > 5000000000) {
        // >50亿
        scanRowsDisplay += ' <span class="badge bg-warning">50亿+</span>';
      } else if (rowsNumber > 1000000000) {
        // >10亿
        scanRowsDisplay += ' <span class="badge bg-info">10亿+</span>';
      } else if (rowsNumber > 100000000) {
        // >1亿
        scanRowsDisplay += ' <span class="badge bg-primary">1亿+</span>';
      } else if (rowsNumber > 10000000) {
        // >1000万
        scanRowsDisplay += ' <span class="badge bg-success">1000万+</span>';
      } else if (rowsNumber > 1000000) {
        // >100万
        scanRowsDisplay += ' <span class="badge bg-secondary">100万+</span>';
      }

      let memoryUsageDisplay = item.MemoryUsage;
      let memoryUsageNumber = 0;
      let units = "B"; // 默认单位
      // 提取单位和数值
      if (typeof memoryUsageDisplay === "string") {
        const match = memoryUsageDisplay.match(/(\d+\.?\d*)\s*([KMGTP]?B)/i);
        if (match) {
          units = match[2].toUpperCase(); // 提取单位并转为大写
          memoryUsageNumber = parseInt(match[1], 10);
        } else {
          memoryUsageNumber = parseInt(memoryUsageDisplay, 10) || 0;
        }
      } else {
        memoryUsageNumber = memoryUsageDisplay;
      }
      memoryUsageDisplay = memoryUsageNumber;
      // 根据单位添加标签
      switch (units) {
        case "TB":
          memoryUsageDisplay += ' <span class="badge bg-danger">TB</span>';
          break;
        case "GB":
          memoryUsageDisplay += ' <span class="badge bg-warning">GB</span>';
          break;
        case "MB":
          memoryUsageDisplay += ' <span class="badge bg-info">MB</span>';
          break;
        case "KB":
          memoryUsageDisplay += ' <span class="badge bg-primary">KB</span>';
          break;
        case "B":
          memoryUsageDisplay += ' <span class="badge bg-secondary">B</span>';
          break;
        default:
          // 未知单位不添加标签
          break;
      }

      let cpuTimeDisplay = item.CPUTime;
      let seconds = 0;
      // 提取秒数
      if (typeof cpuTimeDisplay === "string") {
        // 处理"835.167 s"或"85 s"这样的格式
        const match = cpuTimeDisplay.match(/(\d+\.?\d*)\s*s/i);
        if (match) {
          seconds = parseFloat(match[1]);
        } else {
          // 如果不是"数字 s"格式，尝试直接转换为数字
          seconds = parseFloat(cpuTimeDisplay) || 0;
        }
      } else {
        // 如果不是字符串，直接使用数值
        seconds = cpuTimeDisplay;
      }
      // 转换为hour:min:sec格式
      const hours = Math.floor(seconds / 3600);
      const mins = Math.floor((seconds % 3600) / 60);
      const secs = Math.floor(seconds % 60);
      // 格式化时间显示
      let formattedTime = "";
      if (hours > 0) {
        formattedTime = `${hours}hour:${mins.toString().padStart(2, "0")}min:${secs.toString().padStart(2, "0")}s`;
      } else if (mins > 0) {
        formattedTime = `${mins}min:${secs.toString().padStart(2, "0")}s`;
      } else {
        formattedTime = `${secs}s`;
      }
      // 根据时长添加不同的Bootstrap样式
      let timeClass = "";
      if (seconds > 3600) {
        // >1小时
        timeClass = "bg-danger";
      } else if (seconds > 1800) {
        // >30分钟
        timeClass = "bg-warning";
      } else if (seconds > 300) {
        // >5分钟
        timeClass = "bg-info";
      } else if (seconds > 60) {
        // >1分钟
        timeClass = "bg-primary";
      } else {
        // ≤1分钟
        timeClass = "bg-success";
      }
      // 更新显示内容（只显示带样式的标签）
      cpuTimeDisplay = `<span class="badge ${timeClass}">${formattedTime}</span>`;

      let execTimeDisplay = item.ExecTime;
      let seconds1 = 0;
      // 提取秒数
      if (typeof execTimeDisplay === "string") {
        // 处理"835.167 s"或"85 s"这样的格式
        const match = execTimeDisplay.match(/(\d+\.?\d*)\s*s/i);
        if (match) {
          seconds1 = parseFloat(match[1]);
        } else {
          // 如果不是"数字 s"格式，尝试直接转换为数字
          seconds1 = parseFloat(execTimeDisplay) || 0;
        }
      } else {
        // 如果不是字符串，直接使用数值
        seconds1 = execTimeDisplay;
      }
      // 转换为hour:min:sec格式
      const hours1 = Math.floor(seconds1 / 3600);
      const mins1 = Math.floor((seconds1 % 3600) / 60);
      const secs1 = Math.floor(seconds1 % 60);
      // 格式化时间显示
      let formattedTime1 = "";
      if (hours1 > 0) {
        formattedTime1 = `${hours1}hour:${mins1.toString().padStart(2, "0")}min:${secs1.toString().padStart(2, "0")}s`;
      } else if (mins1 > 0) {
        formattedTime1 = `${mins1}min:${secs1.toString().padStart(2, "0")}s`;
      } else {
        formattedTime1 = `${secs1}s`;
      }
      // 根据时长添加不同的Bootstrap样式
      let timeClass1 = "";
      if (seconds1 > 3600) {
        // >1小时
        timeClass1 = "bg-danger";
      } else if (seconds1 > 1800) {
        // >30分钟
        timeClass1 = "bg-warning";
      } else if (seconds1 > 300) {
        // >5分钟
        timeClass1 = "bg-info";
      } else if (seconds1 > 60) {
        // >1分钟
        timeClass1 = "bg-primary";
      } else {
        // ≤1分钟
        timeClass1 = "bg-success";
      }
      // 更新显示内容（只显示带样式的标签）
      execTimeDisplay = `<span class="badge ${timeClass1}">${formattedTime1}</span>`;

      row.innerHTML = `
			    <td>${index + 1}</td>
                <td data-tooltip="${item.StartTime}">${item.ConnectionId}</td>
                <td>${scanBytesDisplay}</td>
                <td data-tooltip="${item.QueryId}">${scanRowsDisplay}</td>
                <td>${memoryUsageDisplay}</td>
                <td>${cpuTimeDisplay}</td>
                <td>${execTimeDisplay}</td>
            `;

      fragment.appendChild(row);
    });
  }

  // 一次性更新DOM
  tbody.innerHTML = "";
  tbody.appendChild(fragment);
}

// 辅助函数：格式化字节为易读单位
function formatBytes(bytes, decimals = 2) {
  if (bytes === 0) return "0 Bytes";
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
}

// 辅助函数：格式化数字为中文单位
function formatNumber(num) {
  if (num >= 100000000) {
    return (num / 100000000).toFixed(1) + "亿";
  }
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + "万";
  }
  return num.toString();
}

/**
 * 显示错误提示
 */
function showErrorToast(message) {
  // 这里可以替换为您的UI通知系统
  console.error("错误提示:", message);
  alert(`错误: ${message}`);
}

/**
 * 清理右侧面板资源
 */
function cleanupRightPanel() {
  if (rightRefreshInterval) {
    clearInterval(rightRefreshInterval);
    rightRefreshInterval = null;
  }
  isRightPanelActive = false;
  console.log("右侧面板资源已清理");
}

// 导出公共接口
window.rightPanel = {
  init: initRightPanel,
  setupRefresh: setupRightAutoRefresh,
  cleanup: cleanupRightPanel,
  refresh: fetchResourceData,
};
