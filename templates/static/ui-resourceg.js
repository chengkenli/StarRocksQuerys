/* resource group */
function fetchData_resourcegroup() {
  return apiFetch("/resourcegroup").then((data) => {
    renderTable_resourcegroup(data.data);
    updateCounter_resourcegroup(data);
    return data;
  });
}

/* 刷新resource group 数据 */
window.fetchData_resourcegroup = function () {
  return apiFetch("/resourcegroup")
    .then((data) => {
      renderTable_resourcegroup(data.data);
      updateCounter_resourcegroup(data);
      return data;
    })
    .catch((error) => {
      console.error("获取数据失败:", error);
    });
};

/* 解析resource group数据 */
function renderTable_resourcegroup(data) {
  console.debug("resource group:", data);
  const tableBody = document.getElementById("resourcegroup-data");
  tableBody.innerHTML = "";

  if (!data || data.length === 0) {
    const row = document.createElement("tr");
    row.innerHTML = `<td colspan="9" class="text-center py-2 text-muted">没有查询数据</td>`;
    tableBody.appendChild(row);
    return;
  }
  data.forEach((item, index) => {
    const row = document.createElement("tr");

    // 正文
    row.innerHTML = `
					<td class="py-1">${index + 1}</td>
                    <td class="py-1" data-tooltip="未设置">${item.Name}</td>
                    <td class="py-1" data-tooltip="未设置">${item.User}</td>
                    <td class="py-1" data-tooltip="权重:${item.CpuWeight}/硬隔离:${item.ExclusiveCpuCores}">${item.CpuWeight}/${item.ExclusiveCpuCores}</td>
                    <td class="py-1" data-tooltip="未设置">${item.MemLimit}</td>
                    <td class="py-1" data-tooltip="未设置">${item.ConcurrencyLimit}</td>
                    <td class="py-1" data-tooltip="未设置">${item.BigQueryCPU}</td>
                    <td class="py-1" data-tooltip="未设置">${item.BigQueryRows}</td>
					<td class="py-1" data-tooltip="未设置">${item.BigQueryMemLimit}</td>
					<td class="py-1" data-tooltip="未设置">${item.Command}</td>
                `;

    tableBody.appendChild(row);
  });
}
// 动态加载指标
function updateCounter_resourcegroup({ total }) {
  document.getElementById("resourcegroup-total").textContent = "TOTAL:" + total;

  // 更新resource group标签的徽章
  const badge = document.getElementById("resourcegroup-count-badge");
  if (total > 0) {
    badge.textContent = total;
    badge.style.display = "inline-block"; // 显示徽章
  } else {
    badge.style.display = "none"; // 隐藏徽章（如果total为0）
  }
}
