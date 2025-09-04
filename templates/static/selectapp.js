/**
 * 初始化集群选择下拉框
 */
function initClusterSelect() {
  const selectElement = document.getElementById("dynamicSelect");
  if (!selectElement) return;

  // 初始化下拉框状态
  function resetSelectState() {
    selectElement.innerHTML = '<option value="" selected disabled>选择集群...</option>';
  }

  /**
   * 获取集群列表数据
   */
  async function fetchClusterList() {
    try {
      console.debug("get starrocks app list...");
      const startTime = performance.now();

      // 显示加载状态
      selectElement.innerHTML = '<option value="" selected disabled>加载中...</option>';

      const data = await apiFetch("/appids");

      // 清空现有选项
      selectElement.innerHTML = "";

      // 添加默认选项
      const defaultOption = document.createElement("option");
      defaultOption.value = "";
      defaultOption.textContent = "switch starrocks...";
      defaultOption.disabled = true;
      selectElement.appendChild(defaultOption);

      // 填充数据选项
      data.forEach((item) => {
        const option = document.createElement("option");
        option.value = item.value;
        option.textContent = item.text || item.name;
        selectElement.appendChild(option);
      });

      const duration = (performance.now() - startTime).toFixed(2);
      console.debug(`starrocks app list ${duration}ms`);
    } catch (error) {
      console.error("starrocks app list:", error);
      handleFetchError(error);
      throw error;
    }
  }

  /**
   * 处理获取数据失败的情况
   */
  function handleFetchError(error) {
    let errorMessage = "加载失败，点击重试";

    if (error.message.includes("401")) {
      errorMessage = "认证失败，请刷新页面";
    }

    selectElement.innerHTML = `<option value="" selected disabled>${errorMessage}</option>`;

    // 如果不是认证错误，添加重试功能
    if (!error.message.includes("401")) {
      selectElement.addEventListener("click", handleSelectClick, {
        once: true,
      });
    }
  }

  /**
   * 处理下拉框点击事件
   */
  function handleSelectClick() {
    fetchClusterList().catch(() => {});
  }

  /**
   * 处理下拉框变更事件
   */
  function handleSelectChange() {
    const selectedAppId = this.value;
    console.log("change starrocks:", selectedAppId);

    // 如果慢查询面板是活动的，刷新数据
    if (isRight2PanelActive) {
      fetchResource2Data().catch(() => {});
    }
  }

  /**
   * 获取当前选中的集群ID
   */
  function getSelectedAppId() {
    return selectElement.value;
  }

  // 初始化下拉框状态
  resetSelectState();

  // 添加事件监听器
  selectElement.addEventListener("click", handleSelectClick, { once: true });
  selectElement.addEventListener("change", handleSelectChange);

  // 返回公共方法
  return {
    refresh: fetchClusterList,
    getSelectedAppId: getSelectedAppId,
  };
}
