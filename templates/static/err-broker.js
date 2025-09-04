function initBrokerError() {
  // 获取DOM元素
  const errorBroker = document.getElementById("errorBroker");
  const modalOverlay = document.getElementById("modalOverlay");
  const closeModalBtn = document.getElementById("closeModal");

  // 定义定时器变量
  let refreshTimer = null;

  // 打开模态框
  errorBroker.addEventListener("click", function () {
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
      const data = await apiFetch("/broker-err", "POST");
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
