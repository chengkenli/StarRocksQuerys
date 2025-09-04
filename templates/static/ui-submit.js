/* submit task */
function fetchData_submit() {
  return apiFetch("/submit").then((data) => {
    renderTable_submit(data.data);
    updateCounter_submit(data);
    return data;
  });
}

/* 刷新submit task 数据 */
window.fetchData_submit = function () {
  return apiFetch("/submit")
    .then((data) => {
      renderTable_submit(data.data);
      updateCounter_submit(data);
      return data;
    })
    .catch((error) => {
      console.error("获取数据失败:", error);
    });
};

/* 解析submit task数据 */
function renderTable_submit(data) {
  console.debug("submit task:", data);
  const tableBody = document.getElementById("submit-data");
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
    const progressInt = parseInt(item.ProgressInt, 10);
    if (progressInt >= 99) {
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

    // 格式化item.STATE
    let State = "";
    if (item.STATE === "RUNNING") {
      State = ` <span class="badge bg-warning text-white">${item.STATE}</span>`;
    } else if (item.STATE === "PENDING") {
      State = ` <span class="badge bg-info text-white">${item.STATE}</span>`;
    } else if (item.STATE === "MERGED") {
      State = ` <span class="badge bg-success text-white">${item.STATE}</span>`;
    } else {
      State = item.STATE;
    }

    row.innerHTML = `
					<td class="py-1">${index + 1}</td>
                    <td class="py-1" data-tooltip="${item.QueryIdName}">${item.QUERY_ID}</td>
                    <td class="py-1" data-tooltip="${item.DEFINITION}">${item.TASK_NAME}</td>
                    <td class="py-1" data-tooltip="${item.STATE}">${State}</td>
                    <td class="py-1" data-tooltip="${item.DATABASE}">${item.DATABASE}</td>
                    <td class="py-1" data-tooltip="${item.PROGRESS}">${progress}</td>
                    <td class="py-1" data-tooltip="${item.TimeMsg}">${item.GetHour}</td>
					<td class="py-1" data-tooltip="cancel">${item.Command}</td>
                `;

    tableBody.appendChild(row);
  });
}

// 动态加载指标
function updateCounter_submit({ running, pending, merged, total }) {
  document.getElementById("submit-total").textContent = "TOTAL:" + total;
  document.getElementById("submit-running").textContent = "RUNNING:" + running;
  document.getElementById("submit-pending").textContent = "PENDING:" + pending;
  document.getElementById("submit-merged").textContent = "MERGED:" + merged;

  // 更新submit task标签的徽章
  const badge = document.getElementById("submit-count-badge");
  if (total > 0) {
    badge.textContent = total;
    badge.style.display = "inline-block"; // 显示徽章
  } else {
    badge.style.display = "none"; // 隐藏徽章（如果total为0）
  }
}

// cancel submit task
function cancelSubmitTask() {
  document.addEventListener("click", async function (event) {
    const killButton = event.target.closest("#cancel-task");
    if (killButton) {
      event.preventDefault();
      const task = killButton.getAttribute("data-task");
      const selectedIds = [task];
      // 使用美观的确认对话框
      const isConfirmed = await bsConfirm(`cancel task ${task}？`);
      if (isConfirmed) {
        // 创建并显示加载指示器
        const loadingIndicator = createLoadingIndicator();
        document.body.appendChild(loadingIndicator);
        // 禁用按钮防止重复点击
        killButton.disabled = true;
        try {
          const data = await apiFetch("/cancel-task", "POST", { taskname: selectedIds });
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
