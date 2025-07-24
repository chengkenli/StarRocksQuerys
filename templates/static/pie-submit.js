/**
 * 初始化饼图模态框功能
 */
function initSubmitPie() {
  // 1. 获取DOM元素 - 确保只获取特定的饼图按钮
  const pieSubmit = document.getElementById("pieSubmit");
  if (!pieSubmit) {
    console.warn("未找到ID为'pieSubmit'的按钮");
    return; // 如果按钮不存在则退出
  }

  const modalOverlay = document.getElementById("pie-modalOverlay");
  const closeModalBtn = document.getElementById("pie-closeModal");
  const detailsContainer = document.getElementById("pie-DetailsContainer");

  // 2. 定义状态变量
  let refreshTimer = null;
  let pieChart = null;
  let isModalOpen = false;

  // 3. 为饼图按钮添加专属点击处理
  pieSubmit.addEventListener("click", handlePieButtonClick);

  // 4. 关闭按钮事件
  closeModalBtn.addEventListener("click", closeModal);

  // 5. 点击模态框外部关闭
  modalOverlay.addEventListener("click", function (e) {
    if (e.target === modalOverlay) {
      closeModal();
    }
  });

  /**
   * 处理饼图按钮点击事件
   * @param {Event} e
   */
  function handlePieButtonClick(e) {
    e.stopPropagation(); // 阻止事件冒泡
    e.preventDefault(); // 防止默认行为

    if (isModalOpen) {
      closeModal();
      return;
    }

    openPieModal();
  }

  /**
   * 打开饼图模态框
   */
  function openPieModal() {
    isModalOpen = true;
    modalOverlay.style.display = "flex";

    // 修改这里：移除固定尺寸，改用CSS控制
    detailsContainer.innerHTML = `
    <div class="pie-container">
      <div class="item" id="pieChartContainer"></div>
    </div>
    <div class="loading-indicator">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">加载中...</span>
      </div>
    </div>
  `;

    // 初始化饼图
    fetchAndRenderPieChart();
    refreshTimer = setInterval(fetchAndRenderPieChart, 60000);
  }

  /**
   * 关闭模态框
   */
  function closeModal() {
    isModalOpen = false;
    modalOverlay.style.display = "none";

    // 清除定时器
    if (refreshTimer) {
      clearInterval(refreshTimer);
      refreshTimer = null;
    }

    // 销毁饼图实例
    if (pieChart) {
      pieChart.dispose();
      pieChart = null;
    }
  }

  /**
   * 从接口获取数据并渲染饼图
   */
  async function fetchAndRenderPieChart() {
    try {
      console.debug("正在加载饼图数据...");
      const startTime = performance.now();

      // 显示加载状态
      const loadingElement = detailsContainer.querySelector(".loading-indicator");
      if (loadingElement) {
        loadingElement.style.display = "block";
      }

      // 使用apiFetch从接口获取数据
      const pieData = await apiFetch("/pie", "POST");
      console.debug("API返回数据:", pieData);

      // 验证数据格式
      if (!pieData || !Array.isArray(pieData.data)) {
        throw new Error("无效的饼图数据格式");
      }

      // 初始化或获取饼图实例
      if (!pieChart) {
        const container = document.getElementById("pieChartContainer");
        if (!container) {
          throw new Error("找不到饼图容器");
        }

        pieChart = echarts.init(container, "shine", { renderer: "canvas" });

        // 窗口大小变化时重新调整饼图大小
        window.addEventListener("resize", function () {
          if (pieChart) {
            pieChart.resize();
          }
        });
      }

      // 饼图配置
      const pieOption = {
        animation: true,
        legend: {
          right: "0",
          orient: "vertical",
          data: pieData.categories || pieData.data.map((item) => item.name),
        },
        series: [
          {
            name: pieData.title || "NULL",
            type: "pie",
            roseType: "radius",
            radius: ["40%", "70%"],
            center: ["30%", "50%"],
            data: pieData.data,
            itemStyle: {
              borderColor: "#fff",
              borderWidth: 2,
            },
            label: {
              show: true,
              color: "auto",
              position: "outside",
              formatter: "{b}\n{d}%",
            },
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowColor: "rgba(0, 0, 0, 0.5)",
              },
            },
          },
        ],
        title: {
          text: pieData.title || "NULL",
          subtext: pieData.subtext || "数据来源: StarRocks",
          left: "center",
        },
        toolbox: {
          feature: {
            saveAsImage: {
              title: "保存为图片",
              pixelRatio: 2,
            },
          },
        },
        tooltip: {
          trigger: "item",
          formatter: "{a} <br/>{b}: {c} ({d}%)",
          textStyle: {
            // 直接在JS中定义字体样式
            color: "#ffffff", // 深色字体
          },
        },
      };

      // 渲染饼图
      pieChart.setOption(pieOption);

      // 隐藏加载指示器
      if (loadingElement) {
        loadingElement.style.display = "none";
      }

      const duration = (performance.now() - startTime).toFixed(2);
      console.debug(`饼图加载完成，耗时 ${duration}ms`);
    } catch (error) {
      console.error("加载饼图失败:", error);

      // 显示错误信息
      detailsContainer.innerHTML = `
        <div class="alert alert-danger mt-3">
          <h5>加载饼图数据失败</h5>
          <p>${error.message}</p>
          <button class="btn btn-sm btn-outline-secondary" onclick="initSubmitPie()">重试</button>
        </div>
      `;

      // 出错时清除定时器
      if (refreshTimer) {
        clearInterval(refreshTimer);
        refreshTimer = null;
      }
    }
  }
}
