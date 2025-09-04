let API_TOKEN = "";

// 修改后的verifyToken函数，返回Promise以支持异步验证
function verifyToken() {
  return new Promise((resolve) => {
    // 确保DOM元素存在
    const tokenInput = document.getElementById("tokenInput");
    const errorMessage = document.getElementById("errorMessage");
    const submitBtn = document.getElementById("submitBtn");
    const mainContent = document.getElementById("mainContent");
    const tokenModal = document.getElementById("tokenModal");

    if (!tokenInput || !errorMessage || !submitBtn || !mainContent || !tokenModal) {
      console.error("无法找到必要的DOM元素");
      resolve(false);
      return;
    }

    function showError(message) {
      errorMessage.textContent = message;
      errorMessage.style.display = "block";
      tokenInput.style.borderColor = "#e74c3c";
    }

    function hideError() {
      errorMessage.style.display = "none";
      tokenInput.style.borderColor = "#ddd";
    }

    async function validateToken(token) {
      if (!token) {
        showError("token不能为空");
        return false;
      }
      if (token.length < 6) {
        showError("token长度必须至少6个字符");
        return false;
      }

      try {
        // 等待 fetch 完成
        const response = await fetch("/api/verify-token", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ token: token }),
        });
        const data = await response.json();

        if (data.valid === true) {
          return true;
        } else {
          showError("校验失败");
          return false;
        }
      } catch (error) {
        showError("token校验失败");
        return false;
      }
    }

    async function handleTokenSubmission() {
      const token = tokenInput.value.trim();

      const valid = await validateToken(token); // 等待异步校验
      if (valid) {
        API_TOKEN = token;
        tokenModal.style.display = "none";
        mainContent.style.display = "block";
        resolve(true);
      }
    }

    // 清除之前的监听器
    submitBtn.removeEventListener("click", handleTokenSubmission);
    tokenInput.removeEventListener("keypress", handleKeyPress);

    // 点击按钮提交
    submitBtn.addEventListener("click", handleTokenSubmission);

    // 回车键提交
    function handleKeyPress(e) {
      if (e.key === "Enter") {
        handleTokenSubmission();
      }
    }
    tokenInput.addEventListener("keypress", handleKeyPress);

    // 输入时自动隐藏错误
    tokenInput.addEventListener("input", hideError);

    // 自动聚焦输入框
    tokenInput.focus();

    // 显示模态框
    tokenModal.style.display = "flex";
  });
}

function getToken() {
  return API_TOKEN;
}
