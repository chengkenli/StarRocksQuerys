// bead.js - 悬浮珠功能实现
function initBeadMonitor() {
  const bead = document.getElementById("bead");
  if (!bead) {
    console.warn("悬浮珠元素未找到");
    return;
  }

  let isDragging = false,
    offsetX,
    offsetY;

  // 拖动逻辑
  bead.onmousedown = function (e) {
    isDragging = true;
    const rect = bead.getBoundingClientRect();
    offsetX = e.clientX - rect.left;
    offsetY = e.clientY - rect.top;
    bead.style.transition = "none";
    e.preventDefault();
  };

  document.onmousemove = function (e) {
    if (!isDragging) return;
    let x = e.clientX - offsetX;
    let y = e.clientY - offsetY;

    // 限制不超出视窗
    x = Math.max(0, Math.min(x, window.innerWidth - bead.offsetWidth));
    y = Math.max(0, Math.min(y, window.innerHeight - bead.offsetHeight));

    bead.style.left = x + "px";
    bead.style.top = y + "px";
    bead.style.right = "auto";
    bead.style.bottom = "auto";
  };

  document.onmouseup = function () {
    if (!isDragging) return;
    isDragging = false;
    bead.style.transition = "all 0.4s cubic-bezier(0.18, 0.89, 0.32, 1.28)";

    const rect = bead.getBoundingClientRect();
    const centerX = rect.left + bead.offsetWidth / 2;
    const centerY = rect.top + bead.offsetHeight / 2;

    // 计算到各边的距离
    const distances = {
      left: centerX,
      right: window.innerWidth - centerX,
      top: centerY,
      bottom: window.innerHeight - centerY,
    };

    // 找出最近的边
    const closestEdge = Object.keys(distances).reduce(function (a, b) {
      return distances[a] < distances[b] ? a : b;
    });

    // 根据最近的边设置位置
    switch (closestEdge) {
      case "left":
        bead.style.left = "0";
        bead.style.right = "auto";
        break;
      case "right":
        bead.style.left = "auto";
        bead.style.right = "0";
        break;
      case "top":
        bead.style.top = "0";
        bead.style.bottom = "auto";
        break;
      case "bottom":
        bead.style.top = "auto";
        bead.style.bottom = "0";
        break;
    }
  };

  // 数据更新函数
  function updateBeadData() {
    apiFetch("/beads")
      .then((data) => {
        document.getElementById("errquery").textContent = "QUERY:" + data.errquery;
        document.getElementById("errbroker").textContent = "BROKER:" + data.errbroker;
        document.getElementById("errsubmit").textContent = "SUBMIT:" + data.errsubmit;
      })
      .catch((error) => {
        console.error("获取数据失败:", error);
        document.getElementById("errquery").textContent = "Failed";
        document.getElementById("errbroker").textContent = "Failed";
        document.getElementById("errsubmit").textContent = "Failed";
      });
  }

  // 初始化数据并设置定时器
  updateBeadData();
  var beadInterval = setInterval(updateBeadData, 5000);

  // 确保悬浮珠初始位置正确
  bead.style.left = "auto";
  bead.style.right = "10px";
  bead.style.top = "10px";
  bead.style.bottom = "auto";

  // 返回清理函数
  return function () {
    clearInterval(beadInterval);
    document.onmousemove = null;
    document.onmouseup = null;
  };
}
