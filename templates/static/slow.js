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
let right2RefreshInterval = null;
let isRight2PanelActive = false;

/**
 * 初始化右侧面板
 */
function initRight2Panel() {
  if (!document.getElementById("slow-data")) {
    console.warn("右侧面板容器未找到，初始化中止");
    return;
  }

  isRight2PanelActive = true;
  console.log("正在初始化右侧资源面板...");

  // 初始加载数据
  fetchResource2Data()
    .then(() => console.log("右侧面板初始化完成"))
    .catch((error) => console.error("右侧面板初始化失败:", error));
}

/**
 * 设置右侧面板自动刷新
 * @param {number} interval - 刷新间隔(毫秒)，默认5000
 */
function setupRight2AutoRefresh(interval = 5000) {
  if (!isRight2PanelActive) {
    console.warn("右侧面板未激活，请先调用initRight2Panel()");
    return;
  }

  // 清除现有定时器
  if (right2RefreshInterval) {
    clearInterval(right2RefreshInterval);
    console.log("已清除现有右侧刷新定时器");
  }

  // 设置新定时器
  right2RefreshInterval = setInterval(async () => {
    try {
      await fetchResource2Data();
    } catch (error) {
      console.error("右侧自动刷新失败:", error);
      // 错误时尝试重新连接
      if (error.status === 401) {
        try {
          await verifyToken();
          await fetchResource2Data();
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
async function fetchResource2Data() {
  if (!isRight2PanelActive) return;

  try {
    console.debug("正在获取资源数据...");
    const startTime = performance.now();

    const data = await apiFetch("/slowquery");
    renderResource2Table(data);

    const duration = (performance.now() - startTime).toFixed(2);
    console.debug(`慢查询资源数据获取完成，耗时 ${duration}ms`);
    console.log("传递的集群名称", getAppID());
    return data;
  } catch (error) {
    console.error("获取资源数据失败:", error);
    showErrorToast("获取资源数据失败", error);
    throw error;
  }
}

/**
 * 渲染资源表格
 * @param {Array} data - 资源数据数组
 */
function renderResource2Table(data) {
  const tbody = document.getElementById("slow-data");
  if (!tbody) return;

  // 使用文档片段提高性能
  const fragment = document.createDocumentFragment();

  if (!data || data.length === 0) {
    const row = document.createElement("tr");
    row.innerHTML = `<td colspan="8" class="text-center text-muted">没有可用的资源数据</td>`;

    updateClear(0);
    updateWarn(0);
    updateKill(0);
    updateTotal(0);

    fragment.appendChild(row);
  } else {
    data.forEach((item, index) => {
      const row = document.createElement("tr");

      updateClear(item.clear);
      updateWarn(item.warn);
      updateKill(item.kill);
      updateTotal(item.clear + item.warn + item.kill);

      let stateText;
      if (item.state === "clear") {
        stateText = '<span class="badge bg-success">clear</span>';
      } else if (item.state === "warn") {
        stateText = '<span class="badge bg-warning">warn</span>';
      } else if (item.state === "kill") {
        stateText = '<span class="badge bg-danger">kill</span>';
      } else {
        stateText = '<span class="badge bg-light">total</span>';
      }

      row.innerHTML = `
				<td>${index + 1}</td>
                <td>${item.id}</td>
                <td>${item.name}</td>
                <td>${item.starttime}</td>
                <td>${item.exectime}</td>
                <td>${item.type}</td>
                <td>${stateText}</td>
            `;

      fragment.appendChild(row);
    });
  }

  // 一次性更新DOM
  tbody.innerHTML = "";
  tbody.appendChild(fragment);
}

/**
 * 清理右侧面板资源
 */
function cleanupRight2Panel() {
  if (right2RefreshInterval) {
    clearInterval(right2RefreshInterval);
    right2RefreshInterval = null;
  }
  isRight2PanelActive = false;
  console.log("右侧面板资源已清理");
}

function updateClear(count) {
  document.getElementById("slow-clear").textContent = "clear:" + count;
}
function updateWarn(count) {
  document.getElementById("slow-warn").textContent = "warn:" + count;
}
function updateKill(count) {
  document.getElementById("slow-kill").textContent = "kill:" + count;
}
function updateTotal(count) {
  document.getElementById("slow-total").textContent = "total:" + count;
}

// 导出公共接口
window.rightPanel = {
  init: initRight2Panel,
  setupRefresh: setupRight2AutoRefresh,
  cleanup: cleanupRight2Panel,
  refresh: fetchResource2Data,
};
