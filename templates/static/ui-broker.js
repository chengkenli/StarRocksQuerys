/* broker load */
function fetchData_brokerload() {
  return apiFetch("/broker").then((data) => {
    renderTable_brokerload(data.data);
    updateCounter_brokerload(data);
    return data;
  });
}

/* 刷新broker load 数据 */
window.fetchData_brokerload = function () {
  return apiFetch("/broker")
    .then((data) => {
      renderTable_brokerload(data.data);
      updateCounter_brokerload(data);
      return data;
    })
    .catch((error) => {
      console.error("获取数据失败:", error);
    });
};

/* 解析broker load数据 */
function renderTable_brokerload(data) {
  console.debug("broker load:", data);
  const tableBody = document.getElementById("broker-data");
  tableBody.innerHTML = "";

  if (!data || data.length === 0) {
    const row = document.createElement("tr");
    row.innerHTML = `<td colspan="9" class="text-center py-2 text-muted">没有查询数据</td>`;
    tableBody.appendChild(row);
    return;
  }
  data.forEach((item, index) => {
    const row = document.createElement("tr");

    // 格式化item.Progress
    let progress = "";
    const progressInt = parseInt(item.Progress, 10);
    if (progressInt === 99) {
      progress = ` <span class="badge bg-success">${progressInt}%</span>`;
    } else if (progressInt >= 80 && progressInt < 99) {
      progress = ` <span class="badge bg-primary text-white">${progressInt}%</span>`;
    } else if (progressInt >= 50 && progressInt < 80) {
      progress = ` <span class="badge bg-info text-dark">${progressInt}%</span>`;
    } else if (progressInt >= 1 && progressInt < 50) {
      progress = ` <span class="badge bg-warning text-dark">${progressInt}%</span>`;
    } else if (progressInt === 0) {
      progress = ` <span class="badge bg-danger text-white">${progressInt}%</span>`;
    } else {
      progress = `${progressInt}`;
    }

    let state = "";
    state = `<span class="badge text-warning-emphasis">${item.State}</span>`;

    // 正文
    row.innerHTML = `
					<td class="py-1">${index + 1}</td>
                    <td class="py-1" data-tooltip="${item.LabelName}">${item.Label}</td>
                    <td class="py-1" data-tooltip="${item.State}">${state}</td>
                    <td class="py-1" data-tooltip="${item.ProgressMsg}">${progress}</td>
                    <td class="py-1" data-tooltip="${item.Type}">${item.Type}</td>
                    <td class="py-1" data-tooltip="${item.ScanInfo}">${item.ScanRows}</td>
                    <td class="py-1" data-tooltip="${item.JobDetails}">${item.Backends}</td>
                    <td class="py-1" data-tooltip="${item.TimeMsg}">${item.GetHour}</td>
					<td class="py-1" data-tooltip="cancel">${item.Command}</td>
                `;

    tableBody.appendChild(row);
  });
}
// 动态加载指标
function updateCounter_brokerload({ loading, pending, queueing, prepared, total }) {
  document.getElementById("broker-total").textContent = "TOTAL:" + total;
  document.getElementById("broker-loading").textContent = "RUNNING:" + loading;
  document.getElementById("broker-pending").textContent = "PENDING:" + pending;
  document.getElementById("broker-queueing").textContent = "QUEUEING:" + queueing;
  document.getElementById("broker-prepared").textContent = "PREPARED:" + prepared;

  // 更新broker load标签的徽章
  const badge = document.getElementById("broker-count-badge");
  if (total > 0) {
    badge.textContent = total;
    badge.style.display = "inline-block"; // 显示徽章
  } else {
    badge.style.display = "none"; // 隐藏徽章（如果total为0）
  }
}

// cancel broker load
function cancelBrokerLabel() {
  document.addEventListener("click", async function (event) {
    const killButton = event.target.closest("#cancel-label");
    if (killButton) {
      event.preventDefault();
      const label = killButton.getAttribute("data-label");
      const selectedIds = [label];
      // 使用美观的确认对话框
      const isConfirmed = await bsConfirm(`cancel label ${label}？`);
      if (isConfirmed) {
        // 创建并显示加载指示器
        const loadingIndicator = createLoadingIndicator();
        document.body.appendChild(loadingIndicator);
        // 禁用按钮防止重复点击
        killButton.disabled = true;
        try {
          const data = await apiFetch("/cancel-labelId", "POST", { labelId: selectedIds });
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
