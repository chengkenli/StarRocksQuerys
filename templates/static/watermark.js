// 创建水印元素
const watermarkElement = document.createElement("div");
watermarkElement.id = "global-watermark"; // 添加ID便于管理

// 设置水印样式（已覆盖全屏）
watermarkElement.style.cssText = `
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;    /* 使用视口单位确保全覆盖 */
  height: 100vh;
  pointer-events: none;
  z-index: 9999;
  display: flex;
  flex-wrap: wrap;  /* 允许水印换行 */
  align-content: flex-start; /* 从顶部开始排列 */
  font-size: 24px;
  color: rgba(0,0,0,0.1);
  transform: rotate(-30deg);
  font-weight: bold;
  overflow: hidden;
`;

// 创建水印文本节点（重复多个实现全屏覆盖）
function createWatermarkContent(text) {
  const container = document.createElement("div");
  container.style.cssText = `
    width: 300px;
    height: 200px;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-shrink: 0;
  `;

  const span = document.createElement("span");
  span.textContent = text;
  container.appendChild(span);

  return container;
}

// 填充水印（根据屏幕尺寸计算需要的水印数量）
function fillWatermark() {
  // 清空现有水印
  watermarkElement.innerHTML = "";

  // 计算屏幕需要的水印数量
  const cols = Math.ceil(window.innerWidth / 300) + 1;
  const rows = Math.ceil(window.innerHeight / 200) + 1;
  const total = cols * rows;

  // 添加水印
  for (let i = 0; i < total; i++) {
    watermarkElement.appendChild(createWatermarkContent(API_APPID));
  }
}

// 初始设置
document.body.appendChild(watermarkElement);
fillWatermark();

// 窗口大小改变时重新填充
window.addEventListener("resize", fillWatermark);

// 定时更新内容（如果需要）
setInterval(() => {
  fillWatermark();
}, 1000);
