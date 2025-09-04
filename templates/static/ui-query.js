/* query */
function fetchData() {
  console.log("query fetchdata", getAppID());
  return apiFetch("/query").then((data) => {
    renderTable(data.data);
    updateCounters(data);
    return data;
  });
}
// api.js - 确保函数是全局的
window.fetchData = function () {
  return apiFetch("/query")
    .then((data) => {
      renderTable(data.data);
      updateCounters(data);
      return data;
    })
    .catch((error) => {
      console.error("获取数据失败:", error);
    });
};

function renderTable(data) {
  console.debug("query global:", data);
  const tableBody = document.getElementById("query-data");
  tableBody.innerHTML = "";

  if (!data || data.length === 0) {
    const row = document.createElement("tr");
    row.innerHTML = `<td colspan="9" class="text-center py-2 text-muted">没有查询数据</td>`;
    tableBody.appendChild(row);
    return;
  }

  data.forEach((item, index) => {
    const row = document.createElement("tr");
    let timeClass = "";
    let timeText = item.time;
    if (item.time > 1500) {
      timeClass = "critical-query-cell";
      timeText += ' <span class="badge bg-danger">高消耗</span>';
    } else if (item.time > 600) {
      timeClass = "slow-query-cell";
      timeText += ' <span class="badge bg-warning text-dark">慢查询</span>';
    } else if (item.time > 300) {
      timeClass = "slow-query-cell";
      timeText += ' <span class="badge bg-info text-dark">慢查询</span>';
    } else if (item.time === 0 && item.user === "cndlopsns") {
      timeClass = "slow-query-cell";
      timeText += ' <span class="badge bg-info text-white">connection</span>';
    }
    // 处理用户角色显示
    let userDisplay = item.user;
    // 处理管理员显示
    if (item.admin) {
      userDisplay += ' <span class="badge bg-white text-success">A</span>';
    }
    // 处理二级重点账号显示
    if (item.point) {
      userDisplay += ' <span class="badge bg-white text-secondary">P0</span>';
    }
    // 处理短查询保障
    if (item.shortlist) {
      userDisplay += ' <span class="badge bg-white text-dark">S+P0</span>';
    }
    // 处理白名单显示
    if (item.whitelist) {
      userDisplay += ' <span class="badge bg-white text-primary">W</span>';
    }
    // 处理追踪用户名
    let clientDisplay = item.host;
    if (item.clientuser !== null && item.clientuser !== undefined) {
      clientDisplay += `<span class="badge bg-white text-dark border border-gray-300 px-2 py-1">${item.clientuser}</span>`;
    }
    // 处理stmt内容，用于tooltip显示
    let stmtContent = item.info || "无";
    if (stmtContent.length > 200) {
      // 限制tooltip显示的内容长度
      stmtContent = stmtContent.substring(0, 200) + "...";
    }

    // <td class="py-1">${userDisplay}</td>
    row.innerHTML = `
					<td class="py-1">${index + 1}</td>
                    <td class="py-1" data-tooltip="${stmtContent}">${item.id}</td>
					<td class="py-1" data-tooltip="${item.feip}">${userDisplay}</td>
					<td class="py-1" data-tooltip="${item.OperationalMsg}">${item.Operational}</td>
                    <td class="py-1" data-tooltip="${item.ctxip}">${clientDisplay}</td>   
                    <td class="py-1 ${timeClass}" data-tooltip="${item.gethour}">${timeText}</td> 
                    <td class="py-1 text-center">${
                      item.isPending
                        ? '<i data-tooltip="pending" class="bi bi-hourglass-top text-danger"></i>'
                        : '<i data-tooltip="running" class="bi bi-check-circle text-success"></i>'
                    }</td>
                    <td class="py-1" data-tooltip="disconnect ${item.user} all connections">${item.warehouse}</td>
					<td class="py-1" data-tooltip="kill Pid">${item.command}</td>
                `;

    tableBody.appendChild(row);
  });
}

// 动态加载指标
function updateCounters({ run, sleep, pend, count, fe }) {
  document.getElementById("query-run").textContent = "RUNNING:" + run;
  document.getElementById("query-sleep").textContent = "SLEEP:" + sleep;
  document.getElementById("query-pend").textContent = "PENDING:" + pend;
  document.getElementById("query-count").textContent = "TOTAL:" + count;
  document.querySelector(".monitor-title .tooltip").textContent = fe;

  // 更新query global标签的徽章
  const badge = document.getElementById("query-count-badge");
  if (run > 0) {
    badge.textContent = run;
    badge.style.display = "inline-block"; // 显示徽章
  } else {
    badge.style.display = "none"; // 隐藏徽章（如果run为0）
  }
}

function initializeUIEvents() {
  // 清理睡眠按钮事件 - 修改为使用apiFetch
  document.getElementById("kill-info").addEventListener("click", async function () {
    const isConfirmed = await bsConfirm(`Kill all Sleep ConnectionId？`);
    if (isConfirmed) {
      const loadingIndicator = createLoadingIndicator();
      document.body.appendChild(loadingIndicator);
      // 禁用按钮防止重复点击
      const killAllButton = this;
      killAllButton.disabled = true;
      try {
        const data = await apiFetch("/kills", "POST");
        alertMsg(data.message);
        if (typeof fetchData === "function") fetchData(); // 刷新数据
      } catch (error) {
        console.error("Error:", error);
        alertMsg("清理失败", "danger");
      } finally {
        // 无论成功或失败，都移除加载指示器并恢复按钮状态
        loadingIndicator.remove();
        killAllButton.disabled = false;
      }
    }
  });

  // 清理全体语句按钮事件 - 修改为使用apiFetch
  document.getElementById("kill-all").addEventListener("click", async function () {
    const isConfirmed = await bsConfirm(`Kill all ConnectionId？`);
    if (isConfirmed) {
      // 创建并显示加载指示器
      const loadingIndicator = createLoadingIndicator();
      document.body.appendChild(loadingIndicator);
      // 禁用按钮防止重复点击
      const killAllButton = this;
      killAllButton.disabled = true;
      try {
        const data = await apiFetch("/killall", "POST");
        alertMsg(data.message);
        if (typeof fetchData === "function") fetchData(); // 刷新数据
      } catch (error) {
        console.error("Error:", error);
        alertMsg("清理失败", "danger");
      } finally {
        // 无论成功或失败，都移除加载指示器并恢复按钮状态
        loadingIndicator.remove();
        killAllButton.disabled = false;
      }
    }
  });

  // 清理选中的语句按钮事件 - 修改为使用apiFetch
  document.getElementById("kill-one").addEventListener("click", async function () {
    const selectedIds = [];
    document.querySelectorAll(".query-checkbox:checked").forEach((checkbox) => {
      selectedIds.push(checkbox.getAttribute("data-id"));
    });

    if (selectedIds.length === 0) {
      alertMsg("请先选择需要清退的语句ID~", "warning");
      return;
    }
    const isConfirmed = await bsConfirm(`提醒！确定要清理勾选的语句吗？`);

    if (isConfirmed) {
      const loadingIndicator = createLoadingIndicator();
      document.body.appendChild(loadingIndicator);
      // 禁用按钮防止重复点击
      const killAllButton = this;
      killAllButton.disabled = true;
      try {
        const data = await apiFetch("/killone", "POST");
        alertMsg(data.message);
        if (typeof fetchData === "function") fetchData(); // 刷新数据
      } catch (error) {
        console.error("Error:", error);
        alertMsg("清理失败", "danger");
      } finally {
        // 无论成功或失败，都移除加载指示器并恢复按钮状态
        loadingIndicator.remove();
        killAllButton.disabled = false;
      }
    }
  });

  // 查杀慢查询按钮事件 - 修改为使用apiFetch
  document.getElementById("kill-slow").addEventListener("click", async function () {
    const isConfirmed = await bsConfirm(`Kill all WarnQuery ConnectionId？`);

    if (isConfirmed) {
      const loadingIndicator = createLoadingIndicator();
      document.body.appendChild(loadingIndicator);
      // 禁用按钮防止重复点击
      const killAllButton = this;
      killAllButton.disabled = true;
      try {
        const data = await apiFetch("/killw", "POST");
        alertMsg(data.message);
        if (typeof fetchData === "function") fetchData(); // 刷新数据
      } catch (error) {
        console.error("Error:", error);
        alertMsg("清理失败", "danger");
      } finally {
        // 无论成功或失败，都移除加载指示器并恢复按钮状态
        loadingIndicator.remove();
        killAllButton.disabled = false;
      }
    }
  });

  // 查杀高消耗查询按钮事件 - 修改为使用apiFetch
  document.getElementById("kill-critical").addEventListener("click", async function () {
    const isConfirmed = await bsConfirm(`Kill all SlowQuery ConnectionId？`);
    if (isConfirmed) {
      const loadingIndicator = createLoadingIndicator();
      document.body.appendChild(loadingIndicator);
      // 禁用按钮防止重复点击
      const killAllButton = this;
      killAllButton.disabled = true;
      try {
        const data = await apiFetch("/killc", "POST");
        alertMsg(data.message);
        if (typeof fetchData === "function") fetchData(); // 刷新数据
      } catch (error) {
        console.error("Error:", error);
        alertMsg("清理失败", "danger");
      } finally {
        // 无论成功或失败，都移除加载指示器并恢复按钮状态
        loadingIndicator.remove();
        killAllButton.disabled = false;
      }
    }
  });
}

// connection id 跳转地址
function fetchAndRedirect(url) {
  console.debug(url);
  fetch(url, {
    method: "GET",
    redirect: "manual", // 防止浏览器自动跟随重定向
    credentials: "include", // 如果需要携带 Cookie
  })
    .then((response) => {
      console.debug(response);
      if (response.redirected) {
        // 如果服务器返回了 302 或 301 重定向
        // window.location.href = response.url;
        window.open(response.url, "_blank");
      } else if (response.ok) {
        // 如果服务器返回的是 JSON 数据（如 {"url": "..."}）
        return response.json();
      } else {
        console.error("Request failed", response.statusText);
        return Promise.reject("Request failed");
      }
    })
    .then((data) => {
      console.debug(data);
      // 假设返回的 JSON 数据包含 { url: "..." }
      if (data && data.url) {
        // window.location.href = data.url;
        window.open(data.url, "_blank");
      }
      window.open(data, "_blank");
    })
    .catch((error) => {
      console.error("Error:", error);
      // 可选：回退到手动点击链接
      window.open(url, "_blank");
    });
}

// 1. 定义事件处理函数
function setupDisconnectButtonHandler() {
  // 使用事件委托方式监听click事件
  document.addEventListener("click", async function (event) {
    // 检查是否点击的是kill-connectid按钮或其后代元素
    const killButton = event.target.closest("#kill-disconnect");
    if (killButton) {
      // 阻止默认行为（如需要）
      event.preventDefault();

      const username = killButton.getAttribute("data-disconnect");
      const userid = [username];

      const isConfirmed = await bsConfirm(`Kill user ${userid}？`);

      if (isConfirmed) {
        const loadingIndicator = createLoadingIndicator();
        document.body.appendChild(loadingIndicator);
        // 禁用按钮防止重复点击
        const killAllButton = this;
        killAllButton.disabled = true;
        try {
          const data = await apiFetch("/disconnect", "POST", { user: userid });
          alertMsg(data.message);
          if (typeof fetchData === "function") fetchData(); // 刷新数据
        } catch (error) {
          console.error("Error:", error);
          alertMsg("清理失败", "danger");
        } finally {
          // 无论成功或失败，都移除加载指示器并恢复按钮状态
          loadingIndicator.remove();
          killAllButton.disabled = false;
        }
      }
    }
  });
}

function setupKillButtonHandler() {
  document.addEventListener("click", async function (event) {
    const killButton = event.target.closest("#kill-connectid");
    if (killButton) {
      event.preventDefault();
      const connectId = killButton.getAttribute("data-connectid");
      const selectedIds = [connectId];
      // 使用美观的确认对话框
      const isConfirmed = await bsConfirm(`Kill connectId ${connectId}？`);
      if (isConfirmed) {
        // 创建并显示加载指示器
        const loadingIndicator = createLoadingIndicator();
        document.body.appendChild(loadingIndicator);
        // 禁用按钮防止重复点击
        killButton.disabled = true;
        try {
          const data = await apiFetch("/killone", "POST", { ids: selectedIds });
          alertMsg(data.message, "success");
          if (typeof fetchData === "function") fetchData();
        } catch (error) {
          alertMsg("清理失败：" + error.message, "danger");
        } finally {
          // 无论成功或失败，都移除加载指示器并恢复按钮状态
          loadingIndicator.remove();
          killButton.disabled = false;
        }
      }
    }
  });
}

// 创建加载指示器的辅助函数
function createLoadingIndicator() {
  const loader = document.createElement("div");
  loader.style.position = "fixed";
  loader.style.top = "0";
  loader.style.left = "0";
  loader.style.width = "100%";
  loader.style.height = "5px";
  loader.style.backgroundColor = "#ffff00";
  loader.style.zIndex = "9999";
  loader.style.animation = "loadingAnimation 2s infinite";
  // 添加动画样式
  const style = document.createElement("style");
  style.textContent = `
    @keyframes loadingAnimation {
      0% { width: 0%; left: 0; }
      50% { width: 100%; left: 0; }
      100% { width: 0%; left: 100%; }
    }
  `;
  document.head.appendChild(style);

  return loader;
}

/**
 * 显示美观的提示消息
 * @param {string} message - 要显示的消息内容
 * @param {string} [type='info'] - 类型：'info', 'success', 'danger', 'warning'
 */
function alertMsg(message, type = "info") {
  const toastEl = document.getElementById("liveToast");
  const toastMessage = document.getElementById("toastMessage");
  // 设置消息内容
  toastMessage.textContent = message;
  // 根据类型设置背景色（Bootstrap原生颜色类）
  toastEl.className = `toast bg-${type} text-white`;
  // 初始化Toast并显示
  const toast = new bootstrap.Toast(toastEl, {
    autohide: true, // 自动隐藏
    delay: 3000, // 3秒后消失
  });
  toast.show();
}

/**
 * 美观的确认对话框（替换原生confirm）
 * @param {string} message 提示消息
 * @returns {Promise<boolean>} 返回Promise，用户点击确定时resolve(true)
 */
function bsConfirm(message) {
  return new Promise((resolve) => {
    const modal = new bootstrap.Modal(document.getElementById("confirmModal"));
    const confirmMsg = document.getElementById("confirmMessage");
    const confirmBtn = document.getElementById("confirmAction");
    // 设置消息内容
    confirmMsg.textContent = message;
    // 清除旧事件监听（避免重复绑定）
    const cleanUp = () => {
      confirmBtn.removeEventListener("click", confirmHandler);
      modal._element.removeEventListener("hidden.bs.modal", cancelHandler);
    };
    // 确认处理
    const confirmHandler = () => {
      cleanUp();
      modal.hide();
      resolve(true);
    };
    // 取消处理
    const cancelHandler = () => {
      cleanUp();
      resolve(false);
    };
    // 绑定新事件
    confirmBtn.addEventListener("click", confirmHandler);
    modal._element.addEventListener("hidden.bs.modal", cancelHandler);
    // 显示对话框
    modal.show();
  });
}
