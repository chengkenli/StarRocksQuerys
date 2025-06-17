window.setupAutoRefresh = function (interval = 5000) {
  if (window.refreshInterval) {
    clearInterval(window.refreshInterval);
  }
  window.refreshInterval = setInterval(window.fetchData, interval);
};

async function apiFetch(endpoint, method = "GET", body = null) {
  try {
    console.log("api传递的集群名称", getAppID());
    const response = await fetch(`/api${endpoint}`, {
      method,
      headers: {
        "Content-Type": "application/json",
        Authorization: getToken(), // 获取当前token
        AppID: getAppID(), // 获取当前AppID
      },
      body: body ? JSON.stringify(body) : null,
    });

    if (response.status === 401) {
      // 提示用户token无效
      alert("Token无效或已过期，请重新登录");
      // 清除当前无效token
      clearToken();
      // 重定向到登录页面或执行其他登录流程
      window.location.reload(true); // 强制刷新当前页面（不缓存）
      console.debug("强制刷新");
      // 抛出错误终止后续执行
      throw new Error("Unauthorized: Token无效");
    }

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error("API请求失败:", error);
    throw error; // 重新抛出错误以便上层处理
  }
}

function clearToken() {
  // 清除token
  console.debug("当前token:", getToken());
  localStorage.removeItem("token");
  API_TOKEN = "";
  console.debug("清理token:", getToken());
}

function fetchData() {
  console.log("传递的集群名称", getAppID());
  return apiFetch("/query").then((data) => {
    renderTable(data.data); // 调用ui.js中的函数
    updateCounters(data); // 调用ui.js中的函数
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
