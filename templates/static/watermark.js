// 创建水印元素并保持引用
const watermarkElement = document.createElement("div");
watermarkElement.style.cssText = `
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 9999;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 40px;
  color: rgba(0,0,0,0.1);
  transform: rotate(-30deg);
  font-weight: bold;
`;
document.body.appendChild(watermarkElement);

// 实时更新函数
function updateWatermark() {
  watermarkElement.textContent = API_APPID; // 直接同步最新值
}

// 启动更新检测（每秒同步一次）
setInterval(updateWatermark, 1000);

// 初始更新
updateWatermark();
