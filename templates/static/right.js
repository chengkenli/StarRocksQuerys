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

      // 格式化数据
      const scanBytes = formatBytes(item.ScanBytes);
      const scanRows = formatNumber(item.ScanRows);
      const memoryUsage = formatBytes(item.MemoryUsage);

      // 根据资源消耗添加警告类和标签
      let scanBytesDisplay = item.ScanBytes;
      if (item.ScanBytes > 1 * 1024 * 1024 * 1024 * 1024) {
        // >1TB
        scanBytesDisplay += ' <span class="badge bg-danger">1TB</span>';
      } else if (item.ScanBytes > 536870912000) {
        // >500GB
        scanBytesDisplay += ' <span class="badge bg-warning">500GB</span>';
      }

      let scanRowsDisplay = item.ScanRows;
      if (item.ScanRows > 10000000000) {
        // >100亿
        scanRowsDisplay += ' <span class="badge bg-danger">100亿</span>';
      } else if (item.ScanRows > 5000000000) {
        // >50亿
        scanRowsDisplay += ' <span class="badge bg-warning">50亿</span>';
      }

      let memoryUsageDisplay = item.MemoryUsage;
      if (item.MemoryUsage > 200 * 1024 * 1024 * 1024) {
        // >200GB
        memoryUsageDisplay += ' <span class="badge bg-danger">200GB</span>';
        row.classList.add("high-memory-usage");
      }

      if (item.CPUTime > 60) {
        // >60秒
        row.classList.add("high-cpu-usage");
      }

      row.innerHTML = `
			    <td>${index + 1}</td>
                <td>${item.ID}</td>
                <td>${scanBytesDisplay}</td>
                <td>${scanRowsDisplay}</td>
                <td>${memoryUsageDisplay}</td>
                <td>${item.CPUTime}</td>
                <td>${item.ExecTime}</td>
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
