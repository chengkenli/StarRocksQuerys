function renderTable(data) {
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
    if (item.user === "cndlopsns" || item.user === "root") {
      userDisplay += ' <span class="badge bg-light text-success">管理员</span>';
    }
    // 处理短查询保障
    if (item.shortlist) {
      userDisplay += ' <span class="badge bg-white text-dark">短查询保障</span>';
    }
    // 处理白名单显示
    if (item.whitelist) {
      userDisplay += ' <span class="badge bg-info text-white">白名单</span>';
    }
    // 处理追踪用户名
    let clientDisplay = item.host;
    if (item.clientuser !== null && item.clientuser !== undefined) {
      clientDisplay += ` <span class="badge bg-white text-dark">${item.clientuser}</span>`;
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
                    <td class="py-1" data-tooltip="${item.ctxip}">${clientDisplay}</td>   
                    <td class="py-1">${item.db}</td>
                    <td class="py-1 ${timeClass}" data-tooltip="${item.gethour}">${timeText}</td> 
                    <td class="py-1">${item.state}</td>
                    <td class="py-1 text-center">${item.isPending ? '<i class="bi bi-hourglass-top text-danger"></i>' : '<i class="bi bi-check-circle text-success"></i>'}</td>
                    <td class="py-1">${item.warehouse}</td>
					<td class="py-1">${item.command}</td>
                `;

    tableBody.appendChild(row);
  });
}

// 动态加载指标
function updateCounters({ run, sleep, pend, count, fe }) {
  document.getElementById("query-run").textContent = "running:" + run;
  document.getElementById("query-sleep").textContent = "sleep:" + sleep;
  document.getElementById("query-pend").textContent = "pending:" + pend;
  document.getElementById("query-count").textContent = "total:" + count;
  document.querySelector(".monitor-title .tooltip").textContent = fe;
}

function initializeUIEvents() {
  //fetchData(); // 初始加载数据
  //setInterval(fetchData, 5000); // 设置自动刷新

  // 清理睡眠按钮事件 - 修改为使用apiFetch
  document.getElementById("kill-info").addEventListener("click", function () {
    if (confirm("(◔‿◔)确定要清理所有睡眠连接吗？")) {
      apiFetch("/kills", "POST")
        .then((data) => {
          alert(data.message);
          fetchData(); // 刷新数据
        })
        .catch((error) => {
          console.error("Error:", error);
          alert("清理失败");
        });
    }
  });

  // 清理全体语句按钮事件 - 修改为使用apiFetch
  document.getElementById("kill-all").addEventListener("click", function () {
    if (confirm("(◔‿◔)高危！确定要清退所有语句吗？")) {
      apiFetch("/killall", "POST")
        .then((data) => {
          alert(data.message);
          fetchData(); // 刷新数据
        })
        .catch((error) => {
          console.error("Error:", error);
          alert("清理失败");
        });
    }
  });

  // 清理选中的语句按钮事件 - 修改为使用apiFetch
  document.getElementById("kill-one").addEventListener("click", function () {
    const selectedIds = [];
    document.querySelectorAll(".query-checkbox:checked").forEach((checkbox) => {
      selectedIds.push(checkbox.getAttribute("data-id"));
    });

    if (selectedIds.length === 0) {
      alert("请先选择需要清退的语句ID~");
      return;
    }
    if (confirm("提醒！确定要清理勾选的语句吗？")) {
      apiFetch("/killone", "POST", { ids: selectedIds })
        .then((data) => {
          alert(data.message);
          fetchData(); // 刷新数据
        })
        .catch((error) => {
          console.error("Error:", error);
          alert("清理失败");
        });
    }
  });

  // 查杀慢查询按钮事件 - 修改为使用apiFetch
  document.getElementById("kill-slow").addEventListener("click", function () {
    if (confirm("(◔‿◔)确定要查杀所有慢查询吗？")) {
      apiFetch("/killw", "POST")
        .then((data) => {
          alert(data.message);
          fetchData(); // 刷新数据
        })
        .catch((error) => {
          console.error("Error:", error);
          alert("查杀失败");
        });
    }
  });

  // 查杀高消耗查询按钮事件 - 修改为使用apiFetch
  document.getElementById("kill-critical").addEventListener("click", function () {
    if (confirm("(◔‿◔)确定要查杀所有高消耗查询吗？")) {
      apiFetch("/killc", "POST")
        .then((data) => {
          alert(data.message);
          fetchData(); // 刷新数据
        })
        .catch((error) => {
          console.error("Error:", error);
          alert("查杀失败");
        });
    }
  });

  // 退出按钮
  document.getElementById("exit-info").addEventListener("click", function () {
    if (confirm("没有开发这个功能~")) {
      window.open("", "_self").close();
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
function setupKillButtonHandler() {
  // 使用事件委托方式监听click事件
  document.addEventListener("click", function (event) {
    // 检查是否点击的是kill-connectid按钮或其后代元素
    const killButton = event.target.closest("#kill-connectid");
    if (killButton) {
      // 阻止默认行为（如需要）
      event.preventDefault();

      const connectId = killButton.getAttribute("data-connectid");
      const selectedIds = [connectId];

      if (confirm(`(◔‿◔)确定要清理连接 ${connectId} 吗？`)) {
        apiFetch("/killone", "POST", { ids: selectedIds })
          .then((data) => {
            console.info(data.message);
            alert(data.message);
            if (typeof fetchData === "function") {
              fetchData(); // 刷新数据
            }
          })
          .catch((error) => {
            console.error("Error:", error);
            alert("清理失败：" + error.message);
          });
      }
    }
  });
}
