async function apiFetch(endpoint, method = "GET", body = null) {
  try {
    const response = await fetch(`/api${endpoint}`, {
      method,
      headers: {
        "Content-Type": "application/json",
        Authorization: getToken(), // 获取当前token
        AppID: getAppID(), // 获取当前AppID
        Begintime: begintime,
        Endtime: endtime,
      },
      body: body ? JSON.stringify(body) : null,
    });

    if (response.status === 401) {
      // 提示用户token无效
      alertMsg("Token无效或已过期，请重新登录", "warning");
      // 清除当前无效token
      clearToken();
      // 重定向到登录页面或执行其他登录流程
      window.location.reload(true); // 强制刷新当前页面（不缓存）
      // 抛出错误终止后续执行
      throw new Error("Unauthorized: Token无效");
    }

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    throw error; // 重新抛出错误以便上层处理
  }
}

function clearToken() {
  // 清除token
  localStorage.removeItem("token");
  API_TOKEN = "";
}
