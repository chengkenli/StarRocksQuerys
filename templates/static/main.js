// 修改后的主入口函数
async function initializeApp() {
  // 先验证token
  const isTokenValid = await verifyToken();

  console.debug("main:", isTokenValid);
  if (isTokenValid) {
    console.warn("verify token success:", API_TOKEN);

    initBeadMonitor(); //初始化悬浮珠功能模块
    initErrorReport(); //初始化查询报错展示
    initBrokerError(); //初始化broker load报错展示
    initSubmitError(); //初始化submit task报错展示
    initStreamError(); //初始化stream load报错展示
    initSubmitPie(); //初始化查看24小时连接数饼图

    setupKillButtonHandler(); //定义查杀按钮事件委托
    setupDisconnectButtonHandler(); //定义释放连接按钮事件委托

    initDashboard(); //初始化仪表盘
    setupDashboard(5000); //仪表盘定时5秒请求一次

    initializeUIEvents();
    initClusterSelect(); // 初始化集群select标签

    fetchData(); // 左刷新
    setupAutoRefresh(5000); // 设置自动刷新

    initRightPanel(); // 右侧刷新
    setupRightAutoRefresh(5000); // 设置右侧自动刷新

    initRight2Panel(); // 慢查询刷新
    setupRight2AutoRefresh(60000); // 设置慢查询自动刷新

    fetchData_brokerload(); //broker load刷新
    fetchData_submit(); //submit task刷新
    fetchData_resourcegroup(); //resource group 刷新
    cancelBrokerLabel(); //定义cancel broker load事件
    cancelSubmitTask(); //定义cancel submit task事件
    //
  } else {
    console.warn("get token fail:", API_TOKEN);
  }
}

// 修改事件监听
window.addEventListener("DOMContentLoaded", function () {
  initializeApp();
});

// 确保在页面卸载时清理定时器
window.addEventListener("beforeunload", function () {
  // 清理左侧面板定时器
  if (window.refreshInterval) {
    clearInterval(window.refreshInterval);
  }

  // 清理右侧面板定时器
  if (window.rightRefreshInterval) {
    clearInterval(window.rightRefreshInterval);
  }

  // 调用右侧面板的清理函数(如果存在)
  if (typeof cleanupRightPanel === "function") {
    cleanupRightPanel();
  }
});
