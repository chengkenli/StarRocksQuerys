function initErrorReport() {
  // 获取DOM元素
  const errorReportBtn = document.getElementById("errorReportBtn");
  const modalOverlay = document.getElementById("modalOverlay");
  const closeModalBtn = document.getElementById("closeModal");

  // 定义定时器变量
  let refreshTimer = null;

  // 打开模态框
  errorReportBtn.addEventListener("click", function () {
    modalOverlay.style.display = "flex";
    const tableBody = document.getElementById("err-data");
    tableBody.innerHTML = `<tr><td colspan="6" class="text-center py-2 text-muted">正在加载报错信息...</td></tr>`;
    fetchErrorData(); // 立即加载一次数据

    // 启动定时器，每60秒刷新一次
    refreshTimer = setInterval(fetchErrorData, 600000); // 60秒刷新一次
  });

  // 关闭模态框
  closeModalBtn.addEventListener("click", function () {
    modalOverlay.style.display = "none";
    // 关闭模态框时清除定时器
    if (refreshTimer) {
      clearInterval(refreshTimer);
      refreshTimer = null;
    }
  });

  // 点击模态框外部关闭
  modalOverlay.addEventListener("click", function (e) {
    if (e.target === modalOverlay) {
      modalOverlay.style.display = "none";
      // 关闭模态框时清除定时器
      if (refreshTimer) {
        clearInterval(refreshTimer);
        refreshTimer = null;
      }
    }
  });

  // 封装fetchErrorData
  async function fetchErrorData() {
    try {
      console.debug("正在加载报错信息...");

      const startTime = performance.now();
      const data = await apiFetch("/showerr", "POST");
      renderShowerr(data);

      const duration = (performance.now() - startTime).toFixed(2);
      console.debug(`加载报错信息完成，耗时 ${duration}ms`);
    } catch (error) {
      const tableBody = document.getElementById("err-data");
      tableBody.innerHTML = `<tr><td colspan="6" class="text-center py-2 text-danger">加载失败: ${error.message}</td></tr>`;
      console.error("加载报错信息失败:", error);

      // 出错时也清除定时器
      if (refreshTimer) {
        clearInterval(refreshTimer);
        refreshTimer = null;
      }
    }
  }
}

// 显示报错详情（使用表格形式）
function renderShowerr(data) {
  console.debug("加载报错信息执行中...");
  const tableBody = document.getElementById("err-data");
  tableBody.innerHTML = ""; // 清空现有内容

  console.debug(data);
  if (!data || data.length === 0) {
    const row = document.createElement("tr");
    row.innerHTML = `<td colspan="6" class="text-center py-2 text-muted">没有查询数据</td>`;
    tableBody.appendChild(row);
    return;
  }

  // 添加行数据
  data.forEach((item, index) => {
    const row = document.createElement("tr");

    // 处理错误消息，超过80字符截断并添加省略号
    let errMsg = item.errmsg || "无";
    if (errMsg.length > 120) {
      errMsg = errMsg.substring(0, 120) + "...";
    }

    // 处理stmt内容，用于tooltip显示
    let stmtContent = item.stmt || "无";
    if (stmtContent.length > 200) {
      // 限制tooltip显示的内容长度
      stmtContent = stmtContent.substring(0, 200) + "...";
    }

    let errmsgDisplay;
    if (item.category == 1) {
      errmsgDisplay = ' <span class="badge bg-white text-danger">平台侧</span>';
    } else if (item.category == 2) {
      errmsgDisplay = ' <span class="badge bg-white text-success">用户侧</span>';
    } else {
      errmsgDisplay = ' <span class="badge bg-white text-warning">未可知</span>';
    }

    row.innerHTML = `
      <td class="py-1" data-tooltip="未开放...">${index + 1}</td>
      <td class="py-1" data-tooltip="未开放...">${item.starttime || "无"}</td>
      <td class="py-1" data-tooltip="未开放...">${item.user || "无"}</td>
      <td class="py-1" data-tooltip="${stmtContent}">${item.queryid || "无"}</td>
      <td class="py-1" style="color: red;" data-tooltip="${item.errmsg || "无"}">${errmsgDisplay}${errMsg}</td>
    `;
    tableBody.appendChild(row);
    console.debug("数据渲染完毕!");
  });
}
