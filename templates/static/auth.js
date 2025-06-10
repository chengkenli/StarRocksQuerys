// 全局变量存储Token
/* let API_TOKEN = "";

function verifyToken() {
  let token;
  do {
    token = prompt("请输入访问Token");
    if (token === null) {
      alert("必须输入Token才能访问本系统");
      continue;
    }
    if (token.trim() === "") {
      alert("Token不能为空");
      continue;
    }

    API_TOKEN = token;
    document.getElementById("mainContent").style.display = "block";
    return true;
  } while (true);
} */

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

    if (
      !tokenInput ||
      !errorMessage ||
      !submitBtn ||
      !mainContent ||
      !tokenModal
    ) {
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
          console.log("Token 有效");
          return true;
        } else {
          console.log("Token 无效");
          showError("校验失败");
          return false;
        }
      } catch (error) {
        console.error("请求失败:", error);
        showError("token校验失败");
        return false;
      }
    }

    async function handleTokenSubmission() {
      const token = tokenInput.value.trim();
      console.debug("当前token:", token);

      const valid = await validateToken(token); // 等待异步校验
      if (valid) {
        console.debug("校验通过");
        API_TOKEN = token;
        tokenModal.style.display = "none";
        mainContent.style.display = "block";
        console.log("Token验证成功:", API_TOKEN);
        resolve(true);
      } else {
        console.debug("校验失败");
      }
    }

    /*    function handleTokenSubmission() {
      const token = tokenInput.value.trim();
      console.debug("当前token:", token);
      
      const vaild = validateToken(token);
      console.debug("校验结果:", vaild);
      if (vaild) {
        console.debug("校验通过:", token);
        API_TOKEN = token;
        tokenModal.style.display = "none";
        mainContent.style.display = "block";
        console.log("Token验证成功:", API_TOKEN);
        resolve(true);
      } else {
        console.debug("校验失败:", token);
      }
    } */

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
  console.log("get token:", API_TOKEN);
  return API_TOKEN;
}
